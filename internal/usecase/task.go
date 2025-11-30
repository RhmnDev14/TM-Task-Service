package usecase

import (
	"context"
	"fmt"
	"task-service/internal/dto"
	"task-service/internal/entity"
	"task-service/internal/helper"
	"task-service/internal/repository"
)

type TaskUc interface {
	CreateTask(ctx context.Context, req dto.TaskDto) error
	GetListTask(ctx context.Context, userID int) ([]entity.Task, error)
	UpdateTaskProgress(ctx context.Context, taskID uint, req dto.UpdateTaskProgressDto) error
}

type taskUc struct {
	repo repository.TaskRepo
}

func NewTaskUc(repo repository.TaskRepo) *taskUc {
	return &taskUc{repo: repo}
}

func (u *taskUc) CreateTask(ctx context.Context, req dto.TaskDto) error {
	userIDVal := ctx.Value("user_id")
	if userIDVal == nil {
		helper.ErrorHandle(fmt.Errorf("user_id not found in context"))
		return fmt.Errorf(helper.InternalServerError)
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		helper.ErrorHandle(fmt.Errorf("invalid user_id type in context"))
		return fmt.Errorf(helper.InternalServerError)
	}
	param := entity.TaskToEntity(req, int(userID))
	return u.repo.CreateTask(ctx, param)
}

func (u *taskUc) GetListTask(ctx context.Context, userID int) ([]entity.Task, error) {
	tasks, err := u.repo.GetListTask(ctx, userID)
	if err != nil {
		helper.ErrorHandle(fmt.Errorf("failed to get list task: %w", err))
		return nil, fmt.Errorf(helper.InternalServerError)
	}
	return tasks, nil
}

func (u *taskUc) UpdateTaskProgress(ctx context.Context, taskID uint, req dto.UpdateTaskProgressDto) error {
	userIDVal := ctx.Value("user_id")
	if userIDVal == nil {
		helper.ErrorHandle(fmt.Errorf("user_id not found in context"))
		return fmt.Errorf(helper.InternalServerError)
	}
	return u.repo.UpdateTaskProgress(ctx, taskID, req.Progress)
}
