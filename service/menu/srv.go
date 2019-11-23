package menu

import (
	"github.com/jinzhu/gorm"
)

// Srv Menu
type Srv struct {
	db *gorm.DB
}

// New Menu Srv
func New(db *gorm.DB) *Srv {
	db.AutoMigrate(&BoxMenu{})
	return &Srv{db}
}

// Menu Tree
type Menu struct {
	*BoxMenu
	Children []*Menu `json:"children,omitempty"`
}

const sqlSubMenu = "father_id = ? AND (access = '' OR access in (?))"

// Menu with perms
func (s *Srv) Menu(fid uint, perms ...string) *Menu {
	root := &Menu{}
	var list []BoxMenu
	if err := s.db.Where(sqlSubMenu, fid, perms).Find(&list).Error; err != nil {
		return nil
	}
	if len(list) == 0 {
		return nil
	}
	root.Children = make([]*Menu, 0, len(list))
	for i := range list {
		sub := list[i]
		m := &Menu{&sub, nil}
		menu := s.Menu(sub.ID, perms...)
		if menu != nil {
			m.Children = menu.Children
		}
		root.Children = append(root.Children, m)
	}
	return root
}
