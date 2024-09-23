package models // Suggest using 'entities' instead of 'entitys'

type BranchOffice struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	TotalCounter uint   `json:"total_counter"`
}
