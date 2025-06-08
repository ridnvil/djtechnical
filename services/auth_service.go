package services

import (
	"DeallsJobsTest/models"
	"DeallsJobsTest/utils"
	"gorm.io/gorm"
	"strconv"
)

func GenerateUserAdmin(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.User{}).Where("role = ?", "admin").Where("username = ?", "admin").Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	hashPassword, errhash := utils.HashPassword("admin123")
	if errhash != nil {
		return errhash
	}

	adminUser := models.User{
		Username:  "admin",
		Password:  hashPassword,
		Role:      "admin",
		FullName:  "Administrator",
		CreatedBy: nil,
		UpdatedBy: nil,
		RequestIP: "127.0.0.1",
	}

	if err := db.FirstOrCreate(&adminUser).Error; err != nil {
		return err
	}

	return nil
}

func Generate100Employees(db *gorm.DB) error {
	for i := 1; i <= 100; i++ {
		idString := strconv.Itoa(i)
		username := "employee" + idString
		hashPassword, errhash := utils.HashPassword("employee123")
		if errhash != nil {
			return errhash
		}

		employee := models.User{
			Username:  username,
			Password:  hashPassword,
			Role:      "employee",
			FullName:  "Employee " + idString,
			CreatedBy: nil,
			UpdatedBy: nil,
			RequestIP: "127.0.0.1",
		}

		if err := db.Table("users").FirstOrCreate(&employee).Error; err != nil {
			return err
		}
	}
	return nil
}
