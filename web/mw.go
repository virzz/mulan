package web

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/virzz/vlog"
)

func LogMw(c *gin.Context) {
	c.Next()
	args := []any{
		"remote_ip", c.RemoteIP(),
		"client_ip", c.ClientIP(),
		"referer", c.Request.Referer(),
		"useragent", c.Request.UserAgent(),
		"status", c.Writer.Status(),
	}
	if requestid := requestid.Get(c); requestid != "" {
		args = append(args, "requestid", requestid)
	}
	vlog.Info("AccessLog", args...)
}
