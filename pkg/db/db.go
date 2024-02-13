package db

import (
	"context"
	"log/slog"
	"math/rand"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectToDB(params map[string]string) *redis.Client {
	dbname, _ := strconv.Atoi(params["dbname"])
	conn := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: params["password"],
		DB:       dbname,
	})
	ctx := context.Background()
	pong := conn.Ping(ctx)
	slog.Info(pong.Result())
	return conn
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

func GetLink(conn *redis.Client, link string) (string, error) {
	ctx := context.Background()
	original, err := conn.Get(ctx, link).Result()
	if err != nil {
		return "", err
	}
	slog.Info(link + "->" + original)
	return original, err
}

func RegisterLink(conn *redis.Client, inputData string) (string, error) {
	val, err := GetLink(conn, inputData)
	if val != "" {
		slog.Info(inputData + "==" + val)
		return val, err
	}
	ctx := context.Background()
	newLink := generateLink()
	err = conn.Set(ctx, newLink, inputData, 24*time.Hour*30).Err()
	slog.Info(inputData + "->" + newLink)
	return newLink, err
}
