package data

import (
	"educ-gpt/models"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var (
	ErrFailedLoadDB          = errors.New("failed to load db")
	ErrFailedLoadAutoMigrate = errors.New("failed to load auto migrate")
)

var ms = []interface{}{
	&models.User{},
	&models.Role{},
	&models.UserRoles{},
	&models.Topic{},
	&models.Theme{},
	&models.Problem{},
	&models.UserTheme{},
	&models.Token{},
	&models.Dialog{},
	&models.DialogItem{},
}

var db *gorm.DB

var dbConfig *DBconfig

func SetDBConfig(config *DBconfig) {
	dbConfig = config
}

func SwitchToMock() error {
	mockGormDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open in-memory SQLite database: %w", err)
	}

	if err := mockGormDB.AutoMigrate(ms...); err != nil {
		return fmt.Errorf("failed to auto-migrate in-memory SQLite database: %w", err)
	}

	db = mockGormDB

	log.Print("DB switched to in-memory SQLite version")

	SeedMock()

	return nil
}

func DB() *gorm.DB {
	if db == nil {
		database, err := gorm.Open(postgres.Open(dbConfig.String()), &gorm.Config{})
		if err != nil {
			log.Fatal(fmt.Errorf("%w: %w", ErrFailedLoadDB, err))
		}

		if err = database.AutoMigrate(ms...); err != nil {
			log.Fatal(fmt.Errorf("%w: %w", ErrFailedLoadAutoMigrate, err))
		}

		db = database

		log.Print("DB Initialized")
		Seed()
	}

	return db
}

type DBconfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLmode  string
}

func (c DBconfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLmode)
}
