package data

import (
	"educ-gpt/models"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
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
	mockDB, _, err := sqlmock.New()
	if err != nil {
		return err
	}

	dialector := postgres.New(postgres.Config{
		Conn: mockDB,
	})

	mockGormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return err
	}

	db = mockGormDB

	log.Print("DB switched to mock version")

	return nil
}

func DB() *gorm.DB {
	if db == nil {
		database, err := gorm.Open(postgres.Open(dbConfig.String()), &gorm.Config{})
		if err != nil {
			log.Fatal(fmt.Errorf("%w: %w", ErrFailedLoadDB, err))
		}

		err = database.AutoMigrate(ms...)

		if err != nil {
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
