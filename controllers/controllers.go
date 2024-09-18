package controllers

import (
	"api-server/models"
	"api-server/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// User Handlers

// GetUsersHandler handles the request to get a list of users with pagination
func GetUsersHandler(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset := (page - 1) * limit

	users, err := repository.GetAllUsers(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	totalCount, err := repository.GetUsersCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users count",
		})
		return
	}

	totalPages := (totalCount + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
		"total_count": totalCount,
		"users":       users,
	})
}

// GetUserHandler retrieves a single user by ID
func GetUserHandler(c *gin.Context) {
	id := c.Param("id")

	// Convert id to integer if necessary
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := repository.GetUserByID(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// CreateUserHandler creates a new user
func CreateUserHandler(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := repository.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
}

// UpdateUserHandler updates an existing user by ID
func UpdateUserHandler(c *gin.Context) {
	id := c.Param("id")

	// Convert id to integer if necessary
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err = repository.UpdateUser(uint(userID), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUserHandler deletes a user by ID
func DeleteUserHandler(c *gin.Context) {
	id := c.Param("id")

	// Convert id to integer if necessary
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = repository.DeleteUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// BranchOffice Handlers

func GetBranchOfficesHandler(c *gin.Context) {
	// Implement logic to retrieve branch offices
	c.JSON(http.StatusOK, gin.H{"message": "Get all branch offices"})
}

func GetBranchOfficeHandler(c *gin.Context) {
	id := c.Param("id")
	// Implement logic to retrieve a branch office by ID
	c.JSON(http.StatusOK, gin.H{"message": "Get branch office", "id": id})
}

func CreateBranchOfficeHandler(c *gin.Context) {
	// Implement logic to create a new branch office
	c.JSON(http.StatusOK, gin.H{"message": "Create branch office"})
}

func UpdateBranchOfficeHandler(c *gin.Context) {
	id := c.Param("id")
	// Implement logic to update a branch office by ID
	c.JSON(http.StatusOK, gin.H{"message": "Update branch office", "id": id})
}

func DeleteBranchOfficeHandler(c *gin.Context) {
	id := c.Param("id")
	// Implement logic to delete a branch office by ID
	c.JSON(http.StatusOK, gin.H{"message": "Delete branch office", "id": id})
}

// BranchCounter Handlers

func GetBranchCountersHandler(c *gin.Context) {
	// Implement logic to retrieve branch counters
	c.JSON(http.StatusOK, gin.H{"message": "Get all branch counters"})
}

func GetBranchCounterHandler(c *gin.Context) {
	id := c.Param("id")
	// Implement logic to retrieve a branch counter by ID
	c.JSON(http.StatusOK, gin.H{"message": "Get branch counter", "id": id})
}

func CreateBranchCounterHandler(c *gin.Context) {
	// Implement logic to create a new branch counter
	c.JSON(http.StatusOK, gin.H{"message": "Create branch counter"})
}

func UpdateBranchCounterHandler(c *gin.Context) {
	id := c.Param("id")
	// Implement logic to update a branch counter by ID
	c.JSON(http.StatusOK, gin.H{"message": "Update branch counter", "id": id})
}

func DeleteBranchCounterHandler(c *gin.Context) {
	id := c.Param("id")
	// Implement logic to delete a branch counter by ID
	c.JSON(http.StatusOK, gin.H{"message": "Delete branch counter", "id": id})
}

// CompanyProfile Handlers

func GetCompanyProfilesHandler(c *gin.Context) {
	// Implement logic to retrieve company profiles
	c.JSON(http.StatusOK, gin.H{"message": "Get all company profiles"})
}

func GetCompanyProfileHandler(c *gin.Context) {
	id := c.Param("id")
	// Implement logic to retrieve a company profile by ID
	c.JSON(http.StatusOK, gin.H{"message": "Get company profile", "id": id})
}

func CreateCompanyProfileHandler(c *gin.Context) {
	// Implement logic to create a new company profile
	c.JSON(http.StatusOK, gin.H{"message": "Create company profile"})
}

func UpdateCompanyProfileHandler(c *gin.Context) {
	id := c.Param("id")
	// Implement logic to update a company profile by ID
	c.JSON(http.StatusOK, gin.H{"message": "Update company profile", "id": id})
}

func DeleteCompanyProfileHandler(c *gin.Context) {
	id := c.Param("id")
	// Implement logic to delete a company profile by ID
	c.JSON(http.StatusOK, gin.H{"message": "Delete company profile", "id": id})
}
