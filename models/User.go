package models

type User struct {
	ID       uint // Consider if uint is preferred, int could also work
	FullName string
	Email    string
	Password string
	Role     string
	Likes    uint // Changed from 'Like' to 'Likes' (plural for clarity)
	Dislikes uint // Changed from 'Dislike' to 'Dislikes'
}
