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

	// User routes
	userRoutes := r.Group("/users")
	{
		userRoutes.GET("", controllers.GetUsersHandler)
		userRoutes.GET("/:id", controllers.GetUserHandler)
		userRoutes.POST("", controllers.CreateUserHandler)
		userRoutes.PUT("/:id", controllers.UpdateUserHandler)
		userRoutes.DELETE("/:id", controllers.DeleteUserHandler)
	}

	// BranchOffice routes
	branchOfficeRoutes := r.Group("/branch_offices")
	{
		branchOfficeRoutes.GET("", controllers.GetBranchOfficesHandler)
		branchOfficeRoutes.GET("/:id", controllers.GetBranchOfficeHandler)
		branchOfficeRoutes.POST("", controllers.CreateBranchOfficeHandler)
		branchOfficeRoutes.PUT("/:id", controllers.UpdateBranchOfficeHandler)
		branchOfficeRoutes.DELETE("/:id", controllers.DeleteBranchOfficeHandler)
	}

	// BranchCounter routes
	branchCounterRoutes := r.Group("/branch_counters")
	{
		branchCounterRoutes.GET("", controllers.GetBranchCountersHandler)
		branchCounterRoutes.GET("/:id", controllers.GetBranchCounterHandler)
		branchCounterRoutes.POST("", controllers.CreateBranchCounterHandler)
		branchCounterRoutes.PUT("/:id", controllers.UpdateBranchCounterHandler)
		branchCounterRoutes.DELETE("/:id", controllers.DeleteBranchCounterHandler)
	}

	// CompanyProfile routes
	companyProfileRoutes := r.Group("/company_profiles")
	{
		companyProfileRoutes.GET("", controllers.GetCompanyProfileHandler)
		companyProfileRoutes.PUT("/:id", controllers.UpdateCompanyProfileHandler)
	}
}
