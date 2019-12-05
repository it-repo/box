package box

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
	"github.com/ddosakura/sola/v2/middleware/auth"
	"github.com/ddosakura/sola/v2/middleware/router"
	"github.com/it-repo/box/service/route"
	"github.com/jinzhu/gorm"
)

// Route Middleware
func Route(db *gorm.DB, r *router.Router) {
	s := route.New(db)
	r.Bind("/routes", func(c sola.Context) error {
		perms := toStringArray(auth.Claims(c, "perms"))
		tree := s.Route(perms...)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			"msg":  "SUCCESS",
			"data": tree,
		})
	})
}
