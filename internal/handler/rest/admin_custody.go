package rest

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/helper"
	"github.com/azmiagr/garudahacks-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) GetAdminReceiveOrderDetail(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order_id", err)
		return
	}

	result, err := r.service.AdminCustodyService.GetReceiveOrderDetail(user, orderID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get receive order detail", result)
}

func (r *Rest) SubmitAdminPostHandoff(c *gin.Context) {
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

	result, err := r.service.AdminCustodyService.SubmitPostHandoff(user, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to submit post handoff", result)
}

func (r *Rest) CreateAdminSupplementalNeed(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order_id", err)
		return
	}

	var req model.CreateSupplementalNeedRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind request body", err)
		return
	}

	result, err := r.service.AdminCustodyService.CreateSupplementalNeed(user, orderID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "success to create supplemental need", result)
}

func (r *Rest) UploadAdminDistributionProof(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	req, err := bindUploadDistributionProofForm(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind request form", err)
		return
	}

	result, err := r.service.AdminCustodyService.UploadDistributionProof(user, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "success to upload distribution proof", result)
}

func (r *Rest) CompleteAdminDistribution(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid order_id", err)
		return
	}

	var req model.CompleteDistributionRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind request body", err)
		return
	}

	result, err := r.service.AdminCustodyService.CompleteDistribution(user, orderID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to complete distribution", result)
}

func bindUploadDistributionProofForm(c *gin.Context) (model.UploadDistributionProofRequest, error) {
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		return model.UploadDistributionProofRequest{}, err
	}

	itemID, err := uuid.Parse(strings.TrimSpace(c.PostForm("item_id")))
	if err != nil {
		return model.UploadDistributionProofRequest{}, err
	}

	latitude, err := strconv.ParseFloat(strings.TrimSpace(c.PostForm("latitude")), 64)
	if err != nil {
		return model.UploadDistributionProofRequest{}, err
	}

	longitude, err := strconv.ParseFloat(strings.TrimSpace(c.PostForm("longitude")), 64)
	if err != nil {
		return model.UploadDistributionProofRequest{}, err
	}

	distributedQuantity := 0
	if strings.TrimSpace(c.PostForm("distributed_quantity")) != "" {
		distributedQuantity, err = strconv.Atoi(strings.TrimSpace(c.PostForm("distributed_quantity")))
		if err != nil {
			return model.UploadDistributionProofRequest{}, err
		}
	}

	blurFaceEnabled := true
	if strings.TrimSpace(c.PostForm("blur_face_enabled")) != "" {
		blurFaceEnabled, err = strconv.ParseBool(strings.TrimSpace(c.PostForm("blur_face_enabled")))
		if err != nil {
			return model.UploadDistributionProofRequest{}, err
		}
	}

	capturedFromCamera := true
	if strings.TrimSpace(c.PostForm("captured_from_camera")) != "" {
		capturedFromCamera, err = strconv.ParseBool(strings.TrimSpace(c.PostForm("captured_from_camera")))
		if err != nil {
			return model.UploadDistributionProofRequest{}, err
		}
	}

	var capturedAt time.Time
	if strings.TrimSpace(c.PostForm("captured_at")) != "" {
		capturedAt, err = time.Parse(time.RFC3339, strings.TrimSpace(c.PostForm("captured_at")))
		if err != nil {
			return model.UploadDistributionProofRequest{}, err
		}
	}

	photo, _ := c.FormFile("photo")

	return model.UploadDistributionProofRequest{
		OrderID:             orderID,
		ItemID:              itemID,
		Photo:               photo,
		RecipientNote:       c.PostForm("recipient_note"),
		DistributedQuantity: distributedQuantity,
		Latitude:            latitude,
		Longitude:           longitude,
		BlurFaceEnabled:     blurFaceEnabled,
		CapturedFromCamera:  capturedFromCamera,
		CapturedAt:          capturedAt,
	}, nil
}
