package captcha

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/virzz/mulan/code"
	"github.com/virzz/mulan/req"
	"github.com/virzz/mulan/rsp"
)

func GinHandler(debug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _code, data, err := CreateB64()
		if err != nil {
			zap.L().Error("Failed to create base64 captcha", zap.Error(err))
			c.AbortWithStatusJSON(200, rsp.C(code.CaptchaGenerate))
			return
		}
		if debug && c.GetHeader("X-Debug-Captcha") != "" {
			c.Header("Captcha", _code)
		}
		c.JSON(200, rsp.S(req.Captcha{UUID: id, Code: data}))
	}
}
