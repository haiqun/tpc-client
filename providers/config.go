package providers

import (
	"github.com/spf13/viper"
	"path/filepath"
	"sync"
)

//@TODO 因为同一个包里面的 init 调用顺序是看编译器排序的，但变量的声明赋值是在 init 之前的,由此来保障配置文件在最前面
// 配置文件对象

var Config *viper.Viper

var OnecC sync.Once
// 初始化函数
func GetConfig() *viper.Viper {
	OnecC.Do(func() {
		Config = getConfig()
	})
	return  Config
}

func getConfig() *viper.Viper {
	Config := viper.New()
	// 设置配置位置
	Config.AddConfigPath(filepath.Join(RootPath, "config"))
	Config.SetConfigName("config")
	Config.SetConfigType("toml")

	if err := Config.ReadInConfig(); err != nil {
		panic("读取配置文件失败!,原因:" + err.Error())
	}
	return Config
}