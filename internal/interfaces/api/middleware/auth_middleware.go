package middleware

import (
	"net/http"
	"pt-xyz-multifinance/internal/domain/entity"
	"pt-xyz-multifinance/internal/usecase"
	"pt-xyz-multifinance/pkg/response"
	"strconv"
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

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims", "Invalid user_id in token")
			c.Abort()
			return
		}

		username, ok := claims["username"].(string)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims", "Invalid username in token")
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			response.Error(c, http.StatusUnauthorized, "Invalid token claims", "Invalid role in token")
			c.Abort()
			return
		}

		userID := uint64(userIDFloat)
		c.Set("user_id", userID)
		c.Set("username", username)
		c.Set("role", role)

		if role == string(entity.RoleCustomer) {
			customer, err := authUseCase.GetCustomerFromUser(c.Request.Context(), userID)
			if err == nil {
				c.Set("customer_id", customer.ID)
			}
		}

		c.Next()
	}
}

func RequireRole(roles ...entity.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "Role not found in context", "Please login first")
			c.Abort()
			return
		}

		currentRole := entity.UserRole(userRole.(string))

		for _, role := range roles {
			if currentRole == role {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "Insufficient permissions", "You don't have permission to access this resource")
		c.Abort()
	}
}

func CustomerOwnershipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "Role not found", "Please login first")
			c.Abort()
			return
		}

		if userRole.(string) == string(entity.RoleAdmin) {
			c.Next()
			return
		}

		if userRole.(string) == string(entity.RoleCustomer) {
			customerID, exists := c.Get("customer_id")
			if !exists {
				response.Error(c, http.StatusForbidden, "Customer ID not found", "Customer data not available")
				c.Abort()
				return
			}

			paramID := c.Param("id")
			if paramID != "" {
				requestedID, err := strconv.ParseUint(paramID, 10, 64)
				if err != nil {
					response.Error(c, http.StatusBadRequest, "Invalid ID parameter", "ID must be a valid number")
					c.Abort()
					return
				}

				if requestedID != customerID.(uint64) {
					response.Error(c, http.StatusForbidden, "Access denied", "You can only access your own data")
					c.Abort()
					return
				}
			}

			paramCustomerID := c.Param("customer_id")
			if paramCustomerID != "" {
				requestedCustomerID, err := strconv.ParseUint(paramCustomerID, 10, 64)
				if err != nil {
					response.Error(c, http.StatusBadRequest, "Invalid customer ID parameter", "Customer ID must be a valid number")
					c.Abort()
					return
				}

				if requestedCustomerID != customerID.(uint64) {
					response.Error(c, http.StatusForbidden, "Access denied", "You can only access your own data")
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return RequireRole(entity.RoleAdmin)
}

func CustomerOnly() gin.HandlerFunc {
	return RequireRole(entity.RoleCustomer)
}
