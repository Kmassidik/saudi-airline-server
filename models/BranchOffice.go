package models

type BranchOfficeCreateRequest struct {
	Name         string `json:"name"`
	Address      string `json:"address"`
	TotalCounter uint   `json:"total_counter"`
}

type BranchOfficeResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	TotalCounter uint   `json:"total_counter"`
}

type BranchOfficeListResponse struct {
	Offices []BranchOfficeResponse `json:"offices"`
}
