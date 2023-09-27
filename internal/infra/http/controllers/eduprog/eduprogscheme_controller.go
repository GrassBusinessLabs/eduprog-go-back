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
	"sort"

	"log"
	"net/http"
	"strconv"
)

type EduprogschemeController struct {
	eduprogschemeService app.EduprogschemeService
	eduprogcompService   app.EduprogcompService
	disciplineService    app.DisciplineService
	eduprogService       app.EduprogService
}

func NewEduprogschemeController(ess app.EduprogschemeService, ecs app.EduprogcompService, ds app.DisciplineService, es app.EduprogService) EduprogschemeController {
	return EduprogschemeController{
		eduprogschemeService: ess,
		eduprogcompService:   ecs,
		disciplineService:    ds,
		eduprogService:       es,
	}
}

func (c EduprogschemeController) SplitEduprogschemeComponent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eduprogcompIdStr := r.URL.Query().Get("eduprogcompId")
		eduprogcompId, err := strconv.ParseUint(eduprogcompIdStr, 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, fmt.Errorf("error parsing 'eduprogcompId' parameter"))
			return
		}

		semNumStr := r.URL.Query().Get("semNum")
		semNum, err := strconv.ParseUint(semNumStr, 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, fmt.Errorf("error parsing 'semNum' parameter"))
			return
		}

		eduprogscheme, err := c.eduprogschemeService.SplitEduprogschemeComponent(eduprogcompId, semNum)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogscheme[0].EduprogcompId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcompsToShow, err := c.eduprogcompService.ShowListByEduprogId(eduprogcomp.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDtoCollection(eduprogscheme, eduprogcompsToShow))
	}
}

func (c EduprogschemeController) ExpandOrShrinkComponent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eduprogcompIdStr := r.URL.Query().Get("eduprogcompId")
		eduprogcompId, err := strconv.ParseUint(eduprogcompIdStr, 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, fmt.Errorf("error parsing 'eduprogcompId' parameter"))
			return
		}

		semNumStr := r.URL.Query().Get("semNum")
		semNum, err := strconv.ParseUint(semNumStr, 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, fmt.Errorf("error parsing 'semNum' parameter"))
			return
		}

		direction := r.URL.Query().Get("direction")
		if direction != "LEFT" && direction != "RIGHT" && direction != "" {
			controllers.BadRequest(w, errors.New("direction only 'RIGHT' or LEFT"))
			return
		}

		eduprogscheme, err := c.eduprogschemeService.ExpandOrShrinkEduprogschemeComponent(eduprogcompId, semNum, direction)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogscheme[0].EduprogcompId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcompsToShow, err := c.eduprogcompService.ShowListByEduprogId(eduprogcomp.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDtoCollection(eduprogscheme, eduprogcompsToShow))
	}
}

