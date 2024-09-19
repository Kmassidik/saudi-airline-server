package controllers

import (
	"api-server/helpers"
	"api-server/models"
	"api-server/repository"
	"api-server/repository/validation"
	"api-server/services"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	users, err := services.GetAllUsers(limit, offset)
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

	user, err := services.GetUserByID(uint(userID))
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

// CreateUserHandler creates a new user and handles image upload
func CreateUserHandler(c *gin.Context) {
	// Parse the form data, including files
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	// Extract form data
	fullName := c.Request.FormValue("full_name")
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	role := c.Request.FormValue("role")

	// Create a user struct with the extracted form data
	user := models.User{
		FullName: fullName,
		Email:    email,
		Password: password,
		Role:     role,
	}

	// Handle file upload if present
	fileHeader, err := c.FormFile("image")
	if err != nil {
		if err != http.ErrMissingFile {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
			return
		}
		// No file uploaded, continue with user creation without image
	} else {
		// Generate a unique filename and save the file
		fileName := helpers.GenerateFileName()
		filePath := filepath.Join("public/images", fileName)
		file, err := fileHeader.Open()
		if err != nil {
			log.Println("Error opening file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer file.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			log.Println("Error creating file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			log.Println("Error saving file:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Update user image path
		user.Image = fileName
	}

	// Validate user details before processing
	if err := validation.ValidateUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service layer to handle user creation
	if err := services.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
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

	err = services.UpdateUser(uint(userID), &user)
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

	err = services.DeleteUser(uint(userID))
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
