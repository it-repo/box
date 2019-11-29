package route

import (
	"github.com/jinzhu/gorm"
)

// BoxRoute Model
type BoxRoute struct {
	gorm.Model
	FatherID uint
	Name     string `gorm:"not null" json:"name"`
	Desc     string `gorm:"not null" json:"desc,omitempty"`
	Sort     int    `gorm:"not null" json:"sort"`
	Path     string `gorm:"not null" json:"path"`
	Perm     string `gorm:"not null" json:"perm,omitempty"`

	Component  string `json:"component,omitempty"`
	AlwaysShow bool   `json:"alwaysShow"`
	Hidden     bool   `json:"hidden"`
	Redirect   string `json:"redirect,omitempty"`
	Title      string `json:"title,omitempty"`
	Icon       string `json:"icon,omitempty"`
	NoCache    bool   `json:"noCache"`
	Breadcrumb bool   `json:"breadcrumb"`

	Affix bool `json:"affix"`
}
