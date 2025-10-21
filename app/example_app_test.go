package app_test

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/virzz/mulan/app"
	"github.com/virzz/mulan/db"
	"github.com/virzz/mulan/rdb"
	"github.com/virzz/mulan/web"
)

type Config struct {
	DB   db.Config  `json:"db" yaml:"db"`
	HTTP web.Config `json:"http" yaml:"http"`
	RDB  rdb.Config `json:"rdb" yaml:"rdb"`
}

var (
	Version string = "1.0.0"
	Commit  string = "dev"
	BuildAt string = time.Now().Format(time.RFC3339)

	Conf = &Config{}
)

func Example() {
	meta := &app.Meta{
		ID:          "com.virzz.mulan.example",
		Name:        "example",
		Description: "ExampleService",
		Version:     Version,
		Commit:      Commit,
		BuildAt:     BuildAt,
	}
	std := app.New(meta, Conf)

	routerFunc := func(api gin.IRouter) {
		api.Handle("GET", "/", func(c *gin.Context) {
			c.String(200, "Hello, World!")
		})
	}
	webInfo := &web.Info{Name: meta.Name, Version: meta.Version, Commit: meta.Commit}
	webSrv := web.New(&Conf.HTTP, webInfo, routerFunc)

	std.AddService(webSrv)

	if err := std.Execute(context.Background()); err != nil {
		panic(err)
	}
}
