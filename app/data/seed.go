package data

import (
	"educ-gpt/logger"
	"educ-gpt/models"
	"log"
)

func Seed() {
	rolesSeed()
}

func rolesSeed() {
	sRoles := []*models.Role{
		{Name: "ADMIN"},
		{Name: "USER"},
	}

	for _, role := range sRoles {
		result := db.FirstOrCreate(role, models.Role{Name: role.Name})
		if result.Error != nil {
			log.Fatalf("Failed to create or find role %s: %v", role.Name, result.Error)
		}
	}

	logger.Logger().Debug("Roles seed completed")
}
