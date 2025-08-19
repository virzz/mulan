package auth

import (
	"crypto/rand"
	"encoding"
	"encoding/hex"
	"encoding/json"
	"io"
)

type (
	DataStringSlice []string
	DataMap         map[string]any
	Data            interface {
		Token() string
		ID() uint64
		Account() string
		State() uint16
		Roles() []string
		Items() DataMap
		Get(string) any
		New() string
		Set(string, any) Data
		SetToken(string) Data
		SetID(uint64) Data
		SetAccount(string) Data
		SetState(uint16) Data
		SetValues(string, any) Data
		SetRoles([]string) Data
		Delete(string) Data
		Clear() Data
	}
)

var _ encoding.TextUnmarshaler = (*DataStringSlice)(nil)
var _ encoding.TextUnmarshaler = (*DataMap)(nil)

func (d DataStringSlice) MarshalBinary() ([]byte, error)  { return json.Marshal(d) }
func (d *DataStringSlice) UnmarshalText(buf []byte) error { return d.UnmarshalBinary(buf) }
func (d *DataStringSlice) UnmarshalJSON(buf []byte) error { return d.UnmarshalText(buf) }
func (d *DataStringSlice) UnmarshalBinary(buf []byte) error {
	v := []string{}
	err := json.Unmarshal(buf, &v)
	if err != nil {
		return err
	}
	*d = DataStringSlice(v)
	return nil
}
func (d DataMap) MarshalBinary() ([]byte, error)    { return json.Marshal(d) }
func (d *DataMap) UnmarshalBinary(buf []byte) error { return json.Unmarshal(buf, d) }
func (d *DataMap) UnmarshalJSON(buf []byte) error   { return d.UnmarshalText(buf) }
func (d *DataMap) UnmarshalText(buf []byte) error {
	_d := make(map[string]any)
	if err := json.Unmarshal(buf, &_d); err != nil {
		return err
	}
	*d = _d
	return nil
}

type DefaultData struct {
	Token_   string          `json:"token" redis:"token"`
	ID_      uint64          `json:"id" redis:"id"`
	Account_ string          `json:"account" redis:"account"`
	State_   uint16          `json:"state" redis:"state"`
	Roles_   DataStringSlice `json:"roles" redis:"roles"`
	Items_   DataMap         `json:"items" redis:"items"`
}

var _ Data = (*DefaultData)(nil)

func New() string {
	k := make([]byte, 20)
	io.ReadFull(rand.Reader, k)
	return hex.EncodeToString(k)
}
func (d *DefaultData) New() string {
	_d := &DefaultData{}
	_d.Token_ = New()
	*d = *_d
	return d.Token_
}
func (d *DefaultData) ID() uint64      { return d.ID_ }
func (d *DefaultData) Token() string   { return d.Token_ }
func (d *DefaultData) Account() string { return d.Account_ }
func (d *DefaultData) State() uint16   { return d.State_ }
func (d *DefaultData) Roles() []string { return []string(d.Roles_) }
func (d *DefaultData) Items() DataMap  { return d.Items_ }
func (d *DefaultData) SetToken(v string) Data {
	d.Token_ = v
	return d
}
func (d *DefaultData) SetID(v uint64) Data {
	d.ID_ = v
	return d
}
func (d *DefaultData) SetAccount(v string) Data {
	d.Account_ = v
	return d
}
func (d *DefaultData) SetState(v uint16) Data {
	d.State_ = v
	return d
}
func (d *DefaultData) SetRoles(v []string) Data {
	d.Roles_ = DataStringSlice(v)
	return d
}
func (d *DefaultData) SetValues(k string, v any) Data {
	if d.Items_ == nil {
		d.Items_ = make(DataMap)
	}
	d.Items_[k] = v
	return d
}
func (d *DefaultData) Set(key string, val any) Data {
	switch key {
	case "id":
		if v, ok := val.(uint64); ok {
			d.ID_ = v
		}
	case "account":
		if v, ok := val.(string); ok {
			d.Account_ = v
		}
	case "roles":
		if v, ok := val.([]string); ok {
			d.Roles_ = DataStringSlice(v)
		}
	default:
		if d.Items_ == nil {
			d.Items_ = make(DataMap)
		}
		d.Items_[key] = val
	}
	return d
}
func (d *DefaultData) Get(key string) any {
	switch key {
	case "id":
		return d.ID_
	case "account":
		return d.Account_
	case "roles":
		return []string(d.Roles_)
	default:
		if d.Items_ != nil {
			if v, ok := d.Items_[key]; ok {
				return v
			}
		}
	}
	return nil
}
func (d *DefaultData) Delete(key string) Data {
	switch key {
	case "id":
		d.ID_ = 0
	case "account":
		d.Account_ = ""
	case "roles":
		d.Roles_ = nil
	default:
		if d.Items_ != nil {
			delete(d.Items_, key)
		}
	}
	return d
}
func (d *DefaultData) Clear() Data {
	_d := new(DefaultData)
	*d = *_d
	return d
}
