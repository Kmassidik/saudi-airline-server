package routes

import (
	"api-server/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Welcome to API
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to API Saudi Airlines",
		})
	})

	// BranchOffice routes
	branchOfficeRoutes := r.Group("/branch_offices")
	{
		branchOfficeRoutes.GET("", controllers.GetBranchOfficesHandler)
		branchOfficeRoutes.GET("/option-list", controllers.GetBranchOfficesOptionHandler)
		branchOfficeRoutes.GET("/:id", controllers.GetBranchOfficeHandler)
		branchOfficeRoutes.POST("", controllers.CreateBranchOfficeHandler)
		branchOfficeRoutes.PUT("/:id", controllers.UpdateBranchOfficeHandler)
		branchOfficeRoutes.DELETE("/:id", controllers.DeleteBranchOfficeHandler)
	}

	// User routes
	userRoutes := r.Group("/users")
	{
		userRoutes.GET("", controllers.GetUsersHandler)
		userRoutes.GET("/:id", controllers.GetUserHandler)
		userRoutes.POST("", controllers.CreateUserHandler)
		userRoutes.PUT("/:id", controllers.UpdateUserHandler)
		userRoutes.DELETE("/:id", controllers.DeleteUserHandler)
		userRoutes.GET("/branch-office/:id", controllers.GetUsersByBranchOffice)
	}

	// BranchCounter routes
	branchCounterRoutes := r.Group("/branch_counters")
	{
		branchCounterRoutes.GET("/:branch_id", controllers.GetBranchCounterHandlerByBranchId)
		branchCounterRoutes.POST("", controllers.CreateBranchCounterHandler)
		branchCounterRoutes.DELETE("/:id", controllers.DeleteBranchCounterHandler)
	}

	// CompanyProfile routes
	companyProfileRoutes := r.Group("/company_profiles")
	{
		companyProfileRoutes.GET("", controllers.GetCompanyProfileHandler)
		companyProfileRoutes.PUT("", controllers.UpdateCompanyProfileHandler)
	}

	// Vote User routes
	votedUserRoutes := r.Group("/voted-user")
	{
		votedUserRoutes.POST("/:userId", controllers.VotedUserHandler)
	}

	// Dashboard
	dashboardRoutes := r.Group("/dashboard")
	{
		dashboardRoutes.GET("/total-data", controllers.TotalDataDashboard)
		dashboardRoutes.GET("/graph-data/:branchOfficeId", controllers.TotalDataWithBranchIdHandler)
		dashboardRoutes.GET("/total-vote-office", controllers.TotalLikeDislikeBranchOfficeHandler)
		dashboardRoutes.GET("/vote-data-officer", controllers.TotalLikeDislikeOfficerHandler)
		dashboardRoutes.PATCH("/update/:branchId", controllers.UpdateDataDashboardHandler)
	}

	// Authentication
	r.POST("/login", controllers.LoginWebServerHandler)
	r.POST("/login-mobile", controllers.LoginMobileHandler)
}
