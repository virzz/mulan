package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/virzz/vlog"

	"github.com/virzz/mulan/auth"
	"github.com/virzz/mulan/auth/apikey"
	"github.com/virzz/mulan/rdb"
)

func New(conf *Config, routers Routers, mwBefore, mwAfter []gin.HandlerFunc) (*http.Server, error) {
	os.MkdirAll("logs", 0755)
	logFile, err := os.OpenFile(filepath.Join("logs", "gin.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	if conf.Debug {
		gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
		gin.DefaultErrorWriter = io.MultiWriter(logFile, os.Stderr)
		engine.Use(gin.Logger())
	} else {
		gin.DefaultWriter = logFile
		gin.DefaultErrorWriter = logFile
	}
	engine.GET(conf.Prefix+"/version", func(c *gin.Context) { c.String(200, conf.version+" "+conf.commit) })
	engine.GET(conf.Prefix+"/health", func(c *gin.Context) { c.Status(200) })
	if conf.Metrics {
		m := ginmetrics.GetMonitor()
		m.SetMetricPath("/metrics")
		m.SetSlowTime(10)
		m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
		m.Use(engine)
	}
	if conf.System != "" {
		systemGroup := engine.Group("/system", apikey.Mw("system", conf.System))
		systemGroup.POST("/system/upgrade", handleSystemUpgrade)
		systemGroup.POST("/system/upload", handleSystemUpload)
		if conf.Pprof {
			pprof.Register(systemGroup, "/pprof")
		}
	}
	if conf.RequestID {
		engine.Use(requestid.New())
	}
	// CORS
	c := cors.DefaultConfig()
	c.AddAllowHeaders(conf.Headers...)
	if len(conf.Origins) > 0 {
		c.AllowAllOrigins = false
		c.AllowOrigins = conf.Origins
	} else {
		c.AllowAllOrigins = true
	}
	engine.Use(cors.New(c))
	// Auth: Session
	if conf.Auth {
		engine.Use(auth.Init(rdb.R()))
	}
	// Register Router
	api := engine.Group(conf.Prefix)
	// Register Before Middleware
	if len(mwBefore) > 0 {
		api.Use(mwBefore...)
	}
	// Register Routers
	for _, register := range routers {
		register(api)
	}
	// Register After Middleware
	if len(mwAfter) > 0 {
		api.Use(mwAfter...)
	}
	engine.NoRoute(func(c *gin.Context) {
		buf, err := io.ReadAll(c.Request.Body)
		if err != nil {
			vlog.Error("404 Not Found", "method", c.Request.Method, "path", c.Request.RequestURI, "method", c.Request.Method)
		} else {
			vlog.Error("404 Not Found", "method", c.Request.Method, "path", c.Request.RequestURI, "body", string(buf))
		}
	})
	fmt.Fprintln(os.Stderr, "HTTP Server Listening on : "+conf.GetEndpoint())
	return &http.Server{Addr: conf.Addr(), Handler: engine}, nil
}
