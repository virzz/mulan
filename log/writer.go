package log

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap/zapcore"
)

type HTTPWriter struct {
	buf    *bytes.Buffer
	url    string
	client *http.Client
}

func (w *HTTPWriter) Write(p []byte) (n int, err error) {
	defer w.Sync()
	return w.buf.Write(p)
}

func (w *HTTPWriter) Sync() error {
	for {
		line, err := w.buf.ReadBytes('\n')
		if err != nil {
			break
		}
		w.post(line)
	}
	return nil
}

func (w *HTTPWriter) post(line []byte) error {
	rsp, err := w.client.Post(w.url, "application/json", bytes.NewBuffer(line))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer rsp.Body.Close()
	_, _ = io.Copy(io.Discard, rsp.Body)
	return nil
}

func newHTTPWriter(url string) zapcore.WriteSyncer {
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConnsPerHost: 0,
		MaxConnsPerHost:     0,
		MaxIdleConns:        0,
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 30}
	return &HTTPWriter{url: url, client: client, buf: bytes.NewBuffer(nil)}
}
