package services

import (
	"api-server/models"
	"api-server/repository"
)

func VotedUser(voteType string, data *models.User) error {
	return repository.VotedUserLike(voteType, data)
}
