package web

import (
	"github.com/gin-gonic/gin"

	"github.com/virzz/mulan/captcha"
	"github.com/virzz/mulan/code"
	"github.com/virzz/mulan/rsp"
)

var versionHandler gin.HandlerFunc = func(c *gin.Context) { c.Status(200) }

func SetVersionHandler(version, commit string) {
	versionHandler = func(c *gin.Context) { c.String(200, version+" "+commit) }
}

func ErrCodeHandler(c *gin.Context) { c.JSON(200, rsp.S(code.Codes)) }

func CaptchaHandler(debug bool) gin.HandlerFunc { return captcha.GinHandler(debug) }
