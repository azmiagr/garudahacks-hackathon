package rest

import (
	"net/http"

	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/helper"
	"github.com/azmiagr/garudahacks-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) GetDonorDonationTransactions(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	var req model.DonorDonationTransactionParam
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind query", err)
		return
	}

	result, err := r.service.DonorTransactionService.GetTransactions(user, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get donor donation transactions", result)
}

func (r *Rest) GetDonorDonationTransactionDetail(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	donationID, err := uuid.Parse(c.Param("donation_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid donation_id", err)
		return
	}

	result, err := r.service.DonorTransactionService.GetTransactionDetail(user, donationID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get donor donation transaction detail", result)
}
