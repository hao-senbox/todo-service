package task

import (
	"context"
	"fmt"
	"todo-service/helper"
	"todo-service/pkg/constants"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	TaskService TaskService
}

func NewTaskHandler(taskService TaskService) *TaskHandler {
	return &TaskHandler{
		TaskService: taskService,
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.TaskService.CreateTask(ctx, req, userID.(string))
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Create task successfully", data, 0)
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.TaskService.GetTasks(ctx)
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Get tasks successfully", data, 0)
}

func (h *TaskHandler) GetTaskById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.SendError(c, 400, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.TaskService.GetTaskById(ctx, id)
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Get task by id successfully", data, 0)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.SendError(c, 400, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.TaskService.UpdateTask(ctx, req, id)
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Update task successfully", nil, 0)
}

func (h *TaskHandler) GetMyTask(c *gin.Context) {
	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.TaskService.GetMyTask(ctx, userID.(string))
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Get my task successfully", data, 0)
}

func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		helper.SendError(c, 400, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
		return
	}

	var req UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.TaskService.UpdateTaskStatus(ctx, req, id, userID.(string))
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Update task status successfully", nil, 0)
}