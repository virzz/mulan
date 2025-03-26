package web

import (
	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/server"
)

func WrapMCP(s *server.MCPServer) (string, string, gin.HandlerFunc) {
	sse := server.NewSSEServer(s, server.WithBasePath("/mcp"))
	return sse.CompleteSsePath(), sse.CompleteMessagePath(), func(c *gin.Context) {
		sse.ServeHTTP(c.Writer, c.Request)
	}
}

func RegisterMCP(r gin.IRoutes, s *server.MCPServer) {
	sse, msg, h := WrapMCP(s)
	r.GET(sse, h)
	r.POST(msg, h)
}
