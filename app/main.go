package main

import (
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

	data.SetDBConfig(&data.DBconfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Name:     os.Getenv("DB_NAME"),
		SSLmode:  os.Getenv("DB_SSL"),
	})

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
