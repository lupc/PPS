package pkg

import (
	"os"
	"pps/util"
	"time"

	"github.com/shirou/gopsutil/process"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type ConfigStart struct {
	ProcessName         string        `yaml:"ProcessName"`         //进程名称（唯一标识）
	RunCmd              string        `yaml:"RunCmd"`              //执行语句
	StopCmd             string        `yaml:"StopCmd"`             //停止语句
	CheckInterval       time.Duration `yaml:"CheckInterval"`       //检测间隔
	RunDelay            time.Duration `yaml:"RunDelay"`            //延时启动
	StopWhenMenoryUsage float64       `yaml:"StopWhenMenoryUsage"` //指定内存占用超过该值停止程序（单位mb）
	IsEnable            bool          `yaml:"IsEnable"`            //是否启用
}
type Config struct {
	Configs []*ConfigStart //支持多个进程守护
}

func GetConfig(path string) (cfg *Config) {
	defer util.LogRecoverToError("加载配置出错")
	// defer func() {
	// 	err := recover() //内置函数，可以捕捉到函数异常
	// 	if err != nil {
	// 		//这里是打印错误，还可以进行报警处理，例如微信，邮箱通知
	// 		util.GetLogger().Error("加载配置出错", zap.Error(err.(error)))
	// 	}
	// }()

	// var ccc *ConfigStart = nil
	// if ccc.IsEnable {
	// 	return
	// }

	data, err := os.ReadFile(path)
	if err != nil {
		util.GetLogger().Error("打开配置文件出错", zap.Any("path", path), zap.Error(err))
		if os.IsNotExist(err) {
			//创建一个空的
			cfg = new(Config)
			cfg.Configs = append(cfg.Configs, &ConfigStart{
				ProcessName:   "xxx.exe",
				RunCmd:        "start_xxx.bat",
				CheckInterval: time.Second * 5,
				RunDelay:      time.Second,
				IsEnable:      false,
			})
			data, err := yaml.Marshal(cfg)
			if err != nil {
				util.GetLogger().Error("序列化配置出错", zap.Any("path", path), zap.Error(err))
			} else {
				newFile, err := os.Create(path)

				if err != nil {
					util.GetLogger().Error("创建配置文件出错", zap.Any("path", path), zap.Error(err))
				} else {
					util.GetLogger().Info("成功创建配置文件", zap.Any("path", path))
					newFile.Write(data)
				}
				newFile.Close()
			}
		}
		return
	}
	// if err != nil {
	// 	util.GetLogger().Error("读取配置文件出错", zap.Any("path", path), zap.Error(err))
	// 	return
	// }
	if len(data) == 0 {
		util.GetLogger().Error("配置文件为空", zap.Any("path", path))
		return
	}
	// fmt.Printf("string(data): %v\n", string(data))
	util.GetLogger().Sugar().Infof("成功加载配置:\n%v", string(data))
	var ccc = &Config{}
	err = yaml.Unmarshal(data, ccc)
	// // 创建解析器
	// decoder := yaml.NewDecoder(file)
	// // 解析 YAML 数据
	// var cfg2 *Config
	// err = decoder.Decode(cfg2)
	if err != nil {
		util.GetLogger().Error("解析配置文件出错", zap.Any("path", path), zap.Error(err))
		return
	}
	cfg = ccc
	return
}

func findProcess(name string) (proc *process.Process, err error) {
	pids, err := process.Pids()
	//fmt.Println("pids:", pids)
	if err != nil {
		return
	}
	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err == nil {
			n, err := p.Name()
			if err == nil && n == name {
				proc = p
			}
		}
	}
	return
}
