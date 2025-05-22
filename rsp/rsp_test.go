package rsp_test

import (
	"testing"

	"github.com/go-json-experiment/json"

	"github.com/virzz/mulan/code"
	"github.com/virzz/mulan/rsp"
)

func jsonString(v any) string {
	buf, _ := json.Marshal(v)
	return string(buf)
}
func TestMsg(t *testing.T) {
	t.Log(jsonString(rsp.M("test")))
	t.Log(jsonString(rsp.OK()))
	t.Log(jsonString(rsp.M("aaaaaaaaaaaaaaaa")))
	t.Log(jsonString(rsp.OK()))
	t.Log(jsonString(rsp.E(code.DatabaseUnknown, "DatabaseUnknown")))
	t.Log(jsonString(rsp.OK()))
}
