package eduprog

import (
	"errors"
	"fmt"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/app"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/requests"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/resources"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type EduprogController struct {
	eduprogService             app.EduprogService
	eduprogcompService         app.EduprogcompService
	eduprogcompetenciesService app.EduprogcompetenciesService
	competenciesMatrixService  app.CompetenciesMatrixService
	resultsMatrixService       app.ResultsMatrixService
	specialtiesService         app.SpecialtiesService
	educompRelationsService    app.EducompRelationsService
	disciplineService          app.DisciplineService
	eduprogschemeService       app.EduprogschemeService
}

func NewEduprogController(es app.EduprogService, ecs app.EduprogcompService, epcs app.EduprogcompetenciesService, cms app.CompetenciesMatrixService, rms app.ResultsMatrixService, ss app.SpecialtiesService, errs app.EducompRelationsService, ds app.DisciplineService, edss app.EduprogschemeService) EduprogController {
	return EduprogController{
		eduprogService:             es,
		eduprogcompService:         ecs,
		eduprogcompetenciesService: epcs,
		competenciesMatrixService:  cms,
		resultsMatrixService:       rms,
		specialtiesService:         ss,
		educompRelationsService:    errs,
		disciplineService:          ds,
		eduprogschemeService:       edss,
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

		maxYear := time.Now().Year() + 10
		if eduprog.ApprovalYear <= 1990 || eduprog.ApprovalYear > maxYear {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, fmt.Errorf("approval year cant be less then 1990 and greater than %d", maxYear))
			return
		}

		u := r.Context().Value(controllers.UserKey).(domain.User)
		eduprog.UserId = u.Id

		levelData, err := c.eduprogService.GetOPPLevelData(eduprog.EducationLevel)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, errors.New("eduprog level error: no such level in enumeration, can use only `Початковий рівень (короткий цикл)`, `Перший (бакалаврський) рівень`, `Другий (магістерський) рівень`, `Третій (освітньо-науковий/освітньо-творчий) рівень`; get this value from method `LevelsList`"))
			return
		}
		eduprog.EducationLevel = levelData.Level
		eduprog.Stage = levelData.Stage

		allSpecialties, err := c.specialtiesService.ShowAllSpecialties()
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		check := false
		for i := range allSpecialties {
			if allSpecialties[i].Code == eduprog.SpecialtyCode {
				check = true
				eduprog.Speciality = allSpecialties[i].Name
				eduprog.KFCode = allSpecialties[i].KFCode
				eduprog.KnowledgeField = allSpecialties[i].KnowledgeField
			}
		}
		if !check {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, errors.New("there is no such specialty in enum, only values from `ShowAllSpecialties` can be used"))
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

