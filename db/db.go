package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"

	_ "github.com/lib/pq"
)

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

func WorkWithDB(conn *sql.DB, inputData string) {
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
