package ac

import (
	"github.com/jinzhu/gorm"
)

// BoxUser Model
type BoxUser struct {
	gorm.Model
	Roles []BoxRole `gorm:"many2many:box_user_roles;"`
	Name  string    `gorm:"not null;unique"`
	Nick  string    `gorm:"not null;unique"`
	Pass  string    `gorm:"not null"`
}

// BoxRole Model
type BoxRole struct {
	gorm.Model
	Users []BoxPerm `gorm:"many2many:box_user_roles;"`
	Perms []BoxPerm `gorm:"many2many:box_role_perms;"`
	Name  string    `gorm:"not null;unique"`
	Desc  string    `gorm:"not null"`
}

// BoxPerm Model
type BoxPerm struct {
	gorm.Model
	Roles []BoxRole `gorm:"many2many:box_role_perms;"`
	Name  string    `gorm:"not null;unique"`
	Desc  string    `gorm:"not null"`
}
