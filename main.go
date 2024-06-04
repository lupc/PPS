package main

import (
	"fmt"
	"time"

	"pps/pkg"
	"pps/util"

	"github.com/kardianos/service"
	"github.com/lupc/go_service"
)

func main() {
	srvConfig := &service.Config{
		Name:        "PPS",
		DisplayName: "PPS进程守护服务",
		Description: "PPS进程守护服务",
	}

	_ = go_service.RunWithService(srvConfig, run)
	// if s == nil {
	// 	util.GetLogger().Info("守护服务服务启动失败")
	// } else {
	// 	util.GetLogger().Info("守护服务服务已启动", zap.Any("status", s.Status))
	// }

}

func run() {
	// ctx, cancel := context.WithCancel(context.Background())
	// go WaitTerm(cancel)
	util.GetLogger().Info("守护服务启动..")

	var cfgPath = "./config.yml"
	var cfg = pkg.GetConfig(cfgPath)

	if cfg != nil && len(cfg.Configs) > 0 {
		for _, c := range cfg.Configs {
			if c.IsEnable {
				util.GoWithRecover(fmt.Sprintf("守护进程[%v]出错", c.ProcessName), func() {
					for {
						pkg.ProtectProcess(c)
						time.Sleep(time.Second * 10)
					}
				})
			}
		}
	}
	// <-ctx.Done()
}

// func run2() {
// 	// go func() {
// 	for {
// 		util.GetLogger().Info("run2 test")
// 		time.Sleep(time.Second)
// 	}
// 	// }()
// }

// func WaitTerm(cancel context.CancelFunc) {
// 	sigc := make(chan os.Signal, 1)
// 	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
// 	defer signal.Stop(sigc)
// 	<-sigc
// 	cancel()
// }
