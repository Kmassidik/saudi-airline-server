package services

import (
	"api-server/repository"
)

func TotalDataDashboard() (int, int, int, int, error) {
	return repository.TotalDataDashboard()
}

func TotalDataBranch() (int, int, int, int, error) {
	return repository.TotalDataDashboard()
}
