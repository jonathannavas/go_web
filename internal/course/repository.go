package course

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jonathannavas/go_web/internal/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(course *domain.Course) error
		GetAll(filters Filters, offset, limit int) ([]domain.Course, error)
		Get(id string) (*domain.Course, error)
		Delete(id string) error
		Update(id string, name *string, startDate, endDate *time.Time) error
		Count(filters Filters) (int, error)
	}

	repo struct {
		log *log.Logger
		db  *gorm.DB
	}
)

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(course *domain.Course) error {
	if err := repo.db.Create(course).Error; err != nil {
		repo.log.Println(err)
		return err
	}
	repo.log.Println("Course created with id: ", course.ID)
	return nil
}

func (repo *repo) GetAll(filters Filters, offset, limit int) ([]domain.Course, error) {
	var course []domain.Course

	tx := repo.db.Model(&course)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)

	result := tx.Order("created_at desc").Find(&course)

	if result.Error != nil {
		return nil, result.Error
	}

	return course, nil
}

func (repo *repo) Get(id string) (*domain.Course, error) {
	course := domain.Course{
		ID: id,
	}
	result := repo.db.First(&course)
	if result.Error != nil {
		return nil, result.Error
	}
	return &course, nil
}

func (repo *repo) Delete(id string) error {
	course := domain.Course{
		ID: id,
	}
	result := repo.db.Delete(&course)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *repo) Update(id string, name *string, startDate, endDate *time.Time) error {
	values := make(map[string]interface{})
	if name != nil {
		values["name"] = *name
	}

	if startDate != nil {
		values["start_date"] = *startDate
	}

	if endDate != nil {
		values["end_date"] = *endDate
	}

	result := repo.db.Model(&domain.Course{}).Where("id = ?", id).Updates(values)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(name) like ?", filters.Name)
	}
	return tx
}

func (repo *repo) Count(filters Filters) (int, error) {
	var count int64
	tx := repo.db.Model(domain.Course{})
	tx = applyFilters(tx, filters)
	if err := tx.Count(&count).Error; err != nil {
		return 0, nil
	}
	return int(count), nil
}
