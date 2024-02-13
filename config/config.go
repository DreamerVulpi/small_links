package config

import (
	"log/slog"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Database struct {
	Password string `mapstructure:"password"`
	DBname   string `mapstructure:"dbname"`
	Port     string `mapstructure:"port"`
}

type Server struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type AppConfig struct {
	S        Server
	DB       Database
	Path     string
	NameFile string
	TypeFile string
}

func LoadConfig(config *AppConfig) error {
	v := viper.New()
	v.SetConfigName(config.NameFile)
	v.SetConfigType(config.TypeFile)
	v.AddConfigPath(config.Path)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		slog.Warn("failed to read the configuration file: %s", err)
		return err
	} else {
		config.DB.Password = v.GetString("database.password")
		config.DB.DBname = v.GetString("database.dbname")
		config.DB.Port = v.GetString("database.port")
		config.S.Host = v.GetString("server.host")
		config.S.Port = v.GetString("server.port")
		return err
	}
}
