package apikey

import (
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/virzz/mulan/rsp"
	"github.com/virzz/mulan/rsp/code"
)

func Mw(name string, apikeys ...string) func(*gin.Context) {
	if name == "" {
		name = "apikey"
	}
	if name == "disable" || len(apikeys) == 0 {
		return func(c *gin.Context) { c.Next() }
	}
	return func(c *gin.Context) {
		apikey := c.GetHeader(name)
		if name == "Authorization" {
			apikey = strings.TrimPrefix(apikey, "Bearer ")
		}
		if apikey == "" {
			apikey = c.Query(name)
			if apikey == "" {
				apikey, _ = c.Cookie(name)
			}
		}
		apikey = strings.TrimSpace(apikey)
		if apikey != "" && slices.Contains(apikeys, apikey) {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(401, rsp.C(code.Unauthorized))
	}
}
