package services

import (
	"api-server/models"
	"api-server/repository"
)

func TotalDataDashboard() (int, int, int, int, error) {
	return repository.TotalDataDashboard()
}

func TotalDataBranchOfficeDashboard() ([]models.BranchData, error) {
	return repository.TotalDataBranchDashboard()
}

func UpdateDataDashboard(branchId uint, voteType string) error {
	return repository.UpdateDashboard(branchId, voteType)
}

// GetAllBranchOffices retrieves all branch offices with pagination

func GetAllOfficers(limit uint, offset uint) ([]models.DashboardUsers, error) {
	return repository.DataOfficerDashboard(limit, offset)
}
