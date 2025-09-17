package todo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Name           string             `json:"name" bson:"name"`
	OrganizationID string             `json:"organization_id" bson:"organization_id"`
	Description    *string            `json:"description" bson:"description"`
	DueDate        time.Time          `json:"due_date" bson:"due_date"`
	Urgent         bool               `json:"urgent" bson:"urgent"`
	Link           *string            `json:"link" bson:"link"`
	Progress       int                `json:"progress" bson:"progress"`
	Status         string             `json:"status" bson:"status"`
	Stage          *string            `json:"stage" bson:"stage"`
	QRCode         string             `json:"qrcode" bson:"qrcode"`
	Options        *string            `json:"options" bson:"options"`
	Pictures       []string           `json:"pictures" bson:"pictures"`
	ImageTask      string             `json:"image_task" bson:"image_task"`
	TaskUsers      TaskUsers          `json:"task_users" bson:"task_users"`
	Feedback       *string            `json:"feedback" bson:"feedback"`
	CreatedBy      string             `json:"created_by" bson:"created_by"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt      *string            `json:"deleted_at" bson:"deleted_at"`
	DeletedBy      *string            `json:"deleted_by" bson:"deleted_by"`
}

type TaskUsers struct {
	Teachers []string `json:"teachers" bson:"teachers"`
	Students []string `json:"students" bson:"students"`
	Staffs   []string `json:"staffs" bson:"staffs"`
}
