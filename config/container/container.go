package container

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/config"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/app"
	auth2 "github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/auth"
	eduprog2 "github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
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
	app.CompetenciesBaseService
	app.EduprogcompetenciesService
	app.CompetenciesMatrixService
	app.ResultsMatrixService
	app.SpecialtiesService
}

type Controllers struct {
	auth.AuthController
	auth.UserController
	eduprog.EduprogController
	eduprog.EduprogcompController
	eduprog.EduprogschemeController
	eduprog.DisciplineController
	eduprog.EducompRelationsController
	eduprog.CompetenciesBaseController
	eduprog.EduprogcompetenciesController
	eduprog.CompetenciesMatrixController
	eduprog.ResultsMatrixController
	eduprog.SpecialtyController
}

func New(conf config.Configuration) Container {
	tknAuth := jwtauth.New("HS256", []byte(conf.JwtSecret), nil)
	sess := getDbSess(conf)

	userRepository := auth2.NewUserRepository(sess)
	sessionRepository := auth2.NewSessRepository(sess)
	eduprogRepository := eduprog2.NewEduprogRepository(sess)
	eduprogcompRepository := eduprog2.NewEduprogcompRepository(sess)
	eduprogschemeRepository := eduprog2.NewEduprogschemeRepository(sess)
	disciplineRepository := eduprog2.NewDisciplineRepository(sess)
	educompRelationsRepository := eduprog2.NewEducompRelationsRepository(sess)
	competencyBaseRepository := eduprog2.NewCompetenciesBaseRepository(sess)
	competencyMatrixRepository := eduprog2.NewCompetenciesMatrixRepository(sess)
	eduprogcompetenciesRepository := eduprog2.NewEduprogcompetenciesRepository(sess)
	resultsMatrixRepository := eduprog2.NewResultsMatrixRepository(sess)
	specialtiesRepository := eduprog2.NewSpecialtiesRepository(sess)

	userService := app.NewUserService(userRepository)
	authService := app.NewAuthService(sessionRepository, userService, conf, tknAuth)
	eduprogService := app.NewEduprogService(eduprogRepository)
	eduprogcompService := app.NewEduprogcompService(eduprogcompRepository)
	eduprogschemeService := app.NewEduprogschemeService(eduprogschemeRepository)
	disciplineService := app.NewDisciplineService(disciplineRepository)
	educompRelationsService := app.NewEducompRelationsService(educompRelationsRepository)
	competencyBaseService := app.NewCompetenciesBaseService(competencyBaseRepository)
	competencyMatrixService := app.NewCompetenciesMatrixService(competencyMatrixRepository)
	eduprogcompetenciesService := app.NewEduprogcompetenciesService(eduprogcompetenciesRepository)
	resultMatrixService := app.NewResultsMatrixService(resultsMatrixRepository)
	specialtiesService := app.NewSpecialtiesService(specialtiesRepository)

	authController := auth.NewAuthController(authService, userService)
	userController := auth.NewUserController(userService)
	eduprogController := eduprog.NewEduprogController(eduprogService, eduprogcompService, eduprogcompetenciesService, competencyMatrixService, resultMatrixService, specialtiesService, educompRelationsService)
	eduprogcompController := eduprog.NewEduprogcompController(eduprogcompService, eduprogService, eduprogController)
	eduprogschemeController := eduprog.NewEduprogschemeController(eduprogschemeService, eduprogcompService, disciplineService)
	disciplineController := eduprog.NewDisciplineController(disciplineService)
	educompRelationsController := eduprog.NewEducompRelationsController(educompRelationsService, eduprogschemeService, eduprogcompService)
	competencyBaseController := eduprog.NewCompetenciesBaseController(competencyBaseService, eduprogService)
	competencyMatrixController := eduprog.NewCompetenciesMatrixController(competencyMatrixService)
	eduprogcompetenciesController := eduprog.NewEduprogcompetenciesController(eduprogcompetenciesService, competencyBaseService, eduprogService)
	resultsMatrixController := eduprog.NewResultsMatrixController(resultMatrixService)
	specialtiesController := eduprog.NewSpecialtiesController(specialtiesService)

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
			competencyBaseService,
			eduprogcompetenciesService,
			competencyMatrixService,
			resultMatrixService,
			specialtiesService,
		},
		Controllers: Controllers{
			authController,
			userController,
			eduprogController,
			eduprogcompController,
			eduprogschemeController,
			disciplineController,
			educompRelationsController,
			competencyBaseController,
			eduprogcompetenciesController,
			competencyMatrixController,
			resultsMatrixController,
			specialtiesController,
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
