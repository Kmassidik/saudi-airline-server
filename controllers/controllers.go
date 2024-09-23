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

	// Use the service layer to get the users
	users, err := services.GetAllUsers(limit, offset)
	if err != nil {
		c.Error(err) // Pass error to the middleware
		return
	}

	// Fetch the total user count
	totalCount, err := repository.GetUsersCount()
	if err != nil {
		c.Error(err) // Pass error to the middleware
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

func GetUserHandler(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := services.GetUserByID(uint(userID))
	if err != nil {
		c.Error(err) // Pass error to the middleware
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func CreateUserHandler(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	fullName := c.Request.FormValue("full_name")
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	role := c.Request.FormValue("role")

	user := models.User{
		FullName: fullName,
		Email:    email,
		Password: password,
		Role:     role,
	}

	// Handle file upload
	fileHeader, err := c.FormFile("image")
	if err == nil {
		// If there's a file, handle the file upload
		fileName := helpers.GenerateFileName()
		filePath := filepath.Join("public/images", fileName)
		file, err := fileHeader.Open()
		if err != nil {
			c.Error(err) // Pass error to the middleware
			return
		}
		defer file.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			c.Error(err) // Pass error to the middleware
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			c.Error(err) // Pass error to the middleware
			return
		}

		user.Image = fileName
	}

	// Validate user
	if err := validation.ValidateUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to create user
	if err := services.CreateUser(&user); err != nil {
		c.Error(err) // Pass error to the middleware
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func UpdateUserHandler(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse the form for multipart data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	// Retrieve form data
	fullName := c.Request.FormValue("full_name")
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")
	role := c.Request.FormValue("role")

	// Find existing user
	user, err := services.GetUserByID(uint(userID))
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Store the current image path
	oldImagePath := user.Image

	// Update fields
	user.FullName = fullName
	user.Email = email
	if password != "" {
		user.Password = password // Only update if password is provided
	}
	user.Role = role

	// Handle file upload
	fileHeader, err := c.FormFile("image")
	if err == nil {
		// If there's a file, handle the file upload
		fileName := helpers.GenerateFileName() // Ensure this generates just the filename
		filePath := filepath.Join("public/images", fileName)
		file, err := fileHeader.Open()
		if err != nil {
			c.Error(err)
			return
		}
		defer file.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			c.Error(err)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			c.Error(err)
			return
		}

		// Update the user's image field
		user.Image = fileName // Store just the filename
	}

	// Validate user
	if err := validation.ValidateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to update user
	if err := services.UpdateUser(uint(userID), user); err != nil {
		c.Error(err)
		return
	}

	// Delete the old image if it exists and a new image was uploaded
	if user.Image != oldImagePath && oldImagePath != "" {
		oldImageFullPath := filepath.Join("public/images", oldImagePath) // Use the correct path
		if err := os.Remove(oldImageFullPath); err != nil {
			log.Println("Error deleting old image:", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUserHandler(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Find existing user
	user, err := services.GetUserByID(uint(userID))
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Store the image path for deletion
	imagePath := user.Image

	// Delete the user from the database
	if err := services.DeleteUser(uint(user.ID)); err != nil {
		c.Error(err) // Pass error to the middleware
		return
	}

	// Delete the user's image file from the filesystem
	if imagePath != "" {
		oldImageFullPath := filepath.Join("public/images", imagePath) // Adjust the path as needed
		if err := os.Remove(oldImageFullPath); err != nil {
			log.Println("Error deleting user image:", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// CompanyProfile Handlers
func GetCompanyProfileHandler(c *gin.Context) {
	// Get the company profile
	company, err := services.GetCompanyProfile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get company profile"})
		return
	}

	if company == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company profile not found"})
		return
	}

	c.JSON(http.StatusOK, company)
}

func UpdateCompanyProfileHandler(c *gin.Context) {
	// Get the current company profile
	company, err := services.GetCompanyProfile() // Fetch company profile with ID 1
	if err != nil {
		c.Error(err)
		return
	}

	// Store the current image path
	oldImagePath := company.Logo

	// Parse the form for multipart data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form"})
		return
	}

	// Retrieve form data
	name := c.Request.FormValue("name")
	fileHeader, err := c.FormFile("logo")
	var logo string

	if err == nil {
		// Handle the file upload if a new file is provided
		fileName := filepath.Base(fileHeader.Filename)
		filePath := filepath.Join("public/assets", fileName)
		file, err := fileHeader.Open()
		if err != nil {
			c.Error(err)
			return
		}
		defer file.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			c.Error(err)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			c.Error(err)
			return
		}

		logo = fileName // Store the filename for the logo
	} else {
		// No new image uploaded, use the old image path
		logo = oldImagePath
	}

	// Delete the old image if a new one is uploaded and it's different
	if oldImagePath != "" && logo != "" && oldImagePath != logo {
		oldImageFullPath := filepath.Join("public/assets", oldImagePath) // Use the correct path
		if err := os.Remove(oldImageFullPath); err != nil {
			log.Println("Error deleting old image:", err)
		}
	}

	// Call service to update company profile with ID 1
	if err := services.UpdateCompanyProfile(name, logo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company profile updated successfully"})
}

// BranchOffice Handlers

// GetBranchOfficesHandler retrieves all branch offices with pagination
func GetBranchOfficesHandler(c *gin.Context) {
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

	// Use the service layer to get the branch offices
	branchOffices, err := services.GetAllBranchOffices(limit, offset)
	if err != nil {
		c.Error(err) // Pass error to the middleware
		return
	}

	// Fetch the total branch office count
	totalCount, err := repository.GetBranchOfficesCount()
	if err != nil {
		c.Error(err) // Pass error to the middleware
		return
	}

	totalPages := (totalCount + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"page":           page,
		"limit":          limit,
		"total_pages":    totalPages,
		"total_count":    totalCount,
		"branch_offices": branchOffices,
	})
}

// GetBranchOfficeHandler retrieves a single branch office by ID
func GetBranchOfficeHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch office ID"})
		return
	}

	branchOffice, err := services.GetBranchOfficeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Branch office not found"})
		return
	}

	if branchOffice == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Branch office not found"})
		return
	}

	c.JSON(http.StatusOK, branchOffice)
}

// CreateBranchOfficeHandler creates a new branch office
func CreateBranchOfficeHandler(c *gin.Context) {
	var branchOffice models.BranchOffice

	if err := c.BindJSON(&branchOffice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate branch offices
	if err := validation.ValidateBranchOffices(&branchOffice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to create user
	if err := services.CreateBranchOffice(&branchOffice); err != nil {
		c.Error(err) // Pass error to the middleware
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Branch office created successfully"})
}

// UpdateBranchOfficeHandler updates an existing branch office by ID
func UpdateBranchOfficeHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch office ID"})
		return
	}

	var branchOffice models.BranchOffice
	if err := c.ShouldBindJSON(&branchOffice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := services.UpdateBranchOffice(uint(id), &branchOffice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update branch office"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch office updated successfully"})
}

// DeleteBranchOfficeHandler deletes a branch office by ID
func DeleteBranchOfficeHandler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch office ID"})
		return
	}

	if err := services.DeleteBranchOffice(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete branch office"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch office deleted successfully"})
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
