package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

//	TODO: Refactoring code

func generateLink() string {
	digits := "0123456789"
	symbols := "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var letters = []rune(symbols + digits)
	rslt := make([]rune, 12)
	for i := range rslt {
		rslt[i] = letters[rand.Intn(len(letters))]
	}
	slog.Info("https://shrt.link/" + string(rslt))
	return "https://shrt.link/" + string(rslt)

}

func do(inputData string) {
	connStr := "user=postgres password=1234 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Warn("Connect is FAILDED :C\n")
		panic(err)
	}
	defer db.Close()
	sql := "INSERT INTO links (origlink, customlink) VALUES ('" + inputData + "', '" + generateLink() + "')"
	result, err := db.Exec(sql)
	if err != nil {
		slog.Warn("INSERT INTO links IS FAILDED :C\n")
		slog.Warn(sql + "\n")
		slog.Warn(err.Error())
		panic(err)
	}
	fmt.Println(result.RowsAffected())
	slog.Info("u did it, yo -> " + inputData)
}

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title": "Home Page",
		})
	})
	router.POST("/postform", func(c *gin.Context) {
		inputData := c.Request.FormValue("inputData")
		do(inputData)
	})
	router.Run(":8080")
}
