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

func GetLink(conn *sql.DB, id string) string {
	var status bool
	conn.QueryRow("SELECT EXISTS(SELECT origlink FROM links WHERE customlink = $1)", id).Scan(&status)
	if status {
		var link string
		err := conn.QueryRow("SELECT origlink FROM links WHERE customLink = $1", id).Scan(&link)
		if err != nil {
			slog.Warn("oops, can't get link :C")
		}
		slog.Info("Success. Original link -> " + link)
		return link
	} else {
		slog.Warn("Oops, this link isn't work")
		return ""
	}
}

func RegisterLink(conn *sql.DB, inputData string) string {
	var status bool
	conn.QueryRow("SELECT EXISTS(SELECT origlink FROM links WHERE origlink = $1)", inputData).Scan(&status)
	if status {
		slog.Warn("oops, this link is avaliable")
		var link string
		err := conn.QueryRow("SELECT customlink FROM links WHERE origlink = $1)", inputData).Scan(&link)
		slog.Info("SELECT customlink FROM links WHERE origlink = '" + inputData + "'")
		if err != nil {
			slog.Warn("Oops, can't get link")
		} else {
			slog.Info("Original Link ->" + link)
		}
		return link
	} else {
		link := generateLink()
		_, err := conn.Exec("INSERT INTO links (origlink, customlink) VALUES ($1 , $2)", inputData, link)

		if err != nil {
			slog.Warn(fmt.Sprintf("INSERT INTO links (origlink, customlink) VALUES ($1 , $2)", inputData, link))
			panic(err)
		}
		slog.Info("Original link -> " + inputData)
		return link
	}
}
