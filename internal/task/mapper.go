package task

import (
	"context"
	"log"
	"todo-service/internal/uploader"
	"todo-service/internal/user"
)

func MapTaskToResponse(
	ctx context.Context,
	task *Task,
	userService user.UserService,
	fileService uploader.ImageService,
) *TaskResponse {

	if task == nil {
		return nil
	}

	group := make([]UserRoleResponse, len(task.Group))
	for i, role := range task.Group {
		var userInfor *user.UserInfor
		var err error

		switch role.Role {
		case "teacher":
			userInfor, err = userService.GetTeacherInfor(ctx, role.UserID)
		case "staff":
			userInfor, err = userService.GetStaffInfor(ctx, role.UserID)
		case "parent":
			userInfor, err = userService.GetParentInfor(ctx, role.UserID)
		case "student":
			userInfor, err = userService.GetStudentInfor(ctx, role.UserID)
		default:
			continue
		}

		if err != nil {
			return nil
		}

		group[i] = UserRoleResponse{
			UserInfor: userInfor,
			Role:      role.Role,
			Status:    role.Status,
		}
	}

	leader := make([]UserRoleResponse, len(task.Leader))
	for i, role := range task.Leader {
		var userInfor *user.UserInfor
		var err error

		switch role.Role {
		case "teacher":
			userInfor, err = userService.GetTeacherInfor(ctx, role.UserID)
		case "staff":
			userInfor, err = userService.GetStaffInfor(ctx, role.UserID)
		case "parent":
			userInfor, err = userService.GetParentInfor(ctx, role.UserID)
		default:
			continue
		}

		if err != nil {
			return nil
		}

		leader[i] = UserRoleResponse{
			UserInfor: userInfor,
			Role:      role.Role,
		}
	}

	var fileURL *string
	if task.File != nil {
		pdf, err := fileService.GetPDFKey(ctx, *task.File)
		if err != nil {
			return nil
		}
		log.Println("pdf", pdf)
		if pdf != nil {
			fileURL = &pdf.Url
		}
	}

	createdBy, err := userService.GetUserInfor(ctx, task.CreatedBy)
	if err != nil {
		return nil
	}

	return &TaskResponse{
		ID:             task.ID,
		Title:          task.Title,
		OrganizationID: task.OrganizationID,
		Leader:         leader,
		StartDate:      task.StartDate,
		DueDate:        task.DueDate,
		Group:          group,
		File:           task.File,
		FileURL:        fileURL,
		CreatedBy:      task.CreatedBy,
		CreatedByInfor: createdBy,
		CreatedAt:      task.CreatedAt,
		UpdatedAt:      task.UpdatedAt,
	}
}
