package web

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/virzz/mulan/auth"
	"github.com/virzz/mulan/auth/apikey"
	"github.com/virzz/mulan/rdb"
)

var engine *gin.Engine

func New(conf *Config, router *Routers, mwBefore, mwAfter []gin.HandlerFunc) (*http.Server, error) {
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

	versionFn := func(c *gin.Context) { c.String(200, conf.version+" "+conf.commit) }
	engine.GET("/version", versionFn)
	engine.GET("/health", HealthHandler)
	if conf.Prefix != "" && conf.Prefix != "/" {
		engine.GET(conf.Prefix+"/version", versionFn)
		engine.GET(conf.Prefix+"/health", HealthHandler)
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
	router.Apply(api)
	// Register After Middleware
	if len(mwAfter) > 0 {
		api.Use(mwAfter...)
	}
	zap.L().Info("HTTP Server Listening on : " + conf.GetEndpoint())
	return &http.Server{Addr: conf.Addr(), Handler: engine}, nil
}
