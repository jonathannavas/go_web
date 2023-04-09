package course

import (
	"log"
	"time"

	"github.com/jonathannavas/go_web/internal/domain"
)

type (
	Service interface {
		Create(name, startDate, endDate string) (*domain.Course, error)
		GetAll(filters Filters, offset, limit int) ([]domain.Course, error)
		Get(id string) (*domain.Course, error)
		Delete(id string) error
		Update(id string, name *string, startDate, endDate *string) error
		Count(filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}

	Filters struct {
		Name string
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(name, startDate, endDate string) (*domain.Course, error) {
	log.Println("Create course service")

	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}

	course := domain.Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	if err := s.repo.Create(&course); err != nil {
		return nil, err
	}

	return &course, nil
}

func (s service) GetAll(filters Filters, offset, limit int) ([]domain.Course, error) {
	s.log.Println("Service GetAll")
	courses, err := s.repo.GetAll(filters, offset, limit)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}
	return courses, nil
}

func (s service) Get(id string) (*domain.Course, error) {
	s.log.Println("Get course service by id: ", id)
	course, err := s.repo.Get(id)
	if err != nil {
		s.log.Println(err)
		return nil, err
	}
	return course, nil
}

func (s service) Delete(id string) error {
	s.log.Println("Delete course service by id: ", id)
	return s.repo.Delete(id)
}

func (s service) Update(id string, name *string, startDate *string, endDate *string) error {

	var startDateParsed, endDateParsed *time.Time

	if startDate != nil {
		date, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			s.log.Println(err)
			return err
		}
		startDateParsed = &date
		*startDateParsed = startDateParsed.Add(time.Hour * 24)

	}

	if endDate != nil {
		date, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			s.log.Println(err)
			return err
		}
		endDateParsed = &date
		*endDateParsed = endDateParsed.Add(time.Hour * 24)
	}

	return s.repo.Update(id, name, startDateParsed, endDateParsed)
}

func (s service) Count(filters Filters) (int, error) {
	return s.repo.Count(filters)
}
