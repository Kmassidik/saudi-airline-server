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
type BranchOfficeOptionResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type BranchOfficeListResponse struct {
	Offices []BranchOfficeResponse `json:"offices"`
}

type BranchData struct {
	ID            int    `json:"id"`
	NameOffice    string `json:"name_office"`
	TotalLikes    int    `json:"total_likes"`
	TotalDislikes int    `json:"total_dislikes"`
	BranchID      int    `json:"branch_id"`
}
