package ac

import (
	"github.com/jinzhu/gorm"
)

// Srv AC
type Srv struct {
	db *gorm.DB
}

// New AC Srv
func New(db *gorm.DB) *Srv {
	db.AutoMigrate(&BoxUser{}, &BoxRole{}, &BoxPerm{})
	return &Srv{db}
}

// Login Srv
func (s *Srv) Login(name, pass string) *BoxUser {
	var u BoxUser
	if err := s.db.First(&u, "name = ?", name).Error; err != nil {
		return nil
	}
	if u.Pass != pass {
		return nil
	}
	return &u
}

// Register Srv
func (s *Srv) Register(name, pass string) bool {
	if err := s.db.Create(&BoxUser{
		Name: name,
		Pass: pass,
	}).Error; err != nil {
		return false
	}
	return true
}

func (s *Srv) roles(u *BoxUser) []BoxRole {
	var roles []BoxRole
	if err := s.db.Model(u).Related(&roles, "Roles").Error; err != nil {
		return []BoxRole{}
	}
	return roles
}

// Roles Selector
func (s *Srv) Roles(u *BoxUser) []string {
	roles := s.roles(u)
	result := make([]string, 0, len(roles))
	for _, role := range roles {
		result = append(result, role.Name)
	}
	return result
}

func (s *Srv) perms(roles []BoxRole) []BoxPerm {
	var perms []BoxPerm
	if err := s.db.Model(&roles).Related(&perms, "Perms").Error; err != nil {
		return []BoxPerm{}
	}
	return perms
}

// Perms Selector
func (s *Srv) Perms(u *BoxUser) []string {
	perms := s.perms(s.roles(u))
	result := make([]string, 0, len(perms))
	for _, perm := range perms {
		result = append(result, perm.Name)
	}
	return result
}
