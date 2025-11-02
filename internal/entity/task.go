package entity

import (
	"task-service/internal/dto"
	"time"
)

type Task struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	TaskName   string    `json:"task_name"`
	StoryPoint uint      `json:"story_point"`
	Level      int       `json:"level"`
	WorkSpace  int       `json:"works_pace"`
	Assignee   int       `json:"assigne"`
	CreatedBy  int       `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func TaskToEntity(req dto.TaskDto, idUser int) Task {
	return Task{
		TaskName:   req.TaskName,
		StoryPoint: req.StoryPoint,
		Level:      req.Level,
		WorkSpace:  req.WorkSpace,
		Assignee:   req.Assignee,
		CreatedBy:  idUser,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
