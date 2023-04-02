package user

import (
	"log"
)

type Service interface {
	Create(firstName string, lastName string, email string, phone string) (*User, error)
}

type service struct {
	log  *log.Logger
	repo Repository
}

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(firstName string, lastName string, email string, phone string) (*User, error) {
	log.Println("Create user service")

	user := User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}

	// si recibo un puntero debo enviar la dirección de memoria a donde esta con el simbolo de &

	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
