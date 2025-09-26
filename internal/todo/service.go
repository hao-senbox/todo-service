package todo

import (
	"context"
	"fmt"
	"log"
	"math"
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
	JoinTodo(ctx context.Context, req JoinTodoRequest, userID string, isCreator bool) error
	AddUser(ctx context.Context, req AddUserRequest) error
	GetMyTodo(ctx context.Context, userID string) ([]*TodoResponse, float64, error)
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
		if todo == nil {
			continue
		}
		response := s.buildTodoResponse(ctx, todo)
		if response != nil {
			results = append(results, response)
		}
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

	if req.OrganizationID == "" {
		return nil, fmt.Errorf("organization id is required")
	}

	dueDate, err := time.Parse("2006-01-02 15:04:05", req.DueDate)
	if err != nil {
		return nil, fmt.Errorf("invalid due_date format, expected YYYY-MM-DD HH:MM:SS: %v", err)
	}

	ID := primitive.NewObjectID()

	QRCocde := fmt.Sprintf("SENBOX.ORG[TODO]:%s", ID.Hex())

	todo := &Todo{
		ID:             ID,
		Name:           req.Name,
		OrganizationID: req.OrganizationID,
		Description:    req.Description,
		DueDate:        dueDate,
		Urgent:         req.Urgent,
		Link:           req.Link,
		Progress:       0,
		Status:         "pending",
		Stage:          req.Stage,
		QRCode:         QRCocde,
		Options:        req.Options,
		CreatedBy:      userID,
		Pictures:       []string{},
		TaskUsers: TaskUsers{
			Teachers: []string{},
			Students: []string{},
			Staffs:   []string{},
		},
		Feedback:  nil,
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

	if !req.Urgent {
		req.Urgent = existingTodo.Urgent
	} else if req.Urgent == true {
		req.Urgent = true
	} else {
		req.Urgent = false
	}

	updatedAt := time.Now()

	todo := &Todo{
		ID:             existingTodo.ID,
		Name:           req.Name,
		OrganizationID: existingTodo.OrganizationID,
		Description:    req.Description,
		DueDate:        existingTodo.DueDate,
		Urgent:         req.Urgent,
		Link:           req.Link,
		Progress:       req.Progress,
		Status:         req.Status,
		Stage:          req.Stage,
		QRCode:         existingTodo.QRCode,
		Options:        req.Options,
		CreatedBy:      existingTodo.CreatedBy,
		Pictures:       req.Pictures,
		TaskUsers:      existingTodo.TaskUsers,
		CreatedAt:      existingTodo.CreatedAt,
		UpdatedAt:      updatedAt,
		Feedback:       req.Feedback,
		ImageTask:      req.ImageTask,
		DeletedAt:      existingTodo.DeletedAt,
		DeletedBy:      existingTodo.DeletedBy,
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

func (s *todoService) JoinTodo(ctx context.Context, req JoinTodoRequest, userID string, isCreator bool) error {

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

	if todoExist == nil {
		return fmt.Errorf("todo not found")
	}

	if todoExist.CreatedBy == userID {
		return fmt.Errorf("you cannot join your own todo")
	}

	if len(todoExist.TaskUsers.Students) > 0 {
		for _, student := range todoExist.TaskUsers.Students {
			if student == userID && req.Type == "students" {
				return fmt.Errorf("you have already joined this todo as student")
			}
		}
	}

	if len(todoExist.TaskUsers.Staffs) > 0 {
		for _, staff := range todoExist.TaskUsers.Staffs {
			if staff == userID && req.Type == "staffs" {
				return fmt.Errorf("you have already joined this todo as staff")
			}
		}
	}

	if len(todoExist.TaskUsers.Teachers) > 0 {
		for _, teacher := range todoExist.TaskUsers.Teachers {
			if teacher == userID && req.Type == "teachers" {
				return fmt.Errorf("you have already joined this todo as teacher")
			}
		}
	}

	return s.TodoRepo.JoinTodo(ctx, todoExist.ID, userID, req.Type, isCreator)
}

func (s *todoService) AddUser(ctx context.Context, req AddUserRequest) error {

	if req.TodoID == "" {
		return fmt.Errorf("todo id is required")
	}

	objectID, err := primitive.ObjectIDFromHex(req.TodoID)
	if err != nil {
		return err
	}

	todo, err := s.TodoRepo.GetTodoByID(ctx, objectID)
	if err != nil {
		return err
	}

	if todo == nil {
		return fmt.Errorf("todo not found")
	}

	if req.Type == "" {
		return fmt.Errorf("type is required")
	}

	if len(req.UserIDs) == 0 {
		return fmt.Errorf("user array is required")
	}

	return s.TodoRepo.AddUsers(ctx, objectID, req.UserIDs, req.Type)

}

func (s *todoService) GetMyTodo(ctx context.Context, userID string) ([]*TodoResponse, float64, error) {
	if userID == "" {
		return nil, 0, fmt.Errorf("user id is required")
	}

	myTodo, err := s.TodoRepo.GetMyTodo(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	if len(myTodo) == 0 {
		return nil, 0, fmt.Errorf("todo not found")
	}

	var results []*TodoResponse
	var totalProgress float64

	for _, todo := range myTodo {
		var createdBy TaskUser
		if todo.CreatedBy != "" {
			createdByInfor, err := s.UserService.GetUserInfor(ctx, todo.CreatedBy)
			if err != nil {
				log.Printf("[WARN] failed to get createdBy user info for %s: %v", todo.CreatedBy, err)
			} else {
				createdBy = TaskUser{
					UserID:   createdByInfor.UserID,
					UserName: createdByInfor.UserName,
					Avartar:  createdByInfor.Avartar,
				}
			}
		}

		var taskUsersResp TaskUsersResponse

		for _, teacher := range todo.TaskUsers.Teachers {
			if teacher == "" {
				continue
			}
			info, err := s.UserService.GetTeacherInfor(ctx, teacher)
			if err != nil {
				log.Printf("[WARN] failed to get teacher info for %s: %v", teacher, err)
				continue
			}
			taskUsersResp.Teachers = append(taskUsersResp.Teachers, TaskUser{
				UserID:   info.UserID,
				UserName: info.UserName,
				Avartar:  info.Avartar,
			})
		}

		for _, student := range todo.TaskUsers.Students {
			if student == "" {
				continue
			}
			info, err := s.UserService.GetStudentInfor(ctx, student)
			if err != nil {
				log.Printf("[WARN] failed to get student info for %s: %v", student, err)
				continue
			}
			taskUsersResp.Students = append(taskUsersResp.Students, TaskUser{
				UserID:   info.UserID,
				UserName: info.UserName,
				Avartar:  info.Avartar,
			})
		}

		for _, staff := range todo.TaskUsers.Staffs {
			if staff == "" {
				continue
			}
			info, err := s.UserService.GetStaffInfor(ctx, staff)
			if err != nil {
				log.Printf("[WARN] failed to get staff info for %s: %v", staff, err)
				continue
			}
			taskUsersResp.Staffs = append(taskUsersResp.Staffs, TaskUser{
				UserID:   info.UserID,
				UserName: info.UserName,
				Avartar:  info.Avartar,
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
			FeedBack:    todo.Feedback,
			TaskUsers:   taskUsersResp,
			CreatedAt:   todo.CreatedAt,
			UpdatedAt:   todo.UpdatedAt,
			DeletedAt:   todo.DeletedAt,
			DeletedBy:   todo.DeletedBy,
		}

		results = append(results, todoResp)

		totalProgress += float64(todo.Progress)
	}

	avgProgress := totalProgress / float64(len(myTodo))
	avgProgress = math.Round(avgProgress*100) / 100

	return results, avgProgress, nil
}

func safeCreateTaskUser(info *user.UserInfor) TaskUser {
	if info == nil {
		return TaskUser{}
	}
	return TaskUser{
		UserID:   info.UserID,
		UserName: info.UserName,
		Avartar:  info.Avartar,
	}
}

func (s *todoService) buildTodoResponse(ctx context.Context, todo *Todo) *TodoResponse {

	if todo == nil {
		log.Printf("[ERROR] buildTodoResponse: todo is nil")
		return nil
	}

	var createdBy TaskUser

	var Roles []string
	if todo.CreatedBy != "" {
		if createdByInfor, err := s.UserService.GetUserInfor(ctx, todo.CreatedBy); err != nil {
			log.Printf("[WARN] failed to get createdBy user info for %s: %v", todo.CreatedBy, err)
		} else {

			createdBy = safeCreateTaskUser(createdByInfor)

			teacher, err := s.UserService.GetTeacherInforByOrg(ctx, todo.CreatedBy, todo.OrganizationID)
			if err != nil {
				return nil
			}

			if teacher != nil {
				Roles = append(Roles, "teacher")
			}

			staff, err := s.UserService.GetStaffInforByOrg(ctx, todo.CreatedBy, todo.OrganizationID)
			if err != nil {
				return nil
			}

			if staff != nil {
				Roles = append(Roles, "staff")
			}

			if teacher == nil && staff == nil {
				Roles = append(Roles, "user")
			}

			createdBy.Roles = Roles

		}
	}

	var taskUsersResp TaskUsersResponse

	if todo.TaskUsers.Teachers != nil {
		for _, teacher := range todo.TaskUsers.Teachers {
			if teacher == "" {
				log.Printf("[WARN] teacher id is empty")
				continue
			}
			if info, err := s.UserService.GetTeacherInfor(ctx, teacher); err != nil {
				log.Printf("[WARN] failed to get teacher info for %s: %v", teacher, err)
			} else if info != nil {
				taskUsersResp.Teachers = append(taskUsersResp.Teachers, safeCreateTaskUser(info))
			}
		}
	}

	if todo.TaskUsers.Students != nil {
		for _, student := range todo.TaskUsers.Students {
			if student == "" {
				log.Printf("[WARN] student id is empty")
				continue
			}
			if info, err := s.UserService.GetStudentInfor(ctx, student); err != nil {
				log.Printf("[WARN] failed to get student info for %s: %v", student, err)
			} else if info != nil {
				taskUsersResp.Students = append(taskUsersResp.Students, safeCreateTaskUser(info))
			}
		}
	}

	if todo.TaskUsers.Staffs != nil {
		for _, staff := range todo.TaskUsers.Staffs {
			if staff == "" {
				log.Printf("[WARN] staff id is empty")
				continue
			}
			if info, err := s.UserService.GetStaffInfor(ctx, staff); err != nil {
				log.Printf("[WARN] failed to get staff info for %s: %v", staff, err)
			} else if info != nil {
				taskUsersResp.Staffs = append(taskUsersResp.Staffs, safeCreateTaskUser(info))
			}
		}
	}

	data := &TodoResponse{
		ID:             todo.ID,
		Name:           todo.Name,
		Description:    todo.Description,
		OrganizationID: todo.OrganizationID,
		DueDate:        todo.DueDate,
		Urgent:         todo.Urgent,
		Link:           todo.Link,
		Progress:       todo.Progress,
		Stage:          todo.Stage,
		Status:         todo.Status,
		QRCode:         todo.QRCode,
		Options:        todo.Options,
		CreatedBy:      createdBy,
		Pictures:       todo.Pictures,
		ImageTask:      todo.ImageTask,
		FeedBack:       todo.Feedback,
		TaskUsers:      taskUsersResp,
		CreatedAt:      todo.CreatedAt,
		UpdatedAt:      todo.UpdatedAt,
		DeletedAt:      todo.DeletedAt,
		DeletedBy:      todo.DeletedBy,
	}

	return data
}
