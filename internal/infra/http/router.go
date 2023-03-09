package http

import (
	"encoding/json"
	"fmt"
	"github.com/GrassBusinessLabs/eduprog-go-back/config"
	"github.com/GrassBusinessLabs/eduprog-go-back/config/container"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers/auth"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers/eduprog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Router(cont container.Container) http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.RedirectSlashes, middleware.Logger, cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Route("/api", func(apiRouter chi.Router) {
		// Health
		apiRouter.Route("/ping", func(healthRouter chi.Router) {
			healthRouter.Get("/", PingHandler())
			healthRouter.Handle("/*", NotFoundJSON())
		})

		apiRouter.Route("/v1", func(apiRouter chi.Router) {
			// Public routes
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Route("/auth", func(apiRouter chi.Router) {
					AuthRouter(apiRouter, cont.AuthController, cont.AuthMw)
				})
			})

			// Protected routes
			apiRouter.Group(func(apiRouter chi.Router) {
				apiRouter.Use(cont.AuthMw)

				UserRouter(apiRouter, cont.UserController)
				EduprogRouter(apiRouter, cont.EduprogController)
				EduprogcompRouter(apiRouter, cont.EduprogcompController)
				EduprogschemeRouter(apiRouter, cont.EduprogschemeController)
				DisciplineRouter(apiRouter, cont.DisciplineController)
				EducompRelationsRouter(apiRouter, cont.EducompRelationsController)
				CompetenciesBaseRouter(apiRouter, cont.CompetenciesBaseController)
				EduprogcompetenciesRouter(apiRouter, cont.EduprogcompetenciesController)
				CompetenciesMatrixRouter(apiRouter, cont.CompetenciesMatrixController)
				EduprogresultsRouter(apiRouter, cont.EduprogresultsController)
				apiRouter.Handle("/*", NotFoundJSON())
			})
		})
	})

	router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, config.GetConfiguration().FileStorageLocation))
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
		fs.ServeHTTP(w, r)
	})

	return router
}

func AuthRouter(r chi.Router, ac auth.AuthController, amw func(http.Handler) http.Handler) {
	r.Route("/", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/register",
			ac.Register(),
		)
		apiRouter.Post(
			"/login",
			ac.Login(),
		)
		apiRouter.With(amw).Post(
			"/change-pwd",
			ac.ChangePassword(),
		)
		apiRouter.With(amw).Post(
			"/logout",
			ac.Logout(),
		)
	})
}

func UserRouter(r chi.Router, uc auth.UserController) {
	r.Route("/users", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/",
			uc.FindMe(),
		)
		apiRouter.Put(
			"/",
			uc.Update(),
		)
		apiRouter.Delete(
			"/",
			uc.Delete(),
		)
	})
}

func EduprogRouter(r chi.Router, ec eduprog.EduprogController) {
	r.Route("/eduprogs", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/create",
			ec.Save(),
		)
		apiRouter.Put(
			"/{epId}",
			ec.Update(),
		)
		apiRouter.Get(
			"/",
			ec.ShowList(),
		)
		apiRouter.Get(
			"/{epId}",
			ec.FindById(),
		)
		apiRouter.Get(
			"/credits/{epId}",
			ec.CreditsInfo(),
		)
		apiRouter.Get(
			"/toExcel/{edId}",
			ec.ExportEduprogListToExcel(),
		)
		apiRouter.Get(
			"/cMatrixToExcel/{edId}",
			ec.ExportCompetenciesMatrixToExcel(),
		)
		apiRouter.Delete(
			"/{epId}",
			ec.Delete(),
		)

	})
}

func EduprogcompRouter(r chi.Router, ec eduprog.EduprogcompController) {
	r.Route("/eduprogs/comps", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/create",
			ec.Save(),
		)
		apiRouter.Put(
			"/{epcId}",
			ec.Update(),
		)
		apiRouter.Get(
			"/",
			ec.ShowList(),
		)
		apiRouter.Get(
			"/byEduprogId/{epcId}",
			ec.ShowListByEduprogId(),
		)
		apiRouter.Get(
			"/{epcId}",
			ec.FindById(),
		)
		apiRouter.Delete(
			"/{epcId}",
			ec.Delete(),
		)
	})
}

