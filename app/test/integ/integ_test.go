package integ

import (
	"context"
	"educ-gpt/config/data"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
	"time"
)

var (
	rdb *redis.Client
	db  *gorm.DB

	ctx = context.Background()

	key   = "testKey"
	value = "testValue"
	ttl   = time.Minute * 1
)

func setup() error {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Print(err)
		return err
	}

	rdb = data.Redis()

	data.SetDBConfig(&data.DBconfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Name:     os.Getenv("DB_NAME"),
		SSLmode:  os.Getenv("DB_SSL"),
	})

	db = data.DB()

	err = db.Exec("DROP TABLE IF EXISTS test_entities").Error
	if err != nil {
		log.Print(err)
		return err
	}

	createTableSQL := `
	CREATE TABLE test_entities (
		column1 SERIAL PRIMARY KEY,
		column2 FLOAT,
		column3 TEXT,
		column4 TIMESTAMP,
		column5 BOOLEAN
	)`
	err = db.Exec(createTableSQL).Error
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func tearDown() error {
	_, err := rdb.Del(ctx, key).Result()

	if err != nil {
		log.Print(err)
		return err
	}

	err = db.Exec("DROP TABLE IF EXISTS test_entities").Error
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		os.Exit(1)
	}

	exitCode := m.Run()

	if err := tearDown(); err != nil {
		os.Exit(1)
	}

	os.Exit(exitCode)
}

type TestEntity struct {
	Column1 int       `gorm:"column:column1"` // Целое число
	Column2 float64   `gorm:"column:column2"` // Число с плавающей точкой
	Column3 string    `gorm:"column:column3"` // Строка
	Column4 time.Time `gorm:"column:column4"` // Время
	Column5 bool      `gorm:"column:column5"` // Булево значение
}
