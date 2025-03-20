package main

import (
	"context"
	"educ-gpt/config/data"
	"educ-gpt/config/dic"
	"educ-gpt/http/router"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

// @title Educ-GPT API
// @version 1.0
// @host https://educgpt.ru
// @BasePath /api/v1
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("ERR: godotenv.Load(): Error loading .env file")
	}

	db, _ := data.DB().DB()
	if err := db.Ping(); err != nil {
		log.Fatalf("ERR: DB: %s\n", err)
	}

	if err := data.Redis().Ping(context.Background()).Err(); err != nil {
		log.Fatalf("ERR: Redis: %s\n", err)
	}

	r := router.NewRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", os.Getenv("PORT")),
		Handler: r,
	}

	if os.Getenv("ENV") != "development" {
		dic.InitDaemons()
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("ERR: ListenAndServe: %s\n", err)
	}
}
