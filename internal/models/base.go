package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	IsActive  bool      `gorm:"column:is_active;default:true" json:"is_active,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at,omitempty"`
}

func (base *BaseEntity) BeforeCreate(tx *gorm.DB) (err error) {
	base.CreatedAt = time.Now()
	base.UpdatedAt = time.Now()
	return
}

func (base *BaseEntity) BeforeUpdate(tx *gorm.DB) (err error) {
	base.UpdatedAt = time.Now()
	return
}
