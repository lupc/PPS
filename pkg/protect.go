package pkg

import (
	"fmt"
	"os"
	"pps/util"
	"strings"
	"time"

	"go.uber.org/zap"
)

// 启动协程守护进程
func ProtectProcess(cfg *ConfigStart) {
	if cfg == nil {
		return
	}
	var runCmd = cfg.RunCmd
	var pName = cfg.ProcessName
	var interval = cfg.CheckInterval
	defer util.LogRecoverToError("启动守护出错")
	defer runProcess(runCmd, cfg.RunDelay)

	p, err := findProcess(pName)
	if err != nil {
		util.GetLogger().Error("查找进程出错", zap.Any("processName", pName), zap.Error(err))
		return
	}
	if p == nil {
		util.GetLogger().Error("找不到进程", zap.Any("processName", pName))
		return
	}
	util.GetLogger().Sugar().Infof("开始守护进程[%v]...", pName)

	for {
		cpu, err := p.CPUPercent()
		if err != nil {
			util.GetLogger().Error("获取进程CPU信息出错", zap.Any("processName", pName), zap.Error(err))
			break
		}
		men, err := p.MemoryInfo()
		if err != nil {
			util.GetLogger().Error("获取进程内存信息出错", zap.Any("processName", pName), zap.Error(err))
			break
		}
		menMB := (float64)(men.RSS) / 1024 / 1024
		if menMB > float64(cfg.StopWhenMenoryUsage) {
			runProcess(cfg.StopCmd, 1000) //内存过大停止
		}
		util.GetLogger().Debug("进程状态", zap.Any("name", pName), zap.String("cpu", fmt.Sprintf("%.2f%%", cpu)), zap.Any("memory", fmt.Sprintf("%.2fMB", menMB)))
		time.Sleep(interval)
	}
	//进程已退出,执行重启脚本
	util.GetLogger().Info("进程已停止", zap.Any("process", pName))

	//p,err = os.StartProcess(runCmd, argv []string, attr *ProcAttr)
}

func runProcess(runCmd string, delay time.Duration) {
	if runCmd == "" {
		return
	}
	arr := strings.Split(runCmd, " ")
	if len(arr) == 0 {
		return
	}
	time.Sleep(delay)
	// cmd := exec.Command(arr[0], arr[1:]...)
	// err := cmd.Run()
	// if err != nil {
	// 	util.GetLogger().Error("执行启动命令失败", zap.Any("cmd", runCmd), zap.Error(err))
	// } else {
	// 	util.GetLogger().Info("执行启动命令成功", zap.Any("cmd", runCmd))
	// }
	attr := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	_, err := os.StartProcess(arr[0], arr[1:], attr)
	if err != nil {
		util.GetLogger().Error("执行启动命令失败", zap.Any("cmd", runCmd), zap.Error(err))
	} else {
		util.GetLogger().Info("执行启动命令成功", zap.Any("cmd", runCmd))
	}
}
