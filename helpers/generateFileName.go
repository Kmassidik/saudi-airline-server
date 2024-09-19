package helpers

import (
	"fmt"
	"time"
)

// GenerateFileName creates a unique random file name using a timestamp and ".png" extension
func GenerateFileName() string {
	return fmt.Sprintf("%d.png", time.Now().UnixNano())
}
