package data

import (
	"go-stock/backend/db"
	"go-stock/backend/logger"
	"go-stock/backend/models"
	"gorm.io/gorm"
)

// SeedEUStockData seeds the database with major European stocks and indices
func SeedEUStockData() error {
	stocks := []models.StockInfoEU{
		// Euronext Amsterdam (Netherlands)
		{Code: "ASML", Name: "ASML控股", EName: "ASML Holding N.V.", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},
		{Code: "SHELL", Name: "壳牌", EName: "Shell plc", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},
		{Code: "UNLAR", Name: "联合利华", EName: "Unilever PLC", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},
		{Code: "PHIA", Name: "飞利浦", EName: "Philips NV", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},
		{Code: "ADYEN", Name: "Adyen", EName: "Adyen N.V.", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},
		{Code: "INGA", Name: "荷兰国际集团", EName: "ING Groep N.V.", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},
		{Code: "HEIA", Name: "喜力", EName: "Heineken N.V.", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},
		{Code: "STLA", Name: "斯特兰蒂斯", EName: "Stellantis N.V.", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},

		// Euronext Paris (France)
		{Code: "RMS", Name: "路威酩轩", EName: "LVMH Moët Hennessy Louis Vuitton", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},
		{Code: "TTE", Name: "道达尔能源", EName: "TotalEnergies SE", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},
		{Code: "SAN", Name: "赛诺菲", EName: "Sanofi S.A.", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},
		{Code: "ORP", Name: "欧莱雅", EName: "L'Oréal S.A.", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},
		{Code: "BNP", Name: "法国巴黎银行", EName: "BNP Paribas S.A.", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},
		{Code: "ACA", Name: "法国农业信贷银行", EName: "Crédit Agricole S.A.", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},
		{Code: "GLE", Name: "法国兴业银行", EName: "Société Générale S.A.", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},
		{Code: "AIR", Name: "空客", EName: "Airbus SE", Exchange: "Euronext", Currency: "EUR", Type: "Stock"},

		// London Stock Exchange (UK)
		{Code: "HSBA", Name: "汇丰控股", EName: "HSBC Holdings plc", Exchange: "LSE", Currency: "GBP", Type: "Stock"},
		{Code: "BP", Name: "英国石油", EName: "BP plc", Exchange: "LSE", Currency: "GBP", Type: "Stock"},
		{Code: "ULVR", Name: "联合利华", EName: "Unilever plc", Exchange: "LSE", Currency: "GBP", Type: "Stock"},
		{Code: "RIO", Name: "力拓", EName: "Rio Tinto plc", Exchange: "LSE", Currency: "GBP", Type: "Stock"},
		{Code: "AZN", Name: "阿斯利康", EName: "AstraZeneca plc", Exchange: "LSE", Currency: "GBP", Type: "Stock"},
		{Code: "GSK", Name: "葛兰素史克", EName: "GSK plc", Exchange: "LSE", Currency: "GBP", Type: "Stock"},
		{Code: "DGE", Name: "帝亚吉欧", EName: "Diageo plc", Exchange: "LSE", Currency: "GBP", Type: "Stock"},
		{Code: "NVO", Name: "诺和诺德", EName: "Novo Nordisk A/S", Exchange: "LSE", Currency: "GBP", Type: "Stock"},
		{Code: "VOD", Name: "沃达丰", EName: "Vodafone Group plc", Exchange: "LSE", Currency: "GBP", Type: "Stock"},
		{Code: "BATS", Name: "英美烟草", EName: "British American Tobacco plc", Exchange: "LSE", Currency: "GBP", Type: "Stock"},
		{Code: "REL", Name: "RELX", EName: "RELX plc", Exchange: "LSE", Currency: "GBP", Type: "Stock"},
		{Code: "LSEG", Name: "伦敦证券交易所", EName: "London Stock Exchange Group plc", Exchange: "LSE", Currency: "GBP", Type: "Stock"},

		// Swiss Exchange (Switzerland)
		{Code: "NOVN", Name: "诺华", EName: "Novartis AG", Exchange: "Swiss", Currency: "CHF", Type: "Stock"},
		{Code: "ROG", Name: "罗氏", EName: "Roche Holding AG", Exchange: "Swiss", Currency: "CHF", Type: "Stock"},
		{Code: "UBSG", Name: "瑞银", EName: "UBS Group AG", Exchange: "Swiss", Currency: "CHF", Type: "Stock"},
		{Code: "NESN", Name: "雀巢", EName: "Nestlé S.A.", Exchange: "Swiss", Currency: "CHF", Type: "Stock"},
		{Code: "ABBN", Name: "ABB", EName: "ABB Ltd", Exchange: "Swiss", Currency: "CHF", Type: "Stock"},
		{Code: "SIKA", Name: "西卡", EName: "Sika AG", Exchange: "Swiss", Currency: "CHF", Type: "Stock"},
		{Code: "ZURN", Name: "苏黎世保险", EName: "Zurich Insurance Group AG", Exchange: "Swiss", Currency: "CHF", Type: "Stock"},
		{Code: "LONN", Name: "历峰集团", EName: "Richemont S.A.", Exchange: "Swiss", Currency: "CHF", Type: "Stock"},
	}

	indices := []models.IndexBasicEU{
		{Code: "DAX", Name: "德国DAX指数", EName: "DAX 40", Exchange: "XETRA", Currency: "EUR"},
		{Code: "CAC", Name: "法国CAC40", EName: "CAC 40", Exchange: "Euronext", Currency: "EUR"},
		{Code: "AEX", Name: "荷兰AEX指数", EName: "AEX Index", Exchange: "Euronext", Currency: "EUR"},
		{Code: "BEL", Name: "比利时BEL20", EName: "BEL 20", Exchange: "Euronext", Currency: "EUR"},
		{Code: "IBEX", Name: "西班牙IBEX35", EName: "IBEX 35", Exchange: "BME", Currency: "EUR"},
		{Code: "FTSE", Name: "英国富时100", EName: "FTSE 100", Exchange: "LSE", Currency: "GBP"},
		{Code: "SMI", Name: "瑞士 SMI", EName: "SMI Index", Exchange: "Swiss", Currency: "CHF"},
		{Code: "ESTX", Name: "欧洲斯托克50", EName: "EURO STOXX 50", Exchange: "Euronext", Currency: "EUR"},
	}

	// Seed stocks
	for _, stock := range stocks {
		var existing models.StockInfoEU
		result := db.Dao.Where("code = ? AND exchange = ?", stock.Code, stock.Exchange).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Dao.Create(&stock).Error; err != nil {
				logger.SugaredLogger.Warnf("Failed to seed EU stock %s: %v", stock.Code, err)
			}
		}
	}

	// Seed indices
	for _, index := range indices {
		var existing models.IndexBasicEU
		result := db.Dao.Where("code = ?", index.Code).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Dao.Create(&index).Error; err != nil {
				logger.SugaredLogger.Warnf("Failed to seed EU index %s: %v", index.Code, err)
			}
		}
	}

	logger.SugaredLogger.Infof("Seeded EU stocks and indices completed")
	return nil
}