package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/tanjl855/blog/pkg/utils"
)

func InitConfig() bool {
	//viper.SetConfigFile("configs/config")
	viper.SetConfigType("json")
	viper.AddConfigPath(utils.ConfigFilePath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}
	changeCfg()
	// 加载动态配置
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("changing config")
		changeCfg()
	})

	return true
}

func changeCfg() {
	if ok := viper.IsSet("env"); ok {

	}
}
