package web_test

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mulan-ext/log"
	"github.com/virzz/mulan/web"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	log.New(true, "test")

	cfg := &web.Config{
		Host:  "127.0.0.1",
		Port:  3003,
		Pprof: false,
	}
	webInfo := &web.Info{
		Name:    "test",
		Version: "1.0.0",
		Commit:  "dev",
		BuildAt: time.Now().Format(time.RFC3339),
	}
	webSrv := web.New(cfg, webInfo, func(api gin.IRouter) {
		api.Handle("GET", "/aaa", func(c *gin.Context) {
			zap.L().Info("Hello, World!")
			c.String(200, "Hello, World!")
		})
	})
	// webSrv.Build()

	go func() {
		err := webSrv.Serve()
		if err != nil && err != http.ErrServerClosed {
			zap.L().Error("Failed to run http server", zap.Error(err))
		}
	}()

	r, err := http.Get("http://127.0.0.1:3003/aaa")
	if err != nil {
		t.Fatal("Failed to get /aaa", "err", err.Error())
	}
	body, err := httputil.DumpResponse(r, true)
	if err != nil {
		t.Fatal("Failed to dump response", "err", err.Error())
	}
	fmt.Println(string(body))

	<-time.After(5 * time.Second)
	webSrv.Close()
}
