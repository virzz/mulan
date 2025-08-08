package web_test

import (
	"log"
	"testing"
	"time"

	"github.com/virzz/mulan/web"
)

func TestNew(t *testing.T) {
	cfg := &web.Config{
		Host:    "127.0.0.1",
		Port:    3003,
		Pprof:   true,
		Metrics: true,
	}
	httpSrv, err := web.New(cfg, nil)
	if err != nil {
		t.Fatal(err)
		return
	}

	go func() {
		err := httpSrv.ListenAndServe()
		if err != nil {
			log.Println("Failed to run http server", "err", err.Error())
		}
	}()

	<-time.After(3 * time.Second)
	httpSrv.Close()
}
