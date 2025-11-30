package repository

import (
	"context"
	"task-service/internal/entity"

	"gorm.io/gorm"
)

type TaskRepo interface {
	CreateTask(ctx context.Context, req entity.Task) error
	GetListTask(ctx context.Context, userID int) ([]entity.Task, error)
	UpdateTaskProgress(ctx context.Context, taskID uint, progress int) error
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

func (r *taskRepo) GetListTask(ctx context.Context, userID int) ([]entity.Task, error) {
	var tasks []entity.Task
	err := r.db.WithContext(ctx).Where("created_by = ?", userID).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepo) UpdateTaskProgress(ctx context.Context, taskID uint, progress int) error {
	return r.db.WithContext(ctx).Model(&entity.Task{}).Where("id = ?", taskID).Update("progress", progress).Error
}
