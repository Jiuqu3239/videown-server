package model

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&VideoMetadata{},
		&Activity{},
	)
}
