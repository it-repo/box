package box

import (
	"net/http"
	"strconv"

	"github.com/ddosakura/sola/v2/middleware/logger"

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
		c.Set(CtxBoxACRT, t)
		c.Set(CtxBoxACRL, l)
		c.Set(CtxBoxACR, rules)
		return nil
	}
}

// RequestAC Check
type RequestAC func(arc, h sola.Handler) sola.Handler

// AC Middleware & Srv
func AC(db *gorm.DB, key string) (sola.Middleware, *router.Router, RequestAC) {
	_sign := auth.Sign(auth.AuthJWT, []byte(key))
	_auth := auth.Auth(auth.AuthJWT, []byte(key))
	s := ac.New(db)
	r := router.New()
	r.BindFunc("POST /login", auth.NewFunc(_sign, login, loginSuccess))
	r.BindFunc("/logout", auth.CleanFunc(success))
	r.BindFunc("/info", auth.NewFunc(_auth, nil, userInfo))
	r.BindFunc("POST /register", register) // TODO: remove

	init := func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			c.Set(CtxBoxAC, s)
			return next(c)
		}
	}
	requestAC := func(acr, h sola.Handler) sola.Handler {
		return auth.NewFunc(_auth, nil, func(c sola.Context) error {
			acr(c)
			t := c.Get(CtxBoxACRT).(ac.Type)
			l := c.Get(CtxBoxACRL).(ac.Logical)
			acr := c.Get(CtxBoxACR).([]string)
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

	acr1 := ACR(ac.TypeRole, ac.LogicalOR, "admin")
	r.BindFunc("GET /user/:id", requestAC(acr1, getUser))
	r.BindFunc("GET /user", requestAC(acr1, getUsers))
	r.BindFunc("DELETE /user/:id", requestAC(acr1, delUser))
	r.BindFunc("POST /user", requestAC(acr1, postUser))
	r.BindFunc("PUT /user/:id", requestAC(acr1, putUser))

	return init, r, requestAC
}

// ReqUser Form
type ReqUser struct {
	Username string
	Password string
}

func login(next sola.Handler) sola.Handler {
	return func(c sola.Context) error {
		s := c.Get(CtxBoxAC).(*ac.Srv)

		var a ReqUser
		if err := c.GetJSON(&a); err != nil {
			return err
		}

		if a.Username == "" || a.Password == "" {
			return fail(c)
		}
		u := s.Login(a.Username, a.Password)
		if u == nil {
			return fail(c)
		}

		roles := s.Roles(u)
		perms := s.Perms(u)
		c.Set(auth.CtxClaims, map[string]interface{}{
			"id":    u.ID,
			"name":  u.Name,
			"roles": roles,
			"perms": perms,
		})
		logger.Action(c, u.ID, "login", u.Nick+" login system!")
		return next(c)
	}
}

func loginSuccess(c sola.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "SUCCESS",
		"data": map[string]interface{}{
			"token": c.Get(auth.CtxToken),
		},
	})
}

// TODO: refresh roles/perms

func userInfo(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	id := auth.Claims(c, "id").(float64)
	u := s.SelectByID(uint(id))
	if u == nil {
		return fail(c)
	}
	roles := s.Roles(u)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "SUCCESS",
		"data": map[string]interface{}{
			"name":         u.Name,
			"nick":         u.Nick,
			"avatar":       u.Avatar,
			"introduction": u.Desc, // TODO: 前端兼容，暂时改名
			"roles":        roles,
		},
	})
}

// TODO: need change to getJSON
func register(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
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

// CRUD

func acSucc(c sola.Context, v interface{}) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "SUCCESS",
		"data": v,
	})
}

func getUsers(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	r := c.Request()
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = 5
	}
	list := s.SelectUserList(page, size)
	return acSucc(c, list)
}

func getUser(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	id, err := strconv.Atoi(router.Param(c, "id"))
	if err != nil {
		return err
	}
	u := s.SelectUser(uint(id))
	return acSucc(c, u)
}

func delUser(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	id, e := strconv.Atoi(router.Param(c, "id"))
	if e != nil {
		return fail(c)
	}
	s.DeleteUser(id)
	return acSucc(c, nil)
}

func postUser(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	var a ReqUser
	if e := c.GetJSON(&a); e != nil {
		return e
	}
	if e := s.PostUser(a.Username, a.Password); e != nil {
		return fail(c)
	}
	return acSucc(c, nil)
}

// UserInfo -
type UserInfo struct {
	Password string
	Nick     string
	Avatar   string
	Desc     string
}

func putUser(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	var a UserInfo
	if e := c.GetJSON(&a); e != nil {
		return fail(c)
	}
	id, e := strconv.Atoi(router.Param(c, "id"))
	if e != nil {
		return fail(c)
	}
	if e := s.PutUser(id, a.Nick, a.Password, a.Avatar, a.Desc); e != nil {
		return fail(c)
	}
	return acSucc(c, nil)
}
