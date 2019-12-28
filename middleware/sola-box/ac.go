package box

import (
	"net/http"
	"strconv"
	"strings"

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

// ACRequest Check
type ACRequest func(arc, h sola.Handler) sola.Handler

// AC Middleware & Srv
func AC(db *gorm.DB, k string, r *router.Router) (sola.Middleware, ACRequest) {
	r.Bind("/logout", auth.CleanFunc(success))

	s := ac.New(db)
	r.Use(func(next sola.Handler) sola.Handler {
		return func(c sola.Context) error {
			c.Set(CtxBoxAC, s)
			return next(c)
		}
	})
	jwtSign, jwtAuth := auth.NewJWT([]byte(k))

	{
		sub := r.Sub(nil)
		sub.Use(auth.LoadAuthCache)
		h := sola.MergeFunc(loginSuccess, login, jwtSign)
		sub.Bind("POST /login", h)
	}
	r.Use(jwtAuth)
	r.Bind("/info", userInfo)

	acRequest := func(acr, h sola.Handler) sola.Handler {
		if h == nil {
			panic(ErrNoHandler)
		}
		return func(c sola.Context) error {
			if acr == nil {
				return h(c)
			}
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
		}
	}

	acr1 := ACR(ac.TypeRole, ac.LogicalOR, "admin")
	r.Bind("GET /user/:id", acRequest(acr1, getUser))
	r.Bind("GET /user", acRequest(acr1, getUsers))
	r.Bind("DELETE /user/:id", acRequest(acr1, delUser))
	r.Bind("POST /user", acRequest(acr1, postUser))
	r.Bind("PUT /user/:id", acRequest(acr1, putUser))

	r.Bind("GET /role/:id", acRequest(acr1, getRole))
	r.Bind("GET /role", acRequest(acr1, getRoles))
	r.Bind("DELETE /role/:id", acRequest(acr1, delRole))
	r.Bind("POST /role", acRequest(acr1, postRole))
	r.Bind("PUT /role/:id", acRequest(acr1, putRole))

	r.Bind("GET /perm/:id", acRequest(acr1, getPerm))
	r.Bind("GET /perm", acRequest(acr1, getPerms))
	r.Bind("DELETE /perm/:id", acRequest(acr1, delPerm))
	r.Bind("POST /perm", acRequest(acr1, postPerm))
	r.Bind("PUT /perm/:id", acRequest(acr1, putPerm))

	return jwtAuth, acRequest
}

// ReqUser Form
type ReqUser struct {
	Username string
	Password string
	Rid      []uint
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

// CRUD

func acSucc(c sola.Context, v interface{}) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "SUCCESS",
		"data": v,
	})
}

//user
func getUsers(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	page, size := getPageSize(c)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "SUCCESS",
		"data": map[string]interface{}{
			"user":  s.SelectUserList(page, size),
			"total": s.SelectUsertotal(),
		},
	})
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
	List := router.Param(c, "id")
	id := strings.Split(List, ",")
	s.DeleteUser(id)
	return acSucc(c, nil)
}

func postUser(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	var a ReqUser
	if e := c.GetJSON(&a); e != nil {
		return e
	}
	if e := s.PostUser(a.Username, a.Password, a.Rid); e != nil {
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
	Rid      []uint
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
	if e := s.PutUser(uint(id), a.Nick, a.Password, a.Avatar, a.Desc, a.Rid); e != nil {
		return fail(c)
	}
	return acSucc(c, nil)
}

//role

//getRoles 获取角色信息
func getRoles(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	page, size := getPageSize(c)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "SUCCESS",
		"data": map[string]interface{}{
			"total": s.GetRoleCount(),
			"roles": s.GetListRole(page, size),
		},
	})
}

//getRole 根据id查询角色
func getRole(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	id, err := strconv.Atoi(router.Param(c, "id"))
	if err != nil {
		return fail(c)
	}
	r := s.GetRole(uint(id))
	return acSucc(c, r)
}

// ReqRole -
type ReqRole struct {
	Name string
	Desc string
	Rid  []uint
}

func postRole(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	var a ReqRole
	if e := c.GetJSON(&a); e != nil {
		return e
	}
	if e := s.PostRole(a.Name, a.Desc, a.Rid); e != nil {
		return e
	}
	return acSucc(c, nil)
}
func putRole(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	id, err := strconv.Atoi(router.Param(c, "id"))
	if err != nil {
		return err
	}
	var a ReqRole
	if e := c.GetJSON(&a); e != nil {
		return e
	}
	s.PutRole(a.Name, a.Desc, uint(id), a.Rid)
	return acSucc(c, nil)
}
func delRole(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	list := router.Param(c, "id")
	id := strings.Split(list, ",")
	if e := s.DeleteRole(id); e != nil {
		return e
	}
	return acSucc(c, nil)
}

// getPerms -
func getPerms(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	page, size := getPageSize(c)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "SUCCESS",
		"data": map[string]interface{}{
			"total": s.GetPermCount(),
			"roles": s.GetListPerm(page, size),
		},
	})
}

// getPerm -
func getPerm(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	id, err := strconv.Atoi(router.Param(c, "id"))
	if err != nil {
		return err
	}
	return acSucc(c, s.GetPerm(id))
}

func postPerm(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	var a ac.BoxPerm
	if err := c.GetJSON(&a); err != nil {
		return err
	}
	s.PostPerm(a)
	return acSucc(c, nil)
}

func delPerm(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	list := router.Param(c, "id")
	id := strings.Split(list, ",")
	s.DelPerm(id)
	return acSucc(c, nil)
}

func putPerm(c sola.Context) error {
	s := c.Get(CtxBoxAC).(*ac.Srv)
	id, err := strconv.Atoi(router.Param(c, "id"))
	if err != nil {
		return nil
	}
	var a ac.BoxPerm
	if err := c.GetJSON(&a); err != nil {
		return err
	}
	s.PutPerm(id, a)
	return acSucc(c, nil)
}
