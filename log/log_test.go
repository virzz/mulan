package log_test

import (
	"testing"

	"github.com/virzz/mulan/log"
	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	log.NewWithConfig(log.Config{
		Level: "debug",
		Http: []log.Http{
			{
				URL:   "http://localhost:3003/log",
				Level: "debug",
			},
		},
	})
	zap.L().Info("test")
	zap.L().Info("testaaaaaaaaaaa")
	zap.L().Info("testaaaaaargverasrfesaaaaaa")
}
