package task

import (
	"todo-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, taskHandler *TaskHandler) {
	taskGroup := r.Group("/api/v1/tasks", middleware.Secured())
	{
		taskGroup.POST("", taskHandler.CreateTask)
		taskGroup.GET("", taskHandler.GetTasks)
		taskGroup.GET("/:id", taskHandler.GetTaskById)
		taskGroup.PUT("/:id", taskHandler.UpdateTask)
		taskGroup.DELETE("/:id", taskHandler.DeleteTask)
		taskGroup.GET("/my-task", taskHandler.GetMyTask)
		taskGroup.POST("/update-status/:id", taskHandler.UpdateTaskStatus)
	}
}
