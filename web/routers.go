package web

import (
	"github.com/gin-gonic/gin"
)

type (
	Controller interface {
		Authed(*gin.RouterGroup)
		UnAuth(*gin.RouterGroup)
	}

	RegisterFunc func(*gin.RouterGroup)
)

type Routers []RegisterFunc

func (rs *Routers) Register(f RegisterFunc) {
	*rs = append(*rs, f)
}

func (rs Routers) Apply(g *gin.RouterGroup) {
	for _, f := range rs {
		f(g)
	}
}
