package db

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type (
	StringSlice = datatypes.JSONSlice[string]

	Modeler interface {
		TableName() string
		Defaults() []Modeler
		GetID() uint64
		Unique() (string, string)
	}

	Model struct {
		ID        uint64    `gorm:"primaryKey;autoIncrement;column:id" json:"-"`
		UUID      uuid.UUID `gorm:"type:uuid;unique;column:uuid;default:gen_random_uuid()" json:"uuid"`
		CreatedAt int64     `gorm:"autoCreateTime;column:created_at" json:"created_at"`
		UpdatedAt int64     `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
	}
)

func (m *Model) Defaults() []Modeler      { return []Modeler{} }
func (m *Model) GetID() uint64            { return m.ID }
func (m *Model) Unique() (string, string) { return "", "" }
func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	if m.UUID == uuid.Nil {
		m.UUID = uuid.Must(uuid.NewV7())
	}
	return
}
