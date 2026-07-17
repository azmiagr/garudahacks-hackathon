package rest

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (r *Rest) Logout(c *gin.Context) {
	tokenValue, exists := c.Get("token")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "token is required", nil)
		return
	}

	token, ok := tokenValue.(string)
	if !ok || token == "" {
		response.Error(c, http.StatusUnauthorized, "token is required", nil)
		return
	}

	result, err := r.service.AuthService.Logout(token)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to logout", result)
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

func (r *Rest) RequestRegisterOtp(c *gin.Context) {
	var req model.RequestAdminRegisterOtpRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind request body", err)
		return
	}

	result, err := r.service.AuthService.RequestRegisterOtp(req)
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

func (r *Rest) CompleteDonorRegister(c *gin.Context) {
	var req model.CompleteDonorRegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind request body", err)
		return
	}

	result, err := r.service.AuthService.CompleteDonorRegister(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "success to complete donor registration", result)
}

func (r *Rest) CompleteStoreRegister(c *gin.Context) {
	req, err := bindCompleteStoreRegisterForm(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind request form", err)
		return
	}

	result, err := r.service.AuthService.CompleteStoreRegister(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "success to complete store registration", result)
}

func bindCompleteStoreRegisterForm(c *gin.Context) (model.CompleteStoreRegisterRequest, error) {
	registrationID, err := uuid.Parse(strings.TrimSpace(c.PostForm("registration_id")))
	if err != nil {
		return model.CompleteStoreRegisterRequest{}, err
	}

	latitude, err := strconv.ParseFloat(strings.TrimSpace(c.PostForm("latitude")), 64)
	if err != nil {
		return model.CompleteStoreRegisterRequest{}, err
	}

	longitude, err := strconv.ParseFloat(strings.TrimSpace(c.PostForm("longitude")), 64)
	if err != nil {
		return model.CompleteStoreRegisterRequest{}, err
	}

	ktpImage, _ := c.FormFile("ktp_image")

	categories := append([]string{}, c.PostFormArray("categories")...)
	categories = append(categories, c.PostFormArray("categories[]")...)

	return model.CompleteStoreRegisterRequest{
		RegistrationID:  registrationID,
		StoreName:       c.PostForm("store_name"),
		OwnerName:       c.PostForm("owner_name"),
		NIB:             c.PostForm("nib"),
		NPWP:            c.PostForm("npwp"),
		KTPImage:        ktpImage,
		BankName:        c.PostForm("bank_name"),
		BankAccountNo:   c.PostForm("bank_account_no"),
		BankAccountName: c.PostForm("bank_account_name"),
		Categories:      categories,
		CategoriesJSON:  c.PostForm("categories_json"),
		Address:         c.PostForm("address"),
		Latitude:        latitude,
		Longitude:       longitude,
	}, nil
}
