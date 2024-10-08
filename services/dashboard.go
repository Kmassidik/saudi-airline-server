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

func UpdateDataDashboard(branchId uint, voteType string) error {
	return repository.UpdateDashboard(branchId, voteType)
}
