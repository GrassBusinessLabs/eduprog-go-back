package eduprog

import (
	"errors"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/app"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/requests"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/resources"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

type EduprogController struct {
	eduprogService             app.EduprogService
	eduprogcompService         app.EduprogcompService
	eduprogcompetenciesService app.EduprogcompetenciesService
	competenciesMatrixService  app.CompetenciesMatrixService
	eduprogresultsService      app.EduprogresultsService
	resultsMatrixService       app.ResultsMatrixService
}

func NewEduprogController(es app.EduprogService, ecs app.EduprogcompService, epcs app.EduprogcompetenciesService, cms app.CompetenciesMatrixService, ers app.EduprogresultsService, rms app.ResultsMatrixService) EduprogController {
	return EduprogController{
		eduprogService:             es,
		eduprogcompService:         ecs,
		eduprogcompetenciesService: epcs,
		competenciesMatrixService:  cms,
		eduprogresultsService:      ers,
		resultsMatrixService:       rms,
	}
}

func (c EduprogController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		eduprog, err := requests.Bind(r, requests.CreateEduprogRequest{}, domain.Eduprog{})
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}
		u := r.Context().Value(controllers.UserKey).(domain.User)
		eduprog.UserId = u.Id

		levelData, err := c.eduprogService.GetOPPLevelData(eduprog.EducationLevel)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, errors.New("eduprog level error: no such level in enumeration, can use only `Початковий рівень (короткий цикл)`, `Перший (бакалаврський) рівень`, `Другий (магістерський) рівень`, `Третій (освітньо-науковий/освітньо-творчий) рівень`; get this value from method `LevelsList`"))
			return
		}
		eduprog.EducationLevel = string(levelData.Level)
		eduprog.Stage = levelData.Stage

		specialty, err := strconv.ParseUint(eduprog.Speciality, 10, 64)
		if err != nil {
			log.Printf("CompetenciesBaseController: %s", err)
			controllers.BadRequest(w, errors.New("only numeric from 11 to 293"))
			return
		}
		if specialty < 11 || specialty > 293 {
			controllers.BadRequest(w, errors.New("from 11 to 293 (currently only 141 and 121 have ZK, FK and PR)"))
			return
		}

		eduprog, err = c.eduprogService.Save(eduprog)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var eduprogDto resources.EduprogDto
		controllers.Created(w, eduprogDto.DomainToDto(eduprog))
	}
}

func (c EduprogController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprog, err := requests.Bind(r, requests.UpdateEduprogRequest{}, domain.Eduprog{})
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		u := r.Context().Value(controllers.UserKey).(domain.User)
		eduprog.UserId = u.Id

		levelData, err := c.eduprogService.GetOPPLevelData(eduprog.EducationLevel)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, errors.New("eduprog level error: no such level in enumeration, can use only `Початковий рівень (короткий цикл)`, `Перший (бакалаврський) рівень`, `Другий (магістерський) рівень`, `Третій (освітньо-науковий/освітньо-творчий) рівень`; use method LevelsList"))
			return
		}
		eduprog.EducationLevel = levelData.Level
		eduprog.Stage = levelData.Stage

		eduprog, err = c.eduprogService.Update(eduprog, id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogDto resources.EduprogDto
		controllers.Success(w, eduprogDto.DomainToDto(eduprog))
	}
}

func (c EduprogController) GetOPPLevelsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		levels, err := c.eduprogService.GetOPPLevelsList()
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogDto resources.EduprogDto
		controllers.Success(w, eduprogDto.OPPLevelDomainToDtoCollection(levels))
	}
}

func (c EduprogController) ShowList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eduprogs, err := c.eduprogService.ShowList()
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		//comps, err := c.eduprogcompService.SortComponentsByMnS()
		//if err != nil {
		//	log.Printf("EduprogController: %s", err)
		//	InternalServerError(w, err)
		//	return
		//}

		var eduprogsDto resources.EduprogDto
		controllers.Success(w, eduprogsDto.DomainToDtoCollection(eduprogs))
	}
}

func (c EduprogController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprog, err := c.eduprogService.FindById(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		comps, err := c.eduprogcompService.SortComponentsByMnS(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogDto resources.EduprogDto
		controllers.Success(w, eduprogDto.DomainToDtoWithComps(eduprog, comps))
	}
}

func (c EduprogController) CreditsInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprog, err := c.eduprogService.FindById(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		comps, err := c.eduprogcompService.SortComponentsByMnS(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		creditsDto, err := c.GetCreditsInfo(comps, eduprog.EducationLevel)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		controllers.Success(w, creditsDto)
	}
}

func (c EduprogController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "epId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		err = c.eduprogService.Delete(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		controllers.Ok(w)
	}
}

func (c EduprogController) GetCreditsInfo(comps domain.Components, edLevel string) (resources.CreditsDto, error) {
	var creditsDto resources.CreditsDto

	levelData, err := c.eduprogService.GetOPPLevelData(edLevel)
	if err != nil {
		log.Printf("EduprogController: %s", err)
		return creditsDto, err
	}

	for _, comp := range comps.Selective {
		creditsDto.SelectiveCredits += comp.Credits
	}
	for _, comp := range comps.Mandatory {
		creditsDto.MandatoryCredits += comp.Credits
	}
	creditsDto.MandatoryCreditsForLevel = levelData.MandatoryCredits
	creditsDto.SelectiveCreditsForLevel = levelData.SelectiveCredits
	creditsDto.TotalCredits = creditsDto.SelectiveCredits + creditsDto.MandatoryCredits
	creditsDto.TotalFreeCredits = (creditsDto.MandatoryCreditsForLevel + creditsDto.SelectiveCreditsForLevel) - creditsDto.TotalCredits
	creditsDto.MandatoryFreeCredits = creditsDto.MandatoryCreditsForLevel - creditsDto.MandatoryCredits
	creditsDto.SelectiveFreeCredits = creditsDto.SelectiveCreditsForLevel - creditsDto.SelectiveCredits

	return creditsDto, nil
}
