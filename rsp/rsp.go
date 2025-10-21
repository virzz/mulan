package rsp

import (
	e "github.com/virzz/mulan/rsp/apperr"
)

type IDRsp[T uint64 | string] struct {
	ID   T      `json:"id,omitempty"`
	UUID string `json:"uuid,omitempty"`
}

type Rsp[T any] struct {
	*e.AppError `json:",inline"`
	Data        T `json:"data"`
}

func S[T any](data T) *Rsp[T] { return &Rsp[T]{AppError: e.Success, Data: data} }

type Resp = Rsp[any]
type RspNone = Rsp[struct{}]

func UnImplemented() *RspNone    { return &RspNone{AppError: e.UnImplemented} }
func E(err *e.AppError) *RspNone { return &RspNone{AppError: err} }
func OK() *RspNone               { return &RspNone{AppError: e.Success} }
func M(msg string) *RspNone      { return &RspNone{AppError: &e.AppError{Code: e.Success.Code, Msg: msg}} }
func C(err error) *RspNone {
	return &RspNone{AppError: &e.AppError{Code: e.UnknownErr.Code, Msg: err.Error()}}
}
