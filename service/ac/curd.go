package ac

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
