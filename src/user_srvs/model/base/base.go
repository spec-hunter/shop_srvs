package base

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint      `gorm:"primaryKey;"`
	CreatedAt time.Time `gorm:"column:created_at;check:created_at <= updated_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;check:created_at <= updated_at"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool `gorm:"column:is_deleted;default:false;"`
}
