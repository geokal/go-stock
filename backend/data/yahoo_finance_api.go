package data

import (
	"encoding/json"
	"fmt"
	"go-stock/backend/db"
	"go-stock/backend/logger"
	"go-stock/backend/models"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// YahooFinance API - no API key required, global coverage including EU/UK stocks
type YahooFinanceApi struct {
	client *resty.Client
}

type YahooKLineData struct {
	Day    string `json:"day"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
}

type YahooChartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Symbol           string  `json:"symbol"`
				RegularMarketPrice   float64 `json:"regularMarketPrice"`
				Currency         string `json:"currency"`
				ExchangeTimezone string `json:"exchangeTimezoneName"`
			} `json:"meta"`
			Timestamp []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []int64   `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error *struct {
			Code    string `json:"code"`
			Message string `json:"description"`
		} `json:"error"`
	} `json:"chart"`
}

func NewYahooFinanceApi() *YahooFinanceApi {
	return &YahooFinanceApi{
		client: resty.New(),
	}
}

// ConvertStockCodeToYahoo converts eu:/uk:/ch: prefix to Yahoo ticker format
// Stocks: eu:ASML -> ASML.AS (Euronext Amsterdam), uk:HSBA -> HSBA.L (LSE), ch:NOVN -> NOVN.S (Swiss)
// Indices: eu:DAX -> ^GDAXI, eu:CAC -> ^FCHI, uk:FTSE -> ^FTSE, ch:SMI -> ^SSMI
func ConvertStockCodeToYahoo(stockCode string) (string, string) {
	code := strings.ToUpper(strings.TrimSpace(stockCode))

	// Handle index codes (all caps, no dots)
	if strings.HasPrefix(code, "EU:") {
		ticker := code[3:]
		// Map common European indices to Yahoo Finance symbols
		switch ticker {
		case "DAX":
			return "^GDAXI", "XETRA"
		case "CAC", "CAC40":
			return "^FCHI", "Euronext"
		case "AEX":
			return "^AEX", "Euronext"
		case "BEL", "BEL20":
			return "^BFX", "Euronext"
		case "IBEX", "IBEX35":
			return "^IBEX", "BME"
		case "ESTX", "ESTX50":
			return "^STOXX50E", "Euronext"
		case "MIB":
			return "^MIB", "LSE"
		}
		// Otherwise treat as stock
		return ticker + ".AS", "Euronext"
	}
	if strings.HasPrefix(code, "UK:") {
		ticker := code[3:]
		// Map common UK indices
		switch ticker {
		case "FTSE", "FTSE100":
			return "^FTSE", "LSE"
		case "FTSE250":
			return "^MCX", "LSE"
		}
		// Otherwise treat as stock
		return ticker + ".L", "LSE"
	}
	if strings.HasPrefix(code, "CH:") {
		ticker := code[3:]
		// Map common Swiss indices
		switch ticker {
		case "SMI", "SMI20":
			return "^SSMI", "Swiss"
		case "SPI":
			return "^SPI", "Swiss"
		}
		// Otherwise treat as stock
		return ticker + ".S", "Swiss"
	}
	return code, "Unknown"
}

