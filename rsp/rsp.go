package rsp

import (
	"fmt"

	c "github.com/virzz/mulan/code"
)

type IDRsp struct {
	ID   uint64 `json:"id,omitempty"`
	UUID string `json:"uuid,omitempty"`
}

type Rsp struct {
	c.APICode
	Data any `json:"data"`
}

func (r *Rsp) WithCode(code c.APICode) *Rsp {
	r.APICode = code
	return r
}
func (r *Rsp) WithData(v any) *Rsp {
	r.Data = v
	return r
}
func (r *Rsp) WithMsg(v string) *Rsp {
	r.Msg = v
	return r
}
func (r *Rsp) WithItem(total int64, items any) *Rsp {
	r.Data = &Items{Items: items, Total: total}
	return r
}
func (r *Rsp) WithItemExt(total int64, items any, ext any) *Rsp {
	r.Data = &Items{Items: items, Total: total, Ext: ext}
	return r
}

func C(code c.APICode) *Rsp        { return &Rsp{APICode: code} }
func S(data any) *Rsp              { return &Rsp{APICode: c.Success, Data: data} }
func M(msg string) *Rsp            { return (&Rsp{APICode: c.Success}).WithMsg(msg) }
func SM(data any, msg string) *Rsp { return (&Rsp{APICode: c.Success, Data: data}).WithMsg(msg) }
func OK() *Rsp                     { return &Rsp{APICode: c.Success} }
func UnImplemented() *Rsp          { return &Rsp{APICode: c.UnImplemented} }

func E(code c.APICode, msg any) *Rsp {
	switch m := msg.(type) {
	case string:
		return (&Rsp{APICode: code}).WithMsg(m)
	case error:
		return (&Rsp{APICode: code}).WithMsg(m.Error())
	default:
		return (&Rsp{APICode: code}).WithMsg(fmt.Sprintf("%v", msg))
	}
}

func New(code int, msg string, data any) *Rsp {
	return &Rsp{c.NewCode(code, msg), data}
}
