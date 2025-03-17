package auth

import (
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	DefaultKey = "github.com/virzz/mulan/auth"
	TokenKey   = "github.com/virzz/mulan/auth/token"
	paramKey   = "token"
)

func HasRole(c *gin.Context, role string) bool {
	return slices.Contains(c.GetStringSlice("roles"), role)
}

func HasRoles(c *gin.Context, roles ...string) bool {
	if len(roles) == 0 {
		return true
	}
	_roles := c.GetStringSlice("roles")
	roleSet := make(map[string]struct{}, len(_roles))
	for _, role := range _roles {
		roleSet[role] = struct{}{}
	}
	for _, role := range roles {
		if _, exists := roleSet[role]; !exists {
			return false
		}
	}
	return true
}

func Default(c *gin.Context) *Session {
	return c.MustGet(DefaultKey).(*Session)
}

func Init(store *redis.Client, data ...Data) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query(paramKey)
		if len(token) == 0 {
			token = c.PostForm(paramKey)
			if len(token) == 0 {
				token, _ = c.Cookie(paramKey)
				if len(token) == 0 {
					token = strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
				}
			}
		}
		if len(token) == 0 {
			token = New()
		}
		var _data Data
		if len(data) > 0 {
			_data = data[0]
			_data.Clear().SetToken(token)
		} else {
			_data = &DefaultData{Token_: token}
		}
		sess := NewSession(c, store, _data)
		c.Set(DefaultKey, sess)
		c.Set(TokenKey, token)
		if !sess.IsNil {
			data := sess.Data()
			roles := data.Roles()
			c.Set("id", data.ID())
			c.Set("account", data.Account())
			c.Set("state", data.State())
			c.Set("roles", roles)
			c.Set("is_admin", slices.Contains(roles, "admin"))
			for k, v := range data.Items() {
				c.Set(k, v)
			}
		}
		c.Next()
	}
}
