package rest

import (
	"net/http"

	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *Rest) GetPublicDashboardMap(c *gin.Context) {
	var req model.PublicDashboardParam
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind query", err)
		return
	}

	result, err := r.service.PublicDashboardService.GetPublicMap(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get public dashboard map", result)

}

func (r *Rest) GetPublicDashboardSummary(c *gin.Context) {
	var req model.PublicDashboardParam
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind query", err)
		return
	}

	result, err := r.service.PublicDashboardService.GetSummary(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get public dashboard summary", result)

}

func (r *Rest) GetPublicDashboardDistributions(c *gin.Context) {
	var req model.PublicDistributionParam
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind query", err)
		return
	}

	result, err := r.service.PublicDashboardService.GetDistributions(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get public dashboard distributions", result)
}

func (r *Rest) GetPublicDashboardTransparency(c *gin.Context) {
	var req model.PublicTransparencyParam
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind query", err)
		return
	}

	result, err := r.service.PublicDashboardService.GetTransparency(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get public dashboard transparency", result)
}
