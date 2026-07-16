package rest

import (
	"net/http"
	"strconv"

	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/helper"
	"github.com/azmiagr/garudahacks-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) GetStoreOrders(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	result, err := r.service.StoreCustodyService.GetStoreOrders(user, bindStoreOrderListParam(c))
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get store orders", result)
}

func (r *Rest) GetStoreOrderDetail(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	orderID, ok := bindOrderIDParam(c)
	if !ok {
		return
	}

	result, err := r.service.StoreCustodyService.GetStoreOrderDetail(user, orderID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get store order detail", result)
}

func (r *Rest) AcceptStoreOrder(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	orderID, ok := bindOrderIDParam(c)
	if !ok {
		return
	}

	result, err := r.service.StoreCustodyService.AcceptOrder(user, model.StoreOrderActionRequest{
		OrderID: orderID,
	})
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to accept order", result)
}

func (r *Rest) MarkStoreOrderReady(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	orderID, ok := bindOrderIDParam(c)
	if !ok {
		return
	}

	result, err := r.service.StoreCustodyService.MarkOrderReady(user, model.StoreOrderActionRequest{
		OrderID: orderID,
	})
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to mark order ready", result)
}

func (r *Rest) GenerateStoreHandoffToken(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	orderID, ok := bindOrderIDParam(c)
	if !ok {
		return
	}

	result, err := r.service.StoreCustodyService.GenerateStoreHandoffToken(user, model.StoreOrderActionRequest{
		OrderID: orderID,
	})
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to generate store handoff token", result)
}

func (r *Rest) SubmitCourierStoreHandoff(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	var req model.SubmitCustodyHandshakeRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind request body", err)
		return
	}

	result, err := r.service.StoreCustodyService.SubmitHandoff(user, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to submit store handoff", result)
}

func bindOrderIDParam(c *gin.Context) (uuid.UUID, bool) {
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order_id", err)
		return uuid.Nil, false
	}

	return orderID, true
}

func bindStoreOrderListParam(c *gin.Context) model.StoreOrderListParam {
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	return model.StoreOrderListParam{
		Status: c.Query("status"),
		Limit:  limit,
		Offset: offset,
	}
}
