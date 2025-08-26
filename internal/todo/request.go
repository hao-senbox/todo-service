package todo

type CreateTodoRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	DueDate     string  `json:"due_date"`
	Urgent      bool    `json:"urgent"`
	Link        *string `json:"link"`
	Stage       *string `json:"stage"`
	Options     *string `json:"options"`
	CreatedBy   string  `json:"created_by"`
	ImageTask   string  `json:"image_task"`
}

type UpdateTaskProgressRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	DueDate     string   `json:"due_date"`
	Urgent      bool     `json:"urgent"`
	Status      string   `json:"status"`
	Link        *string  `json:"link"`
	Stage       *string  `json:"stage"`
	Options     *string  `json:"options"`
	QRCode      string   `json:"qrcode"`
	CreatedBy   string   `json:"created_by"`
	Pictures    []string `json:"pictures"`
	ImageTask   string   `json:"image_task"`
	Progress    int      `json:"progress"`
}

type JoinTodoRequest struct {
	QRCode string `json:"qrcode"`
	Type   string `json:"type"`
}

type AddUserRequest struct {
	TodoID  string   `json:"todo_id"`
	UserIDs []string `json:"user_ids"`
	Type    string   `json:"type"`
}
