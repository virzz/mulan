package db

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Modeler interface {
	TableName() string
	Defaults() []Modeler
	GetID() uint64
	Unique() (string, string)
}

type Model struct {
	ID        uint64            `gorm:"primaryKey;autoIncrement;column:id" json:"-"`
	UUID      datatypes.BinUUID `gorm:"type:varchar(36);unique;column:uuid" json:"uuid"`
	CreatedAt int64             `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt int64             `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.UUID.IsNil() {
		m.UUID = datatypes.NewBinUUIDv4()
	}
	return
}
func (m *Model) Defaults() []Modeler { return []Modeler{} }

func (m *Model) GetID() uint64 { return m.ID }

type StringSlice = datatypes.JSONSlice[string]

var models = []Modeler{}

func Register(model ...Modeler) {
	models = append(models, model...)
}

func Models() []any {
	var result = make([]any, 0, len(models))
	for i, model := range models {
		result[i] = model
	}
	return result
}
