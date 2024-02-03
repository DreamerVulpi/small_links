package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/dreamervulpi/small_links/config"
	"github.com/dreamervulpi/small_links/db"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Data struct {
	HostName string
	Line     string
}

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
	router.Static("/css", "templates/css")
	router.Static("/src", "src")
	router.LoadHTMLFiles("templates/index.html")

	router.POST("/postform", func(c *gin.Context) {
		inputData := c.Request.FormValue("inputData")
		str := db.RegisterLink(conn, inputData)
		data := Data{
			HostName: conf.S.Host + ":" + conf.S.Port,
			Line:     str,
		}
		slog.Info(data.HostName)
		tmpl, _ := template.ParseFiles("templates/index.html")
		tmpl.Execute(c.Writer, data)
	})

	router.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title": "Home Page",
		})
	})

	router.GET("/:id", func(c *gin.Context) {
		slog.Info(c.Param("id"))
		c.Redirect(http.StatusFound, db.GetLink(conn, c.Param("id")))
	})

	router.Run(":" + conf.S.Port)
}