func (c EduprogschemeController) SetComponentToEdprogscheme() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		eduprogscheme, err := requests.Bind(r, requests.SetComponentToEdprogschemeRequest{}, domain.Eduprogscheme{})
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}
		if eduprogscheme.CreditsPerSemester < 3 {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, errors.New("min 3 credits per semester"))
			return
		}

		eduprogschemes, err := c.eduprogschemeService.ShowSchemeByEduprogId(eduprogscheme.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogscheme.EduprogcompId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		discipline, err := c.disciplineService.FindById(eduprogscheme.DisciplineId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		totalCompCredits := eduprogscheme.CreditsPerSemester

		for i := range eduprogschemes {
			if eduprogschemes[i].EduprogcompId == eduprogscheme.EduprogcompId {
				totalCompCredits = totalCompCredits + eduprogschemes[i].CreditsPerSemester
				if eduprogschemes[i].SemesterNum == eduprogscheme.SemesterNum {
					log.Printf("EduprogschemeController: %s", err)
					controllers.BadRequest(w, errors.New("this component already exists in this semester"))
					return
				}
			}
		}

		if discipline.Rows < eduprogscheme.Row || eduprogscheme.Row == 0 {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, fmt.Errorf("max row num for this discipline is %d, and it cant be zero", discipline.Rows))
			return
		}

		if totalCompCredits > eduprogcomp.Credits {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, errors.New("too much credits per semester"))
			return
		}

		eduprogscheme, err = c.eduprogschemeService.SetComponentToEdprogscheme(eduprogscheme)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Created(w, eduprogschemeDto.DomainToDto(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) UpdateComponentInEduprogscheme() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogscheme, err := requests.Bind(r, requests.UpdateComponentInEduprogschemeRequest{}, domain.Eduprogscheme{})
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}
		if eduprogscheme.CreditsPerSemester < 3 {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, errors.New("min 3 credits per semester"))
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogscheme.EduprogcompId)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		discipline, err := c.disciplineService.FindById(eduprogscheme.DisciplineId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		totalCompCredits := eduprogscheme.CreditsPerSemester

		if discipline.Rows < eduprogscheme.Row || eduprogscheme.Row == 0 {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, fmt.Errorf("max row num for this discipline is %d, and it cant be zero", discipline.Rows))
			return
		}

		if totalCompCredits > eduprogcomp.Credits {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, errors.New("too much credits per semester"))
			return
		}

		eduprogscheme.Id = id
		eduprogscheme, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(eduprogscheme, id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDto(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) ExpandComponentInEduprogscheme() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		expandTo := r.URL.Query().Get("expandTo")
		if expandTo != "LEFT" && expandTo != "RIGHT" {
			controllers.BadRequest(w, errors.New("expandTo param only RIGHT or LEFT"))
			return
		}

		eduprogschemeComponent, err := c.eduprogschemeService.FindById(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.ShowSchemeByEduprogId(eduprogschemeComponent.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var schemeComponentsList []domain.Eduprogscheme

		for i := range eduprogscheme {
			if eduprogscheme[i].EduprogcompId == eduprogschemeComponent.EduprogcompId {
				schemeComponentsList = append(schemeComponentsList, eduprogscheme[i])
			}
		}

		var createdSchemeComponent domain.Eduprogscheme
		var componentWithMaxSemester domain.Eduprogscheme
		var componentWithMinSemester domain.Eduprogscheme

		if expandTo == "LEFT" { // EXPANDING SCHEME COMP TO LEFT (-semester)
			if len(schemeComponentsList) < 2 {
				if schemeComponentsList[0].SemesterNum > 1 {
					if (schemeComponentsList[0].CreditsPerSemester / 2) < 3 {
						log.Printf("EduprogschemeController: %s", err)
						controllers.BadRequest(w, errors.New("неможливо розтягнути цей компонент, бо має бути мінімум 3 кредити на семестр"))
						return
					}

					createdSchemeComponent = schemeComponentsList[0]
					schemeComponentsList[0].SemesterNum = createdSchemeComponent.SemesterNum - 1
					schemeComponentsList[0].CreditsPerSemester = createdSchemeComponent.CreditsPerSemester / 2
					createdSchemeComponent.CreditsPerSemester = createdSchemeComponent.CreditsPerSemester / 2

					if schemeComponentsList[0].CreditsPerSemester <= 0 {
						log.Printf("EduprogschemeController: %s", err)
						controllers.BadRequest(w, errors.New("too much credits"))
						return
					}

					_, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(schemeComponentsList[0], schemeComponentsList[0].Id)
					if err != nil {
						log.Printf("EduprogschemeController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}

					createdSchemeComponent, err = c.eduprogschemeService.SetComponentToEdprogscheme(createdSchemeComponent)
					if err != nil {
						log.Printf("EduprogschemeController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
				} else if schemeComponentsList[0].SemesterNum == 1 {
					log.Printf("EduprogschemeController: %s", err)
					controllers.BadRequest(w, errors.New("cannot expand to left, component is in 1st semester"))
					return
				}
			} else if len(schemeComponentsList) > 1 {
				var minSemester uint64 = 8
				var totalCreditsInScheme float64 = 0
				for i := range schemeComponentsList {
					if minSemester > schemeComponentsList[i].SemesterNum {
						minSemester = schemeComponentsList[i].SemesterNum
						componentWithMinSemester = schemeComponentsList[i]
					}
					totalCreditsInScheme = totalCreditsInScheme + schemeComponentsList[i].CreditsPerSemester
				}
				if componentWithMinSemester.SemesterNum > 1 {
					creditsInEachSemester := totalCreditsInScheme / float64(len(schemeComponentsList)+1)
					if creditsInEachSemester < 3 {
						log.Printf("EduprogschemeController: %s", err)
						controllers.BadRequest(w, errors.New("неможливо розтягнути цей компонент, бо має бути мінімум 3 кредити на семестр"))
						return
					}

					createdSchemeComponent = schemeComponentsList[len(schemeComponentsList)-1]
					createdSchemeComponent.CreditsPerSemester = creditsInEachSemester

					for _, component := range schemeComponentsList {
						component.CreditsPerSemester = creditsInEachSemester
						component.SemesterNum = component.SemesterNum - 1
						_, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(component, component.Id)
						if err != nil {
							log.Printf("EduprogschemeController: %s", err)
							controllers.InternalServerError(w, err)
							return
						}
					}

					createdSchemeComponent, err = c.eduprogschemeService.SetComponentToEdprogscheme(createdSchemeComponent)
					if err != nil {
						log.Printf("EduprogschemeController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
				} else if componentWithMinSemester.SemesterNum == 1 {
					log.Printf("EduprogschemeController: %s", err)
					controllers.BadRequest(w, errors.New("cannot expand to left, component is in 1st semester"))
					return
				}
			}
		} else if expandTo == "RIGHT" { // EXPANDING SCHEME COMP TO RIGHT (+semester)
			if len(schemeComponentsList) < 2 {
				if schemeComponentsList[0].SemesterNum < 8 {
					if (schemeComponentsList[0].CreditsPerSemester / 2) < 3 {
						log.Printf("EduprogschemeController: %s", err)
						controllers.BadRequest(w, errors.New("неможливо розтягнути цей компонент, бо має бути мінімум 3 кредити на семестр"))
						return
					}
					createdSchemeComponent = schemeComponentsList[0]
					createdSchemeComponent.CreditsPerSemester = schemeComponentsList[0].CreditsPerSemester / 2
					schemeComponentsList[0].CreditsPerSemester = schemeComponentsList[0].CreditsPerSemester - createdSchemeComponent.CreditsPerSemester
					createdSchemeComponent.SemesterNum = createdSchemeComponent.SemesterNum + 1

					if schemeComponentsList[0].CreditsPerSemester <= 0 {
						log.Printf("EduprogschemeController: %s", err)
						controllers.BadRequest(w, errors.New("too much credits"))
						return
					}

					_, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(schemeComponentsList[0], schemeComponentsList[0].Id)
					if err != nil {
						log.Printf("EduprogschemeController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}

					createdSchemeComponent, err = c.eduprogschemeService.SetComponentToEdprogscheme(createdSchemeComponent)
					if err != nil {
						log.Printf("EduprogschemeController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
				}
			} else if len(schemeComponentsList) > 1 {
				var maxSemester uint64 = 0
				var totalCreditsInScheme float64 = 0
				for i := range schemeComponentsList {
					if maxSemester < schemeComponentsList[i].SemesterNum {
						maxSemester = schemeComponentsList[i].SemesterNum
						componentWithMaxSemester = schemeComponentsList[i]
					}
					totalCreditsInScheme = totalCreditsInScheme + schemeComponentsList[i].CreditsPerSemester
				}
				if componentWithMaxSemester.SemesterNum < 8 {
					creditsInEachSemester := totalCreditsInScheme / float64(len(schemeComponentsList)+1)
					if creditsInEachSemester < 3 {
						log.Printf("EduprogschemeController: %s", err)
						controllers.BadRequest(w, errors.New("неможливо розтягнути цей компонент, бо має бути мінімум 3 кредити на семестр"))
						return
					}

					createdSchemeComponent = componentWithMaxSemester
					createdSchemeComponent.CreditsPerSemester = creditsInEachSemester
					createdSchemeComponent.SemesterNum = createdSchemeComponent.SemesterNum + 1

					for _, component := range schemeComponentsList {
						component.CreditsPerSemester = creditsInEachSemester
						_, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(component, component.Id)
						if err != nil {
							log.Printf("EduprogschemeController: %s", err)
							controllers.InternalServerError(w, err)
							return
						}
					}

					createdSchemeComponent, err = c.eduprogschemeService.SetComponentToEdprogscheme(createdSchemeComponent)
					if err != nil {
						log.Printf("EduprogschemeController: %s", err)
						controllers.InternalServerError(w, err)
						return
					}
				}
			}
		}

		eduprogschemeToShow, err := c.eduprogschemeService.ShowSchemeByEduprogId(eduprogschemeComponent.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcompsToShow, err := c.eduprogcompService.ShowListByEduprogId(eduprogschemeComponent.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDtoCollection(eduprogschemeToShow, eduprogcompsToShow))
	}
}

func (c EduprogschemeController) ShrinkComponentInEduprogscheme() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		shrinkTo := r.URL.Query().Get("shrinkTo")
		if shrinkTo != "LEFT" && shrinkTo != "RIGHT" {
			controllers.BadRequest(w, errors.New("shrinkTo param only RIGHT or LEFT"))
			return
		}

		eduprogschemeComponent, err := c.eduprogschemeService.FindById(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.ShowSchemeByEduprogId(eduprogschemeComponent.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogschemeComponent.EduprogcompId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogschemeComponent.Eduprogcomp = eduprogcomp

		var schemeComponentsList []domain.Eduprogscheme

		for i := range eduprogscheme {
			if eduprogscheme[i].EduprogcompId == eduprogschemeComponent.EduprogcompId {
				schemeComponentsList = append(schemeComponentsList, eduprogscheme[i])
			}
		}

		if len(schemeComponentsList) < 2 {
			controllers.BadRequest(w, errors.New("nothing to shrink, element is in one semester"))
			return
		}

		var componentWithMaxSemester domain.Eduprogscheme

		if shrinkTo == "LEFT" {
			if len(schemeComponentsList) == 2 {
				err = c.eduprogschemeService.Delete(schemeComponentsList[1].Id)
				if err != nil {
					log.Printf("EduprogschemeController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}

				schemeComponentsList[0].CreditsPerSemester = schemeComponentsList[0].CreditsPerSemester + schemeComponentsList[1].CreditsPerSemester

				_, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(schemeComponentsList[0], schemeComponentsList[0].Id)
				if err != nil {
					log.Printf("EduprogschemeController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}
			} else if len(schemeComponentsList) > 2 {
				var maxSemester uint64 = 0
				for i := range schemeComponentsList {
					if maxSemester < schemeComponentsList[i].SemesterNum {
						maxSemester = schemeComponentsList[i].SemesterNum
						componentWithMaxSemester = schemeComponentsList[i]
					}
				}

				err = c.eduprogschemeService.Delete(componentWithMaxSemester.Id)
				if err != nil {
					log.Printf("EduprogschemeController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}

				creditsInEachSemester := componentWithMaxSemester.CreditsPerSemester / float64(len(schemeComponentsList)-1)
				for i := range schemeComponentsList {
					if schemeComponentsList[i].Id != componentWithMaxSemester.Id {
						schemeComponentsList[i].CreditsPerSemester = schemeComponentsList[i].CreditsPerSemester + creditsInEachSemester
						_, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(schemeComponentsList[i], schemeComponentsList[i].Id)
						if err != nil {
							log.Printf("EduprogschemeController: %s", err)
							controllers.InternalServerError(w, err)
							return
						}
					}
				}

			}
		} else if shrinkTo == "RIGHT" {
			if len(schemeComponentsList) == 2 {
				err = c.eduprogschemeService.Delete(schemeComponentsList[1].Id)
				if err != nil {
					log.Printf("EduprogschemeController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}

				schemeComponentsList[0].CreditsPerSemester = schemeComponentsList[0].CreditsPerSemester + schemeComponentsList[1].CreditsPerSemester
				schemeComponentsList[0].SemesterNum = schemeComponentsList[0].SemesterNum + 1

				_, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(schemeComponentsList[0], schemeComponentsList[0].Id)
				if err != nil {
					log.Printf("EduprogschemeController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}
			} else if len(schemeComponentsList) > 2 {
				var maxSemester uint64 = 0
				for i := range schemeComponentsList {
					if maxSemester < schemeComponentsList[i].SemesterNum {
						maxSemester = schemeComponentsList[i].SemesterNum
						componentWithMaxSemester = schemeComponentsList[i]
					}
				}

				err = c.eduprogschemeService.Delete(componentWithMaxSemester.Id)
				if err != nil {
					log.Printf("EduprogschemeController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}

				creditsInEachSemester := componentWithMaxSemester.CreditsPerSemester / float64(len(schemeComponentsList)-1)
				for i := range schemeComponentsList {
					if schemeComponentsList[i].Id != componentWithMaxSemester.Id {
						schemeComponentsList[i].CreditsPerSemester = schemeComponentsList[i].CreditsPerSemester + creditsInEachSemester
						schemeComponentsList[i].SemesterNum = schemeComponentsList[i].SemesterNum + 1
						_, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(schemeComponentsList[i], schemeComponentsList[i].Id)
						if err != nil {
							log.Printf("EduprogschemeController: %s", err)
							controllers.InternalServerError(w, err)
							return
						}
					}
				}
			}
		}

		eduprogschemeToShow, err := c.eduprogschemeService.ShowSchemeByEduprogId(eduprogcomp.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcompsToShow, err := c.eduprogcompService.ShowListByEduprogId(eduprogcomp.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDtoCollection(eduprogschemeToShow, eduprogcompsToShow))
	}
}

func (c EduprogschemeController) MoveComponentInEduprogscheme() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		moveCompReq, err := requests.Bind(r, requests.MoveComponentInEduprogschemeRequest{}, domain.Eduprogscheme{})
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogschemeComponent, err := c.eduprogschemeService.FindById(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.ShowSchemeByEduprogId(eduprogschemeComponent.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var schemeComponentsList []domain.Eduprogscheme

		for i := range eduprogscheme {
			if eduprogscheme[i].EduprogcompId == eduprogschemeComponent.EduprogcompId {
				schemeComponentsList = append(schemeComponentsList, eduprogscheme[i])
			}
		}

		for i := 0; i < len(schemeComponentsList); i++ {
			schemeComponentsList[i].SemesterNum = moveCompReq.SemesterNum
			schemeComponentsList[i].DisciplineId = moveCompReq.DisciplineId
			schemeComponentsList[i].Row = moveCompReq.Row
			moveCompReq.SemesterNum++
		}

		if moveCompReq.SemesterNum <= 9 {
			for i := range schemeComponentsList {
				_, err = c.eduprogschemeService.UpdateComponentInEduprogscheme(schemeComponentsList[i], schemeComponentsList[i].Id)
				if err != nil {
					log.Printf("EduprogschemeController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}
			}
		} else {
			controllers.BadRequest(w, errors.New("cant move this schemecomp here, it comes out of semesters"))
			return
		}

		eduprogschemeToShow, err := c.eduprogschemeService.ShowSchemeByEduprogId(eduprogschemeComponent.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcompsToShow, err := c.eduprogcompService.ShowListByEduprogId(eduprogschemeComponent.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDtoCollection(eduprogschemeToShow, eduprogcompsToShow))
	}
}

func (c EduprogschemeController) DeleteFullCompFromScheme() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogschemecomp, err := c.eduprogschemeService.FindById(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogschemecomps, err := c.eduprogschemeService.ShowSchemeByEduprogId(eduprogschemecomp.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for _, eduprogscheme := range eduprogschemecomps {
			if eduprogscheme.EduprogcompId == eduprogschemecomp.EduprogcompId {
				err = c.eduprogschemeService.Delete(eduprogscheme.Id)
				if err != nil {
					log.Printf("EduprogschemeController: %s", err)
					controllers.InternalServerError(w, err)
					return
				}
			}
		}

		eduprogschemeToShow, err := c.eduprogschemeService.ShowSchemeByEduprogId(eduprogschemecomp.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcompsToShow, err := c.eduprogcompService.ShowListByEduprogId(eduprogschemecomp.EduprogId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDtoCollection(eduprogschemeToShow, eduprogcompsToShow))
	}
}

func (c EduprogschemeController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.FindById(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.FindById(eduprogscheme.EduprogcompId)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDto(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) FindBySemesterNum() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sNum, err := strconv.ParseUint(chi.URLParam(r, "sNum"), 10, 16)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		id, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.FindBySemesterNum(uint16(sNum), id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.ShowListByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDtoCollection(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) ShowSchemeByEduprogId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "sNum"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.ShowSchemeByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcomp, err := c.eduprogcompService.ShowListByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var eduprogschemeDto resources.EduprogschemeDto
		controllers.Success(w, eduprogschemeDto.DomainToDtoCollection(eduprogscheme, eduprogcomp))
	}
}

func (c EduprogschemeController) ShowFreeComponents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "sNum"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		sortOrder := r.URL.Query().Get("order")
		if sortOrder != "Az" && sortOrder != "Za" && sortOrder != "" {
			controllers.BadRequest(w, errors.New("order param only Az (alphabetic) or Za (naoborot)"))
			return
		}

		eduprogcomps, err := c.eduprogcompService.ShowListByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogscheme, err := c.eduprogschemeService.ShowSchemeByEduprogId(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for i := range eduprogcomps {
			eduprogcomps[i].FreeCredits = eduprogcomps[i].Credits
			for i2 := range eduprogscheme {
				if eduprogcomps[i].Id == eduprogscheme[i2].EduprogcompId {
					eduprogcomps[i].FreeCredits = eduprogcomps[i].FreeCredits - eduprogscheme[i2].CreditsPerSemester
				}
			}
		}

		var result []domain.Eduprogcomp

		for i := range eduprogcomps {
			if eduprogcomps[i].FreeCredits > 0 {
				result = append(result, eduprogcomps[i])
			}
		}

		if sortOrder == "Az" {
			sort.Slice(result, func(i, j int) bool {
				return result[i].Name < result[j].Name
			})
		} else if sortOrder == "Za" {
			sort.Slice(result, func(i, j int) bool {
				return result[i].Name > result[j].Name
			})
		} else {
			result = c.eduprogService.SortByCode(result)
		}

		var eduprogcompDto resources.EduprogcompDtoWithFreeCredits
		controllers.Success(w, eduprogcompDto.DomainToDtoCollectionWithFreeCredits(result))
	}
}

func (c EduprogschemeController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "essId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		err = c.eduprogschemeService.Delete(id)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		controllers.Ok(w)
	}
}
