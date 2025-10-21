package req_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/virzz/mulan/req"
)

type User struct {
	Name string `json:"name" form:"name"`
	Info string `json:"info" form:"info"`
	ID   uint64 `json:"id" form:"id" uri:"id"`
}

func TestBind(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/u/:id", func(c *gin.Context) {
		var obj User
		err := req.Bind(c, &obj)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, obj)
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/u/123?info=hello,world", bytes.NewBufferString(`{"name": "test"}`))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)
	t.Log(w.Body.String())
}

func BenchmarkBind(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/u/:id", func(c *gin.Context) {
		var obj User
		req.Bind(c, &obj)
	})
	w := httptest.NewRecorder()
	buf := bytes.NewBufferString(`{"name": "test"}`)
	r := httptest.NewRequest(http.MethodPost, "/u/123?info=test", buf)
	r.Header.Set("Content-Type", "application/json")

	for b.Loop() {
		router.ServeHTTP(w, r)
	}
}
