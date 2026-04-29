//go:build linux
// +build linux

package main

import (
	"context"
	"fmt"
	"go-stock/backend/data"
	"go-stock/backend/db"
	"go-stock/backend/logger"
	"log"
	"time"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gen2brain/beeep"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	defer PanicHandler()
	runtime.EventsOn(ctx, "frontendError", func(optionalData ...interface{}) {
		logger.SugaredLogger.Errorf("Frontend error: %v\n", optionalData)
	})
	logger.SugaredLogger.Infof("Version:%s", Version)
	// Perform your setup here
	a.ctx = ctx

	// 应用启动时自动创建已启用的定时任务
	a.InitCronTasks()

	preCacheTradingDays()

	// 监听设置更新事件
	runtime.EventsOn(ctx, "updateSettings", func(optionalData ...interface{}) {
		config := data.GetSettingConfig()
		//logger.SugaredLogger.Infof("updateSettings config:%+v", config)
		if config.DarkTheme {
			runtime.WindowSetBackgroundColour(ctx, 27, 38, 54, 1)
			runtime.WindowSetDarkTheme(ctx)
		} else {
			runtime.WindowSetBackgroundColour(ctx, 255, 255, 255, 1)
			runtime.WindowSetLightTheme(ctx)
		}
		runtime.WindowReloadApp(ctx)
	})

	// 创建 Linux 托盘通知
	go func() {
		ts := data.GetTrayStrings()
		err := beeep.Notify(ts.AppTitle, ts.NotifyStarted, "")
		if err != nil {
			log.Fatalf("系统通知失败：%v", err)
		}
	}()

	logger.SugaredLogger.Infof(" application startup Version:%s", Version)
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	defer PanicHandler()

	// 记录当前窗口大小，供下次启动时还原
	if a.ctx != nil {
		w, h := runtime.WindowGetSize(ctx)
		logger.SugaredLogger.Infof(" window size: %dx%d", w, h)
		if w > 0 && h > 0 {
			cfg := data.GetSettingConfig()
			cfg.WindowWidth = w
			cfg.WindowHeight = h
			data.UpdateConfig(cfg)
			logger.SugaredLogger.Infof("save window size: %dx%d", w, h)
		}
	}

	// 在 Linux 上使用 MessageDialog 显示确认窗口
	ts := data.GetTrayStrings()
	dialog, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:         runtime.QuestionDialog,
		Title:        ts.AppTitle,
		Message:      ts.DialogClose,
		Buttons:      []string{ts.ButtonOK, ts.ButtonCancel},
		Icon:         icon2,
		CancelButton: ts.ButtonCancel,
	})

	if err != nil {
		logger.SugaredLogger.Errorf("dialog error:%s", err.Error())
		return false
	}

	logger.SugaredLogger.Debugf("dialog:%s", dialog)
	if dialog == ts.ButtonCancel || dialog == "No" {
		return true // 如果选择了取消，不关闭应用
	} else {
		// 在 Linux 上应用退出时执行清理工作
		if a.cron != nil {
			a.cron.Stop() // 停止定时任务
		}
		return false // 如果选择了确定，继续关闭应用
	}
}

// OnSecondInstanceLaunch 处理第二实例启动时的通知
func OnSecondInstanceLaunch(secondInstanceData options.SecondInstanceData) {
	ts := data.GetTrayStrings()
	err := beeep.Notify(ts.AppTitle, ts.NotifyRunning, "")
	if err != nil {
		logger.SugaredLogger.Error(err)
	}
	time.Sleep(time.Second * 3)
}

// MonitorStockPrices 监控股票价格
func MonitorStockPrices(a *App) {
	// 检查是否至少有一个市场开市
	isAStockOpen := isTradingTime(time.Now())
	isHKStockOpen := IsHKTradingTime(time.Now())
	isUSStockOpen := IsUSTradingTime(time.Now())

	// 如果所有市场都不在交易时间，则提前返回
	if !isAStockOpen && !isHKStockOpen && !isUSStockOpen {
		logger.SugaredLogger.Debugf("当前所有市场均未开市，跳过价格监控")
		return
	}

	logger.SugaredLogger.Debugf("市场状态 - A股: %v, 港股: %v, 美股: %v", isAStockOpen, isHKStockOpen, isUSStockOpen)

	dest := &[]data.FollowedStock{}
	db.Dao.Model(&data.FollowedStock{}).Find(dest)
	total := float64(0)

	// 股票信息处理逻辑
	stockInfos := GetStockInfos(*dest...)
	for _, stockInfo := range *stockInfos {
		if strutil.HasPrefixAny(stockInfo.Code, []string{"SZ", "SH", "sh", "sz"}) && (!isTradingTime(time.Now())) {
			continue
		}
		if strutil.HasPrefixAny(stockInfo.Code, []string{"hk", "HK"}) && (!IsHKTradingTime(time.Now())) {
			continue
		}
		if strutil.HasPrefixAny(stockInfo.Code, []string{"us", "US", "gb_"}) && (!IsUSTradingTime(time.Now())) {
			continue
		}

		total += stockInfo.ProfitAmountToday
		price, _ := convertor.ToFloat(stockInfo.Price)

		if stockInfo.PrePrice != price {
			go runtime.EventsEmit(a.ctx, "stock_price", stockInfo)
		}
	}

	// 计算总收益并更新状态
	if total != 0 {
		ts := data.GetTrayStrings()
		title := ts.AppTitle + " " + time.Now().Format(time.DateTime) + fmt.Sprintf("  %.2f¥", total)

		// 发送通知显示实时数据
		err := beeep.Notify(ts.AppTitle, title, "")
		if err != nil {
			logger.SugaredLogger.Errorf("%s: %v", ts.NotifyProfit, err)
		}
	}

	// 触发实时利润事件
	go runtime.EventsEmit(a.ctx, "realtime_profit", fmt.Sprintf("  %.2f", total))
}

// getFrameless 返回是否使用无边框窗口
func getFrameless() bool {
	return false
}

// getScreenResolution 返回屏幕分辨率
func getScreenResolution() (int, int, int, int, error) {
	// Linux 上使用简单的默认值
	// 可以通过 xrandr 或其他工具获取实际分辨率
	return 1412, 834, 900, 600, nil
}
