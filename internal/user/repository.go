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
	Delete(id string) error
	Update(id string, firstName *string, lastName *string, email *string, phone *string) error
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

func (repo *repo) Delete(id string) error {
	user := User{
		ID: id,
	}
	result := repo.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *repo) Update(id string, firstName *string, lastName *string, email *string, phone *string) error {

	values := make(map[string]interface{})

	if firstName != nil {
		values["first_name"] = firstName
	}

	if lastName != nil {
		values["last_name"] = lastName
	}

	if email != nil {
		values["email"] = email
	}

	if phone != nil {
		values["phone"] = phone
	}

	result := repo.db.Model(&User{}).Where("id = ?", id).Updates(values)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
