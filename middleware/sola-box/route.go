package box

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ddosakura/sola/v2"
	"github.com/ddosakura/sola/v2/middleware/auth"
	"github.com/ddosakura/sola/v2/middleware/router"
	"github.com/it-repo/box/service/ac"
	"github.com/it-repo/box/service/route"
	"github.com/jinzhu/gorm"
)

// Ctx
const (
	CtxBoxRoute = "box.route"
)

// Route Middleware
func Route(db *gorm.DB, r *router.Router) {
	s := route.New(db)
	r.Use(func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			c.Set(CtxBoxRoute, s)
			return next(c)
		}
	})
	r.Bind("/routes", func(c sola.Context) error {
		perms := toStringArray(auth.Claims(c, "perms"))
		tree := s.Route(perms...)
		return c.JSON(http.StatusOK, map[string]interface{}{
			"code": 0,
			"msg":  "SUCCESS",
			"data": tree,
		})
	})

	acRequest := func(a, b sola.Handler) sola.Handler {
		return b
	}
	acr1 := ACR(ac.TypeRole, ac.LogicalOR, "admin")
	r.Bind("GET /router/:id", acRequest(acr1, getRouter))
	r.Bind("GET /router", acRequest(acr1, getRouters))
	r.Bind("DELETE /router/:id", acRequest(acr1, delRouter))
	r.Bind("POST /router", acRequest(acr1, postRouter))
	r.Bind("PUT /router/:id", acRequest(acr1, putRouter))
}

func getRouters(c sola.Context) error {
	s := c.Get(CtxBoxRoute).(*route.Srv)
	r := c.Request()
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		return err
	}
	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "SUCCESS",
		"data": map[string]interface{}{
			"total":  s.GetRouterCount(),
			"Routes": s.GetRouterlist(page, size),
		},
	})
}

func getRouter(c sola.Context) error {
	s := c.Get(CtxBoxRoute).(*route.Srv)
	id, err := strconv.Atoi(router.Param(c, "id"))
	if err != nil {
		return err
	}

	return acSucc(c, s.GetRouter(id))
}

func postRouter(c sola.Context) error {
	var a route.BoxRoute
	s := c.Get(CtxBoxRoute).(*route.Srv)
	if err := c.GetJSON(&a); err != nil {
		return err
	}
	s.PostRouter(a)
	return acSucc(c, nil)
}

func delRouter(c sola.Context) error {
	s := c.Get(CtxBoxRoute).(*route.Srv)
	list := router.Param(c, "id")
	id := strings.Split(list, ",")
	s.DelRouter(id)
	return acSucc(c, nil)
}

func putRouter(c sola.Context) error {
	s := c.Get(CtxBoxRoute).(*route.Srv)
	id, err := strconv.Atoi(router.Param(c, "id"))
	if err != nil {
		return err
	}
	var a route.BoxRoute
	if err := c.GetJSON(&a); err != nil {
		return err
	}
	s.PutRouter(id, a)
	return acSucc(c, nil)
}
