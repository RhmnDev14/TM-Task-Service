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
