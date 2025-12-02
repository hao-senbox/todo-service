package task

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Title          string             `json:"title" bson:"title"`
	OrganizationID string             `json:"organization_id" bson:"organization_id"`
	Leader         []Leader           `json:"leader" bson:"leader"`
	StartDate      time.Time          `json:"start_date" bson:"start_date"`
	DueDate        time.Time          `json:"due_date" bson:"due_date"`
	Group          []UserRole         `json:"group" bson:"group"`
	File           *string            `json:"file" bson:"file"`
	FileURL        *string            `json:"file_url" bson:"file_url"`
	CreatedBy      string             `json:"created_by" bson:"created_by"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserRole struct {
	UserID string `json:"user_id" bson:"user_id"`
	Role   string `json:"role" bson:"role"`
	Status string `json:"status" bson:"status"`
}

type Leader struct {
	UserID string `json:"user_id" bson:"user_id"`
	Role   string `json:"role" bson:"role"`
}
