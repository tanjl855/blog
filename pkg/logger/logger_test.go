package logger

import (
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestSetupLogger(t *testing.T) {
	// windows test
	viper.SetConfigFile("../.././configs/log.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	pwd, err := os.Getwd()
	if err != nil {
		Logger.Errorf("get pwd error: %v", err)
		t.Error(err)
	}
	SetupLogger()
	Logger.Infof("Here is test: %v", pwd)
}
