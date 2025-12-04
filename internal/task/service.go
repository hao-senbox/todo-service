package task

import (
	"context"
	"fmt"
	"log"
	"time"
	"todo-service/internal/uploader"
	"todo-service/internal/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskService interface {
	CreateTask(ctx context.Context, req CreateTaskRequest, userID string) (*string, error)
	GetTasks(ctx context.Context, role string, status string) ([]*TaskResponse, error)
	GetTaskById(ctx context.Context, id string) (*TaskResponse, error)
	UpdateTask(ctx context.Context, req UpdateTaskRequest, id string) error
	DeleteTask(ctx context.Context, id string) error

	GetMyTask(ctx context.Context, userID string) ([]*TaskResponse, error)
	UpdateTaskStatus(ctx context.Context, req []*UpdateTaskStatusRequest, id string, userID string) error
}

type taskService struct {
	TaskRepo    TaskRepository
	UserGateway user.UserService
	FileGateway uploader.ImageService
}

func NewTaskService(
	taskRepo TaskRepository,
	userGateway user.UserService,
	fileGateway uploader.ImageService,
) TaskService {
	return &taskService{
		TaskRepo:    taskRepo,
		UserGateway: userGateway,
		FileGateway: fileGateway,
	}
}

func (s *taskService) CreateTask(ctx context.Context, req CreateTaskRequest, userID string) (*string, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	if req.OrganizationID == "" {
		return nil, fmt.Errorf("organization_id is required")
	}

	if req.StartDate == "" {
		return nil, fmt.Errorf("start_date is required")
	}

	if req.DueDate == "" {
		return nil, fmt.Errorf("due_date is required")
	}

	if len(req.Group) == 0 {
		return nil, fmt.Errorf("group is required")
	}

	startDate, err := time.Parse("2006-01-02 15:04:05", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format, expected YYYY-MM-DD HH:MM:SS: %v", err)
	}

	dueDate, err := time.Parse("2006-01-02 15:04:05", req.DueDate)
	if err != nil {
		return nil, fmt.Errorf("invalid due_date format, expected YYYY-MM-DD HH:MM:SS: %v", err)
	}

	id := primitive.NewObjectID()

	var status string
	now := time.Now().Truncate(24 * time.Hour)
	start := startDate.Truncate(24 * time.Hour)

	if now.Before(start) {
		status = "not_started"
	} else {
		status = "progressing"
	}

	group := make([]UserRole, len(req.Group))
	for i, role := range req.Group {
		group[i] = UserRole{
			UserID: role.UserID,
			Role:   role.Role,
			Status: status,
		}
	}

	leader := make([]Leader, len(req.Leader))
	for i, reqLeader := range req.Leader {
		leader[i] = Leader{
			UserID: reqLeader.UserID,
			Role:   reqLeader.Role,
		}
	}

	task := &Task{
		ID:             id,
		Title:          req.Title,
		OrganizationID: req.OrganizationID,
		Leader:         leader,
		StartDate:      startDate,
		DueDate:        dueDate,
		Group:          group,
		File:           req.File,
		CreatedBy:      userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = s.TaskRepo.CreateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	idParse := id.Hex()
	return &idParse, nil
}

func (s *taskService) GetTasks(ctx context.Context, role string, status string) ([]*TaskResponse, error) {
	tasks, err := s.TaskRepo.GetTasks(ctx, role, status)
	if err != nil {
		return nil, err
	}

	results := make([]*TaskResponse, len(tasks))
	for i, task := range tasks {
		results[i] = MapTaskToResponse(ctx, task, s.UserGateway, s.FileGateway)
	}

	return results, nil
}

func (s *taskService) GetTaskById(ctx context.Context, id string) (*TaskResponse, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format, expected ObjectID: %v", err)
	}

	task, err := s.TaskRepo.GetTaskById(ctx, objectID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, fmt.Errorf("task not found")
	}

	return MapTaskToResponse(ctx, task, s.UserGateway, s.FileGateway), nil
}

func (s *taskService) UpdateTask(ctx context.Context, req UpdateTaskRequest, id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format, expected ObjectID: %v", err)
	}

	task, err := s.TaskRepo.GetTaskById(ctx, objectID)
	if err != nil {
		return err
	}

	if task == nil {
		return fmt.Errorf("task not found")
	}

	if req.Title != nil {
		task.Title = *req.Title
	}

	if req.File != nil {
		if task.File != nil && *task.File != *req.File {
			err = s.FileGateway.DeletePDFKey(ctx, *task.File)
			if err != nil {
				return err
			}
		}
		task.File = req.File
	}

	if req.Group != nil {
		task.Group = *req.Group
	}

	if req.Leader != nil {
		task.Leader = *req.Leader
	}

	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02 15:04:05", *req.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start_date format, expected YYYY-MM-DD HH:MM:SS: %v", err)
		}
		task.StartDate = startDate
	}

	if req.DueDate != nil {
		dueDate, err := time.Parse("2006-01-02 15:04:05", *req.DueDate)
		if err != nil {
			return fmt.Errorf("invalid due_date format, expected YYYY-MM-DD HH:MM:SS: %v", err)
		}
		task.DueDate = dueDate
	}

	task.UpdatedAt = time.Now()

	return s.TaskRepo.UpdateTask(ctx, objectID, task)
}

func (s *taskService) DeleteTask(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format, expected ObjectID: %v", err)
	}

	task, err := s.TaskRepo.GetTaskById(ctx, objectID)
	if err != nil {
		return err
	}

	if task == nil {
		return fmt.Errorf("task not found")
	}

	if task.File != nil {
		err = s.FileGateway.DeletePDFKey(ctx, *task.File)
		if err != nil {
			return err
		}
	}

	err = s.TaskRepo.DeleteTask(ctx, objectID)
	if err != nil {
		return err
	}

	return nil
}

func (s *taskService) GetMyTask(ctx context.Context, userID string) ([]*TaskResponse, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	task, err := s.TaskRepo.GetMyTask(ctx, userID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, fmt.Errorf("task not found")
	}

	results := make([]*TaskResponse, len(task))
	for i, task := range task {
		results[i] = MapTaskToResponse(ctx, task, s.UserGateway, s.FileGateway)
	}

	return results, nil
}

func (s *taskService) 	UpdateTaskStatus(ctx context.Context, req []*UpdateTaskStatusRequest, id string, userID string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format, expected ObjectID: %v", err)
	}

	task, err := s.TaskRepo.GetTaskById(ctx, objectID)
	if err != nil {
		return err
	}

	if task == nil {
		return fmt.Errorf("task not found")
	}

	// Update status for each group item in the request
	for _, groupUpdate := range req {
		found := false
		for i, group := range task.Group {
			if group.UserID == groupUpdate.UserID && group.Role == groupUpdate.Role {
				task.Group[i].Status = groupUpdate.Status
				found = true
				break
			}
		}
		if !found {
			log.Printf("user_id %s not found in group", groupUpdate.UserID)
		}
	}

	task.UpdatedAt = time.Now()

	return s.TaskRepo.UpdateTask(ctx, objectID, task)
}
