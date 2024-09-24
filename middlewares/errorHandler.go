package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// ErrorHandler is a middleware that captures errors from all routes and handles them centrally
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process the request

		// Check if any errors occurred
		errs := c.Errors
		if len(errs) > 0 {
			err := errs[0].Err

			// Log the error for debugging
			log.Printf("Encountered error: %v", err)

			// Handle specific PostgreSQL errors
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Table {
				case "users":
					handlePostgresErrorUser(c, pqErr)
				case "branch_offices":
					handlePostgresErrorBranchOffice(c, pqErr)
				case "branch_counters":
					handlePostgresErrorBranchCounter(c, pqErr)
				}
				return
			}

			// Handle generic errors
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "An internal error occurred",
			})
			return
		}
	}
}

// handlePostgresError handles PostgreSQL specific errors and sends appropriate responses
func handlePostgresErrorUser(c *gin.Context, pqErr *pq.Error) {
	log.Printf("PostgreSQL error: %v", pqErr)

	switch pqErr.Code.Name() {
	case "unique_violation":
		c.JSON(http.StatusConflict, gin.H{"error": "Duplicate key: the email already exists"})
	case "foreign_key_violation":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Foreign key violation, Branch ID not found!"})
	case "not_null_violation":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field"})
	case "check_violation":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data failed validation check"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error: " + pqErr.Message,
		})
	}
}

// handlePostgresError handles PostgreSQL specific errors and sends appropriate responses
func handlePostgresErrorBranchOffice(c *gin.Context, pqErr *pq.Error) {
	log.Printf("PostgreSQL error: %v", pqErr)

	switch pqErr.Code.Name() {
	case "unique_violation":
		c.JSON(http.StatusConflict, gin.H{"error": "Duplicate key: the email already exists"})
	case "foreign_key_violation":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Foreign key violation"})
	case "not_null_violation":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field"})
	case "check_violation":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data failed validation check"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error: " + pqErr.Message,
		})
	}
}

// handlePostgresError handles PostgreSQL specific errors and sends appropriate responses
func handlePostgresErrorBranchCounter(c *gin.Context, pqErr *pq.Error) {
	log.Printf("PostgreSQL error: %v", pqErr)

	switch pqErr.Code.Name() {
	case "unique_violation":
		c.JSON(http.StatusConflict, gin.H{"error": "Duplicate key: the email already exists"})
	case "foreign_key_violation":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Foreign key violation, please check user_id or branch_id "})
	case "not_null_violation":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field"})
	case "check_violation":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data failed validation check"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error: " + pqErr.Message,
		})
	}
}
