package todo

import (
	"time"
	"todo-service/internal/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoResponse struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description *string            `json:"description" bson:"description"`
	DueDate     time.Time          `json:"due_date" bson:"due_date"`
	Urgent      bool               `json:"urgent" bson:"urgent"`
	Link        *string            `json:"link" bson:"link"`
	Progress    int                `json:"progress" bson:"progress"`
	Stage       *string            `json:"stage" bson:"stage"`
	QRCode      string             `json:"qrcode" bson:"qrcode"`
	Options     *string            `json:"options" bson:"options"`
	CreatedBy   TaskUser           `json:"created_by" bson:"created_by"`
	Pictures    []string           `json:"pictures" bson:"pictures"`
	ImageTask   string             `json:"image_task" bson:"image_task"`
	TaskUsers   TaskUsersResponse  `json:"task_users" bson:"task_users"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt   *string            `json:"deleted_at" bson:"deleted_at"`
	DeletedBy   *string            `json:"deleted_by" bson:"deleted_by"`
}

type TaskUsersResponse struct {
	Teachers []TaskUser `json:"teachers" bson:"teachers"`
	Students []TaskUser `json:"students" bson:"students"`
	Staff    []TaskUser `json:"staff" bson:"staff"`
}

type TaskUser struct {
	UserID         string   `json:"user_id"`
	UserName       string   `json:"user_name"`
	Avartars       []user.Avatar `json:"avatars"`
}
