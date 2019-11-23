package box

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
	"github.com/ddosakura/sola/v2/middleware/auth"
	"github.com/ddosakura/sola/v2/middleware/router"
	"github.com/it-repo/box/service/menu"
	"github.com/jinzhu/gorm"
)

// Menu Middleware
func Menu(db *gorm.DB) *router.Router {
	s := menu.New(db)
	r := router.New()
	r.BindFunc("/menu", func(c sola.Context) error {
		perms := toStringArray(auth.Claims(c, "perms"))
		tree := s.Menu(0, perms...)
		if tree == nil || tree.Children == nil {
			tree = &menu.Menu{
				Children: []*menu.Menu{},
			}
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			"msg":  "SUCCESS",
			"data": tree.Children,
		})
	})
	return r
}
