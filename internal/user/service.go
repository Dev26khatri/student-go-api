package user

import (
	"strings"
	"student-go-service/internal/dto"
	"student-go-service/internal/package/jwt"
	"student-go-service/internal/package/utils"
	"time"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}
func (s *Service) Register(req dto.RegisterRequest) (*User, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	hashedPassword, err := utils.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:     req.Name,
		Email:    email,
		Password: hashedPassword,
		Role:     "User",
	}

	if err := s.repository.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(req dto.LoginRequest, jwtSeceret string, jwtExpiration time.Duration) (*User, string, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.repository.FindByEmail(email)

	if err != nil {
		return nil, "", utils.ErrInvalidCredentials
	}
	if err := utils.Compare(user.Password, req.Password); err != nil {
		return nil, "", utils.ErrInvalidCredentials
	}
	token, err := jwt.GenerateToken(
		user.ID,
		user.Role,
		jwtSeceret,
		jwtExpiration,
	)
	if err != nil {
		return nil, "", err
	}
	return user, token, nil

}
