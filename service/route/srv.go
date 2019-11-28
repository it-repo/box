package route

import (
	"github.com/jinzhu/gorm"
)

// Srv Route
type Srv struct {
	db *gorm.DB
}

// New Route Srv
func New(db *gorm.DB) *Srv {
	db.AutoMigrate(&BoxRoute{})
	return &Srv{db}
}

// Route Tree
type Route struct {
	*BoxRoute
	Children []*Route `json:"children,omitempty"`
}

// const sqlSubRoute = "father_id = ? AND (perm = '' OR perm in (?))"
// // Route with perms
// func (s *Srv) Route(fid uint, perms ...string) *Route {
// 	root := &Route{}
// 	var l []BoxRoute
// 	if err := s.db.Where(sqlSubRoute, fid, perms).Find(&l).Error; err != nil {
// 		return nil
// 	}
// 	if len(l) == 0 {
// 		return nil
// 	}
// 	root.Children = make([]*Route, 0, len(l))
// 	for i := range l {
// 		sub := l[i]
// 		m := &Route{&sub, nil}
// 		route := s.Route(sub.ID, perms...)
// 		if route != nil {
// 			m.Children = route.Children
// 		}
// 		root.Children = append(root.Children, m)
// 	}
// 	return root
// }

const sqlSubRoute = "perm = '' OR perm in (?)"

// Route with perms
func (s *Srv) Route(perms ...string) []*Route {
	var l []BoxRoute
	if err := s.db.Where(sqlSubRoute, perms).Find(&l).Error; err != nil {
		return []*Route{}
	}
	if len(l) == 0 {
		return []*Route{}
	}
	m := make(map[uint][]*Route)
	for i := range l {
		r := l[i]
		if m[r.FatherID] == nil {
			m[r.FatherID] = make([]*Route, 0)
		}
		route := &Route{BoxRoute: &r}
		m[r.FatherID] = append(m[r.FatherID], route)
	}
	return getRoutes(m, 0)
}

func getRoutes(m map[uint][]*Route, fid uint) []*Route {
	for _, r := range m[fid] {
		r.Children = getRoutes(m, r.ID)
	}
	return m[fid]
}
