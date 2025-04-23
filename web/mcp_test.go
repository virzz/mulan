package web

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func helloHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := req.Params.Arguments["name"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}
	return mcp.NewToolResultText("Hello, " + name), nil
}

func TestWrapMCP(t *testing.T) {
	s := server.NewMCPServer("test-server", "0.1.0")
	// Add tool handler
	s.AddTool(
		mcp.NewTool("hello_world",
			mcp.WithDescription("Say hello to someone"),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Name of the person to greet"),
			),
		),
		helloHandler,
	)

	gin.SetMode(gin.DebugMode)
	router := gin.New()
	g := router.Group("/api")
	RegisterMCP(g, s, "/api")

	server := httptest.NewServer(router)
	defer server.Close()

	serverURL := server.URL + "/api/mcp/sse"

	client, err := client.NewSSEMCPClient(serverURL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start the client
	if err := client.Start(ctx); err != nil {
		t.Fatalf("Failed to start client: %v", err)
	}

	// Initialize
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{Name: "test-client", Version: "1.0.0"}
	result, err := client.Initialize(ctx, initRequest)
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	if result.ServerInfo.Name != "test-server" {
		t.Errorf(
			"Expected server name 'test-server', got '%s'",
			result.ServerInfo.Name,
		)
	}

	// Test Ping
	if err := client.Ping(ctx); err != nil {
		t.Errorf("Ping failed: %v", err)
	}

	// Test ListTools
	toolsRequest := mcp.ListToolsRequest{}
	toolResult, err := client.ListTools(ctx, toolsRequest)
	if err != nil {
		t.Errorf("ListTools failed: %v", err)
	}
	for _, tool := range toolResult.Tools {
		t.Logf("Tool: %s", tool.Name)
	}

	// Test CallTool
	callToolRequest := mcp.CallToolRequest{}
	callToolRequest.Params.Name = "hello_world"
	callToolRequest.Params.Arguments = map[string]any{"name": "John"}
	callToolResult, err := client.CallTool(ctx, callToolRequest)
	if err != nil {
		t.Errorf("CallTool failed: %v", err)
	}
	for _, content := range callToolResult.Content {
		switch c := content.(type) {
		case *mcp.TextContent:
			t.Logf("CallTool result: %s", c.Text)
		case *mcp.ImageContent:
			t.Logf("CallTool result: %s %s", c.MIMEType, c.Data)
		}
	}
}
