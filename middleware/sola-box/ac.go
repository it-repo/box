package box

import (
	"net/http"

	"github.com/ddosakura/sola/v2"
	"github.com/ddosakura/sola/v2/middleware/auth"
	"github.com/ddosakura/sola/v2/middleware/router"
	"github.com/it-repo/box/service/ac"
	"github.com/jinzhu/gorm"
)

// Ctx
const (
	CtxBoxAC   = "box.ac"
	CtxBoxACR  = "box.ac.rule"
	CtxBoxACRT = "box.ac.rule.type"
	CtxBoxACRL = "box.ac.rule.logical"
)

// ACR - AC Rule Builder
func ACR(t ac.Type, l ac.Logical, rules ...string) sola.Handler {
	return func(c sola.Context) error {
		c[CtxBoxACRT] = t
		c[CtxBoxACRL] = l
		c[CtxBoxACR] = rules
		return nil
	}
}

// RequestAC Check
type RequestAC func(arc, h sola.Handler) sola.Handler

// AC Middleware & Srv
func AC(db *gorm.DB, key string) (sola.Middleware, RequestAC) {
	_sign := auth.Sign(auth.AuthJWT, []byte(key))
	_auth := auth.Auth(auth.AuthJWT, []byte(key))
	s := ac.New(db)
	r := router.New()
	r.BindFunc("POST /login", auth.NewFunc(_sign, login, loginSuccess))
	r.BindFunc("/logout", auth.CleanFunc(success))
	r.BindFunc("POST /register", register) // TODO: remove

	routes := sola.Merge(func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			c[CtxBoxAC] = s
			return next(c)
		}
	}, r.Routes())
	requestAC := func(acr, h sola.Handler) sola.Handler {
		return auth.NewFunc(_auth, nil, func(c sola.Context) error {
			acr(c)
			t := c[CtxBoxACRT].(ac.Type)
			l := c[CtxBoxACRL].(ac.Logical)
			acr := c[CtxBoxACR].([]string)
			var rule *ac.Rule
			if l == ac.LogicalAND {
				rule = ac.NewRule(acr, true)
			} else {
				rule = ac.NewRule(acr, false)
			}
			var rules []string
			if t == ac.TypeRole {
				rules = toStringArray(auth.Claims(c, "roles"))
			} else {
				rules = toStringArray(auth.Claims(c, "perms"))
			}
			if rule.Check(rules) {
				return h(c)
			}
			return c.String(http.StatusForbidden, "Forbidden")
		})
	}
	return routes, requestAC
}

func login(next sola.Handler) sola.Handler {
	return func(c sola.Context) error {
		s := c[CtxBoxAC].(*ac.Srv)
		r := c.Request()
		name := r.PostFormValue("name")
		pass := r.PostFormValue("pass")

		if name == "" || pass == "" {
			return fail(c)
		}
		u := s.Login(name, pass)
		if u == nil {
			return fail(c)
		}

		roles := s.Roles(u)
		perms := s.Perms(u)
		c[auth.CtxClaims] = map[string]interface{}{
			"id":    u.ID,
			"name":  u.Name,
			"roles": roles,
			"perms": perms,
		}
		return next(c)
	}
}

func loginSuccess(c sola.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "SUCCESS",
		"data": c[auth.CtxClaims],
	})
}

func register(c sola.Context) error {
	s := c[CtxBoxAC].(*ac.Srv)
	r := c.Request()
	name := r.PostFormValue("name")
	pass := r.PostFormValue("pass")

	if name == "" || pass == "" {
		return fail(c)
	}

	if s.Register(name, pass) {
		return success(c)
	}
	return fail(c)
}
