package middleware

import "github.com/gin-gonic/gin"

func LoggingMiddleware() gin.HandlerFunc {
    return gin.Logger()
}
