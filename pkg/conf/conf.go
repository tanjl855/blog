package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func InitConfig() bool {
	viper.SetConfigFile("../../configs/log.json")
	viper.SetConfigType("json")
	//viper.AddConfigPath("../../configs/")
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
