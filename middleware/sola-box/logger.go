package box

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ddosakura/sola/v2"
	"github.com/ddosakura/sola/v2/middleware/logger"
	xlog "github.com/it-repo/box/service/log"
	"github.com/jinzhu/gorm"
)

// Logger Builder
func Logger(bufSize int, db *gorm.DB) (sola.Middleware, sola.Handler) {
	srv := xlog.New(db)
	return logger.New(bufSize, func(l *logger.Log) {
			if !l.IsAction {
				log.Printf(l.Format, l.V...)
				return
			}
			defer func() {
				_ = recover()
			}()
			uid := l.V[0].(uint)
			name := l.V[1].(string)
			desc := l.V[2].(string)
			srv.Insert(uid, name, desc)
		}), func(c sola.Context) error {
			r := c.Request()
			qs := r.URL.Query()
			page, err := strconv.Atoi(qs.Get("page"))
			if err != nil || page < 1 {
				page = 1
			}
			size, err := strconv.Atoi(qs.Get("size"))
			if err != nil || size < 1 {
				size = 5
			}
			return c.JSON(http.StatusOK, map[string]interface{}{
				"code": 0,
				"msg":  "SUCCESS",
				"data": map[string]interface{}{
					"items": srv.List(page, size),
					"total": srv.Total(page, size),
				},
			})
		}
}
