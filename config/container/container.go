package container

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/config"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/app"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers"
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
	controllers.AuthController
	controllers.UserController
	controllers.EduprogController
	controllers.EduprogcompController
	controllers.EduprogschemeController
	controllers.DisciplineController
	controllers.EducompRelationsController
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

	authController := controllers.NewAuthController(authService, userService)
	userController := controllers.NewUserController(userService)
	eduprogController := controllers.NewEduprogController(eduprogService, eduprogcompService)
	eduprogcompController := controllers.NewEduprogcompController(eduprogcompService)
	eduprogschemeController := controllers.NewEduprogschemeController(eduprogschemeService, eduprogcompService)
	disciplineController := controllers.NewDisciplineController(disciplineService)
	educompRelationsController := controllers.NewEducompRelationsController(educompRelationsService)

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