func (c EduprogController) CreateDuplicateOf() http.HandlerFunc {
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
			controllers.InternalServerError(w, err)
			return
		}

		eduprog.ChildOf = eduprog.Id
		eduprog.Name = eduprog.Name + " [КОПІЯ]"
		eduprog, err = c.eduprogService.Save(eduprog)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcomps, err := c.eduprogcompService.ShowListByEduprogId(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for _, eduprogcomp := range eduprogcomps {
			eduprogcomp.EduprogId = eduprog.Id
			_, err = c.eduprogcompService.Save(eduprogcomp)
			if err != nil {
				log.Printf("EduprogController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
		}

		eduprogcompetenices, err := c.eduprogcompetenciesService.ShowCompetenciesByEduprogId(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for _, eduprogcompetency := range eduprogcompetenices {
			eduprogcompetency.EduprogId = eduprog.Id
			_, err = c.eduprogcompetenciesService.AddCompetencyToEduprog(eduprogcompetency)
			if err != nil {
				log.Printf("EduprogController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
		}

		disciplines, err := c.disciplineService.ShowDisciplinesByEduprogId(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for _, discipline := range disciplines {
			discipline.EduprogId = eduprog.Id
			_, err = c.disciplineService.Save(discipline)
			if err != nil {
				log.Printf("EduprogController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
		}

		eduprogscheme, err := c.eduprogschemeService.ShowSchemeByEduprogId(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		copiedDisciplines, err := c.disciplineService.ShowDisciplinesByEduprogId(eduprog.Id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		copiedEduprogcomps, err := c.eduprogcompService.ShowListByEduprogId(eduprog.Id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for _, schemeElement := range eduprogscheme {
			schemeElement.EduprogId = eduprog.Id

			for _, copiedDiscipline := range copiedDisciplines {
				schemeDiscipline, err := c.disciplineService.FindById(schemeElement.DisciplineId)
				if err != nil {
					log.Printf("EduprogController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}
				if copiedDiscipline.Name == schemeDiscipline.Name {
					schemeElement.DisciplineId = copiedDiscipline.Id
				}
			}

			for _, copiedEduprogcomp := range copiedEduprogcomps {
				schemeComp, err := c.eduprogcompService.FindById(schemeElement.EduprogcompId)
				if err != nil {
					log.Printf("EduprogController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}
				if copiedEduprogcomp.Name == schemeComp.Name {
					schemeElement.EduprogcompId = copiedEduprogcomp.Id
				}
			}

			_, err = c.eduprogschemeService.SetComponentToEdprogscheme(schemeElement)
			if err != nil {
				log.Printf("EduprogController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
		}

		competenciesMatrix, err := c.competenciesMatrixService.ShowByEduprogId(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		copiedCompetencies, err := c.eduprogcompetenciesService.ShowCompetenciesByEduprogId(eduprog.Id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for _, matrix := range competenciesMatrix {
			matrix.EduprogId = eduprog.Id

			for _, eduprogcomp := range copiedEduprogcomps {
				matrixComp, err := c.eduprogcompService.FindById(matrix.ComponentId)
				if err != nil {
					log.Printf("EduprogController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}
				if matrixComp.Name == eduprogcomp.Name {
					matrix.ComponentId = eduprogcomp.Id
				}
			}

			for _, copiedCompetency := range copiedCompetencies {
				matrixCompetency, err := c.eduprogcompetenciesService.FindById(matrix.CompetencyId)
				if err != nil {
					log.Printf("EduprogController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}
				if matrixCompetency.CompetencyId == copiedCompetency.CompetencyId {
					matrix.CompetencyId = copiedCompetency.Id
				}
			}

			_, err = c.competenciesMatrixService.CreateRelation(matrix)
			if err != nil {
				log.Printf("EduprogController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
		}

		resultsMatrix, err := c.resultsMatrixService.ShowByEduprogId(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for _, matrix := range resultsMatrix {
			matrix.EduprogId = eduprog.Id

			for _, eduprogcomp := range copiedEduprogcomps {
				matrixComp, err := c.eduprogcompService.FindById(matrix.ComponentId)
				if err != nil {
					log.Printf("EduprogController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}
				if matrixComp.Name == eduprogcomp.Name {
					matrix.ComponentId = eduprogcomp.Id
				}
			}

			for _, copiedCompetency := range copiedCompetencies {
				matrixCompetency, err := c.eduprogcompetenciesService.FindById(matrix.EduprogresultId)
				if err != nil {
					log.Printf("EduprogController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}
				if matrixCompetency.CompetencyId == copiedCompetency.CompetencyId {
					matrix.EduprogresultId = copiedCompetency.Id
				}
			}

			_, err = c.resultsMatrixService.CreateRelation(matrix)
			if err != nil {
				log.Printf("EduprogController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
		}

		educompRelations, err := c.educompRelationsService.ShowByEduprogId(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for _, relation := range educompRelations {
			relation.EduprogId = eduprog.Id

			for _, eduprogcomp := range copiedEduprogcomps {
				baseComp, err := c.eduprogcompService.FindById(relation.BaseCompId)
				if err != nil {
					log.Printf("EduprogController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}
				if baseComp.Name == eduprogcomp.Name {
					relation.BaseCompId = eduprogcomp.Id
				}
			}

			for _, eduprogcomp := range copiedEduprogcomps {
				childComp, err := c.eduprogcompService.FindById(relation.ChildCompId)
				if err != nil {
					log.Printf("EduprogController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}
				if childComp.Name == eduprogcomp.Name {
					relation.ChildCompId = eduprogcomp.Id
				}
			}

			_, err = c.educompRelationsService.CreateRelation(relation)
			if err != nil {
				log.Printf("EduprogController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
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

		allSpecialties, err := c.specialtiesService.ShowAllSpecialties()
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		check := false
		for i := range allSpecialties {
			if allSpecialties[i].Code == eduprog.SpecialtyCode {
				check = true
				eduprog.Speciality = allSpecialties[i].Name
				eduprog.KFCode = allSpecialties[i].KFCode
				eduprog.KnowledgeField = allSpecialties[i].KnowledgeField
			}
		}
		if !check {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, errors.New("there is no such specialty in enum, only values from `ShowAllSpecialties` can be used"))
			return
		}

		eduprog.Id = id
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

		comps.Mandatory = sortByCode(comps.Mandatory)

		var eduprogDto resources.EduprogDto
		controllers.Success(w, eduprogDto.DomainToDtoWithComps(eduprog, comps, comps.Selective))
		//controllers.Success(w, eduprogDto.DomainToDtoWithComps(eduprog, comps))
	}
}

func sortByCode(eduprogcomps []domain.Eduprogcomp) []domain.Eduprogcomp {
	sort.Slice(eduprogcomps, func(i, j int) bool {
		// Parse the Code field as integers and compare them
		codeI, errI := strconv.ParseUint(eduprogcomps[i].Code, 10, 64)
		codeJ, errJ := strconv.ParseUint(eduprogcomps[j].Code, 10, 64)
		if errI != nil || errJ != nil {
			return eduprogcomps[i].Code < eduprogcomps[j].Code
		}
		return codeI < codeJ
	})
	return eduprogcomps
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

	for i := range comps.Selective {
		for _, comp := range comps.Selective[i].CompsInBlock {
			creditsDto.SelectiveCredits += comp.Credits
		}

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
