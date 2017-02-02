package app

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/jmoiron/sqlx"
	"github.com/valyala/fasthttp"
	"zxq.co/ripple/rippleapi/common"
)

type router struct {
	r *fasthttprouter.Router
}

func (r router) Method(path string, f func(md common.MethodData) common.CodeMessager, privilegesNeeded ...int) {
	r.r.GET(path, Method(f, privilegesNeeded...))
}
func (r router) POSTMethod(path string, f func(md common.MethodData) common.CodeMessager, privilegesNeeded ...int) {
	r.r.POST(path, Method(f, privilegesNeeded...))
}
func (r router) Peppy(path string, a func(c *fasthttp.RequestCtx, db *sqlx.DB)) {
	r.r.GET(path, PeppyMethod(a))
}
func (r router) GET(path string, handle fasthttp.RequestHandler) {
	r.r.GET(path, handle)
}
