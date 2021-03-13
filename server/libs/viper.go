package libs

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/fsnotify/fsnotify"
)

var Viper *viper.Viper

func ViperInit(config string) {
	//var config string

	v := viper.New()
	v.SetConfigFile(config)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	Viper = v
}
