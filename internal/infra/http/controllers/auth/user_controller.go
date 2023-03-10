package auth

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/app"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/requests"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/resources"
	"log"
	"net/http"
)

type UserController struct {
	userService app.UserService
}

func NewUserController(us app.UserService) UserController {
	return UserController{
		userService: us,
	}
}

func (c UserController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.RegisterRequest{}, domain.User{})
		if err != nil {
			log.Printf("UserController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		user, err = c.userService.Save(user)
		if err != nil {
			log.Printf("UserController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var userDto resources.UserDto
		controllers.Created(w, userDto.DomainToDto(user))
	}
}

func (c UserController) FindMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(controllers.UserKey).(domain.User)
		controllers.Success(w, resources.UserDto{}.DomainToDto(user))
	}
}

func (c UserController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.UpdateUserRequest{}, domain.User{})
		if err != nil {
			log.Printf("UserController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		u := r.Context().Value(controllers.UserKey).(domain.User)
		u.Name = user.Name
		user, err = c.userService.Update(u)
		if err != nil {
			log.Printf("UserController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var userDto resources.UserDto
		controllers.Success(w, userDto.DomainToDto(user))
	}
}

func (c UserController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(controllers.UserKey).(domain.User)

		err := c.userService.Delete(u.Id)
		if err != nil {
			log.Printf("UserController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		controllers.Ok(w)
	}
}
