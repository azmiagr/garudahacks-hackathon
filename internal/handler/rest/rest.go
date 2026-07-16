package rest

import (
	"fmt"
	"os"

	"github.com/azmiagr/garudahacks-hackathon/internal/service"
	"github.com/azmiagr/garudahacks-hackathon/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	router     *gin.Engine
	service    *service.Service
	middleware middleware.Interface
}

func NewRest(service *service.Service, middleware middleware.Interface) *Rest {
	return &Rest{
		router:     gin.Default(),
		service:    service,
		middleware: middleware,
	}
}

func (r *Rest) MountEndpoint() {
	r.router.Use(r.middleware.Cors())

	baseURL := r.router.Group("/api/v1")
	baseURL.GET("/dashboard/summary", r.GetPublicDashboardSummary)
	baseURL.GET("/dashboard/map", r.GetPublicDashboardMap)
	baseURL.GET("/dashboard/distributions", r.GetPublicDashboardDistributions)
	baseURL.GET("/dashboard/transparency", r.GetPublicDashboardTransparency)

	auth := baseURL.Group("/auth")
	auth.POST("/login", r.Login)
	adminRegister := auth.Group("/register/admin")
	adminRegister.POST("/request-otp", r.RequestAdminRegisterOtp)
	adminRegister.POST("/verify-otp", r.VerifyAdminRegisterOtp)
	adminRegister.POST("/password", r.SetAdminRegisterPassword)
	adminRegister.POST("/profile", r.CompleteAdminRegister)

	admin := baseURL.Group("/admin")
	admin.Use(r.middleware.AuthenticateUser, r.middleware.OnlyAdmin())
	admin.GET("/dashboard", r.GetAdminDashboardHome)
	admin.POST("/events", r.CreateAdminEvent)

}

func (r *Rest) Run() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")

	r.router.Run(fmt.Sprintf("%s:%s", addr, port))
}
