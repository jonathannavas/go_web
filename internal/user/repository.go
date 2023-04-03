package user

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	GetAll() ([]User, error)
	Get(id string) (*User, error)
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

func (repo *repo) GetAll() ([]User, error) {
	var user []User
	result := repo.db.Model(&user).Order("created_at desc").Find(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *repo) Get(id string) (*User, error) {
	user := User{
		ID: id,
	}
	result := repo.db.First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
