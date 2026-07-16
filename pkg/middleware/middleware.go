package middleware

import (
	"github.com/azmiagr/garudahacks-hackathon/internal/service"
	"github.com/azmiagr/garudahacks-hackathon/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type Interface interface {
	Cors() gin.HandlerFunc
	AuthenticateUser(c *gin.Context)
	OnlyRoles(allowedRoles ...string) gin.HandlerFunc
	OnlyAdmin() gin.HandlerFunc
	OnlyDonor() gin.HandlerFunc
	OnlyStore() gin.HandlerFunc
	OnlyCourier() gin.HandlerFunc
}

type middleware struct {
	service *service.Service
	jwtAuth jwt.Interface
}

func Init(service *service.Service, jwtAuth jwt.Interface) Interface {
	return &middleware{
		service: service,
		jwtAuth: jwtAuth,
	}
}
