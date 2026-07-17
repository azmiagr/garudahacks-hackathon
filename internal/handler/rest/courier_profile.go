package rest

import (
	"net/http"

	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/helper"
	"github.com/azmiagr/garudahacks-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *Rest) GetCourierProfile(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	result, err := r.service.CourierProfileService.GetProfile(user)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get courier profile", result)
}

func (r *Rest) UpdateCourierProfilePreferences(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	var req model.UpdateCourierProfilePreferencesRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind request body", err)
		return
	}

	result, err := r.service.CourierProfileService.UpdatePreferences(user, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to update courier profile preferences", result)
}
