package data

import (
	"educ-gpt/models"
	"log"
)

func Seed() {
	rolesSeed()
	topicsSeed()
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
