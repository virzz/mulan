package web

import (
	"github.com/gin-gonic/gin"

	"github.com/virzz/mulan/rsp"
	"github.com/virzz/mulan/rsp/code"
)

var versionHandler gin.HandlerFunc = func(c *gin.Context) { c.Status(200) }

func SetVersionHandler(name, version, commit string) {
	versionHandler = func(c *gin.Context) { c.String(200, name+" "+version+" "+commit) }
}

func ErrCodeHandler(c *gin.Context) { c.JSON(200, rsp.S(code.Codes)) }
