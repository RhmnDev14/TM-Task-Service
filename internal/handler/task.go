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
	h.rg.Handle(helper.Task, h.middleware.PrivateMiddleware(http.HandlerFunc(h.handleTask)))
	h.rg.Handle(helper.TaskWithID, h.middleware.PrivateMiddleware(http.HandlerFunc(h.handleTaskWithID)))
}

func (h *TaskHandler) handleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateTask(w, r)
	case http.MethodGet:
		h.GetListTask(w, r)
	default:
		helper.WriteJSON(w, http.StatusMethodNotAllowed, helper.MethodNotAllowed)
	}
}

func (h *TaskHandler) handleTaskWithID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		h.UpdateTaskProgress(w, r)
	default:
		helper.WriteJSON(w, http.StatusMethodNotAllowed, helper.MethodNotAllowed)
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
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

func (h *TaskHandler) GetListTask(w http.ResponseWriter, r *http.Request) {
	userIDVal := r.Context().Value("user_id")
	if userIDVal == nil {
		h.writeUnauthorized(w)
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok {
		h.writeInternalServerError(w, helper.InternalServerError)
		return
	}

	tasks, err := h.taskUc.GetListTask(r.Context(), int(userID))
	if err != nil {
		h.writeInternalServerError(w, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, dto.Response{
		Status:  http.StatusOK,
		Message: helper.Succes,
		Data:    tasks,
	})
}

func (h *TaskHandler) writeMethodNotAllowed(w http.ResponseWriter) {
	helper.WriteJSON(w, http.StatusMethodNotAllowed, helper.MethodNotAllowed)
}

func (h *TaskHandler) writeBadRequest(w http.ResponseWriter, message string) {
	helper.WriteJSON(w, http.StatusBadRequest, dto.Response{
		Status:  http.StatusBadRequest,
		Message: message,
	})
}

func (h *TaskHandler) writeUnauthorized(w http.ResponseWriter) {
	helper.WriteJSON(w, http.StatusUnauthorized, dto.Response{
		Status:  http.StatusUnauthorized,
		Message: helper.Unauthorized,
	})
}

func (h *TaskHandler) writeInternalServerError(w http.ResponseWriter, message string) {
	helper.WriteJSON(w, http.StatusInternalServerError, dto.Response{
		Status:  http.StatusInternalServerError,
		Message: message,
	})
}

func (h *TaskHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	helper.WriteJSON(w, statusCode, data)
}

func (h *TaskHandler) UpdateTaskProgress(w http.ResponseWriter, r *http.Request) {
	taskIDStr := r.URL.Path[len(helper.TaskWithID)-1:] // Assuming URL is /task/{id}
	taskID, err := helper.ParseUint(taskIDStr)
	if err != nil {
		h.writeBadRequest(w, fmt.Sprintf(helper.InvalidTaskID, err))
		return
	}

	var reqBody dto.UpdateTaskProgressDto
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		h.writeBadRequest(w, fmt.Sprintf(helper.InvalidJson, err))
		return
	}

	if err := h.taskUc.UpdateTaskProgress(r.Context(), taskID, reqBody); err != nil {
		h.writeInternalServerError(w, err.Error())
		return
	}

	h.writeJSON(w, http.StatusOK, dto.Response{
		Status:  http.StatusOK,
		Message: helper.Succes,
	})
}
