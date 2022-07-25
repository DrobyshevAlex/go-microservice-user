package repositories

import (
	"context"
	"main/src/entities"
)

type UsersRepository struct {
}

func (r *UsersRepository) GetUserById(ctx context.Context, userId int64) (*entities.User, error) {

	user := &entities.User{
		Id:        userId,
		FirstName: "Firstname",
		LastName:  "LastName",
	}

	return user, nil
}

func NewUsersRepository() *UsersRepository {
	return &UsersRepository{}
}
