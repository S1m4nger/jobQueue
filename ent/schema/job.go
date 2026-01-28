package schema

import (
	"time"

	"gorm.io/gorm"
)

// Job is a GORM model that represents a queued job.
type Job struct {
	ID        string         `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	Type      string         `gorm:"column:type;type:varchar(255);not null" json:"type"`
	Payload   []byte         `gorm:"column:payload;type:blob" json:"payload,omitempty"`
	Status    string         `gorm:"column:status;type:varchar(32);not null;default:'pending'" json:"status"`
	Result    []byte         `gorm:"column:result;type:blob" json:"result,omitempty"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Job) TableName() string {
	return "jobs"
}
