package shop

import (
	"context"
	"fmt"
	"net/http"
	"todo-service/helper"
	"todo-service/pkg/constants"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShopHandler struct {
	ShopService ShopService
}

func NewShopHandler(shopService ShopService) *ShopHandler {
	return &ShopHandler{
		ShopService: shopService,
	}
}

func (h *ShopHandler) CreateShop(c *gin.Context) {
	var req CreateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.ShopService.CreateShop(ctx, req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Shop created successfully", data, 0)
}

func (h *ShopHandler) GetMyShop(c *gin.Context) {
	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.ShopService.GetMyShop(ctx, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "My shop retrieved successfully", data, 0)
}

func (h *ShopHandler) GetShopByID(c *gin.Context) {

	shopIDStr := c.Param("id")
	if shopIDStr == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("shop_id is required"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	data, err := h.ShopService.GetShopByID(ctx, shopIDStr)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Shop retrieved successfully", data, 0)
}

func (h *ShopHandler) UpdateShop(c *gin.Context) {

	var req UpdateShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	shopIDStr := c.Param("id")
	if shopIDStr == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("shop_id is required"), helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.ShopService.UpdateShop(ctx, shopIDStr, req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Shop updated successfully", nil, 0)

}

func (h *ShopHandler) DeleteShop(c *gin.Context) {
	
	shopIDStr := c.Param("id")
	if shopIDStr == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("shop_id is required"), helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.ShopService.DeleteShop(ctx, shopIDStr, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Shop deleted successfully", nil, 0)
}

func (h *ShopHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	shopIDStr := c.Param("shop_id")
	if shopIDStr == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("shop_id is required"), helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	product, err := h.ShopService.CreateProduct(ctx, req, shopIDStr, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Product created successfully", product, 0)
}

func (h *ShopHandler) GetProductsByShop(c *gin.Context) {
	shopIDStr := c.Param("shop_id")
	if shopIDStr == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("shop_id is required"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	products, err := h.ShopService.GetProductsByShop(ctx, shopIDStr, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Products retrieved successfully", products, 0)
}

func (h *ShopHandler) UpdateProduct(c *gin.Context) {
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	productIDStr := c.Param("product_id")
	if productIDStr == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("product_id is required"), helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.ShopService.UpdateProduct(ctx, productIDStr, req, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Product updated successfully", nil, 0)
}

func (h *ShopHandler) DeleteProduct(c *gin.Context) {
	productIDStr := c.Param("product_id")
	if productIDStr == "" {
		helper.SendError(c, http.StatusBadRequest, fmt.Errorf("product_id is required"), helper.ErrInvalidRequest)
		return
	}

	userID, exists := c.Get(constants.UserID)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("user_id not found"), helper.ErrInvalidRequest)
		return
	}

	token, exists := c.Get(constants.Token)
	if !exists {
		helper.SendError(c, http.StatusUnauthorized, fmt.Errorf("token not found"), helper.ErrInvalidRequest)
		return
	}

	ctx := context.WithValue(c, constants.TokenKey, token)

	err := h.ShopService.DeleteProduct(ctx, productIDStr, userID.(string))
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Product deleted successfully", nil, 0)
}

// AddRepairItem adds a product to a repair
func (h *ShopHandler) AddRepairItem(c *gin.Context) {
	var req AddRepairItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	repairIDStr := c.Param("repair_id")
	repairID, err := primitive.ObjectIDFromHex(repairIDStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	repairItem, err := h.ShopService.AddRepairItem(c.Request.Context(), repairID, req)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusCreated, "Item added to repair successfully", repairItem, 0)
}

// GetRepairItems gets all items for a repair with total cost
func (h *ShopHandler) GetRepairItems(c *gin.Context) {
	repairIDStr := c.Param("repair_id")
	repairID, err := primitive.ObjectIDFromHex(repairIDStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	summary, err := h.ShopService.GetRepairItems(c.Request.Context(), repairID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Repair items retrieved successfully", summary, 0)
}

// RemoveRepairItem removes an item from a repair
func (h *ShopHandler) RemoveRepairItem(c *gin.Context) {
	repairItemIDStr := c.Param("repair_item_id")
	repairItemID, err := primitive.ObjectIDFromHex(repairItemIDStr)
	if err != nil {
		helper.SendError(c, http.StatusBadRequest, err, helper.ErrInvalidRequest)
		return
	}

	err = h.ShopService.RemoveRepairItem(c.Request.Context(), repairItemID)
	if err != nil {
		helper.SendError(c, http.StatusInternalServerError, err, helper.ErrInvalidOperation)
		return
	}

	helper.SendSuccess(c, http.StatusOK, "Item removed from repair successfully", nil, 0)
}
