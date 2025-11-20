package repair

import (
	"todo-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, repairHandler *RepairHandler) {
	repairGroup := r.Group("/api/v1/repairs", middleware.Secured())
	{
		repairGroup.POST("", repairHandler.CreateRepair)
		repairGroup.GET("", repairHandler.GetRepairs)
		repairGroup.GET("/:id", repairHandler.GetRepairByID)
		repairGroup.PUT("/:id", repairHandler.UpdateRepair)
		repairGroup.DELETE("/:id", repairHandler.DeleteRepair)

		repairGroup.POST("/:id/assign", repairHandler.AssignRepair)
		repairGroup.POST("/:id/complete", repairHandler.CompleteRepair)
	}
}
