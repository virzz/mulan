package web_test

import (
	"log"
	"testing"
	"time"

	"github.com/virzz/mulan/web"
)

func TestNew(t *testing.T) {
	httpSrv, err := web.New(&web.Config{
		Host:    "127.0.0.1",
		Port:    3003,
		Pprof:   true,
		Metrics: true,
	}, nil, nil, nil)
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

	time.Sleep(10 * time.Second)

}
