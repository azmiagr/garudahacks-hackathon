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
	auth.POST("/logout", r.middleware.AuthenticateUser, r.Logout)
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

	courierRegister := auth.Group("/register/courier")
	courierRegister.POST("/profile", r.CompleteCourierRegister)

	storeRegister := auth.Group("/register/store")
	storeRegister.POST("/profile", r.CompleteStoreRegister)

	admin := baseURL.Group("/admin")
	admin.Use(r.middleware.AuthenticateUser, r.middleware.OnlyAdmin())
	admin.GET("/dashboard", r.GetAdminDashboardHome)
	admin.GET("/profile", r.GetAdminProfile)
	admin.POST("/events", r.CreateAdminEvent)
	admin.POST("/custody/post-handoff", r.SubmitAdminPostHandoff)

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

	store := baseURL.Group("/store")
	store.Use(r.middleware.AuthenticateUser, r.middleware.OnlyStore())
	store.GET("/profile", r.GetStoreProfile)
	store.GET("/orders", r.GetStoreOrders)
	store.GET("/orders/:order_id", r.GetStoreOrderDetail)
	store.POST("/orders/:order_id/accept", r.AcceptStoreOrder)
	store.POST("/orders/:order_id/ready", r.MarkStoreOrderReady)
	store.POST("/orders/:order_id/handoff-token", r.GenerateStoreHandoffToken)
	store.GET("/disbursements/dashboard", r.GetStoreDisbursementDashboard)
	store.GET("/goodness", r.GetStoreGoodness)

	courier := baseURL.Group("/courier")
	courier.Use(r.middleware.AuthenticateUser, r.middleware.OnlyCourier())
	courier.GET("/tasks", r.GetCourierTasks)
	courier.GET("/tasks/:order_id", r.GetCourierTaskDetail)
	courier.POST("/tasks/:order_id/claim", r.ClaimCourierTask)
	courier.POST("/tasks/:order_id/location", r.UpdateCourierLocation)
	courier.POST("/tasks/:order_id/arrived", r.MarkCourierArrived)
	courier.POST("/tasks/:order_id/arrived-post", r.MarkCourierArrivedAtPost)
	courier.POST("/tasks/:order_id/handoff-token", r.GenerateCourierHandoffToken)
	courier.POST("/custody/store-handoff", r.SubmitCourierStoreHandoff)
	courier.GET("/goodness", r.GetCourierGoodness)

}

func (r *Rest) Run() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")

	r.router.Run(fmt.Sprintf("%s:%s", addr, port))
}
