package Usecases

import (
	"errors"
	"strings"
	"task-manager/Domain"
	"task-manager/Repositories"

	"task-manager/Infrastructure"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrUserNotFound = errors.New("user not found")

type UserUseCase interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
	Promote(id string) error
}

type userUseCase struct {
	repo            Repositories.UserRepository
	passwordService Infrastructure.PasswordService
	jwtService      Infrastructure.JWTService
}

func NewUserUseCase(repo Repositories.UserRepository, ps Infrastructure.PasswordService, js Infrastructure.JWTService) UserUseCase {
	return &userUseCase{
		repo:            repo,
		passwordService: ps,
		jwtService:      js,
	}
}

func (uuc *userUseCase) Register(username, password string) error {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return errors.New("username and password are required")
	}

	_, err := uuc.repo.FindByUsername(username)
	if err == nil || !errors.Is(err, Infrastructure.ErrNoDocuments) {
		return errors.New("username already taken")
	}

	count, err := uuc.repo.Count()
	if err != nil {
		return err
	}

	role := "user"
	if count == 0 {
		role = "admin"
	}

	hashed, err := uuc.passwordService.HashPassword(password)
	if err != nil {
		return err
	}

	user := Domain.User{
		ID:       primitive.NewObjectID(),
		UserName: username,
		Password: hashed,
		Role:     role,
	}

	return uuc.repo.Create(user)
}

func (uuc *userUseCase) Login(username, password string) (string, error) {
	username = strings.TrimSpace(username)

	user, err := uuc.repo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, Infrastructure.ErrNoDocuments) {
			return "", errors.New("invalid username or password")
		}
		return "", err
	}

	err = uuc.passwordService.ComparePassword(user.Password, password)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	token, err := uuc.jwtService.GenerateToken(user.ID, user.UserName, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (uuc *userUseCase) Promote(id string) error {
	user, err := uuc.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, Infrastructure.ErrNoDocuments) {
			return ErrUserNotFound
		}
		return err
	}

	if user.Role == "admin" {
		return errors.New("user is already an admin")
	}

	err = uuc.repo.UpdateRole(user.ID, "admin")
	if err != nil && errors.Is(err, errors.New("user not found")) {
		return ErrUserNotFound
	}
	return err
}
