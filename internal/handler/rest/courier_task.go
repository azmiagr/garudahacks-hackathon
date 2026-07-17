package rest

import (
	"net/http"
	"strconv"

	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/helper"
	"github.com/azmiagr/garudahacks-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *Rest) GetCourierTasks(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	result, err := r.service.CourierTaskService.GetCourierTasks(user, bindCourierTaskListParam(c))
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get courier tasks", result)
}

func (r *Rest) GetCourierTaskDetail(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	orderID, ok := bindOrderIDParam(c)
	if !ok {
		return
	}

	result, err := r.service.CourierTaskService.GetCourierTaskDetail(user, orderID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get courier task detail", result)
}

func (r *Rest) ClaimCourierTask(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	orderID, ok := bindOrderIDParam(c)
	if !ok {
		return
	}

	var body struct {
		Latitude  *float64 `json:"lat"`
		Longitude *float64 `json:"lng"`
	}
	_ = c.ShouldBindJSON(&body)

	req := model.CourierTaskClaimRequest{}
	if body.Latitude != nil && body.Longitude != nil {
		req.Latitude = *body.Latitude
		req.Longitude = *body.Longitude
		req.HasCoords = true
	}

	result, err := r.service.CourierTaskService.ClaimTask(user, orderID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to claim task", result)
}

func (r *Rest) UpdateCourierLocation(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	orderID, ok := bindOrderIDParam(c)
	if !ok {
		return
	}

	var req model.CourierLocationPingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind request body", err)
		return
	}

	result, err := r.service.CourierTaskService.UpdateLocation(user, orderID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to update courier location", result)
}

func (r *Rest) MarkCourierArrived(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	orderID, ok := bindOrderIDParam(c)
	if !ok {
		return
	}

	result, err := r.service.CourierTaskService.MarkArrived(user, orderID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to mark courier arrived", result)
}

func (r *Rest) MarkCourierArrivedAtPost(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	orderID, ok := bindOrderIDParam(c)
	if !ok {
		return
	}

	result, err := r.service.CourierTaskService.MarkArrivedAtPost(user, orderID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to mark courier arrived at post", result)
}

func (r *Rest) GenerateCourierHandoffToken(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	orderID, ok := bindOrderIDParam(c)
	if !ok {
		return
	}

	result, err := r.service.CourierTaskService.GenerateHandoffToken(user, orderID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to generate courier handoff token", result)
}

func bindCourierTaskListParam(c *gin.Context) model.CourierTaskListParam {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))
	radiusKm, _ := strconv.ParseFloat(c.Query("radius_km"), 64)

	param := model.CourierTaskListParam{
		Status:   c.Query("status"),
		RadiusKm: radiusKm,
		Limit:    limit,
		Offset:   offset,
	}

	lat, latErr := strconv.ParseFloat(c.Query("lat"), 64)
	lng, lngErr := strconv.ParseFloat(c.Query("lng"), 64)
	if latErr == nil && lngErr == nil {
		param.Latitude = lat
		param.Longitude = lng
		param.HasCoords = true
	}

	return param
}
