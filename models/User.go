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
}

type AllUserResponse struct {
	ID       uint   `json:"ID"`
	FullName string `json:"FullName"`
	Email    string `json:"Email"`
	Role     string `json:"Role"`
	Likes    uint   `json:"Likes"`
	Dislikes uint   `json:"Dislikes"`
	Image    string `json:"Image"`
}
