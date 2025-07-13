package middleware

import (
	"net/http"
	"pt-xyz-multifinance/internal/usecase"
	"pt-xyz-multifinance/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authUseCase usecase.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header required", "Missing Authorization header")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "Invalid authorization format", "Use Bearer token")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := authUseCase.ValidateToken(token)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid token", err.Error())
			c.Abort()
			return
		}

		customerIDFloat, ok := claims["customer_id"].(float64)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims", "Invalid customer_id in token")
			c.Abort()
			return
		}

		nikString, ok := claims["nik"].(string)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims", "Invalid nik in token")
			c.Abort()
			return
		}

		customerID := uint64(customerIDFloat)
		c.Set("customer_id", customerID)
		c.Set("nik", nikString)

		c.Next()
	}
}
