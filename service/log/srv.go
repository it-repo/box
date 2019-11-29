package log

import (
	"github.com/jinzhu/gorm"
)

// Srv Route
type Srv struct {
	db *gorm.DB
}

// New Route Srv
func New(db *gorm.DB) *Srv {
	db.AutoMigrate(&BoxLog{})
	return &Srv{db}
}

// Insert Action
func (s *Srv) Insert(uid uint, name string, desc string) error {
	return s.db.Create(&BoxLog{
		UID:  uid,
		Name: name,
		Desc: desc,
	}).Error
}

// List Action
func (s *Srv) List(page, size int) []BoxLog {
	var list []BoxLog
	db := s.db.Order("id desc").Limit(size).Offset((page - 1) * size)
	if err := db.Find(&list).Error; err != nil {
		return []BoxLog{}
	}
	return list
}

// Total Action
func (s *Srv) Total(page, size int) int {
	count := 0
	if err := s.db.Model(&BoxLog{}).Count(&count).Error; err != nil {
		return 0
	}
	return count
}
