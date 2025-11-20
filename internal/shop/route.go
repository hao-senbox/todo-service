package shop

import (
	"todo-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, shopHandler *ShopHandler) {
	shopGroup := r.Group("/api/v1/shops", middleware.Secured())
	{
		shopGroup.POST("", shopHandler.CreateShop)
		shopGroup.GET("/:id", shopHandler.GetShopByID)
		shopGroup.GET("/my-shop", shopHandler.GetMyShop)
		shopGroup.PUT("/:id", shopHandler.UpdateShop)
        shopGroup.DELETE("/:id", shopHandler.DeleteShop)

		shopGroup.POST("products/:shop_id", shopHandler.CreateProduct)
		shopGroup.GET("products/:shop_id", shopHandler.GetProductsByShop)

		shopGroup.PUT("/products/:product_id", shopHandler.UpdateProduct)
		shopGroup.DELETE("/products/:product_id", shopHandler.DeleteProduct)

		shopGroup.POST("/repairs/:repair_id/items", shopHandler.AddRepairItem)
		shopGroup.GET("/repairs/:repair_id/items", shopHandler.GetRepairItems)
		shopGroup.DELETE("/repair-items/:repair_item_id", shopHandler.RemoveRepairItem)
	}
}

