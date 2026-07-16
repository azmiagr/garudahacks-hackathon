package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/helper"
	"github.com/azmiagr/garudahacks-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *Rest) CreateAdminEvent(c *gin.Context) {
	user, ok := helper.GetLoginUserFromContext(c)
	if !ok {
		return
	}

	req, err := bindCreateAdminEventForm(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind request body", err)
		return
	}

	result, err := r.service.AdminEventService.CreateEvent(user, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "success to create admin event", result)
}

func bindCreateAdminEventForm(c *gin.Context) (model.CreateAdminEventRequest, error) {
	latitude, err := strconv.ParseFloat(c.PostForm("latitude"), 64)
	if err != nil {
		return model.CreateAdminEventRequest{}, err
	}

	longitude, err := strconv.ParseFloat(c.PostForm("longitude"), 64)
	if err != nil {
		return model.CreateAdminEventRequest{}, err
	}

	geofenceRadius, err := strconv.ParseFloat(c.PostForm("geofence_radius"), 64)
	if err != nil {
		return model.CreateAdminEventRequest{}, err
	}

	var items []model.CreateAdminEventItem
	if err := json.Unmarshal([]byte(c.PostForm("items")), &items); err != nil {
		return model.CreateAdminEventRequest{}, err
	}

	photo, err := c.FormFile("photo")
	if err != nil {
		return model.CreateAdminEventRequest{}, err
	}

	return model.CreateAdminEventRequest{
		Name:           c.PostForm("name"),
		Description:    c.PostForm("description"),
		DisasterType:   c.PostForm("disaster_type"),
		Address:        c.PostForm("address"),
		Latitude:       latitude,
		Longitude:      longitude,
		GeofenceRadius: geofenceRadius,
		Photo:          photo,
		Items:          items,
	}, nil
}
