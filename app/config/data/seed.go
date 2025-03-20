package data

import (
	"educ-gpt/models"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	activate_at := time.Time{}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.DefaultCost)

	if err != nil {
		log.Fatal(err)
	}

	user := &models.User{
		Name:       "admin",
		Email:      os.Getenv("ADMIN_EMAIL"),
		Password:   string(hashedPassword),
		ActivateAt: &activate_at,
	}

	var exists bool

	if err := db.Model(&models.User{}).Select("count(*) > 0").Where("name = ? OR email = ?", user.Name, user.Email).Find(&exists).Error; err != nil {
		log.Fatalf("Failed to create %s: %v", user.Name, err)
	}

	if !exists {
		if err := db.FirstOrCreate(&models.User{}, user).Error; err != nil {
			log.Fatalf("Failed to create %s: %v", user.Name, err)
		}
	} else {
		if err := db.Where("name = ? AND email = ?", user.Name, user.Email).First(user).Error; err != nil {
			log.Fatalf("Failed to create %s: %v", user.Name, err)
		}
	}

	if err := db.FirstOrCreate(&models.UserRoles{}, &models.UserRoles{UserID: user.ID, RoleID: 1}).Error; err != nil {
		log.Fatalf("Failed to create %s: %v", user.Name, err)
	}

	log.Print("Admin seed completed")
}

func usersSeed() {
	activate_at := time.Time{}

	users := make([]*models.User, 10)

	for i := 0; i < 10; i++ {
		iItoa := strconv.Itoa(i)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("user_user"), bcrypt.DefaultCost)

		if err != nil {
			log.Fatal(err)
		}

		users[i] = &models.User{
			Name:       "user" + iItoa,
			Email:      "user" + iItoa + "@educgpt.ru",
			Password:   string(hashedPassword),
			ActivateAt: &activate_at,
		}
	}

	for _, user := range users {
		var exists bool

		if err := db.Model(&models.User{}).Select("count(*) > 0").Where("name = ? OR email = ?", user.Name, user.Email).Find(&exists).Error; err != nil {
			log.Fatalf("Failed to create %s: %v", user.Name, err)
		}

		if !exists {
			if err := db.FirstOrCreate(&models.User{}, user).Error; err != nil {
				log.Fatalf("Failed to create %s: %v", user.Name, err)
			}
		} else {
			if err := db.Where("name = ? AND email = ?", user.Name, user.Email).First(user).Error; err != nil {
				log.Fatalf("Failed to create %s: %v", user.Name, err)
			}
		}

		if err := db.FirstOrCreate(&models.UserRoles{}, &models.UserRoles{UserID: user.ID, RoleID: 1}).Error; err != nil {
			log.Fatalf("Failed to create %s: %v", user.Name, err)
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
		if result.Error != nil {
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
		if result.Error != nil {
			log.Fatalf("Failed to create theme %s: %v", topicTheme, result.Error)
		}
	}

	log.Print("Topics seed completed")
}
