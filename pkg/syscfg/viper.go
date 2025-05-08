package syscfg

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	defaultPath   = "/etc/work/"
	defaultPrefix = "application"
	defaultSuffix = "yaml"
)

func LoadSysCfgByYaml(path string, option ...Option) *SysCfg {
	if path == "" {
		path = defaultPath
	}

	viper.SetConfigName(defaultPrefix) // name of config file (without extension)
	viper.SetConfigType(defaultSuffix) // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path)          // path to look for the config file in
	err := viper.ReadInConfig()        // Find and read the config file
	if err != nil {                    // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	global, err := New(viper.GetViper(), option...)
	if err != nil {
		panic(err)
	}

	return global
}
