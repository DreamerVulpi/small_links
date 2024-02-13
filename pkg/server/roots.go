package server

import (
	"log/slog"
	"net/http"

	"github.com/dreamervulpi/small_links/config"
	"github.com/dreamervulpi/small_links/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Data struct {
	Hostname string
	Link     string
}

func Roots(conn *redis.Client, router *gin.Engine, conf config.AppConfig) {
	router.POST("/postform", func(c *gin.Context) {
		inputData := c.Request.FormValue("inputData")
		val, err := db.RegisterLink(conn, inputData)
		if err != nil {
			slog.Warn(err.Error())
		}
		c.HTML(http.StatusAccepted, "result.html", gin.H{
			"Title":    "Home Page",
			"Hostname": conf.S.Host + ":" + conf.S.Port,
			"Link":     val,
		})
	})

	router.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title": "Home Page",
		})
	})

	router.GET("/:id", func(c *gin.Context) {
		slog.Info(c.Param("id"))
		link, err := db.GetLink(conn, c.Param("id"))
		if err != nil {
			slog.Warn(err.Error())
		} else {
			c.Redirect(http.StatusFound, link)
		}
	})
}
