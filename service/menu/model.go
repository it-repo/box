package menu

import (
	"github.com/jinzhu/gorm"
)

// BoxMenu Model
type BoxMenu struct {
	gorm.Model
	FatherID uint
	Name     string `gorm:"not null"`
	Desc     string `gorm:"not null"`
	Sort     int    `gorm:"not null"`
	Path     string `gorm:"not null"`
	Access   string `gorm:"not null"` // perm name
}
