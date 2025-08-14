package todo

import (
	"todo-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, todoHanlder *TodoHandler) {
	todoGroup := r.Group("/api/v1/todos", middleware.Secured())
	{
		todoGroup.GET("", todoHanlder.GetTodos)
		todoGroup.GET("/:id", todoHanlder.GetTodo)
		todoGroup.POST("", todoHanlder.CreateTodo)
		todoGroup.PUT("/:id", todoHanlder.UpdateTodo)
		todoGroup.DELETE("/:id", todoHanlder.DeleteTodo)
		//Add other routes
		todoGroup.POST("/join", todoHanlder.JoinTodo)
		todoGroup.GET("/my-todo", todoHanlder.GetMyTodo)

	}
}