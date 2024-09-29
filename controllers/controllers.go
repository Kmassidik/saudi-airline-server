package controllers

import (
	"api-server/helpers"
	"api-server/models"
	"api-server/repository"
	"api-server/repository/validation"
	"api-server/services"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BranchOffice Handlers

func GetBranchOfficesHandler(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "5")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
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

func GetBranchOfficesOptionHandler(c *gin.Context) {

	// Use the service layer to get the branch offices
	branchOffices, err := services.GetAllBranchOfficesOptionList()
	if err != nil {
		c.Error(err) // Pass error to the middleware
		return
	}

	c.JSON(http.StatusOK, branchOffices)
}

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

func CreateBranchOfficeHandler(c *gin.Context) {
	var input map[string]interface{}

	// Parse request body
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate input using the validation function
	if err := validation.ValidateBranchOffices(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create BranchOffice model from validated input
	branchOffice := models.BranchOfficeCreateRequest{
		Name:         input["name"].(string),
		Address:      input["address"].(string),
		TotalCounter: uint(input["total_counter"].(float64)), // Assuming total_counter comes as float64 from JSON
	}

	// Call service to create branch office
	if err := services.CreateBranchOffice(&branchOffice); err != nil {
		c.Error(err) // Pass error to the middleware
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Branch office created successfully"})
}

func UpdateBranchOfficeHandler(c *gin.Context) {
	var input map[string]interface{}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch office ID"})
		return
	}

	// Parse request body
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate input using the validation function
	if err := validation.ValidateBranchOffices(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	branchOffice := models.BranchOfficeCreateRequest{
		Name:         input["name"].(string),
		Address:      input["address"].(string),
		TotalCounter: uint(input["total_counter"].(float64)), // Assuming total_counter comes as float64 from JSON
	}

	if err := services.UpdateBranchOffice(uint(id), &branchOffice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update branch office"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch office updated successfully"})
}

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

// User Handlers

func GetUsersHandler(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	role := c.DefaultQuery("role", "officier")
	limitStr := c.DefaultQuery("limit", "5")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	offset := (page - 1) * limit

	// Use the service layer to get the users
	users, err := services.GetAllUsers(limit, offset, role)
	if err != nil {
		c.Error(err) // Pass error to the middleware
		return
	}

	// Fetch the total user count
	totalCount, err := services.GetUsersCount(role)
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
		"role":        role,
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
	branchIdStr := c.Request.FormValue("branch_id")

	branchId, err := strconv.Atoi(branchIdStr)
	if err != nil {
		log.Println("Error converting branch_id to uint:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch_id"})
		return
	}

	user := models.User{
		FullName: fullName,
		Email:    email,
		Password: password,
		Role:     role,
		BranchId: uint(branchId),
	}

	// Handle file upload but don't save it yet
	var fileHeader *multipart.FileHeader
	fileHeader, err = c.FormFile("image")
	if err == nil {
		// If there's a file, set the filename but don't save yet
		fileName := helpers.GenerateFileName()
		user.Image = fileName
	}

	// Validate user
	if err := validation.ValidateUser(&user); err != nil {
		log.Println("Error Validation:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service to create user
	if err := services.CreateUser(&user); err != nil {
		c.Error(err) // Pass error to the middleware
		return
	}

	// Only save the image after the user is successfully created
	if fileHeader != nil {
		filePath := filepath.Join("public/images", user.Image)
		file, err := fileHeader.Open()
		if err != nil {
			c.Error(err) // Pass error to the middleware
			log.Println("Error Image:", err)
			return
		}
		defer file.Close()

		dst, err := os.Create(filePath)
		if err != nil {
			log.Println("Error Image:", err)
			c.Error(err) // Pass error to the middleware
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			log.Println("Error Image:", err)
			c.Error(err) // Pass error to the middleware
			return
		}
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
	branchIdStr := c.Request.FormValue("branch_id")

	branchId, err := strconv.ParseUint(branchIdStr, 10, 32)
	if err != nil {
		log.Println("Error converting branch_id to uint:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch_id"})
		return
	}

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
	user.BranchId = uint(branchId)

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

// BranchCounter Handlers

func GetBranchCounterHandlerByBranchId(c *gin.Context) {
	idParam := c.Param("branch_id")

	// Convert the string idParam to an integer
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch office ID"})
		return
	}

	// Check the branch office by ID
	branchOffice, err := services.GetBranchOfficeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Branch office not found"})
		return
	}

	// Call service to retrieve branch counters by branch ID
	counters, err := services.GetBranchCountersByBranchID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return JSON response with counters and total counter
	c.JSON(http.StatusOK, gin.H{
		"list_counter":  counters,
		"name_branch":   branchOffice.Name,
		"total_counter": branchOffice.TotalCounter,
	})
}

func CreateBranchCounterHandler(c *gin.Context) {
	var branchCounter models.BranchCounter
	var input map[string]interface{}

	// Parse request body into a map for validation
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate input using the validation function
	if err := validation.ValidateBranchCounter(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Manually bind validated values to branchCounter struct
	branchCounter.CounterLocation = input["counter_location"].(string)
	branchCounter.UserID = uint(input["user_id"].(float64))
	branchCounter.BranchID = uint(input["branch_id"].(float64))

	// Call service to create BranchCounter
	if err := services.CreateBranchCounter(&branchCounter); err != nil {
		c.Error(err) // Pass error to middleware
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "Branch counter created successfully"})
}

func DeleteBranchCounterHandler(c *gin.Context) {
	id := c.Param("id")

	// Call service to delete the branch counter
	if err := services.DeleteBranchCounter(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting branch counter"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch counter deleted successfully", "id": id})
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

	// Prepend the URL to the image field
	company.Logo = "http://localhost:3000/assets/" + company.Logo

	c.JSON(http.StatusOK, company)
}

func UpdateCompanyProfileHandler(c *gin.Context) {
	// Get the current company profile
	company, err := services.GetCompanyProfile() // Fetch company profile with ID 1
	if err != nil {
		c.Error(err)
		return
	}

	// Store the current image path (e.g., "application_logo.jpg")
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

	// Handle the file if provided
	if err == nil {
		fileName := "application_logo.png"
		logo = fileName // Store the filename for the logo

		// Open the file
		fileData, err := fileHeader.Open()
		if err != nil {
			c.Error(err)
			return
		}
		defer fileData.Close()

		// Save the new file
		newFilePath := filepath.Join("public/assets", logo)
		dst, err := os.Create(newFilePath)
		if err != nil {
			c.Error(err)
			return
		}
		defer dst.Close()

		// Copy the contents to the new file
		if _, err := io.Copy(dst, fileData); err != nil {
			c.Error(err)
			return
		}

		// Delete the old image if it exists
		if oldImagePath != "" {
			oldImageFullPath := filepath.Join("public/assets", oldImagePath)
			if err := os.Remove(oldImageFullPath); err != nil {
				log.Println("Error deleting old image:", err)
			} else {
				log.Println("Old image deleted:", oldImageFullPath)
			}
		}
	} else {
		// No new image uploaded, use the old image path
		logo = oldImagePath
	}

	// Call service to update company profile with ID 1
	if err := services.UpdateCompanyProfile(name, logo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company profile updated successfully"})
}
