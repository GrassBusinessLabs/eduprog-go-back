package container

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/config"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/app"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers/auth"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers/eduprog"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/middlewares"
	"github.com/go-chi/jwtauth/v5"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"

	//"github.com/upper/db/v4/adapter/postgresql"
	//"github.com/upper/db/v4/adapter/sqlite"
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
	app.EduprogschemeService
	app.DisciplineService
	app.EducompRelationsService
}

type Controllers struct {
	auth.AuthController
	auth.UserController
	eduprog.EduprogController
	eduprog.EduprogcompController
	eduprog.EduprogschemeController
	eduprog.DisciplineController
	eduprog.EducompRelationsController
}

func New(conf config.Configuration) Container {
	tknAuth := jwtauth.New("HS256", []byte(conf.JwtSecret), nil)
	sess := getDbSess(conf)

	userRepository := database.NewUserRepository(sess)
	sessionRepository := database.NewSessRepository(sess)
	eduprogRepository := database.NewEduprogRepository(sess)
	eduprogcompRepository := database.NewEduprogcompRepository(sess)
	eduprogschemeRepository := database.NewEduprogschemeRepository(sess)
	disciplineRepository := database.NewDisciplineRepository(sess)
	educompRelationsRepository := database.NewEducompRelationsRepository(sess)

	userService := app.NewUserService(userRepository)
	authService := app.NewAuthService(sessionRepository, userService, conf, tknAuth)
	eduprogService := app.NewEduprogService(eduprogRepository)
	eduprogcompService := app.NewEduprogcompService(eduprogcompRepository)
	eduprogschemeService := app.NewEduprogschemeService(eduprogschemeRepository)
	disciplineService := app.NewDisciplineService(disciplineRepository)
	educompRelationsService := app.NewEducompRelationsService(educompRelationsRepository)

	authController := auth.NewAuthController(authService, userService)
	userController := auth.NewUserController(userService)
	eduprogController := eduprog.NewEduprogController(eduprogService, eduprogcompService)
	eduprogcompController := eduprog.NewEduprogcompController(eduprogcompService)
	eduprogschemeController := eduprog.NewEduprogschemeController(eduprogschemeService, eduprogcompService)
	disciplineController := eduprog.NewDisciplineController(disciplineService)
	educompRelationsController := eduprog.NewEducompRelationsController(educompRelationsService, eduprogschemeService, eduprogcompService)

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
			eduprogschemeService,
			disciplineService,
			educompRelationsService,
		},
		Controllers: Controllers{
			authController,
			userController,
			eduprogController,
			eduprogcompController,
			eduprogschemeController,
			disciplineController,
			educompRelationsController,
		},
	}
}

func getDbSess(conf config.Configuration) db.Session {
	sess, err := postgresql.Open(
		postgresql.ConnectionURL{
			User:     conf.DatabaseUser,
			Host:     conf.DatabaseHost,
			Password: conf.DatabasePassword,
			Database: conf.DatabaseName,
		})
	//sess, err := sqlite.Open(
	//	sqlite.ConnectionURL{
	//		Database: conf.DatabasePath,
	//	})
	if err != nil {
		log.Fatalf("Unable to create new DB session: %q\n", err)
	}
	return sess
}
