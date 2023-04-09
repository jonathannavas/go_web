package enrollment

import (
	"errors"
	"log"

	"github.com/jonathannavas/go_web/internal/course"
	"github.com/jonathannavas/go_web/internal/domain"
	"github.com/jonathannavas/go_web/internal/user"
)

type (
	Service interface {
		Create(userID, courseID string) (*domain.Enrollment, error)
	}

	service struct {
		log           *log.Logger
		repo          Repository
		userService   user.Service
		courseService course.Service
	}
)

func NewService(l *log.Logger, repo Repository, userService user.Service, courseService course.Service) Service {
	return &service{
		log:           l,
		repo:          repo,
		userService:   userService,
		courseService: courseService,
	}
}

func (s service) Create(userID, courseID string) (*domain.Enrollment, error) {
	enrollment := &domain.Enrollment{
		UserID:   userID,
		CourseId: courseID,
		Status:   "P",
	}

	if _, err := s.userService.Get(userID); err != nil {
		return nil, errors.New("user id doesn't exists")
	}

	if _, err := s.courseService.Get(courseID); err != nil {
		return nil, errors.New("course id doesn't exists")
	}

	if err := s.repo.Create(enrollment); err != nil {
		s.log.Println(err)
		return nil, err
	}
	return enrollment, nil
}
