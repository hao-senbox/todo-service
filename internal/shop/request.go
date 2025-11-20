package shop

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateShopRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateShopRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// CreateProductRequest for creating new products
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock" binding:"min=0"`
}

// UpdateProductRequest for updating products
type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Category    *string  `json:"category"`
	ImageURL    *string  `json:"image_url"`
	Stock       *int     `json:"stock"`
	IsActive    *bool    `json:"is_active"`
}

// AddRepairItemRequest for adding items to repair
type AddRepairItemRequest struct {
	ProductID primitive.ObjectID `json:"product_id" binding:"required"`
	Quantity  int                `json:"quantity" binding:"required,min=1"`
}
