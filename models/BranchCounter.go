package models

type BranchCounter struct {
	ID              uint   `json:"id"`
	CounterLocation string `json:"counter_location"`
	UserID          uint   `json:"user_id"`
	BranchID        uint   `json:"branch_id"`
}
