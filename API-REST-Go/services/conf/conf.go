package conf

import (
	"github.com/spf13/viper"
)

var Env *viper.Viper
var Conf *viper.Viper

func Setup(envFileName string, configFileName string) error {
	envName := ".env"
	configName := "conf"
	if envFileName != "" {
		envName = envFileName
	}
	if configFileName != "" {
		configName = configFileName
	}

	// Load environment configuration
	Env = viper.New()
	Env.SetConfigFile(envName) // name of config file with extension
	Env.AddConfigPath(".")     // optionally look for config in the working directory
	err := Env.ReadInConfig()  // Find and read the config file
	if err != nil {
		return err
	}
	// Load yaml configuration
	Conf = viper.New()
	Conf.SetConfigName(configName)
	Conf.SetConfigType("yaml")
	Conf.AddConfigPath(".")
	err = Conf.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
