package dto

type TaskDto struct {
	TaskName   string `json:"task_name"`
	StoryPoint uint   `json:"story_point"`
	Level      int    `json:"level"`
	WorkSpace  int    `json:"works_pace"`
	Assignee   int    `json:"assigne"`
}
