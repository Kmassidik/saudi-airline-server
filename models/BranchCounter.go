package models

type BranchCounter struct {
	ID              uint   `json:"id"`
	CounterLocation string `json:"counter_location"`
	UserID          uint   `json:"user_id"`
	BranchID        uint   `json:"branch_id"`
}

// BranchCounterWithNames includes additional fields for names from related tables
type BranchCounterWithNames struct {
	ID              uint   `json:"id"`
	CounterLocation string `json:"counter_location"`
	FullName        string `json:"full_name"` // Name from users table
	Image           string `json:"image"`
}
