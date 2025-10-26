package application

import (
	jpdomain "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/jobprofile/domain"
	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/domain"
)

type UserService struct {
	userRepo                    domain.UserRepository
	jobProfileRepo              jpdomain.JobProfileRepository
	notificationPreferencesRepo domain.NotificationPreferencesRepository
}

func NewUserService(userRepo domain.UserRepository, jobProfileRepo jpdomain.JobProfileRepository, notificationPreferencesRepo domain.NotificationPreferencesRepository) *UserService {
	return &UserService{
		userRepo:                    userRepo,
		jobProfileRepo:              jobProfileRepo,
		notificationPreferencesRepo: notificationPreferencesRepo,
	}

}
