package services

import (
	"api-server/models"
	"api-server/repository"
)

// GetCompanyProfile retrieves the company profile from the repository
func GetCompanyProfile() (*models.CompanyProfile, error) {
	return repository.GetCompanyProfile()
}

func UpdateCompanyProfile(name string, logo string) error {
	const companyID = 1 // Fixed ID
	return repository.UpdateCompanyProfile(companyID, name, logo)
}
