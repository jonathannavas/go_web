package user

import "log"

type Service interface {
	Create(firstName string, lastName string, email string, phone string) error
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s service) Create(firstName string, lastName string, email string, phone string) error {
	log.Println("Create user service")
	return nil
}
