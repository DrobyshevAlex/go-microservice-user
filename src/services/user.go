package services

import (
	"context"
	"main/src/ampq/requests"
	"main/src/ampq/responses"
	"main/src/entities"
	"main/src/repositories"
)

type UserService struct {
	repository *repositories.UsersRepository
}

func (s *UserService) GetUserById(
	ctx context.Context,
	request *requests.UserGetByIdRequest,
) (*responses.UserResponse, error) {
	user, err := s.repository.GetUserById(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		user = &entities.User{}
	}

	return &responses.UserResponse{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, nil
}
func NewUserService(repository *repositories.UsersRepository) *UserService {
	return &UserService{repository: repository}
}
