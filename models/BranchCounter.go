package models // Suggest using 'entities' instead of 'entitys'

type BranchCounter struct {
	ID              uint
	CounterLocation string
	UserID          uint // Changed 'UserId' to 'UserID' to match Go naming conventions
	BranchID        uint // Changed 'BranchId' to 'BranchID' for consistency
}
