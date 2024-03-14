package config

import (
	"github.com/spf13/viper"
	"github.com/universalmacro/common/singleton"
)

func GetConfig() *ConfigSingleton {
	return configSingleton.Get()
}

var configSingleton = singleton.SingletonFactory(newConfigSingleton, singleton.Eager)

type ConfigSingleton struct {
	*viper.Viper
}

func newConfigSingleton() *ConfigSingleton {
	c := viper.New()
	c.SetConfigName(".env")
	c.SetConfigType("yaml")
	c.AddConfigPath(".")
	c.ReadInConfig()
	return &ConfigSingleton{c}
}
