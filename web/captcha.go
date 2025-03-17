package web

import (
	"github.com/gin-gonic/gin"

	"github.com/virzz/vlog"

	"github.com/virzz/mulan/captcha"
	"github.com/virzz/mulan/code"
	"github.com/virzz/mulan/req"
	"github.com/virzz/mulan/rsp"
)

func CaptchaHandler(debug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _code, data, err := captcha.CreateB64()
		if err != nil {
			vlog.Error("Failed to create base64 captcha", "err", err.Error())
			c.AbortWithStatusJSON(200, rsp.C(code.CaptchaGenerate))
			return
		}
		if debug && c.GetHeader("X-Debug-Captcha") != "" {
			c.Header("Captcha", _code)
		}
		c.JSON(200, rsp.S(req.Captcha{UUID: id, Code: data}))
	}
}
