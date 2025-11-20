package repair

type CreateRepairRequest struct {
	OrganizationID string   `json:"organization_id"`
	JobName        string   `json:"job_name"`
	Location       string   `json:"location"`
	UrgentVote     int      `json:"urgent_vote"`
	Comment        string   `json:"comment"`
	ImageReport    []string `json:"image_report"`
}

type UpdateRepairRequest struct {
	JobName     string   `json:"job_name"`
	Location    string   `json:"location"`
	UrgentVote  int      `json:"urgent_vote"`
	Comment     string   `json:"comment"`
	ImageReport []string `json:"image_report"`
}

type AssignRepairRequest struct {
	AssignedTo string `json:"assigned_to"`
}

type CompleteRepairRequest struct {
	CommentRepair *string  `json:"comment_repair"`
	ImageRepair   []string `json:"image_repair"`
}
