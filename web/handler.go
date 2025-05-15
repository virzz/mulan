package web

import (
	"github.com/gin-gonic/gin"

	"github.com/virzz/mulan/captcha"
	"github.com/virzz/mulan/code"
	"github.com/virzz/mulan/rsp"
)

func HealthHandler(c *gin.Context) { c.Status(200) }

func ErrCodeHandler(c *gin.Context) { c.JSON(200, rsp.S(code.Codes)) }

func CaptchaHandler(debug bool) gin.HandlerFunc { return captcha.GinHandler(debug) }
