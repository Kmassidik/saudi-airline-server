package models

type User struct {
	ID       uint   // FIXME: make Gin skip mapping this field!
	FullName string `form:"full_name"`
	Email    string `form:"email"`
	Password string `form:"password"`
	Role     string `form:"role"`
	Likes    uint   `form:"likes"`
	Dislikes uint   `form:"dislikes"`
	Image    string `form:"image"`
	BranchId uint   `form:"branch_id"`
}

type UserAllResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Likes    uint   `json:"likes"`
	Dislikes uint   `json:"dislikes"`
	Image    string `json:"image"`
	BranchId *uint  `json:"branch_id"`
}

// UserResponse omits the password when retrieving user data (e.g., all users or by ID)
type UserResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Likes    uint   `json:"likes"`
	Dislikes uint   `json:"dislikes"`
	Image    string `json:"image"`
	BranchId *uint  `json:"branch_id"`
	Password string `json:"password"`
}

// UserDetailResponse for more detailed responses, again without exposing the password
type UserDetailResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Likes    uint   `json:"likes"`
	Dislikes uint   `json:"dislikes"`
	Image    string `json:"image"`
	// Add other fields as needed (but no password)
}

// UserResponse omits the password when retrieving user data (e.g., all users or by ID)
type UserByBranchOfiiceResponse struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
}

// UserDetailResponse for more detailed responses, again without exposing the password
type DashboardUsers struct {
	Name     string `json:"full_name"`
	Likes    uint   `json:"likes"`
	Dislikes uint   `json:"dislikes"`
}
