package rest

import (
	"net/http"

	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) GetDonorDashboardMap(c *gin.Context) {
	var req model.DonorDashboardMapParam
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind query", err)
		return
	}

	result, err := r.service.DonorDashboardService.GetDonorMap(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get donor dashboard map", result)
}

func (r *Rest) GetDonorPostDetail(c *gin.Context) {
	postID := c.Param("post_id")
	if postID == "" {
		response.Error(c, http.StatusBadRequest, "failed to bind uri", nil)
		return
	}

	result, err := r.service.DonorDashboardService.GetPostDetail(uuid.MustParse(postID))
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get donor post detail", result)
}
