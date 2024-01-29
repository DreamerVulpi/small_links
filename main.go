package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dreamervulpi/small_links/config"
	"github.com/dreamervulpi/small_links/db"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	var conf config.AppConfig
	conf.Path = "."
	conf.NameFile = "config"
	conf.TypeFile = "yaml"

	config.LoadConfig(&conf)
	slog.Info(conf.DB.Driver)

	data := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", conf.DB.User, conf.DB.Password, conf.DB.DBname, conf.DB.Sslmode)
	slog.Info(data)
	conn, err := sql.Open(conf.DB.Driver, data)
	if err != nil {
		slog.Warn("Connect is FAILDED :C\n")
		panic(err)
	}
	defer conn.Close()

	router := gin.Default()
	router.LoadHTMLFiles("index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title": "Home Page",
		})
	})
	router.POST("/postform", func(c *gin.Context) {
		inputData := c.Request.FormValue("inputData")
		db.WorkWithDB(conn, inputData)
	})
	router.Run(":" + conf.S.Port)
}
