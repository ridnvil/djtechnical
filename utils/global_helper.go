package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"path/filepath"
)

func AutoCreateFolder(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			log.Errorf("failed create uploads dir: %w", err)
		}
	}
	return nil
}

func DeleteAllFilesInFolder(folder string) error {
	files, err := os.ReadDir(folder)
	if err != nil {
		return fmt.Errorf("failed read dir: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			err := os.Remove(filepath.Join(folder, file.Name()))
			if err != nil {
				return fmt.Errorf("failde delete file %s: %w", file.Name(), err)
			}
		}
	}
	return nil
}

func RandomSalary(min, max int) float64 {
	return float64(rand.Intn(max-min+1) + min)
}

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
