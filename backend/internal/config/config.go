package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Init() error {
	v := viper.New()
	v.SetConfigFile("configs/app.yaml")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("read config file err:%v", err)
		}
	}
	if err := v.Unmarshal(&AppEngine); err != nil {
		return fmt.Errorf("unmarshal config err:%v", err)
	}
	for k, v := range AppEngine.Redis {
		client, err := v.NewClient()
		if err != nil {
			panic(fmt.Errorf("redis:%v new client err:%v", v, err))
		}
		Conns.Store(k, client)
	}
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		if err := v.Unmarshal(&AppEngine); err != nil {
			panic(fmt.Errorf("unmarshal config err:%v", err))
		}
		fmt.Println("config change", AppEngine)
	})

	AppEngine.Flags(pflag.CommandLine)

	pflag.Parse()
	v.BindPFlags(pflag.CommandLine)
	return nil
}
