package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"
)

func TestInitConfig(t *testing.T) {
	InitConfig()
	fmt.Println(viper.GetString("info_file_name"))
}
