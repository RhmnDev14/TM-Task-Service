package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"task-service/internal/dto"
	"task-service/internal/helper"
	"task-service/internal/middleware"
	"task-service/internal/usecase"
)

type TaskHandler struct {
	taskUc     usecase.TaskUc
	rg         *http.ServeMux
	middleware *middleware.Middleware
}

func NewTaskHandler(taskUc usecase.TaskUc, rg *http.ServeMux, middleware *middleware.Middleware) *TaskHandler {
	return &TaskHandler{taskUc: taskUc, rg: rg, middleware: middleware}
}

func (h *TaskHandler) SetupRoutes() {
	h.rg.Handle(helper.Task, h.middleware.PrivateMiddleware(http.HandlerFunc(h.CreateTask)))
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helper.WriteJSON(w, http.StatusMethodNotAllowed, helper.MethodNotAllowed)
		return
	}

	var reqBody dto.TaskDto
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		helper.WriteJSON(w, http.StatusBadRequest, fmt.Sprintf(helper.InvalidJson, err))
		return
	}

	if err := h.taskUc.CreateTask(r.Context(), reqBody); err != nil {
		helper.WriteJSON(w, http.StatusInternalServerError, dto.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	helper.WriteJSON(w, http.StatusCreated, dto.Response{
		Status:  http.StatusCreated,
		Message: helper.Succes,
	})
}
