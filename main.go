package main

import (
	"log/slog"

	"github.com/dreamervulpi/small_links/config"
	"github.com/dreamervulpi/small_links/pkg/db"
	"github.com/dreamervulpi/small_links/pkg/server"
)

func main() {
	var conf config.AppConfig
	conf.Path = "."
	conf.NameFile = "config"
	conf.TypeFile = "yaml"
	config.LoadConfig(&conf)

	params := map[string]string{
		"hostname": conf.S.Host,
		"port":     conf.DB.Port,
		"password": conf.DB.Password,
		"dbname":   conf.DB.DBname,
	}
	conn := db.ConnectToDB(params)
	slog.Info(conn.String())

	host := server.Init()

	server.Roots(conn, host, conf)
	server.Run(host, conf.S.Port)
	defer conn.Close()
}
