package util

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lupc/go-myzap"
	"go.uber.org/zap"
)

var logger *zap.Logger

// 获取默认Logger，首次获取会自动创建
func GetLogger() *zap.Logger {
	if logger == nil {
		var myZapConfig = myzap.NewConfigByName("pps")
		myZapConfig.TimeFormat = "2006-01-02 15:04:05.000"
		logger = myZapConfig.BuildLogger().
			WithOptions(zap.WithCaller(false))
	}
	return logger
}

// 捕获panic并写到error日志
func LogRecoverToError(reason string) {
	if err := recover(); err != nil {
		// fmt.Println("recover in func DefaultRecover")
		// fmt.Println(fmt.Sprintf("%T %v", err, err))
		//...handle  打日志等
		var msg = fmt.Sprintf("panic:%v", reason)
		GetLogger().Error(msg, zap.Error(err.(error)))
	}
}

// 运行方法并捕获panic写日志
func RunWithRecover(panicReason string, action func()) {
	defer LogRecoverToError(panicReason)
	action()
}

// 协程运行方法并捕获panic写日志
func GoWithRecover(panicReason string, action func()) {
	go func() {
		RunWithRecover(panicReason, action)
	}()
}

// 成功 高亮显示
func HighlightSuccess(msg string) string {
	colorFormat := color.New(color.FgWhite, color.BgGreen).SprintFunc()
	return colorFormat(" " + msg + " ")
}

// 失败 高亮显示
func HighlightFail(msg string) string {
	colorFormat := color.New(color.FgWhite, color.BgRed).SprintFunc()
	return colorFormat(" " + msg + " ")
}

// 已连接 高亮显示
func HighlightConnected(msg string) string {
	colorFormat := color.New(color.FgGreen, color.BgWhite).SprintFunc()
	return colorFormat(" " + msg + " ")
}

// 已断开 高亮显示
func HighlightDisconnected(msg string) string {
	colorFormat := color.New(color.FgHiBlack, color.BgWhite).SprintFunc()
	return colorFormat(" " + msg + " ")
}

// 错误 高亮显示
func HighlightError(msg string) string {
	colorFormat := color.New(color.FgRed).SprintFunc()
	return colorFormat(" " + msg + " ")
}
