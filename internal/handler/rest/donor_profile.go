package rest

import (
	"net/http"

	"github.com/azmiagr/garudahacks-hackathon/pkg/helper"
	"github.com/azmiagr/garudahacks-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *Rest) GetDonorProfile(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	result, err := r.service.DonorProfileService.GetProfile(user)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get donor profile", result)
}
