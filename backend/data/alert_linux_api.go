//go:build linux
// +build linux

package data

import (
	"go-stock/backend/logger"

	"github.com/gen2brain/beeep"
)

// AlertWindowsApi @Author spark
// @Date 2025/1/8 9:40
// @Desc Linux implementation using beeep for notifications
type AlertWindowsApi struct {
	AppID    string
	Title    string
	Content  string
	Icon     string
}

func NewAlertWindowsApi(AppID string, Title string, Content string, Icon string) *AlertWindowsApi {
	return &AlertWindowsApi{
		AppID:   AppID,
		Title:   Title,
		Content: Content,
		Icon:    Icon,
	}
}

func (a AlertWindowsApi) SendNotification() bool {
	if GetSettingConfig().LocalPushEnable == false {
		return false
	}

	err := beeep.Notify(a.AppID, a.Content, a.Icon)
	if err != nil {
		logger.SugaredLogger.Error(err)
		return false
	}
	return true
}
