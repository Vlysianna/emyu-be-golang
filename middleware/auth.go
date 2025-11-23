package middleware

import (
	"net/http"
	"strings"

	"github.com/emyu/ecommer-be/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.ID)
		c.Set("userEmail", claims.Email)
		c.Set("roleID", claims.RoleID)
		c.Set("roleName", claims.RoleName)
		c.Set("permissions", claims.Permissions)
		c.Next()
	}
}

// RoleMiddleware checks if user has specific role
func RoleMiddleware(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleName, exists := c.Get("roleName")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User role not found"})
			c.Abort()
			return
		}

		hasRole := false
		for _, role := range requiredRoles {
			if roleName == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient role permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// PermissionMiddleware checks if user has specific permission
func PermissionMiddleware(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User permissions not found"})
			c.Abort()
			return
		}

		userPermissions := permissions.([]string)
		hasPermission := false

		for _, reqPerm := range requiredPermissions {
			if utils.HasPermission(userPermissions, reqPerm) {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminMiddleware checks if user is admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleName, exists := c.Get("roleName")
		if !exists || roleName != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetUserID(c *gin.Context) string {
	userID, exists := c.Get("userID")
	if !exists {
		return ""
	}
	return userID.(string)
}

func GetRoleID(c *gin.Context) int {
	roleID, exists := c.Get("roleID")
	if !exists {
		return 0
	}
	return roleID.(int)
}

func GetRoleName(c *gin.Context) string {
	roleName, exists := c.Get("roleName")
	if !exists {
		return ""
	}
	return roleName.(string)
}

func GetPermissions(c *gin.Context) []string {
	permissions, exists := c.Get("permissions")
	if !exists {
		return []string{}
	}
	return permissions.([]string)
}
