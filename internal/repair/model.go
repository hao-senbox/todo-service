package repair

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repair struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	OrganizationID string             `json:"organization_id" bson:"organization_id"`
	JobNumber      int                `json:"job_number" bson:"job_number"`
	JobName        string             `json:"job_name" bson:"job_name"`
	QRCode         string             `json:"qrcode" bson:"qrcode"`
	Location       string             `json:"location" bson:"location"`
	Status         string             `json:"status" bson:"status"`
	AssignedTo     *string            `json:"assigned_to" bson:"assigned_to"`
	// Report by
	DateReport  time.Time `json:"date_report" bson:"date_report"`
	ReportBy    string    `json:"report_by" bson:"report_by"`
	UrgentVote  int       `json:"urgent_vote" bson:"urgent_vote"`
	Comment     string    `json:"comment" bson:"comment"`
	ImageReport []string  `json:"image_report" bson:"image_report"`
	// Repair by
	DateRepair    *time.Time `json:"date_repair" bson:"date_repair"`
	RepairBy      *string    `json:"repair_by" bson:"repair_by"`
	CommentRepair *string    `json:"comment_repair" bson:"comment_repair"`
	ImageRepair   []*string  `json:"image_repair" bson:"image_repair"`
	TotalCost     *float64   `json:"total_cost" bson:"total_cost"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
