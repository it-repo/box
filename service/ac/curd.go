package ac

import "time"

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
	db := s.db.Order("id desc").Limit(size).Offset((page - 1) * size)
	if err := db.Find(&list).Error; err != nil {
		return []UserInfo{}
	}
	// var relist []UserInfo
	// relist=list
	return turn(list)
}

func turn(list []BoxUser) []UserInfo {
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
func (s *Srv) DeleteUser(id int) error {
	list := BoxUser{}
	return s.db.Where("id=?", id).Delete(&list).Error
}

//PostUser -
func (s *Srv) PostUser(name, pass string) error {
	list := BoxUser{
		Name: name,
		Nick: name,
		Pass: pass,
	}
	return s.db.Create(&list).Error
}

//PutUser -
func (s *Srv) PutUser(id int, nick, pass, avatar, desc string) error {
	list := BoxUser{
		Nick:   nick,
		Pass:   pass,
		Avatar: avatar,
		Desc:   desc,
	}
	return s.db.Where("id=?", id).Table("box_users").Updates(&list).Error
}
