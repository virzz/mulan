package code

import (
	"sync"
)

var Codes = map[int]APICode{}

type APICode struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (a *APICode) C(c int) *APICode {
	a.Code = c
	return a
}

func (a *APICode) M(msg string) *APICode {
	a.Msg = msg
	return a
}

func New(system SystemCode, biz BizCode, code int, msg string) APICode {
	return NewCode(code+int(biz)*1000+int(system)*1000*100, msg)
}

var lock sync.Mutex

func NewCode(code int, msg string) APICode {
	c := APICode{Code: code, Msg: msg}
	lock.Lock()
	if _, ok := Codes[code]; ok {
		panic("code exists")
	}
	Codes[code] = c
	lock.Unlock()
	return c
}
