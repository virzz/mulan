package auth

import (
	"github.com/gin-gonic/gin"

	"github.com/virzz/mulan/rsp"
	"github.com/virzz/mulan/rsp/code"
)

func RoleMW(roles ...string) gin.HandlerFunc {
	roleMap := make(map[string]struct{})
	for _, role := range roles {
		roleMap[role] = struct{}{}
	}
	return func(c *gin.Context) {
		for _, r := range c.GetStringSlice("roles") {
			if _, ok := roleMap[r]; ok {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(200, rsp.C(code.Forbidden))
	}
}

func AuthMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, ok := c.Get("id"); ok && c.GetString("account") != "" {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(200, rsp.C(code.Unauthorized))
	}
}
