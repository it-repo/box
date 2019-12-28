package ac

import (
	"time"

	"github.com/jinzhu/gorm"
)

// SelectByID -
func (s *Srv) SelectByID(id uint) *BoxUser {
	u := BoxUser{}
	u.ID = id
	if err := s.db.First(&u).Error; err != nil {
		return nil
	}
	return &u
}

// TODO

// Model -
type Model struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserInfo -
type UserInfo struct {
	Model
	Roles  []BoxRole
	Name   string `json:"name"`
	Nick   string `json:"nick,omitempty"`
	Avatar string `json:"avatar,omitempty"`
	Desc   string `json:"desc,omitempty"`
}

// SelectUserList -
func (s *Srv) SelectUserList(page, size int) []UserInfo {
	var list []BoxUser
	var count int
	db := s.db.Order("id desc").Limit(size).Offset((page - 1) * size)
	if err := db.Find(&list).Count(&count).Error; err != nil {
		return []UserInfo{}
	}
	return turn(list, count)
}

// SelectUsertotal -
func (s *Srv) SelectUsertotal() int {
	var count int
	if err := s.db.Model(&BoxUser{}).Count(&count).Error; err != nil {
		return 0
	}
	return count
}

func turn(list []BoxUser, count int) []UserInfo {
	relist := make([]UserInfo, 0, len(list))
	for _, x := range list {
		relist = append(relist, UserInfo{
			Model: Model{
				ID:        x.ID,
				CreatedAt: x.CreatedAt,
				UpdatedAt: x.UpdatedAt,
			},
			Roles:  x.Roles,
			Name:   x.Name,
			Nick:   x.Nick,
			Avatar: x.Avatar,
			Desc:   x.Desc,
		})
	}
	return relist
}

//SelectUser -
func (s *Srv) SelectUser(id uint) *UserInfo {
	x := BoxUser{}
	x.ID = id
	if err := s.db.First(&x).Error; err != nil {
		return nil
	}
	relist := UserInfo{}
	relist.ID = x.ID
	relist.CreatedAt = x.CreatedAt
	relist.UpdatedAt = x.UpdatedAt
	relist.Roles = x.Roles
	relist.Name = x.Name
	relist.Nick = x.Nick
	relist.Avatar = x.Avatar
	relist.Desc = x.Desc
	return &relist
}

//DeleteUser -
func (s *Srv) DeleteUser(id []string) error {
	list := []BoxUser{}
	for _, x := range id {
		s.db.Where("id=?", x).Delete(&list)
	}
	return nil
}

//PostUser -
func (s *Srv) PostUser(name, pass string, id []uint) error {
	u := BoxUser{
		Name: name,
		Nick: name,
		Pass: pass,
	}
	if e := s.db.Create(&u).Error; e != nil {
		return e
	}
	list := make([]interface{}, 0, len(id))
	for _, x := range id {
		list = append(list, &BoxRole{Model: gorm.Model{ID: x}})
	}
	return s.db.
		Model(&BoxUser{Model: gorm.Model{ID: u.ID}}).
		Association("Roles").
		Append(list...).
		Error
}

//PutUser -
func (s *Srv) PutUser(id uint, nick, pass, avatar, desc string, rid []uint) error {
	u := BoxUser{
		Nick:   nick,
		Pass:   pass,
		Avatar: avatar,
		Desc:   desc,
	}
	if err := s.db.Where("id=?", id).Table("box_users").Updates(&u).Error; err != nil {
		return err
	}

	list := make([]interface{}, 0, len(rid))
	for _, x := range rid {
		list = append(list, &BoxRole{Model: gorm.Model{ID: x}})
	}
	return s.db.
		Model(&BoxUser{Model: gorm.Model{ID: id}}).
		Association("Roles").
		Replace(list...).
		Error
}

//Role

//GetListRole -
func (s *Srv) GetListRole(page, size int) []BoxRole {
	var list []BoxRole
	db := s.db.Limit(size).Offset((page - 1) * size)
	if err := db.Find(&list).Error; err != nil {
		return []BoxRole{}
	}
	return list
}

//GetRoleCount 获取角色数量
func (s *Srv) GetRoleCount() int {
	var count int
	if err := s.db.Model(&BoxRole{}).Count(&count).Error; err != nil {
		return 0
	}
	return count
}

//GetRole 获取角色信息
func (s *Srv) GetRole(id uint) *BoxRole {
	x := BoxRole{}
	x.ID = id
	if err := s.db.First(&x).Error; err != nil {
		return nil
	}
	return &x
}

//PostRole 增加角色
func (s *Srv) PostRole(name, desc string, rid []uint) error {
	r := BoxRole{
		Name: name,
		Desc: desc,
	}
	if err := s.db.Create(&r).Error; err != nil {
		return err
	}
	list := make([]interface{}, 0, len(rid))
	for _, X := range rid {
		list = append(list, &BoxPerm{Model: gorm.Model{ID: X}})
	}
	return s.db.
		Model(&BoxRole{Model: gorm.Model{ID: r.ID}}).
		Association("Perms").
		Append(list...).
		Error
}

//DeleteRole 删除角色
func (s *Srv) DeleteRole(id []string) error {
	x := BoxRole{}
	for _, i := range id {
		s.db.Where("id=?", i).Delete(&x)
	}
	return nil
}

//PutRole 更新角色
func (s *Srv) PutRole(name, desc string, id uint, rid []uint) error {
	r := &BoxRole{
		Name: name,
		Desc: desc,
	}
	if err := s.db.Where("id=?", id).Table("box_roles").Updates(&r).Error; err != nil {
		return err
	}
	list := make([]interface{}, 0, len(rid))
	for _, X := range rid {
		list = append(list, &BoxPerm{Model: gorm.Model{ID: X}})
	}
	return s.db.
		Model(&BoxRole{Model: gorm.Model{ID: id}}).
		Association("Perms").
		Replace(list...).
		Error
}

// GetPermCount 获取总数
func (s *Srv) GetPermCount() int {
	var count int
	if err := s.db.Model(&BoxPerm{}).Count(&count).Error; err != nil {
		return 0
	}
	return count
}

// GetListPerm 获取列表
func (s *Srv) GetListPerm(page, size int) []BoxPerm {
	var list []BoxPerm
	db := s.db.Limit(size).Offset((page - 1) * size)
	if err := db.Find(&list).Error; err != nil {
		return nil
	}
	return list
}

// GetPerm 查询权限
func (s *Srv) GetPerm(id int) *BoxPerm {
	list := BoxPerm{}
	s.db.Where("id=?", id).First(&list)
	return &list
}

//PostPerm 增加权限
func (s *Srv) PostPerm(a BoxPerm) error {
	list := BoxPerm{
		Name: a.Name,
		Desc: a.Desc,
	}
	return s.db.Table("box_perms").Create(&list).Error
}

//DelPerm 删除权限
func (s *Srv) DelPerm(id []string) error {
	var list BoxPerm
	for _, i := range id {
		if err := s.db.Where("id=?", i).Delete(&list).Error; err != nil {
			return err
		}
	}
	return nil
}

//PutPerm 更新权限
func (s *Srv) PutPerm(id int, a BoxPerm) error {
	list := BoxPerm{
		Name: a.Name,
		Desc: a.Desc,
	}
	return s.db.Where("id=?", id).Updates(&list).Error
}
