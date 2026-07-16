package rest

import (
	"net/http"

	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *Rest) Login(c *gin.Context) {
	var req model.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind request body", err)
		return
	}

	result, err := r.service.AuthService.Login(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to login", result)
}

func (r *Rest) RequestAdminRegisterOtp(c *gin.Context) {
	var req model.RequestAdminRegisterOtpRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind request body", err)
		return
	}

	result, err := r.service.AuthService.RequestAdminRegisterOtp(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "success to request register otp", result)
}

func (r *Rest) VerifyAdminRegisterOtp(c *gin.Context) {
	var req model.VerifyAdminRegisterOtpRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind request body", err)
		return
	}

	result, err := r.service.AuthService.VerifyAdminRegisterOtp(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to verify register otp", result)
}

func (r *Rest) SetAdminRegisterPassword(c *gin.Context) {
	var req model.SetAdminRegisterPasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind request body", err)
		return
	}

	result, err := r.service.AuthService.SetAdminRegisterPassword(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to create register password", result)
}

func (r *Rest) CompleteAdminRegister(c *gin.Context) {
	var req model.CompleteAdminRegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind request body", err)
		return
	}

	result, err := r.service.AuthService.CompleteAdminRegister(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "success to complete admin registration", result)
}
