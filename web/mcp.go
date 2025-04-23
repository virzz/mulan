package web

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/server"
)

func WrapMCP(s *server.MCPServer, prefix string) (string, string, gin.HandlerFunc) {
	sse := server.NewSSEServer(s, server.WithBasePath(prefix+"/mcp"))
	ssePath := strings.TrimPrefix(sse.CompleteSsePath(), prefix)
	msgPath := strings.TrimPrefix(sse.CompleteMessagePath(), prefix)
	return ssePath, msgPath, func(c *gin.Context) {
		sse.ServeHTTP(c.Writer, c.Request)
	}
}

func RegisterMCP(r gin.IRoutes, s *server.MCPServer, prefix string) {
	sse, msg, h := WrapMCP(s, prefix)
	r.GET(sse, h)
	r.POST(msg, h)
}
