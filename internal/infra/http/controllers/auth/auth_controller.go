package auth

import (
	"errors"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/app"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/requests"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/resources"
	"log"
	"net/http"
)

type AuthController struct {
	authService app.AuthService
	userService app.UserService
}

func NewAuthController(as app.AuthService, us app.UserService) AuthController {
	return AuthController{
		authService: as,
		userService: us,
	}
}

func (c AuthController) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.RegisterRequest{}, domain.User{})
		if err != nil {
			log.Printf("AuthController: %s", err)
			controllers.BadRequest(w, errors.New("invalid request body"))
			return
		}

		user, token, err := c.authService.Register(user)
		if err != nil {
			log.Printf("AuthController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var authDto resources.AuthDto
		controllers.Success(w, authDto.DomainToDto(token, user))
	}
}

func (c AuthController) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.AuthRequest{}, domain.User{})
		if err != nil {
			log.Printf("AuthController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		u, token, err := c.authService.Login(user)
		if err != nil {
			controllers.Unauthorized(w, err)
			return
		}

		var authDto resources.AuthDto
		controllers.Success(w, authDto.DomainToDto(token, u))
	}
}

func (c AuthController) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := r.Context().Value(controllers.SessKey).(domain.Session)
		err := c.authService.Logout(sess)
		if err != nil {
			log.Print(err)
			controllers.InternalServerError(w, err)
			return
		}

		controllers.NoContent(w)
	}
}

func (c AuthController) ChangePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := requests.Bind(r, requests.ChangePasswordRequest{}, domain.ChangePassword{})
		if err != nil {
			log.Printf("AuthController: %s", err)
			controllers.BadRequest(w, err)
			return
		}
		sess := r.Context().Value(controllers.SessKey).(domain.Session)
		user := r.Context().Value(controllers.UserKey).(domain.User)

		err = c.authService.ChangePassword(user, req, sess)
		if err != nil {
			log.Printf("AuthController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		controllers.Ok(w)
	}
}
