package data

import (
	"educ-gpt/logger"
	"educ-gpt/models"
	"errors"
	"fmt"
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
	&models.UserTopic{},
}

var db *gorm.DB

var dbConfig *DBconfig

func SetDBConfig(config *DBconfig) {
	dbConfig = config
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
		logger.Logger().Debug("DB initialized")

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
