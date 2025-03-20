package data

import (
	"educ-gpt/models"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"time"
)

func Seed() {
	rolesSeed()
	topicsSeed()
	adminSeed()
}

func SeedMock() {
	Seed()
	usersSeed()
}

func adminSeed() {
	now := time.Now()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	}

	user := &models.User{
		Name:       "admin",
		Email:      os.Getenv("ADMIN_EMAIL"),
		Password:   string(hashedPassword),
		ActivateAt: &now,
		CreatedAt:  now,
		Roles: []*models.Role{
			{
				ID:   1,
				Name: "ADMIN",
			},
		},
	}

	result := db.FirstOrCreate(&models.User{}, user)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Fatalf("Failed to create %s: %v", user.Name, result.Error)
	}

	log.Print("Users seed completed")
}

func usersSeed() {
	now := time.Now()

	users := make([]*models.User, 10)

	for i := 0; i < 10; i++ {
		iItoa := strconv.Itoa(i)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("user_user")), bcrypt.DefaultCost)

		if err != nil {
			log.Fatal(err)
		}

		users[i] = &models.User{
			Name:       "user" + iItoa,
			Email:      "user" + iItoa + "@educgpt.ru",
			Password:   string(hashedPassword),
			ActivateAt: &now,
			CreatedAt:  now,
			Roles: []*models.Role{
				{
					ID:   2,
					Name: "USER",
				},
			},
		}
	}

	for _, user := range users {
		result := db.FirstOrCreate(&models.User{}, &user)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Fatalf("Failed to create %s: %v", user.Name, result.Error)
		}
	}

	log.Print("Users seed completed")
}

func rolesSeed() {
	sRoleNames := []string{
		"ADMIN",
		"USER",
	}

	for _, roleName := range sRoleNames {
		result := db.FirstOrCreate(&models.Role{}, &models.Role{Name: roleName})
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Fatalf("Failed to create role %s: %v", sRoleNames, result.Error)
		}
	}

	log.Print("Roles seed completed")
}

func topicsSeed() {
	sTopicThemes := []string{
		"Frontend Developer",
		"Backend Developer",
		"DevOps Engineer",
		"Android Developer",
		"QA Engineer",
		"React Developer",
		"Angular Developer",
		"Vue.js Developer",
		"Node.js Developer",
		"Go Developer",
		"Java Developer",
		"PHP Developer",
		"Python Developer",
		".NET Developer",
		"Flutter Developer",
		"Swift Developer",
		"DBA",
		"Blockchain Developer",
		"Cyber Security",
		"Software Architect",
		"Software Design & Architecture",
		"Computer Science",
		"System Design",
		"AI Engineer",
		"Data Engineer",
		"Machine Learning Engineer",
		"Data Scientist",
		"Rust Developer",
	}

	for _, topicTheme := range sTopicThemes {
		result := db.FirstOrCreate(&models.Topic{}, &models.Topic{Title: topicTheme})
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Fatalf("Failed to create theme %s: %v", topicTheme, result.Error)
		}
	}

	log.Print("Topics seed completed")
}
