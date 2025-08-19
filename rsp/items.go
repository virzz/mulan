package rsp

import "github.com/virzz/mulan/rsp/code"

type IntX interface {
	int | int64 | int32 | int16 | int8 | uint | uint64 | uint32 | uint16 | uint8
}

type Items struct {
	Items any   `json:"items"`
	Total int64 `json:"total"`
	Ext   any   `json:"ext,omitempty"`
}

func Item[T IntX](total T, items any) Rsp {
	return Rsp{APICode: code.Success, Data: &Items{Items: items, Total: int64(total)}}
}
func ItemExt[T IntX](total T, items any, ext any) Rsp {
	return Rsp{APICode: code.Success, Data: &Items{Items: items, Total: int64(total), Ext: ext}}
}
func ItemCode(c code.APICode) Rsp {
	return Rsp{APICode: c, Data: &Items{Items: []struct{}{}, Total: 0}}
}
func ItemNone() Rsp { return ItemCode(code.NotFound) }
func MItem[T IntX](total T, items any, msg string, ext ...any) Rsp {
	r := Rsp{APICode: code.Success, Data: &Items{Items: items, Total: int64(total)}}
	r.WithMsg(msg)
	if len(ext) > 0 {
		r.Data.(*Items).Ext = ext[0]
	}
	return r
}
