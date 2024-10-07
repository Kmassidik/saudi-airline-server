package services

import (
	"api-server/repository"
)

func TotalDataDashboard() (int, int, int, int, error) {
	return repository.TotalDataDashboard()
}

func TotalDataBranchOfficeDashboard(id uint, option string) error {
	return repository.TotalDataBranchDashboard(id, option)
}
