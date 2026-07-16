package middleware

import (
	"github.com/azmiagr/garudahacks-hackathon/internal/service"
	"github.com/azmiagr/garudahacks-hackathon/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type Interface interface {
	Cors() gin.HandlerFunc
	AuthenticateUser(c *gin.Context)
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
