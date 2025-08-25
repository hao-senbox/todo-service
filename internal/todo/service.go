package todo

import (
	"context"
	"fmt"
	"log"
	"time"
	"todo-service/internal/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoService interface {
	GetAllTodo(ctx context.Context, status, name, teacher, student, staff string) ([]*TodoResponse, error)
	GetTodoByID(ctx context.Context, todoID string) (*TodoResponse, error)
	CreateTodo(ctx context.Context, req CreateTodoRequest, userID string) (*string, error)
	UpdateTodo(ctx context.Context, req UpdateTaskProgressRequest, id string) error
	DeleteTodo(ctx context.Context, id string) error
	// Join Todo
	JoinTodo(ctx context.Context, req JoinTodoRequest, userID string) error
	AddUser(ctx context.Context, req AddUserRequest) error
	GetMyTodo(ctx context.Context, userID string) ([]*TodoResponse, error)
}

type todoService struct {
	TodoRepo    TodoRepository
	UserService user.UserService
}

func NewTodoService(TodoRepo TodoRepository, UserService user.UserService) TodoService {
	return &todoService{
		TodoRepo:    TodoRepo,
		UserService: UserService,
	}
}

func (s *todoService) GetAllTodo(ctx context.Context, status, name, teacher, student, staff string) ([]*TodoResponse, error) {
	todos, err := s.TodoRepo.GetAllTodo(ctx, status, name, teacher, student, staff)
	if err != nil {
		return nil, err
	}

	var results []*TodoResponse
	for _, todo := range todos {
		results = append(results, s.buildTodoResponse(ctx, todo))
	}
	return results, nil
}


func (s *todoService) GetTodoByID(ctx context.Context, todoID string) (*TodoResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		return nil, err
	}

	todo, err := s.TodoRepo.GetTodoByID(ctx, objectID)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, fmt.Errorf("todo not found")
	}

	return s.buildTodoResponse(ctx, todo), nil
}


