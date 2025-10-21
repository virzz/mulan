package web

import (
	"github.com/gin-gonic/gin"

	"github.com/virzz/mulan/rsp"
	"github.com/virzz/mulan/rsp/apperr"
)

type Handler interface {
	Register(gin.IRouter)
}

func VersionHandler(info *Info) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(200, info.Name+" "+info.Version+" "+info.Commit)
	}
}

func ErrCodeHandler(c *gin.Context) { c.JSON(200, rsp.S(apperr.Errors)) }
