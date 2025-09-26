package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware gán 1 giá trị X-Request-ID vào header và context
// RequestIDMiddleware lấy X-Request-ID từ request (nếu có) và gán vào context + response
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy từ request header
		requestID := c.GetHeader("X-Request-ID")

		// Nếu client không gửi, thì tự sinh
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Gán vào context
		c.Set("X-Request-ID", requestID)

		c.Next()
	}
}
