package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

//	TODO: Refactoring code

type database struct {
	driver   string `mapstructure:"driver"`
	user     string `mapstructure:"user"`
	password string `mapstructure:"password"`
	dbname   string `mapstructure:"dbname"`
	sslmode  string `mapstructure:"sslmode"`
}

type server struct {
	host string `mapstructure:"host"`
	port string `mapstructure:"port"`
}

type AppConfig struct {
	s        server
	db       database
	path     string
	nameFile string
	typeFile string
}

func LoadConfig(config *AppConfig) {
	v := viper.New()
	v.SetConfigName(config.nameFile)
	v.SetConfigType(config.typeFile)
	v.AddConfigPath(config.path)

	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		slog.Warn("failed to read the configuration file: %s", err)
		return
	}
	config.db.driver = v.GetString("database.driver")
	config.db.user = v.GetString("database.user")
	config.db.password = v.GetString("database.password")
	config.db.dbname = v.GetString("database.dbname")
	config.db.sslmode = v.GetString("database.sslmode")
	config.s.host = v.GetString("server.host")
	config.s.port = v.GetString("server.port")
	config.s.port = v.GetString("server.port")
	config.s.port = v.GetString("server.port")
	config.s.port = v.GetString("server.port")
	config.s.port = v.GetString("server.port")
	config.s.port = v.GetString("server.port")
	config.s.port = v.GetString("server.port")
	return
}

func generateLink() string {
	digits := "0123456789"
	symbols := "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var letters = []rune(symbols + digits)
	rslt := make([]rune, 6)
	for i := range rslt {
		rslt[i] = letters[rand.Intn(len(letters))]
	}
	slog.Info(string(rslt))
	return string(rslt)

}

func do(conn *sql.DB, inputData string) {
	// SELECT EXISTS(SELECT origlink FROM links WHERE origlink = 'UWU.ooo')
	var status bool
	conn.QueryRow("SELECT EXISTS(SELECT origlink FROM links WHERE origlink = $1)", inputData).Scan(&status)
	if status {
		slog.Warn("oops, this link is avaliable")
		var link string
		err := conn.QueryRow("SELECT customlink FROM links WHERE origlink = $1)", inputData).Scan(&link)
		slog.Info("SELECT customlink FROM links WHERE origlink = '" + inputData + "'")
		if err != nil {
			slog.Warn("oops, can't get link")
		}
		slog.Info("ur link ->" + link)
	} else {

		result, err := conn.Exec("INSERT INTO links (origlink, customlink) VALUES ($1 , $2)", inputData, generateLink())

		if err != nil {
			slog.Warn(fmt.Sprintf("INSERT INTO links (origlink, customlink) VALUES ($1 , $2)", inputData, generateLink()))
			panic(err)
		}
		fmt.Println(result.RowsAffected())
		slog.Info("u did it, yo -> " + inputData)
	}
}

func main() {
	var config AppConfig
	config.path = "."
	config.nameFile = "config"
	config.typeFile = "yaml"

	LoadConfig(&config)
	slog.Info(config.db.driver)

	// FIXME: THIS IS BULLSSHIT
	data := "user=" + config.db.user + " password= " + config.db.password + " dbname= " + config.db.dbname + " sslmode=" + config.db.sslmode
	conn, err := sql.Open(config.db.driver, data)
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
		do(conn, inputData)
	})
	router.Run(":" + config.s.port)
}
