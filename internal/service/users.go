package service

import (
	"awesome-go/internal/models"
	"awesome-go/internal/types"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	pre  string = "*89f8dsIo$(){^"
	post string = "^}fdssaJFDdls"
)

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pre+password+post), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func (s *Service) CreateUser(user types.UserForm) (models.User, error) {
	newUser := models.User{
		Name:     user.Name.Value,
		Email:    user.Email.Value,
		Password: hashPassword(user.Password.Value),
	}
	err := gorm.G[models.User](s.db).Create(s.context(), &newUser)
	if err != nil {
		log.Print(err)
		return models.User{}, err
	}
	return newUser, nil
}

func (s *Service) GetUserByEmail(email string) (models.User, error) {
	user, err := gorm.G[models.User](s.db).Where("email = ?", email).First(s.context())
	if err != nil {
		log.Print(err)
		return models.User{}, err
	}

	return user, nil
}

func (s *Service) AuthenticateUser(credentials types.AuthForm) (models.User, error) {
	user, err := s.GetUserByEmail(credentials.Email.Value)
	if err != nil {
		return models.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pre+credentials.Password.Value+post))
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *Service) UpdateUser(user models.User) error {
	return nil
}
