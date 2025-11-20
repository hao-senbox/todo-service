package repair

import (
	"context"
	"fmt"
	"todo-service/helper"
	"todo-service/pkg/constants"

	"github.com/gin-gonic/gin"
)

type RepairHandler struct {
	RepairService RepairService
}

func NewRepairHandler(RepairService RepairService) *RepairHandler {
	return &RepairHandler{
		RepairService: RepairService,
	}
}

func (h *RepairHandler) CreateRepair(c *gin.Context) {
	var req CreateRepairRequest
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

	data, err := h.RepairService.CreateRepair(ctx, req, userID.(string))
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Create repair successfully", data, 0)
}

func (h *RepairHandler) GetRepairs(c *gin.Context) {
	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.RepairService.GetRepairs(ctx)
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, 200, "Get repairs successfully", data, 0)
}

func (h *RepairHandler) GetRepairByID(c *gin.Context) {

	id := c.Param("id")

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, 400, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.RepairService.GetRepairByID(ctx, id)
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, 200, "Get repair by id successfully", data, 0)
}

func (h *RepairHandler) UpdateRepair(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, 400, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
		return
	}

	var req UpdateRepairRequest
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

	err := h.RepairService.UpdateRepair(ctx, req, id, userID.(string))
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, 200, "Update repair successfully", nil, 0)
}

func (h *RepairHandler) DeleteRepair(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, 400, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
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

	err := h.RepairService.DeleteRepair(ctx, id, userID.(string))
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, 200, "Delete repair successfully", nil, 0)
}

func (h *RepairHandler) AssignRepair(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, 400, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
		return
	}

	var req AssignRepairRequest
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

	err := h.RepairService.AssignRepair(ctx, id, userID.(string), req)
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, 200, "Assign repair successfully", nil, 0)
}

func (h *RepairHandler) CompleteRepair(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		helper.SendError(c, 400, fmt.Errorf("id is required"), helper.ErrInvalidRequest)
		return
	}

	var req CompleteRepairRequest
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

	err := h.RepairService.CompleteRepair(ctx, id, req, userID.(string))
	if err != nil {
		helper.SendError(c, 500, err, helper.ErrInvalidOperation)
		return
	}
	helper.SendSuccess(c, 200, "Complete repair successfully", nil, 0)

}
