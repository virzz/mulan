package rsp_test

import (
	"fmt"
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
	fmt.Println(jsonString(rsp.M("test")))
	fmt.Println(jsonString(rsp.OK()))
	fmt.Println(jsonString(rsp.M("aaaaaaaaaaaaaaaa")))
	fmt.Println(jsonString(rsp.OK()))
	fmt.Println(jsonString(rsp.E(code.DatabaseUnknown, "DatabaseUnknown")))
	fmt.Println(jsonString(rsp.OK()))
}
