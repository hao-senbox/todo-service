package shop

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductResponse for product details
type ProductResponse struct {
	ID          primitive.ObjectID `json:"id"`
	ShopID      string             `json:"shop_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Price       float64            `json:"price"`
	Category    string             `json:"category"`
	ImageURL    string             `json:"image_url"`
	Stock       int                `json:"stock"`
	IsActive    bool               `json:"is_active"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// ShopProductsResponse for shop with its products
type ShopProductsResponse struct {
	Shop     Shop              `json:"shop"`
	Products []ProductResponse `json:"products"`
}

// RepairItemResponse for repair items
type RepairItemResponse struct {
	ID       primitive.ObjectID `json:"id"`
	Product  ProductResponse    `json:"product"`
	Quantity int                `json:"quantity"`
	Price    float64            `json:"price"`
	SubTotal float64            `json:"sub_total"` // quantity * price
	AddedAt  time.Time          `json:"added_at"`
}

// RepairItemsSummary for repair items summary
type RepairItemsSummary struct {
	Items     []RepairItemResponse `json:"items"`
	TotalCost float64              `json:"total_cost"`
	ItemCount int                  `json:"item_count"`
}
