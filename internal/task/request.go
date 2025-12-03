package task

type CreateTaskRequest struct {
	Title          string     `json:"title" binding:"required"`
	OrganizationID string     `json:"organization_id" binding:"required"`
	StartDate      string     `json:"start_date" binding:"required"`
	DueDate        string     `json:"due_date" binding:"required"`
	Group          []UserRole `json:"group" binding:"required"`
	Leader         []Leader   `json:"leader" binding:"required"`
	File           *string    `json:"file"`
}

type UpdateTaskRequest struct {
	Title     *string     `json:"title" binding:"required"`
	StartDate *string     `json:"start_date" binding:"required"`
	DueDate   *string     `json:"due_date" binding:"required"`
	Group     *[]UserRole `json:"group" binding:"required"`
	Leader    *[]Leader   `json:"leader" binding:"required"`
	File      *string     `json:"file"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required"`
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}
