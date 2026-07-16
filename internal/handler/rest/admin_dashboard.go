package rest

import (
	"net/http"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *Rest) GetAdminDashboardHome(c *gin.Context) {
	userValue, exists := c.Get("user")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "failed to get user login", nil)
		return
	}

	user, ok := userValue.(*entity.User)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "invalid user login", nil)
		return
	}

	result, err := r.service.AdminDashboardService.GetAdminDashboard(user)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get admin dashboard home", result)
}
