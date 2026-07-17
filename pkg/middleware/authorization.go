package middleware

import (
	"net/http"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
)

func (m *middleware) OnlyRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userValue, exists := c.Get("user")
		if !exists {
			response.Error(c, http.StatusForbidden, "access denied", nil)
			c.Abort()
			return
		}

		user, ok := userValue.(*entity.User)
		if !ok {
			response.Error(c, http.StatusForbidden, "access denied", nil)
			c.Abort()
			return
		}

		userRoleName, err := m.service.UserService.GetUserRoleName(user)
		if err != nil {
			response.Error(c, http.StatusForbidden, "access denied", err)
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if userRoleName == role {
				c.Next()
				return
			}
		}

		response.Error(c, http.StatusForbidden, "access denied", nil)
		c.Abort()
	}
}

func (m *middleware) OnlyAdmin() gin.HandlerFunc {
	return m.OnlyRoles("admin")
}

func (m *middleware) OnlyDonor() gin.HandlerFunc {
	return m.OnlyRoles("donor")
}

func (m *middleware) OnlyStore() gin.HandlerFunc {
	return m.OnlyRoles("store")
}

func (m *middleware) OnlyCourier() gin.HandlerFunc {
	return m.OnlyRoles("relawan", "courier", "kurir")
}
