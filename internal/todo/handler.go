package todo

import (
	"context"
	"fmt"
	"todo-service/helper"
	"todo-service/pkg/constants"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	TodoService TodoService
}

func NewTodoHandler(TodoService TodoService) *TodoHandler {
	return &TodoHandler{
		TodoService: TodoService,
	}
}

func (h *TodoHandler) GetTodos(c *gin.Context) {

	status := c.Query("status")
	name := c.Query("name")
	teacher := c.Query("teacher")
	student := c.Query("student")
	staff := c.Query("staff")
	fmt.Printf("status: %s, name: %s, teacher: %s, student: %s, staff: %s", status, name, teacher, student, staff)
	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.TodoService.GetAllTodo(ctx, status, name, teacher, student, staff)
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Get todos successfully", data, 0)

}

func (h *TodoHandler) GetTodo(c *gin.Context) {

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

	data, err := h.TodoService.GetTodoByID(ctx, id)
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Get detail todo successfully", data, 0)

}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	
	var req CreateTodoRequest
	
	if err := c.BindJSON(&req); err != nil {
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

	data, err := h.TodoService.CreateTodo(ctx, req, userID.(string))
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Create todo successfully", data, 0)

}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {

	id := c.Param("id")

	if id == "" {
		helper.SendError(c, 400, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
		return
	}

	var req UpdateTaskProgressRequest
	
	if err := c.BindJSON(&req); err != nil {
		helper.SendError(c, 400, err, helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}
	
	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.TodoService.UpdateTodo(ctx, req, id)
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Update todo successfully", nil, 0)

}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {

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

	err := h.TodoService.DeleteTodo(ctx, id)
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Delete todo successfully", nil, 0)	
}

func (h *TodoHandler) JoinTodo(c *gin.Context) {

	var req JoinTodoRequest
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

	err := h.TodoService.JoinTodo(ctx, req, userID.(string))
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Join todo successfully", nil, 0)

}

func (h TodoHandler) AddUser(c *gin.Context) {

	var req AddUserRequest
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

	err := h.TodoService.AddUser(ctx, req)
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Add user successfully", nil, 0)

}

func (h *TodoHandler) GetMyTodo(c *gin.Context) {

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

	data, avg, err := h.TodoService.GetMyTodo(ctx, userID.(string))
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Get my todo successfully", data, avg)
}
