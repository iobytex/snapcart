package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Postgresql Postgresql
}

type Postgresql struct {
	PostgresqlHost string
	PostgresqlPort string
	PostgresqlUser string
	PostgresqlPassword string
	PostgresqlDbname string
	PostgresqlSslMode string
	PgDriver string
}

func LoadConfig()  (*viper.Viper,error){
	v := viper.New()
	v.SetConfigName("config-local")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config/")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("Config file not found; ignore error if desired: %v",err)
			return nil,errors.Wrap(err,"Load Config:")
		} else {
			fmt.Printf("Config file was found but another error was produced: %v",err)
			return nil,errors.Wrap(err,"Load Config:")
		}
	}
	return v,nil
}


func GetConfig() (*Config,error) {
	v,vErr := LoadConfig()
	if vErr != nil {
		return nil,errors.Wrap(vErr,"LoadConfig ==> Get Config: ")
	}
	config, err := ParseConfig(v)
	if err != nil {
		return nil,errors.Wrap(err," ParseConfig ==> Get Config")
	}

	return config,nil
}

func ParseConfig(v *viper.Viper)  (*Config,error) {
	var config *Config
	err := v.Unmarshal(&config)
	if err != nil {
		fmt.Printf("unable to decode into struct: %v",err)
		return nil, errors.Wrap(err,"Parse Config:")
	}

	return config,nil
}
