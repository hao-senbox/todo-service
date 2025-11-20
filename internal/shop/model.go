package shop

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Shop represents a shop/store
type Shop struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	OwnerID     string             `json:"owner_id" bson:"owner_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// Product represents items/products sold by shops
type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	ShopID      string             `json:"shop_id" bson:"shop_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	Category    string             `json:"category" bson:"category"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	Stock       int                `json:"stock" bson:"stock"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

// RepairItem represents items added to a repair from shops
type RepairItem struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	RepairID  primitive.ObjectID `json:"repair_id" bson:"repair_id"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Quantity  int                `json:"quantity" bson:"quantity"`
	Price     float64            `json:"price" bson:"price"` // Price at the time of adding
	AddedAt   time.Time          `json:"added_at" bson:"added_at"`
}
