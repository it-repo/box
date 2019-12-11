package route

//GetRouterlist 获取菜单列表
func (s *Srv) GetRouterlist(page, size int) []BoxRoute {
	var list []BoxRoute
	db := s.db.Limit(size).Offset((page - 1) * size)
	if err := db.Find(&list).Error; err != nil {
		return []BoxRoute{}
	}
	return list
}

//GetRouterCount  获取总数
func (s *Srv) GetRouterCount() int {
	var count int
	if err := s.db.Model(&BoxRoute{}).Count(&count).Error; err != nil {
		return 0
	}
	return count
}

//GetRouter 查询菜单
func (s *Srv) GetRouter(id int) *BoxRoute {
	x := BoxRoute{}
	s.db.Where("id=?", id).First(&x)
	return &x
}

//PostRouter 增加菜单
func (s *Srv) PostRouter(a BoxRoute) error {
	list := BoxRoute{
		FatherID:   a.FatherID,
		Name:       a.Name,
		Desc:       a.Desc,
		Sort:       a.Sort,
		Path:       a.Path,
		Perm:       a.Perm,
		Component:  a.Component,
		AlwaysShow: a.AlwaysShow,
		Hidden:     a.Hidden,
		Redirect:   a.Redirect,
		Title:      a.Title,
		Icon:       a.Icon,
		NoCache:    a.NoCache,
		Breadcrumb: a.Breadcrumb,
		Affix:      a.Affix,
	}
	return s.db.Table("Box_routes").Create(&list).Error
}

//DelRouter 删除菜单
func (s *Srv) DelRouter(id []string) error {
	route := BoxRoute{}
	for _, x := range id {
		s.db.Where("id=?", x).Delete(&route)
	}
	return nil
}

//PutRouter 更新菜单
func (s *Srv) PutRouter(id int, a BoxRoute) error {
	list := BoxRoute{
		Name:       a.Name,
		Desc:       a.Desc,
		Sort:       a.Sort,
		Path:       a.Path,
		Perm:       a.Perm,
		Component:  a.Component,
		AlwaysShow: a.AlwaysShow,
		Hidden:     a.Hidden,
		Redirect:   a.Redirect,
		Title:      a.Title,
		Icon:       a.Icon,
		NoCache:    a.NoCache,
		Breadcrumb: a.Breadcrumb,
		Affix:      a.Affix,
	}
	return s.db.Where("id=?", id).Table("Box_routes").Updates(&list).Error
}