// GetKLineData fetches historical K-line data from Yahoo Finance
func (y *YahooFinanceApi) GetKLineData(stockCode, period string, limit int) *[]KLineData {
	yahooTicker, exchange := ConvertStockCodeToYahoo(stockCode)
	if exchange == "Unknown" {
		logger.SugaredLogger.Warnf("YahooFinance: unknown exchange for stock code: %s", stockCode)
		return &[]KLineData{}
	}

	if limit <= 0 {
		limit = 500
	}
	if limit > 1023 {
		limit = 1023
	}

	// Map period string to Yahoo interval
	interval := y.periodToInterval(period)
	if interval == "" {
		logger.SugaredLogger.Warnf("YahooFinance: unsupported period %s", period)
		return &[]KLineData{}
	}

	baseURL := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s", url.PathEscape(yahooTicker))
	params := url.Values{}
	params.Set("interval", interval)
	params.Set("range", y.limitToRange(limit, period))
	params.Set("includePrePost", "false")

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	logger.SugaredLogger.Infof("YahooFinance GetKLineData: %s", reqURL)

	resp, err := y.client.R().
		SetHeader("User-Agent", getRandomUA()).
		SetHeader("Accept", "application/json").
		SetHeader("Accept-Language", "en-US,en;q=0.9").
		Get(reqURL)

	if err != nil {
		logger.SugaredLogger.Errorf("YahooFinance HTTP error: %v", err)
		return &[]KLineData{}
	}

	if resp.StatusCode() != 200 {
		logger.SugaredLogger.Errorf("YahooFinance HTTP %d", resp.StatusCode())
		return &[]KLineData{}
	}

	var result YahooChartResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		logger.SugaredLogger.Errorf("YahooFinance JSON parse error: %v", err)
		return &[]KLineData{}
	}

	if result.Chart.Error != nil {
		logger.SugaredLogger.Errorf("YahooFinance API error: %s - %s", result.Chart.Error.Code, result.Chart.Error.Message)
		return &[]KLineData{}
	}

	if len(result.Chart.Result) == 0 {
		logger.SugaredLogger.Warnf("YahooFinance: no data returned for %s", stockCode)
		return &[]KLineData{}
	}

	chartResult := result.Chart.Result[0]
	timestamps := chartResult.Timestamp
	quotes := chartResult.Indicators.Quote

	if len(timestamps) == 0 || len(quotes) == 0 {
		logger.SugaredLogger.Warnf("YahooFinance: no quote data for %s", stockCode)
		return &[]KLineData{}
	}

	quote := quotes[0]
	klineData := make([]KLineData, 0, len(timestamps))

	for i := 0; i < len(timestamps) && i < limit; i++ {
		day := time.Unix(timestamps[i], 0).Format("2006-01-02")

		open := ""
		if i < len(quote.Open) && !isNaN(quote.Open[i]) {
			open = fmt.Sprintf("%.2f", quote.Open[i])
		}
		high := ""
		if i < len(quote.High) && !isNaN(quote.High[i]) {
			high = fmt.Sprintf("%.2f", quote.High[i])
		}
		low := ""
		if i < len(quote.Low) && !isNaN(quote.Low[i]) {
			low = fmt.Sprintf("%.2f", quote.Low[i])
		}
		close := ""
		if i < len(quote.Close) && !isNaN(quote.Close[i]) {
			close = fmt.Sprintf("%.2f", quote.Close[i])
		}
		volume := ""
		if i < len(quote.Volume) && !isNaNFloat(float64(quote.Volume[i])) {
			volume = fmt.Sprintf("%.0f", float64(quote.Volume[i]))
		}

		klineData = append(klineData, KLineData{
			Day:    day,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		})
	}

	return &klineData
}

// GetStockList searches for EU/UK stocks from the database
func (y *YahooFinanceApi) GetStockList(key string) []StockBasic {
	var result []StockBasic
	var stocksEU []models.StockInfoEU
	db.Dao.Model(&models.StockInfoEU{}).Where("name like ? or e_name like ? or code like ?", "%"+key+"%", "%"+key+"%", "%"+key+"%").Find(&stocksEU)

	for _, item := range stocksEU {
		prefix := "eu:"
		if item.Exchange == "LSE" {
			prefix = "uk:"
		} else if item.Exchange == "Swiss" {
			prefix = "ch:"
		}
		result = append(result, StockBasic{
			TsCode:   prefix + item.Code,
			Name:     item.EName, // English name for display
			Fullname: item.Name,  // Local name
			Symbol:   item.Code,
			Market:   item.Exchange,
		})
	}
	return result
}

// periodToInterval converts period string to Yahoo interval
func (y *YahooFinanceApi) periodToInterval(period string) string {
	switch period {
	case "1", "daily", "day":
		return "1d"
	case "5", "week", "weekly":
		return "1wk"
	case "15":
		return "15m"
	case "30":
		return "30m"
	case "60":
		return "1h"
	case "101", "240":
		return "1d"
	case "102", "1200":
		return "1wk"
	default:
		return "1d"
	}
}

// limitToRange converts limit and period to Yahoo range parameter
func (y *YahooFinanceApi) limitToRange(limit int, period string) string {
	switch period {
	case "1", "daily", "day":
		if limit <= 5 {
			return "5d"
		} else if limit <= 30 {
			return "1mo"
		} else if limit <= 90 {
			return "3mo"
		} else if limit <= 180 {
			return "6mo"
		} else if limit <= 365 {
			return "1y"
		} else {
			return "2y"
		}
	case "5", "week", "weekly":
		if limit <= 4 {
			return "1mo"
		} else if limit <= 12 {
			return "3mo"
		} else if limit <= 24 {
			return "6mo"
		} else if limit <= 52 {
			return "1y"
		} else {
			return "2y"
		}
	default:
		return "1mo"
	}
}

func isNaN(f float64) bool {
	return f != f
}

func isNaNFloat(f float64) bool {
	return f != f
}