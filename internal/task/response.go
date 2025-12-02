package task

import (
	"time"
	"todo-service/internal/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskResponse struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Title          string             `json:"title" bson:"title"`
	OrganizationID string             `json:"organization_id" bson:"organization_id"`
	Leader         []UserRoleResponse `json:"leader" bson:"leader"`
	StartDate      time.Time          `json:"start_date" bson:"start_date"`
	DueDate        time.Time          `json:"due_date" bson:"due_date"`
	Group          []UserRoleResponse `json:"group" bson:"group"`
	File           *string            `json:"file" bson:"file"`
	FileURL        *string            `json:"file_url" bson:"file_url"`
	CreatedBy      string             `json:"created_by" bson:"created_by"`
	CreatedByInfor *user.UserInfor    `json:"created_by_infor" bson:"created_by_infor"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserRoleResponse struct {
	UserInfor *user.UserInfor `json:"user_infor" bson:"user_infor"`
	Role      string          `json:"role" bson:"role"`
	Status    string          `json:"status" bson:"status"`
}
