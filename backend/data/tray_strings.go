//go:build linux || windows || darwin
// +build linux,windows,darwin

package data

// TrayStrings holds all translatable strings for system tray and dialogs.
// Used by all platform-specific app_*.go files.
type TrayStrings struct {
	AppTitle       string // "go-stock"
	AppTooltip     string // systray tooltip
	MenuShowWindow string // tray menu: show window
	MenuQuit       string // tray menu: quit program
	NotifyStarted  string // startup notification
	NotifyRunning  string // second instance notification
	NotifyProfit   string // profit notification
	DialogClose    string // "Are you sure you want to close?"
	ButtonOK       string // confirm button
	ButtonCancel   string // cancel button
}

var TrayZhCN = TrayStrings{
	AppTitle:       "go-stock",
	AppTooltip:     "go-stock：AI赋能股票分析",
	MenuShowWindow: "显示窗口",
	MenuQuit:       "退出程序",
	NotifyStarted:  "应用程序已启动",
	NotifyRunning:  "程序已经在运行了",
	NotifyProfit:   "发送通知失败：%v",
	DialogClose:    "确定关闭吗？",
	ButtonOK:       "确定",
	ButtonCancel:   "取消",
}

var TrayEn = TrayStrings{
	AppTitle:       "go-stock",
	AppTooltip:     "go-stock: AI-Powered Stock Analysis",
	MenuShowWindow: "Show Window",
	MenuQuit:       "Exit Program",
	NotifyStarted:  "Application started",
	NotifyRunning:  "Another instance is already running",
	NotifyProfit:   "Failed to send notification: %v",
	DialogClose:    "Are you sure you want to close?",
	ButtonOK:       "OK",
	ButtonCancel:   "Cancel",
}

// GetTrayStrings returns the appropriate TrayStrings based on the saved language setting.
func GetTrayStrings() TrayStrings {
	cfg := GetSettingConfig()
	if cfg.Language == "en" {
		return TrayEn
	}
	return TrayZhCN
}
