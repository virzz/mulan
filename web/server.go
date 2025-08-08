package web

import (
	"net/http"
	"slices"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/virzz/mulan/auth/apikey"
)

var (
	engine *gin.Engine

	skipPaths = []string{"/health", "/version", "/metrics"}
)

func New(conf *Config, applyFunc func(*gin.RouterGroup)) (*http.Server, error) {
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	ginLog := zap.L().Named("gin")
	engine = gin.New()
	engine.Use(
		ginzap.RecoveryWithZap(ginLog, true),
		ginzap.GinzapWithConfig(ginLog, &ginzap.Config{
			TimeFormat: time.RFC3339,
			UTC:        true,
			Skipper: func(c *gin.Context) bool {
				if c.Request.Response != nil &&
					c.Request.Response.StatusCode == 404 {
					return true
				}
				return slices.Contains(skipPaths, c.Request.URL.Path)
			},
			Context: func(c *gin.Context) []zap.Field {
				ctx := c.Request.Context()
				fields := []zapcore.Field{
					zap.String("referer", c.Request.Referer()),
				}
				if requestid := requestid.Get(c); requestid != "" {
					fields = append(fields, zap.String("requestid", requestid))
				}
				// log trace and span ID
				if span := trace.SpanFromContext(ctx).SpanContext(); span.IsValid() {
					fields = append(fields, zap.String("trace_id", span.TraceID().String()))
					fields = append(fields, zap.String("span_id", span.SpanID().String()))
				}
				return fields
			},
		}),
	)

	engine.GET("/version", versionHandler)
	engine.GET("/health", versionHandler)
	if conf.Prefix != "" && conf.Prefix != "/" {
		engine.GET(conf.Prefix+"/version", versionHandler)
		engine.GET(conf.Prefix+"/health", versionHandler)
	}

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
	// Register Router
	api := engine.Group(conf.Prefix)

	if applyFunc != nil {
		applyFunc(api)
	}

	zap.L().Info("HTTP Server Listening on",
		zap.String("endpoint", conf.GetEndpoint()),
		zap.String("host", conf.Host),
		zap.Int("port", conf.Port),
	)
	return &http.Server{Addr: conf.Addr(), Handler: engine}, nil
}

func Routes() []gin.RouteInfo { return engine.Routes() }
