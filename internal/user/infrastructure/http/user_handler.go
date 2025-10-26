package http

import "github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/internal/user/application"

type UserHandler struct {
	userService *application.UserService
}

func NewUserHandler(userService *application.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
