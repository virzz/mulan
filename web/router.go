package web

import (
	"github.com/gin-gonic/gin"
)

type (
	RegisterFunc func(*gin.RouterGroup)
	Routers      []RegisterFunc
)

func (rs *Routers) Handle(method, path string, f gin.HandlerFunc) {
	*rs = append(*rs, func(g *gin.RouterGroup) { g.Handle(method, path, f) })
}
func (rs *Routers) Register(f RegisterFunc) { *rs = append(*rs, f) }
func (rs Routers) Apply(g *gin.RouterGroup) {
	for _, f := range rs {
		f(g)
	}
}
func Routes() []gin.RouteInfo { return engine.Routes() }
