package gormrepo

import (
	"context"
	"job-queue/internal/domain"
	"time"

	"gorm.io/gorm"
)

type JobModel struct {
	ID        string    `gorm:"primaryKey;column:id" json:"id"`
	Type      string    `gorm:"column:type" json:"type"`
	Payload   []byte    `gorm:"column:payload" json:"payload"`
	Status    string    `gorm:"column:status" json:"status"`
	Result    []byte    `gorm:"column:result" json:"result"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (JobModel) TableName() string {
	return "jobs"
}

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) *JobRepository {
	// auto-migrate the jobs table
	db.AutoMigrate(&JobModel{})
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(ctx context.Context, job domain.Job) error {
	m := JobModel{
		ID:        job.ID,
		Type:      job.Type,
		Payload:   job.Payload,
		Status:    job.Status,
		Result:    job.Result,
		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *JobRepository) GetByID(ctx context.Context, id string) (domain.Job, error) {
	var m JobModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return domain.Job{}, err
	}
	return domain.Job{
		ID:        m.ID,
		Type:      m.Type,
		Payload:   m.Payload,
		Status:    m.Status,
		Result:    m.Result,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

func (r *JobRepository) UpdateStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).Model(&JobModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}).Error
}

func (r *JobRepository) SaveResult(ctx context.Context, id string, result []byte) error {
	return r.db.WithContext(ctx).Model(&JobModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"result":     result,
		"status":     domain.StatusDone,
		"updated_at": time.Now(),
	}).Error
}