func (s *todoService) CreateTodo(ctx context.Context, req CreateTodoRequest, userID string) (*string, error) {

	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if req.DueDate == "" {
		return nil, fmt.Errorf("due_date is required")
	}

	dueDate, err := time.Parse("2006-01-02 15:04:05", req.DueDate)
	if err != nil {
		return nil, fmt.Errorf("invalid due_date format, expected YYYY-MM-DD HH:MM:SS: %v", err)
	}

	if req.ImageTask == "" {
		return nil, fmt.Errorf("image_task is required")
	}

	if !req.Urgent {
		return nil, fmt.Errorf("urgent is required")
	}

	ID := primitive.NewObjectID()

	QRCocde := fmt.Sprintf("SENBOX.ORG[TODO]:%s", ID.Hex())

	todo := &Todo{
		ID:          ID,
		Name:        req.Name,
		Description: req.Description,
		DueDate:     dueDate,
		Urgent:      req.Urgent,
		Link:        req.Link,
		Progress:    0,
		Status:      "pending",
		Stage:       req.Stage,
		QRCode:      QRCocde,
		Options:     req.Options,
		CreatedBy:   userID,
		Pictures:    []string{},
		TaskUsers: TaskUsers{
			Teachers: []string{},
			Students: []string{},
			Staff:    []string{},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ImageTask: req.ImageTask,
		DeletedAt: nil,
		DeletedBy: nil,
	}

	id, err := s.TodoRepo.CreateTodo(ctx, todo)
	if err != nil {
		return nil, err
	}

	return id, nil

}

func (s *todoService) UpdateTodo(ctx context.Context, req UpdateTaskProgressRequest, id string) error {

	if id == "" {
		return fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	existingTodo, err := s.TodoRepo.GetTodoByID(ctx, objectID)
	if err != nil {
		return err
	}

	if existingTodo == nil {
		return fmt.Errorf("todo not found")
	}

	if req.Progress < 0 || req.Progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}

	if existingTodo.Progress >= req.Progress {
		return fmt.Errorf("progress must be greater than the current progress")
	}

	if req.Name == "" {
		req.Name = existingTodo.Name
	}

	if req.Description == nil {
		req.Description = existingTodo.Description
	}

	if req.DueDate == "" {
		req.DueDate = existingTodo.DueDate.Format("2006-01-02 15:04:05")
	}

	if req.Link == nil {
		req.Link = existingTodo.Link
	}

	if req.Stage == nil {
		req.Stage = existingTodo.Stage
	}

	if req.Options == nil {
		req.Options = existingTodo.Options
	}

	var newPictures []string

	for _, picture := range req.Pictures {
		if picture != "" {
			newPictures = append(newPictures, picture)
		}
	}

	for _, picture := range existingTodo.Pictures {
		if picture != "" {
			newPictures = append(newPictures, picture)
		}
	}

	req.Pictures = newPictures

	if req.ImageTask == "" {
		req.ImageTask = existingTodo.ImageTask
	}

	if req.Status == "" {
		req.Status = existingTodo.Status
	}

	updatedAt := time.Now()

	todo := &Todo{
		ID:          existingTodo.ID,
		Name:        req.Name,
		Description: req.Description,
		DueDate:     existingTodo.DueDate,
		Urgent:      existingTodo.Urgent,
		Link:        req.Link,
		Progress:    req.Progress,
		Status:      req.Status,
		Stage:       req.Stage,
		QRCode:      existingTodo.QRCode,
		Options:     req.Options,
		CreatedBy:   existingTodo.CreatedBy,
		Pictures:    req.Pictures,
		TaskUsers:   existingTodo.TaskUsers,
		CreatedAt:   existingTodo.CreatedAt,
		UpdatedAt:   updatedAt,
		ImageTask:   req.ImageTask,
		DeletedAt:   existingTodo.DeletedAt,
		DeletedBy:   existingTodo.DeletedBy,
	}

	return s.TodoRepo.UpdateTodo(ctx, todo)

}

func (s *todoService) DeleteTodo(ctx context.Context, id string) error {

	if id == "" {
		return fmt.Errorf("id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return s.TodoRepo.DeleteTodo(ctx, objectID)

}

func (s *todoService) JoinTodo(ctx context.Context, req JoinTodoRequest, userID string) error {

	if req.QRCode == "" {
		return fmt.Errorf("qrcode is required")
	}

	if userID == "" {
		return fmt.Errorf("user id is required")
	}

	if req.Type == "" {
		return fmt.Errorf("type is required")
	}

	todoExist, err := s.TodoRepo.GetTodoByQRCode(ctx, req.QRCode)
	if err != nil {
		return err
	}

	if todoExist.CreatedBy == userID {
		return fmt.Errorf("you cannot join your own todo")
	}

	if todoExist == nil {
		return fmt.Errorf("todo not found")
	}

	if len(todoExist.TaskUsers.Students) >= 1 {
		for _, student := range todoExist.TaskUsers.Students {
			if student == userID {
				return fmt.Errorf("you have already joined this todo")
			}
		}
	}

	if len(todoExist.TaskUsers.Teachers) >= 1 {
		for _, teacher := range todoExist.TaskUsers.Teachers {
			if teacher == userID {
				return fmt.Errorf("you have already joined this todo")
			}
		}
	}

	if len(todoExist.TaskUsers.Staff) >= 1 {
		for _, staff := range todoExist.TaskUsers.Staff {
			if staff == userID {
				return fmt.Errorf("you have already joined this todo")
			}
		}
	}

	return s.TodoRepo.JoinTodo(ctx, todoExist.ID, userID, req.Type)
}

func (s *todoService) AddUser(ctx context.Context, req AddUserRequest) error {

	if req.TodoID == "" {
		return fmt.Errorf("todo id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(req.TodoID)
	if err != nil {
		return err
	}

	if req.UserID == "" {
		return fmt.Errorf("user id is required")
	}

	if req.Type == "" {
		return fmt.Errorf("type is required")
	}

	return s.TodoRepo.AddUser(ctx, objectID, req.UserID, req.Type)

}

func (s *todoService) GetMyTodo(ctx context.Context, userID string) ([]*TodoResponse, error) {

	if userID == "" {
		return nil, fmt.Errorf("user id is required")
	}

	myTodo, err := s.TodoRepo.GetMyTodo(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(myTodo) == 0 {
		return nil, fmt.Errorf("todo not found")
	}

	var results []*TodoResponse

	for _, todo := range myTodo {
		var createdBy TaskUser
		if todo.CreatedBy != "" {
			createdByInfor, err := s.UserService.GetUserInfor(ctx, todo.CreatedBy)
			if err != nil {
				log.Printf("[WARN] failed to get createdBy user info for %s: %v", todo.CreatedBy, err)
			} else {
				createdBy = TaskUser{
					UserID:         createdByInfor.UserID,
					UserName:       createdByInfor.UserName,
					Avartar:       createdByInfor.Avartar,
				}
			}
		}

		var taskUsersResp TaskUsersResponse

		for _, teacherID := range todo.TaskUsers.Teachers {
			if teacherID == "" {
				continue
			}
			info, err := s.UserService.GetTeacherInfor(ctx, teacherID)
			if err != nil {
				log.Printf("[WARN] failed to get teacher info for %s: %v", teacherID, err)
				continue
			}
			taskUsersResp.Teachers = append(taskUsersResp.Teachers, TaskUser{
				UserID:         info.UserID,
				UserName:       info.UserName,
				Avartar:       info.Avartar,
			})
		}

		for _, studentID := range todo.TaskUsers.Students {
			if studentID == "" {
				continue
			}
			info, err := s.UserService.GetStudentInfor(ctx, studentID)
			if err != nil {
				log.Printf("[WARN] failed to get student info for %s: %v", studentID, err)
				continue
			}
			taskUsersResp.Students = append(taskUsersResp.Students, TaskUser{
				UserID:         info.UserID,
				UserName:       info.UserName,
				Avartar:       info.Avartar,
			})
		}

		for _, staffID := range todo.TaskUsers.Staff {
			if staffID == "" {
				continue
			}
			info, err := s.UserService.GetStaffInfor(ctx, staffID)
			if err != nil {
				log.Printf("[WARN] failed to get staff info for %s: %v", staffID, err)
				continue
			}
			taskUsersResp.Staff = append(taskUsersResp.Staff, TaskUser{
				UserID:         info.UserID,
				UserName:       info.UserName,
				Avartar:       info.Avartar,
			})
		}

		todoResp := &TodoResponse{
			ID:          todo.ID,
			Name:        todo.Name,
			Description: todo.Description,
			DueDate:     todo.DueDate,
			Urgent:      todo.Urgent,
			Link:        todo.Link,
			Progress:    todo.Progress,
			Stage:       todo.Stage,
			QRCode:      todo.QRCode,
			Options:     todo.Options,
			CreatedBy:   createdBy,
			Pictures:    todo.Pictures,
			ImageTask:   todo.ImageTask,
			TaskUsers:   taskUsersResp,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
			DeletedAt:   todo.DeletedAt,
			DeletedBy:   todo.DeletedBy,
		}

		results = append(results, todoResp)
	}

	return results, nil
}

func (s *todoService) buildTodoResponse(ctx context.Context, todo *Todo) *TodoResponse {

	var createdBy TaskUser
	if todo.CreatedBy != "" {
		createdByInfor, err := s.UserService.GetUserInfor(ctx, todo.CreatedBy)
		if err != nil {
			log.Printf("[WARN] failed to get createdBy user info for %s: %v", todo.CreatedBy, err)
		} else {
			createdBy = TaskUser{
				UserID:   createdByInfor.UserID,
				UserName: createdByInfor.UserName,
				Avartar: createdByInfor.Avartar,
			}
		}
	}

	var taskUsersResp TaskUsersResponse

	for _, teacherID := range todo.TaskUsers.Teachers {
		if teacherID == "" {
			continue
		}
		info, err := s.UserService.GetTeacherInfor(ctx, teacherID)
		if err != nil {
			log.Printf("[WARN] failed to get teacher info for %s: %v", teacherID, err)
			continue
		}
		taskUsersResp.Teachers = append(taskUsersResp.Teachers, TaskUser{
			UserID:   info.UserID,
			UserName: info.UserName,
			Avartar: info.Avartar,
		})
	}

	for _, studentID := range todo.TaskUsers.Students {
		if studentID == "" {
			continue
		}
		info, err := s.UserService.GetStudentInfor(ctx, studentID)
		if err != nil {
			log.Printf("[WARN] failed to get student info for %s: %v", studentID, err)
			continue
		}
		taskUsersResp.Students = append(taskUsersResp.Students, TaskUser{
			UserID:   info.UserID,
			UserName: info.UserName,
			Avartar: info.Avartar,
		})
	}

	for _, staffID := range todo.TaskUsers.Staff {
		if staffID == "" {
			continue
		}
		info, err := s.UserService.GetStaffInfor(ctx, staffID)
		if err != nil {
			log.Printf("[WARN] failed to get staff info for %s: %v", staffID, err)
			continue
		}
		taskUsersResp.Staff = append(taskUsersResp.Staff, TaskUser{
			UserID:   info.UserID,
			UserName: info.UserName,
			Avartar: info.Avartar,
		})
	}

	return &TodoResponse{
		ID:          todo.ID,
		Name:        todo.Name,
		Description: todo.Description,
		DueDate:     todo.DueDate,
		Urgent:      todo.Urgent,
		Link:        todo.Link,
		Progress:    todo.Progress,
		Stage:       todo.Stage,
		QRCode:      todo.QRCode,
		Options:     todo.Options,
		CreatedBy:   createdBy,
		Pictures:    todo.Pictures,
		ImageTask:   todo.ImageTask,
		TaskUsers:   taskUsersResp,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
		DeletedAt:   todo.DeletedAt,
		DeletedBy:   todo.DeletedBy,
	}

}