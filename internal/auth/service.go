package auth

import (
	"URLshortener/internal/user"
	"URLshortener/pkg/di"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository di.IUserRepository
}

func NewAuthService(userRepository di.IUserRepository) *AuthService {
	return &AuthService{UserRepository: userRepository}
}

func (service *AuthService) Login(email, password string) (string, uint, error) {
	existedUser, err := service.UserRepository.GetByEmail(email)
	if existedUser == nil || err != nil {
		return "", 0, errors.New(ErrWrongCredentials)
	}
	//проверка введенного пароля и хэшированного в бд
	err = bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password))
	if err != nil {
		return "", 0, errors.New(ErrWrongCredentials)
	}
	return existedUser.Email, existedUser.ID, nil
}

func (service *AuthService) Register(email, password, name string) (string, uint, error) {
	existedUser, _ := service.UserRepository.GetByEmail(email)
	if existedUser != nil {
		return "", 0, errors.New(ErrUserExists)
	}
	//хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", 0, err
	}
	user := &user.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	newUser, err := service.UserRepository.Create(user)
	if err != nil {
		return "", 0, err
	}
	return user.Email, newUser.ID, nil
}
