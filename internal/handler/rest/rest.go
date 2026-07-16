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

	payments := baseURL.Group("/payments")
	payments.POST("/webhook", r.HandleMidtransNotification)

	auth := baseURL.Group("/auth")
	auth.POST("/login", r.Login)
	authRegister := auth.Group("/register")
	authRegister.POST("/request-otp", r.RequestRegisterOtp)
	authRegister.POST("/verify-otp", r.VerifyAdminRegisterOtp)
	authRegister.POST("/password", r.SetAdminRegisterPassword)

	adminRegister := auth.Group("/register/admin")
	adminRegister.POST("/request-otp", r.RequestAdminRegisterOtp)
	adminRegister.POST("/verify-otp", r.VerifyAdminRegisterOtp)
	adminRegister.POST("/password", r.SetAdminRegisterPassword)
	adminRegister.POST("/profile", r.CompleteAdminRegister)

	donorRegister := auth.Group("/register/donor")
	donorRegister.POST("/profile", r.CompleteDonorRegister)

	storeRegister := auth.Group("/register/store")
	storeRegister.POST("/profile", r.CompleteStoreRegister)

	admin := baseURL.Group("/admin")
	admin.Use(r.middleware.AuthenticateUser, r.middleware.OnlyAdmin())
	admin.GET("/dashboard", r.GetAdminDashboardHome)
	admin.GET("/profile", r.GetAdminProfile)
	admin.POST("/events", r.CreateAdminEvent)

	donor := baseURL.Group("/donor")
	donor.Use(r.middleware.AuthenticateUser, r.middleware.OnlyDonor())
	donor.GET("/profile", r.GetDonorProfile)
	donor.GET("/dashboard/map", r.GetDonorDashboardMap)
	donor.GET("/dashboard/posts/:post_id", r.GetDonorPostDetail)
	donor.GET("/donations/transactions", r.GetDonorDonationTransactions)
	donor.GET("/donations/transactions/:donation_id", r.GetDonorDonationTransactionDetail)
	donor.POST("/donations/payments", r.CreateDonationPayment)
	donor.GET("/points", r.GetPointDashboard)
	donor.GET("/points/history", r.GetPointHistory)
	donor.GET("/points/rewards", r.GetRewards)
	donor.POST("/points/rewards/claim", r.ClaimReward)

}

func (r *Rest) Run() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")

	r.router.Run(fmt.Sprintf("%s:%s", addr, port))
}
