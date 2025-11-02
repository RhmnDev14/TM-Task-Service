package repository

import (
	"context"
	"task-service/internal/entity"

	"gorm.io/gorm"
)

type TaskRepo interface {
	CreateTask(ctx context.Context, req entity.Task) error
}

type taskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *taskRepo {
	return &taskRepo{db: db}
}

func (r *taskRepo) CreateTask(ctx context.Context, req entity.Task) error {
	return r.db.WithContext(ctx).Create(&req).Error
}
