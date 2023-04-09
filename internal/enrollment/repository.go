package enrollment

import (
	"log"

	"github.com/jonathannavas/go_web/internal/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(enrollment *domain.Enrollment) error
	}

	repo struct {
		db  *gorm.DB
		log *log.Logger
	}
)

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(enrollment *domain.Enrollment) error {
	if err := repo.db.Create(enrollment).Error; err != nil {
		repo.log.Println(err)
		return err
	}
	repo.log.Println("Enrollment created with id:", enrollment.ID)
	return nil
}
