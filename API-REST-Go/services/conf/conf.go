package conf

import (
	"github.com/spf13/viper"
)

var Env *viper.Viper
var Conf *viper.Viper

var Domain string
var JwtSecret string

func Setup() error {
	// Load environment configuration
	Env = viper.New()
	Env.SetConfigFile(".env") // name of config file with extension
	Env.AddConfigPath(".")    // optionally look for config in the working directory
	err := Env.ReadInConfig() // Find and read the config file
	if err != nil {
		return err
	}
	Domain = Env.GetString("DOMAIN")
	JwtSecret = Env.GetString("JWT_SECRET")
	// Load yaml configuration
	Conf = viper.New()
	Conf.SetConfigName("conf")
	Conf.SetConfigType("yaml")
	Conf.AddConfigPath(".")
	err = Conf.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