func EduprogschemeRouter(r chi.Router, esc eduprog.EduprogschemeController) {
	r.Route("/eduprogs/scheme", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/setCompToSemester",
			esc.SetComponentToEdprogscheme(),
		)
		apiRouter.Put(
			"/{essId}",
			esc.UpdateComponentInEduprogscheme(),
		)
		apiRouter.Get(
			"/{essId}",
			esc.FindById(),
		)
		apiRouter.Get(
			"/bySemester/{sNum}/{edId}",
			esc.FindBySemesterNum(),
		)
		apiRouter.Get(
			"/byEduprogId/{sNum}",
			esc.ShowSchemeByEduprogId(),
		)
		apiRouter.Get(
			"/freeComps/{sNum}",
			esc.ShowFreeComponents(),
		)
		apiRouter.Delete(
			"/{essId}",
			esc.Delete(),
		)
	})
}

func DisciplineRouter(r chi.Router, d eduprog.DisciplineController) {
	r.Route("/eduprogs/scheme/disciplines", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/create",
			d.Save(),
		)
		apiRouter.Put(
			"/{epId}",
			d.Update(),
		)
		apiRouter.Get(
			"/getByEdId/{epId}",
			d.ShowDisciplinesByEduprogId(),
		)
		apiRouter.Get(
			"/{epId}",
			d.FindById(),
		)
		apiRouter.Delete(
			"/{epId}",
			d.Delete(),
		)

	})
}

func EducompRelationsRouter(r chi.Router, ecrc eduprog.EducompRelationsController) {
	r.Route("/eduprogs/compRelations", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/create",
			ecrc.CreateRelation(),
		)
		apiRouter.Get(
			"/{epId}",
			ecrc.ShowByEduprogId(),
		)
		apiRouter.Get(
			"/posRel/{edId}/{compId}",
			ecrc.ShowPossibleRelationsForComp(),
		)
		apiRouter.Delete(
			"/{baseId}/{childId}",
			ecrc.Delete(),
		)
	})

}

func CompetenciesBaseRouter(r chi.Router, cbc eduprog.CompetenciesBaseController) {

	r.Route("/eduprogs/baseCompetencies", func(apiRouter chi.Router) {
		apiRouter.Get(
			"/list",
			cbc.ShowAllCompetencies(),
		)
		apiRouter.Get(
			"/ZK_list",
			cbc.ShowZK(),
		)
		apiRouter.Get(
			"/FK_list",
			cbc.ShowFK(),
		)
		apiRouter.Get(
			"/{cbId}",
			cbc.FindById(),
		)
	})
}

func EduprogcompetenciesRouter(r chi.Router, ecc eduprog.EduprogcompetenciesController) {

	r.Route("/eduprogs/competencies", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/add",
			ecc.AddCompetencyToEduprog(),
		)
		apiRouter.Get(
			"/byEduprogId/{edId}",
			ecc.ShowCompetenciesByEduprogId(),
		)
		apiRouter.Get(
			"/{compId}",
			ecc.FindById(),
		)
		apiRouter.Delete(
			"/{compId}",
			ecc.Delete(),
		)
	})
}

func CompetenciesMatrixRouter(r chi.Router, cmc eduprog.CompetenciesMatrixController) {
	r.Route("/eduprogs/competenciesMatrix", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/create",
			cmc.CreateRelation(),
		)
		apiRouter.Get(
			"/{epId}",
			cmc.ShowByEduprogId(),
		)
		apiRouter.Delete(
			"/{componentId}/{competencyId}",
			cmc.Delete(),
		)
	})

}

func EduprogresultsRouter(r chi.Router, erc eduprog.EduprogresultsController) {

	r.Route("/eduprogs/results", func(apiRouter chi.Router) {
		apiRouter.Post(
			"/add",
			erc.AddEduprogresultToEduprog(),
		)
		apiRouter.Get(
			"/byEduprogId/{edId}",
			erc.ShowEduprogResultsByEduprogId(),
		)
		apiRouter.Get(
			"/{resId}",
			erc.FindById(),
		)
		apiRouter.Delete(
			"/{resId}",
			erc.Delete(),
		)
	})
}

func NotFoundJSON() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode("Resource Not Found")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}

func PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode("Ok")
		if err != nil {
			fmt.Printf("writing response: %s", err)
		}
	}
}
