package container

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/config"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/app"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/middlewares"
	"github.com/go-chi/jwtauth/v5"
	"github.com/upper/db/v4"
	//"github.com/upper/db/v4/adapter/postgresql"
	"github.com/upper/db/v4/adapter/sqlite"
	"log"
	"net/http"
)

type Container struct {
	Middlewares
	Services
	Controllers
}

type Middlewares struct {
	AuthMw func(http.Handler) http.Handler
}

type Services struct {
	app.AuthService
	app.UserService
	app.EduprogService
	app.EduprogcompService
}

type Controllers struct {
	controllers.AuthController
	controllers.UserController
	controllers.EduprogController
	controllers.EduprogcompController
}

func New(conf config.Configuration) Container {
	tknAuth := jwtauth.New("HS256", []byte(conf.JwtSecret), nil)
	sess := getDbSess(conf)

	userRepository := database.NewUserRepository(sess)
	sessionRepository := database.NewSessRepository(sess)
	eduprogRepository := database.NewEduprogRepository(sess)
	eduprogcompRepository := database.NewEduprogcompRepository(sess)

	userService := app.NewUserService(userRepository)
	eduprogService := app.NewEduprogService(eduprogRepository)
	eduprogcompService := app.NewEduprogcompService(eduprogcompRepository)
	authService := app.NewAuthService(sessionRepository, userService, conf, tknAuth)

	authController := controllers.NewAuthController(authService, userService)
	userController := controllers.NewUserController(userService)
	eduprogController := controllers.NewEduprogController(eduprogService, eduprogcompRepository)
	eduprogcompController := controllers.NewEduprogcompController(eduprogcompService)

	authMiddleware := middlewares.AuthMiddleware(tknAuth, authService, userService)

	return Container{
		Middlewares: Middlewares{
			AuthMw: authMiddleware,
		},
		Services: Services{
			authService,
			userService,
			eduprogService,
			eduprogcompService,
		},
		Controllers: Controllers{
			authController,
			userController,
			eduprogController,
			eduprogcompController,
		},
	}
}

func getDbSess(conf config.Configuration) db.Session {
	//sess, err := postgresql.Open(
	//	postgresql.ConnectionURL{
	//		User:     conf.DatabaseUser,
	//		Host:     conf.DatabaseHost,
	//		Password: conf.DatabasePassword,
	//		Database: conf.DatabaseName,
	//	})
	sess, err := sqlite.Open(
		sqlite.ConnectionURL{
			Database: conf.DatabasePath,
		})
	if err != nil {
		log.Fatalf("Unable to create new DB session: %q\n", err)
	}
	return sess
}
