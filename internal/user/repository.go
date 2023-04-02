package user

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(user *User) error {

	user.ID = uuid.New().String()

	// result := repo.db.Debug().Create(user)
	// if result.Error != nil {
	// 	repo.log.Println(result.Error)
	// 	return result.Error
	// }

	if err := repo.db.Debug().Create(user).Error; err != nil {
		repo.log.Println(err)
		return err
	}

	repo.log.Println("User created with id: ", user.ID)
	return nil
}
