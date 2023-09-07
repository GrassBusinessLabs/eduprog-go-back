package app

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/domain"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/database/eduprog"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/resources"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/nguyenthenguyen/docx"
	"github.com/xuri/excelize/v2"
	"log"
	"sort"
	"strconv"
	"time"
)

const (
	SheetName1 = "Перелік компонент"
	SheetName2 = "Матриця компетентностей"
	SheetName3 = "Матриця відповідності ПР"
)

type EduprogService interface {
	Save(eduprog domain.Eduprog, userId uint64) (domain.Eduprog, error)
	Update(ref, req domain.Eduprog) (domain.Eduprog, error)
	CreateDuplicateOf(id, userId uint64, name string, approvalYear int) (domain.Eduprog, error)
	ShowList() (domain.Eduprogs, error)
	FindById(id uint64) (domain.Eduprog, domain.Components, error)
	GetOPPLevelsList() ([]domain.OPPLevelStruct, error)
	GetOPPLevelData(level string) (domain.OPPLevelStruct, error)
	Delete(id uint64) error
	SortByCode(eduprogcomps []domain.Eduprogcomp) []domain.Eduprogcomp
	GetCreditsInfo(eduprogId uint64) (resources.CreditsDto, error)
	ExportEduprogToWord(eduprogId uint64) error
	ExportEduprogToExcel(eduprogId uint64) (string, *bytes.Buffer, error)
	ExportEducompRealtionsToJpg(eduprogId uint64) (string, error)
}

type eduprogService struct {
	eduprogRepo                eduprog.EduprogRepository
	specialtiesService         SpecialtiesService
	eduprogcompService         EduprogcompService
	eduprogcompetenciesService EduprogcompetenciesService
	disciplineService          DisciplineService
	eduprogschemeService       EduprogschemeService
	competenciesMatrixService  CompetenciesMatrixService
	resultsMatrixService       ResultsMatrixService
	educompRelationsService    EducompRelationsService
}

func NewEduprogService(
	er eduprog.EduprogRepository,
	ss SpecialtiesService,
	es EduprogcompService,
	ecs EduprogcompetenciesService,
	ds DisciplineService,
	ess EduprogschemeService,
	cms CompetenciesMatrixService,
	rms ResultsMatrixService,
	ers EducompRelationsService) EduprogService {
	return eduprogService{
		eduprogRepo:                er,
		eduprogcompService:         es,
		specialtiesService:         ss,
		eduprogcompetenciesService: ecs,
		disciplineService:          ds,
		eduprogschemeService:       ess,
		competenciesMatrixService:  cms,
		resultsMatrixService:       rms,
		educompRelationsService:    ers,
	}
}

func (s eduprogService) GetCreditsInfo(eduprogId uint64) (resources.CreditsDto, error) {
	e, _, err := s.FindById(eduprogId)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return resources.CreditsDto{}, err
	}

	creditsDto, err := s.eduprogcompService.GetCreditsInfo(e)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return resources.CreditsDto{}, err
	}

	return creditsDto, nil
}

func (s eduprogService) Save(eduprog domain.Eduprog, userId uint64) (domain.Eduprog, error) {
	var err error

	maxYear := time.Now().Year() + 10
	if eduprog.ApprovalYear <= 1990 || eduprog.ApprovalYear > maxYear {
		log.Printf("EduprogService: %s", fmt.Errorf("approval year cant be less then 1990 and greater than %d", maxYear))
		return domain.Eduprog{}, fmt.Errorf("approval year cant be less then 1990 and greater than %d", maxYear)
	}

	eduprog.UserId = userId

	levelData, err := s.GetOPPLevelData(eduprog.EducationLevel)
	if err != nil {
		log.Printf("EduprogService: %s", errors.New("eduprog level error: no such level in enumeration, can use only `Початковий рівень (короткий цикл)`, `Перший (бакалаврський) рівень`, `Другий (магістерський) рівень`, `Третій (освітньо-науковий/освітньо-творчий) рівень`; get this value from method `LevelsList`"))
		return domain.Eduprog{}, errors.New("eduprog level error: no such level in enumeration, can use only `Початковий рівень (короткий цикл)`, `Перший (бакалаврський) рівень`, `Другий (магістерський) рівень`, `Третій (освітньо-науковий/освітньо-творчий) рівень`; get this value from method `LevelsList`")
	}

	eduprog.EducationLevel = levelData.Level
	eduprog.Stage = levelData.Stage

	allSpecialties, err := s.specialtiesService.ShowAllSpecialties()
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
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
		log.Printf("EduprogService: %s", errors.New("there is no such specialty in enum, only values from `ShowAllSpecialties` can be used"))
		return domain.Eduprog{}, errors.New("there is no such specialty in enum, only values from `ShowAllSpecialties` can be used")
	}

	eduprog, err = s.eduprogRepo.Save(eduprog)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}

	return eduprog, nil
}

func (s eduprogService) Update(ref, req domain.Eduprog) (domain.Eduprog, error) {
	if req.Name != "" && req.Name != ref.Name {
		ref.Name = req.Name
	}
	if req.EducationLevel != "" && req.EducationLevel != ref.EducationLevel {
		levelData, err := s.GetOPPLevelData(req.EducationLevel)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return domain.Eduprog{}, fmt.Errorf("eduprog level error: no such level in enumeration, can use only `Початковий рівень (короткий цикл)`, `Перший (бакалаврський) рівень`, `Другий (магістерський) рівень`, `Третій (освітньо-науковий/освітньо-творчий) рівень`; use method LevelsList")
		}
		ref.EducationLevel = levelData.Level
		ref.Stage = levelData.Stage
	}
	if req.SpecialtyCode != "" && req.SpecialtyCode != ref.SpecialtyCode {
		allSpecialties, err := s.specialtiesService.ShowAllSpecialties()
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return domain.Eduprog{}, err
		}
		check := false
		for i := range allSpecialties {
			if allSpecialties[i].Code == req.SpecialtyCode {
				check = true
				req.Speciality = allSpecialties[i].Name
				req.KFCode = allSpecialties[i].KFCode
				req.KnowledgeField = allSpecialties[i].KnowledgeField
			}
		}
		if !check {
			log.Printf("EduprogService: %s", err)
			return domain.Eduprog{}, fmt.Errorf("there is no such specialty in enum, only values from `ShowAllSpecialties` can be used")
		}
	}
	if req.ApprovalYear != 0 && req.ApprovalYear != ref.ApprovalYear {
		ref.ApprovalYear = req.ApprovalYear
	}

	e, err := s.eduprogRepo.Update(ref)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}

	return e, err
}

func (s eduprogService) CreateDuplicateOf(id, userId uint64, name string, approvalYear int) (domain.Eduprog, error) {
	e, _, err := s.FindById(id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}
	e.ChildOf = e.Id
	e.Name = name
	e.ApprovalYear = approvalYear
	e, err = s.Save(e, userId)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}

	eduprogcomps, err := s.eduprogcompService.ShowListByEduprogId(id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}

	for _, eduprogcomp := range eduprogcomps {
		eduprogcomp.EduprogId = e.Id
		_, err = s.eduprogcompService.Save(eduprogcomp)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return domain.Eduprog{}, err
		}
	}

	eduprogcompetenices, err := s.eduprogcompetenciesService.ShowCompetenciesByEduprogId(id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}

	for _, eduprogcompetency := range eduprogcompetenices {
		eduprogcompetency.EduprogId = e.Id
		_, err = s.eduprogcompetenciesService.AddCompetencyToEduprog(eduprogcompetency)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return domain.Eduprog{}, err
		}
	}

	disciplines, err := s.disciplineService.ShowDisciplinesByEduprogId(id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}

	for _, discipline := range disciplines {
		discipline.EduprogId = e.Id
		_, err = s.disciplineService.Save(discipline)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return domain.Eduprog{}, err
		}
	}

	eduprogscheme, err := s.eduprogschemeService.ShowSchemeByEduprogId(id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}

	copiedDisciplines, err := s.disciplineService.ShowDisciplinesByEduprogId(e.Id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}

	copiedEduprogcomps, err := s.eduprogcompService.ShowListByEduprogId(e.Id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}

	for _, schemeElement := range eduprogscheme {
		schemeElement.EduprogId = e.Id

		for _, copiedDiscipline := range copiedDisciplines {
			schemeDiscipline, err := s.disciplineService.FindById(schemeElement.DisciplineId)
			if err != nil {
				log.Printf("EduprogService: %s", err)
				return domain.Eduprog{}, err
			}
			if copiedDiscipline.Name == schemeDiscipline.Name {
				schemeElement.DisciplineId = copiedDiscipline.Id
			}
		}

		for _, copiedEduprogcomp := range copiedEduprogcomps {
			schemeComp, err := s.eduprogcompService.FindById(schemeElement.EduprogcompId)
			if err != nil {
				log.Printf("EduprogService: %s", err)
				return domain.Eduprog{}, err
			}
			if copiedEduprogcomp.Name == schemeComp.Name {
				schemeElement.EduprogcompId = copiedEduprogcomp.Id
			}
		}

		_, err = s.eduprogschemeService.SetComponentToEdprogscheme(schemeElement)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return domain.Eduprog{}, err
		}
	}

	competenciesMatrix, err := s.competenciesMatrixService.ShowByEduprogId(id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}

	copiedCompetencies, err := s.eduprogcompetenciesService.ShowCompetenciesByEduprogId(e.Id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}

	for _, matrix := range competenciesMatrix {
		matrix.EduprogId = e.Id

		for _, eduprogcomp := range copiedEduprogcomps {
			matrixComp, err := s.eduprogcompService.FindById(matrix.ComponentId)
			if err != nil {
				log.Printf("EduprogService: %s", err)
				return domain.Eduprog{}, err
			}
			if matrixComp.Name == eduprogcomp.Name {
				matrix.ComponentId = eduprogcomp.Id
			}
		}

		for _, copiedCompetency := range copiedCompetencies {
			matrixCompetency, err := s.eduprogcompetenciesService.FindById(matrix.CompetencyId)
			if err != nil {
				log.Printf("EduprogService: %s", err)
				return domain.Eduprog{}, err
			}
			if matrixCompetency.CompetencyId == copiedCompetency.CompetencyId {
				matrix.CompetencyId = copiedCompetency.Id
			}
		}

		_, err = s.competenciesMatrixService.CreateRelation(matrix)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return domain.Eduprog{}, err
		}
	}

	resultsMatrix, err := s.resultsMatrixService.ShowByEduprogId(id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}

	for _, matrix := range resultsMatrix {
		matrix.EduprogId = e.Id

		for _, eduprogcomp := range copiedEduprogcomps {
			matrixComp, err := s.eduprogcompService.FindById(matrix.ComponentId)
			if err != nil {
				log.Printf("EduprogService: %s", err)
				return domain.Eduprog{}, err
			}
			if matrixComp.Name == eduprogcomp.Name {
				matrix.ComponentId = eduprogcomp.Id
			}
		}

		for _, copiedCompetency := range copiedCompetencies {
			matrixCompetency, err := s.eduprogcompetenciesService.FindById(matrix.EduprogresultId)
			if err != nil {
				log.Printf("EduprogService: %s", err)
				return domain.Eduprog{}, err
			}
			if matrixCompetency.CompetencyId == copiedCompetency.CompetencyId {
				matrix.EduprogresultId = copiedCompetency.Id
			}
		}

		_, err = s.resultsMatrixService.CreateRelation(matrix)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return domain.Eduprog{}, err
		}
	}

	educompRelations, err := s.educompRelationsService.ShowByEduprogId(id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, err
	}

	for _, relation := range educompRelations {
		relation.EduprogId = e.Id

		for _, eduprogcomp := range copiedEduprogcomps {
			baseComp, err := s.eduprogcompService.FindById(relation.BaseCompId)
			if err != nil {
				log.Printf("EduprogService: %s", err)
				return domain.Eduprog{}, err
			}
			if baseComp.Name == eduprogcomp.Name {
				relation.BaseCompId = eduprogcomp.Id
			}
		}

		for _, eduprogcomp := range copiedEduprogcomps {
			childComp, err := s.eduprogcompService.FindById(relation.ChildCompId)
			if err != nil {
				log.Printf("EduprogService: %s", err)
				return domain.Eduprog{}, err
			}
			if childComp.Name == eduprogcomp.Name {
				relation.ChildCompId = eduprogcomp.Id
			}
		}

		_, err = s.educompRelationsService.CreateRelation(relation)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return domain.Eduprog{}, err
		}
	}

	return e, err
}

func (s eduprogService) ShowList() (domain.Eduprogs, error) {
	e, err := s.eduprogRepo.ShowList()
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprogs{}, err
	}
	return e, nil
}

func (s eduprogService) FindById(id uint64) (domain.Eduprog, domain.Components, error) {
	e, err := s.eduprogRepo.FindById(id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, domain.Components{}, err
	}
	comps, err := s.eduprogcompService.SortComponentsByMnS(id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.Eduprog{}, domain.Components{}, err
	}

	return e, comps, nil
}

func (s eduprogService) GetOPPLevelsList() ([]domain.OPPLevelStruct, error) {
	e, err := s.eduprogRepo.GetOPPLevelsList()
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return []domain.OPPLevelStruct{}, err
	}
	return e, nil
}

func (s eduprogService) GetOPPLevelData(level string) (domain.OPPLevelStruct, error) {
	e, err := s.eduprogRepo.GetOPPLevelData(level)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return domain.OPPLevelStruct{}, err
	}
	return e, nil
}

func (s eduprogService) Delete(id uint64) error {
	err := s.eduprogRepo.Delete(id)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return err
	}
	return nil
}

func (s eduprogService) SortByCode(eduprogcomps []domain.Eduprogcomp) []domain.Eduprogcomp {
	sort.Slice(eduprogcomps, func(i, j int) bool {
		codeI, errI := strconv.ParseUint(eduprogcomps[i].Code, 10, 64)
		codeJ, errJ := strconv.ParseUint(eduprogcomps[j].Code, 10, 64)
		if errI != nil || errJ != nil {
			return eduprogcomps[i].Code < eduprogcomps[j].Code
		}
		return codeI < codeJ
	})
	return eduprogcomps
}

func (s eduprogService) ExportEduprogToWord(eduprogId uint64) error {
	e, _, err := s.FindById(eduprogId)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return err
	}

	d, err := docx.ReadDocxFile("./opp_template.docx")

	if err != nil {
		panic(err)
	}
	docx1 := d.Editable()

	err = docx1.Replace("oppname", e.Name, -1)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return err
	}
	err = docx1.Replace("level", e.EducationLevel, -1)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return err
	}
	err = docx1.Replace("stupin", e.Stage, -1)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return err
	}
	err = docx1.Replace("specialty", e.SpecialtyCode+" "+e.Speciality, -1)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return err
	}
	err = docx1.Replace("galuz", e.KFCode+" "+e.KnowledgeField, -1)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return err
	}
	err = docx1.Replace("year", strconv.Itoa(e.ApprovalYear), -1)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return err
	}

	err = docx1.WriteToFile("OPP.docx")
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return err
	}

	err = d.Close()
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return err
	}

	return nil
}

func (s eduprogService) ExportEducompRealtionsToJpg(eduprogId uint64) (string, error) {
	e, _, err := s.FindById(eduprogId)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return "", err
	}

	relationships, err := s.educompRelationsService.ShowByEduprogId(eduprogId)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return "", err
	}

	g := graphviz.New()

	graph, err := g.Graph(graphviz.Name("G"))
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return "", err
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Printf("EduprogService: %s", err)
		}
		err := g.Close()
		if err != nil {
			return
		}
	}()

	nodeAttrs := make(map[string]string)
	nodeAttrs["shape"] = "box"
	var nodeName string

	// Add nodes for base components and child components
	for _, r := range relationships {
		baseedcomp, _ := s.eduprogcompService.FindById(r.BaseCompId)
		childedcomp, _ := s.eduprogcompService.FindById(r.ChildCompId)

		if baseedcomp.Type == "ОК" {
			nodeName = fmt.Sprintf("%s%s", baseedcomp.Type, baseedcomp.Code)
		} else if baseedcomp.Type == "ВБ" {
			nodeName = fmt.Sprintf("Блок%s", baseedcomp.BlockNum)
		}
		baseNode, err := graph.CreateNode(nodeName)
		baseNode.SetShape(cgraph.SquareShape)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return "", err
		}
		if childedcomp.Type == "ОК" {
			nodeName = fmt.Sprintf("%s%s", childedcomp.Type, childedcomp.Code)
		} else if childedcomp.Type == "ВБ" {
			nodeName = fmt.Sprintf("Блок%s", childedcomp.BlockNum)
		}
		childNode, err := graph.CreateNode(nodeName)
		childNode.SetShape(cgraph.SquareShape)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return "", err
		}
	}

	// Add edges between base components and child components
	for _, r := range relationships {
		baseedcomp, _ := s.eduprogcompService.FindById(r.BaseCompId)
		childedcomp, _ := s.eduprogcompService.FindById(r.ChildCompId)

		if baseedcomp.Type == "ОК" {
			nodeName = fmt.Sprintf("%s%s", baseedcomp.Type, baseedcomp.Code)
		} else if baseedcomp.Type == "ВБ" {
			nodeName = fmt.Sprintf("Блок%s", baseedcomp.BlockNum)
		}
		baseCompNode, err := graph.Node(nodeName)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return "", err
		}

		if childedcomp.Type == "ОК" {
			nodeName = fmt.Sprintf("%s%s", childedcomp.Type, childedcomp.Code)
		} else if childedcomp.Type == "ВБ" {
			nodeName = fmt.Sprintf("Блок%s", childedcomp.BlockNum)
		}
		childCompNode, err := graph.Node(nodeName)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return "", err
		}

		_, err = graph.CreateEdge("", baseCompNode, childCompNode)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return "", err
		}
	}
	filename := fmt.Sprintf("%s_схема.png", e.Name)
	if err := g.RenderFilename(graph, graphviz.PNG, filename); err != nil {
		log.Printf("EduprogService: %s", err)
		return "", err
	}

	return filename, nil
}

func (s eduprogService) ExportEduprogToExcel(eduprogId uint64) (string, *bytes.Buffer, error) {
	e, _, err := s.FindById(eduprogId)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return "", nil, err
	}

	eduprogcomps, err := s.eduprogcompService.SortComponentsByMnS(eduprogId)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return "", nil, err
	}

	creditsDto, err := s.eduprogcompService.GetCreditsInfo(e)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return "", nil, err
	}

	//------------------------EXPORT EDUPROGCOMPS LOGIC-------------------------------//

	xlsx := excelize.NewFile()
	index, _ := xlsx.NewSheet("Sheet1")
	index2, _ := xlsx.NewSheet("Sheet2")
	index3, _ := xlsx.NewSheet("Sheet3")
	xlsx.SetActiveSheet(index)
	_ = xlsx.SetSheetName("Sheet1", SheetName1)

	style, _ := xlsx.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 12, Family: "Times New Roman"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "left", Color: "#000000", Style: 1},
		},
	})
	styleAlignLeft, _ := xlsx.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 12, Family: "Times New Roman"},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
		Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "left", Color: "#000000", Style: 1},
		},
	})
	styleBold, _ := xlsx.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 12, Bold: true, Family: "Times New Roman"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "left", Color: "#000000", Style: 1},
		},
	})
	styleItalic, _ := xlsx.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 12, Italic: true, Family: "Times New Roman"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "left", Color: "#000000", Style: 1},
		},
	})
	styleBoldAlignLeft, _ := xlsx.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 12, Bold: true, Family: "Times New Roman"},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
		Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "left", Color: "#000000", Style: 1},
		},
	})
	styleError, _ := xlsx.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 16, Bold: true, Family: "Times New Roman", Color: "#FF0000"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "left", Color: "#000000", Style: 1},
		},
	})
	_ = xlsx.SetCellStyle(SheetName1, "A1", "D3", style)
	_ = xlsx.MergeCell(SheetName1, "A3", "D3")
	_ = xlsx.SetColWidth(SheetName1, "A", "A", 10)
	_ = xlsx.SetColWidth(SheetName1, "B", "B", 50)
	_ = xlsx.SetColWidth(SheetName1, "C", "C", 15)
	_ = xlsx.SetColWidth(SheetName1, "D", "D", 20)

	data := [][]interface{}{
		{"Код н/д", "Компоненти освітньої програми (навчальні дисципліни, курсові проекти (роботи), практики, кваліфікаційна робота)", "Кількість кредитів", "Форма підсумкового контролю"},
		{1, 2, 3, 4},
		{"Обов'язкові компоненти ОП"},
	}

	_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", 3), fmt.Sprintf("D%d", 3), styleBold)
	startRow := 1

	for i := startRow; i < len(data)+startRow; i++ {

		_ = xlsx.SetSheetRow(SheetName1, fmt.Sprintf("A%d", i), &data[i-1])

	}

	mandLen := len(eduprogcomps.Mandatory)

	for i := 4; i < mandLen+4; i++ {

		_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
		_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("B%d", i), fmt.Sprintf("B%d", i), styleAlignLeft)

		_ = xlsx.SetSheetRow(SheetName1, fmt.Sprintf("A%d", i), &[]interface{}{
			eduprogcomps.Mandatory[i-4].Type + " " + eduprogcomps.Mandatory[i-4].Code + ".",
			eduprogcomps.Mandatory[i-4].Name,
			eduprogcomps.Mandatory[i-4].Credits,
			eduprogcomps.Mandatory[i-4].ControlType,
		})

	}

	_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("A%d", mandLen+4), fmt.Sprintf("B%d", mandLen+4))
	_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("C%d", mandLen+4), fmt.Sprintf("D%d", mandLen+4))
	_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("C%d", mandLen+4), fmt.Sprintf("D%d", mandLen+4), styleBold)
	_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", mandLen+4), fmt.Sprintf("B%d", mandLen+4), styleBoldAlignLeft)
	_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("A%d", mandLen+4), "Загальний обсяг обов'язкових компонент: ")
	_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("C%d", mandLen+4), fmt.Sprintf("%.1f кредитів", creditsDto.MandatoryCredits))
	_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("A%d", mandLen+5), fmt.Sprintf("D%d", mandLen+5))
	_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("A%d", mandLen+5), "Вибіркові компоненти ОП")
	_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", mandLen+5), fmt.Sprintf("D%d", mandLen+5), styleBold)

	blocksInfo := eduprogcomps.Selective
	var selective []domain.Eduprogcomp
	for i := range eduprogcomps.Selective {
		selective = append(selective, eduprogcomps.Selective[i].CompsInBlock...)
	}
	selLen := len(selective)
	blocksInfoLen := len(blocksInfo)
	var temp = 0
	for i := mandLen + 6; i < blocksInfoLen+selLen+mandLen+6; i++ {
		_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i))
		_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("A%d", i), fmt.Sprintf("Вибірковий блок %s (%s)", blocksInfo[temp].BlockNum, blocksInfo[temp].BlockName))
		_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), styleItalic)
		for _, comps := range blocksInfo[temp].CompsInBlock {
			i++
			_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
			_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("B%d", i), fmt.Sprintf("B%d", i), styleAlignLeft)
			_ = xlsx.SetSheetRow(SheetName1, fmt.Sprintf("A%d", i), &[]interface{}{
				comps.Type + " " + comps.BlockNum + "." + comps.Code + ".",
				comps.Name,
				comps.Credits,
				comps.ControlType,
			})

		}
		temp++
	}

	_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("A%d", blocksInfoLen+selLen+mandLen+6), fmt.Sprintf("B%d", blocksInfoLen+selLen+mandLen+6))
	_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("C%d", blocksInfoLen+selLen+mandLen+6), fmt.Sprintf("D%d", blocksInfoLen+selLen+mandLen+6))
	_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("A%d", blocksInfoLen+selLen+mandLen+7), fmt.Sprintf("B%d", blocksInfoLen+selLen+mandLen+7))
	_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("C%d", blocksInfoLen+selLen+mandLen+7), fmt.Sprintf("D%d", blocksInfoLen+selLen+mandLen+7))
	_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", blocksInfoLen+selLen+mandLen+6), fmt.Sprintf("B%d", blocksInfoLen+selLen+mandLen+6), styleBoldAlignLeft)
	_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("C%d", blocksInfoLen+selLen+mandLen+6), fmt.Sprintf("D%d", blocksInfoLen+selLen+mandLen+6), styleBold)
	_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("A%d", blocksInfoLen+selLen+mandLen+6), "Загальний обсяг вибіркових компонент: ")
	_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("C%d", blocksInfoLen+selLen+mandLen+6), fmt.Sprintf("%.1f кредитів", creditsDto.SelectiveCredits))
	_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", blocksInfoLen+selLen+mandLen+7), fmt.Sprintf("B%d", blocksInfoLen+selLen+mandLen+7), styleBoldAlignLeft)
	_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("C%d", blocksInfoLen+selLen+mandLen+7), fmt.Sprintf("D%d", blocksInfoLen+selLen+mandLen+7), styleBold)
	_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("A%d", blocksInfoLen+selLen+mandLen+7), "ЗАГАЛЬНИЙ ОБСЯГ ОСВІТНЬОЇ ПРОГРАМИ: ")
	_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("C%d", blocksInfoLen+selLen+mandLen+7), fmt.Sprintf("%.1f кредитів", creditsDto.TotalCredits))

	//----------------------------EXPORT COMPETENCIES MATRIX LOGIC----------------------------------//

	eduprogcompetenciesZK, _ := s.eduprogcompetenciesService.ShowCompetenciesByType(eduprogId, "ZK")
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return "", nil, err
	}
	eduprogcompetenciesFK, _ := s.eduprogcompetenciesService.ShowCompetenciesByType(eduprogId, "FK")
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return "", nil, err
	}
	eduprogcompetenciesPR, _ := s.eduprogcompetenciesService.ShowCompetenciesByType(eduprogId, "PR")
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return "", nil, err
	}

	xlsx.SetActiveSheet(index2)
	err = xlsx.SetSheetName("Sheet2", SheetName2)

	styleRotated, _ := xlsx.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 12, Family: "Times New Roman", Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true, TextRotation: 90},
		Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "left", Color: "#000000", Style: 1},
		},
	})

	styleDot, _ := xlsx.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 24, Family: "Times New Roman", Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
			{Type: "top", Color: "#000000", Style: 1},
			{Type: "bottom", Color: "#000000", Style: 1},
			{Type: "right", Color: "#000000", Style: 1},
			{Type: "left", Color: "#000000", Style: 1},
		},
	})

	mandLen = len(eduprogcomps.Mandatory)

	lastLetter := ""
	bufLetter := ""
	_ = xlsx.SetRowHeight(SheetName2, 1, 40)
	for i := 66; i < mandLen+66; i++ {
		if i <= 90 {
			_ = xlsx.SetCellStyle(SheetName2, fmt.Sprintf("%s1", string(rune(i))), fmt.Sprintf("%s1", string(rune(i))), styleRotated)
			_ = xlsx.SetColWidth(SheetName2, string(rune(i)), string(rune(i)), 3)

			_ = xlsx.SetSheetCol(SheetName2, fmt.Sprintf("%s1", string(rune(i))), &[]interface{}{
				eduprogcomps.Mandatory[i-66].Type + " " + eduprogcomps.Mandatory[i-66].Code,
			})
			lastLetter = string(rune(i))
		} else if i > 90 && i <= 116 {
			bufLetter = string(rune(65))
			_ = xlsx.SetCellStyle(SheetName2, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), styleRotated)
			_ = xlsx.SetColWidth(SheetName2, fmt.Sprintf("%s%s", bufLetter, string(rune(i-26))), fmt.Sprintf("%s%s", bufLetter, string(rune(i-26))), 3)

			_ = xlsx.SetSheetCol(SheetName2, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), &[]interface{}{
				eduprogcomps.Mandatory[i-66].Type + " " + eduprogcomps.Mandatory[i-66].Code,
			})
			lastLetter = fmt.Sprintf("%s%s", bufLetter, string(rune(i-26)))
		} else if i > 116 && i <= 142 {
			bufLetter = string(rune(66))
			//	_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
			_ = xlsx.SetCellStyle(SheetName2, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), styleRotated)
			_ = xlsx.SetColWidth(SheetName2, fmt.Sprintf("%s%s", bufLetter, string(rune(i-52))), fmt.Sprintf("%s%s", bufLetter, string(rune(i-52))), 3)

			_ = xlsx.SetSheetCol(SheetName2, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), &[]interface{}{
				eduprogcomps.Mandatory[i-66].Type + " " + eduprogcomps.Mandatory[i-66].Code,
			})

			lastLetter = fmt.Sprintf("%s%s", bufLetter, string(rune(i-52)))
		}

	}

	for i := mandLen + 66; i < mandLen+selLen+66; i++ {

		if i <= 90 {
			_ = xlsx.SetCellStyle(SheetName2, fmt.Sprintf("%s1", string(rune(i))), fmt.Sprintf("%s1", string(rune(i))), styleRotated)
			_ = xlsx.SetColWidth(SheetName2, string(rune(i)), string(rune(i)), 3)

			_ = xlsx.SetSheetCol(SheetName2, fmt.Sprintf("%s1", string(rune(i))), &[]interface{}{
				selective[i-mandLen-66].Type + " " + selective[i-mandLen-66].BlockNum + "." + selective[i-mandLen-66].Code,
			})

			lastLetter = string(rune(i))
		} else if i > 90 && i <= 116 {
			bufLetter = string(rune(65))
			_ = xlsx.SetCellStyle(SheetName2, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), styleRotated)
			_ = xlsx.SetColWidth(SheetName2, fmt.Sprintf("%s%s", bufLetter, string(rune(i-26))), fmt.Sprintf("%s%s", bufLetter, string(rune(i-26))), 3)

			_ = xlsx.SetSheetCol(SheetName2, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), &[]interface{}{
				selective[i-mandLen-66].Type + " " + selective[i-mandLen-66].BlockNum + "." + selective[i-mandLen-66].Code,
			})

			lastLetter = fmt.Sprintf("%s%s", bufLetter, string(rune(i-26)))
		} else if i > 116 && i <= 142 {
			bufLetter = string(rune(66))
			_ = xlsx.SetCellStyle(SheetName2, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), styleRotated)
			_ = xlsx.SetColWidth(SheetName2, fmt.Sprintf("%s%s", bufLetter, string(rune(i-52))), fmt.Sprintf("%s%s", bufLetter, string(rune(i-52))), 3)

			_ = xlsx.SetSheetCol(SheetName2, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), &[]interface{}{
				selective[i-mandLen-66].Type + " " + selective[i-mandLen-66].BlockNum + "." + selective[i-mandLen-66].Code,
			})

			lastLetter = fmt.Sprintf("%s%s", bufLetter, string(rune(i-52)))
		}

	}
	competenicesZKLen := len(eduprogcompetenciesZK)
	competenicesFKLen := len(eduprogcompetenciesFK)
	if competenicesZKLen == 0 {
		_ = xlsx.MergeCell(SheetName2, fmt.Sprintf("A%d", competenicesFKLen+2), fmt.Sprintf("Z%d", competenicesFKLen+2))
		_ = xlsx.SetCellStyle(SheetName2, fmt.Sprintf("A%d", competenicesFKLen+2), fmt.Sprintf("Z%d", competenicesFKLen+2), styleError)
		_ = xlsx.SetCellValue(SheetName2, fmt.Sprintf("A%d", competenicesFKLen+2), "Помилка: у освітньої програми відсутні ЗК")
	} else if competenicesZKLen > 0 {
		for i := 2; i < competenicesZKLen+2; i++ {
			_ = xlsx.SetCellStyle(SheetName2, fmt.Sprintf("A%d", i), fmt.Sprintf("A%d", i), styleBold)
			_ = xlsx.SetRowHeight(SheetName2, i, 15)
			_ = xlsx.SetSheetRow(SheetName2, fmt.Sprintf("A%d", i), &[]interface{}{
				eduprogcompetenciesZK[i-2].Type + " " + strconv.FormatUint(eduprogcompetenciesZK[i-2].Code, 10),
			})
		}
	}

	if competenicesFKLen == 0 {
		_ = xlsx.MergeCell(SheetName2, fmt.Sprintf("A%d", competenicesZKLen+competenicesFKLen+2), fmt.Sprintf("Z%d", competenicesZKLen+competenicesFKLen+2))
		_ = xlsx.SetCellStyle(SheetName2, fmt.Sprintf("A%d", competenicesZKLen+competenicesFKLen+2), fmt.Sprintf("Z%d", competenicesZKLen+competenicesFKLen+2), styleError)
		_ = xlsx.SetCellValue(SheetName2, fmt.Sprintf("A%d", competenicesZKLen+competenicesFKLen+2), "Помилка: у освітньої програми відсутні ФК")
	} else if competenicesFKLen > 0 {
		for i := competenicesZKLen + 2; i < competenicesZKLen+competenicesFKLen+2; i++ {
			_ = xlsx.SetCellStyle(SheetName2, fmt.Sprintf("A%d", i), fmt.Sprintf("A%d", i), styleBold)
			_ = xlsx.SetRowHeight(SheetName2, i, 15)
			_ = xlsx.SetSheetRow(SheetName2, fmt.Sprintf("A%d", i), &[]interface{}{
				eduprogcompetenciesFK[i-competenicesZKLen-2].Type + " " + strconv.FormatUint(eduprogcompetenciesFK[i-competenicesZKLen-2].Code, 10),
			})
		}
	}

	_ = xlsx.SetCellStyle(SheetName2, "B2", fmt.Sprintf("%s%d", lastLetter, competenicesZKLen+competenicesFKLen+1), styleDot)

	competenciesMatrix, _ := s.competenciesMatrixService.ShowByEduprogId(eduprogId)
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return "", nil, err
	}

	for i := 0; i < len(competenciesMatrix); i++ {
		eduprogcomp, _ := s.eduprogcompService.FindById(competenciesMatrix[i].ComponentId)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return "", nil, err
		}
		competency, _ := s.eduprogcompetenciesService.FindById(competenciesMatrix[i].CompetencyId)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return "", nil, err
		}
		edcode, _ := strconv.ParseUint(eduprogcomp.Code, 10, 64)
		if eduprogcomp.Type == "ВБ" {
			edcode = edcode + uint64(len(eduprogcomps.Mandatory))
		}

		if edcode+65 <= 90 {
			if competency.Type == "ЗК" {
				_ = xlsx.SetCellValue(SheetName2, fmt.Sprintf("%s%d", string(rune(edcode+65)), competency.Code+1), "·")
			} else if competency.Type == "ФК" {
				_ = xlsx.SetCellValue(SheetName2, fmt.Sprintf("%s%d", string(rune(edcode+65)), competency.Code+uint64(competenicesZKLen)+1), "·")
			}

		} else if edcode+65 > 90 && edcode+65 <= 116 {
			bufLetter = string(rune(65))

			if competency.Type == "ЗК" {
				_ = xlsx.SetCellValue(SheetName2, fmt.Sprintf("%s%s%d", bufLetter, string(rune(edcode+65-26)), competency.Code+1), "·")
			} else if competency.Type == "ФК" {
				_ = xlsx.SetCellValue(SheetName2, fmt.Sprintf("%s%s%d", bufLetter, string(rune(edcode+65-26)), competency.Code+uint64(competenicesZKLen)+1), "·")
			}
		}

	}

	//----------------------------EXPORT EDUPROGRESULTS MATRIX LOGIC----------------------------------//

	xlsx.SetActiveSheet(index3)
	_ = xlsx.SetSheetName("Sheet3", SheetName3)

	mandLen = len(eduprogcomps.Mandatory)
	lastLetter = ""
	_ = xlsx.SetRowHeight(SheetName3, 1, 40)
	for i := 66; i < mandLen+66; i++ {
		if i <= 90 {
			_ = xlsx.SetCellStyle(SheetName3, fmt.Sprintf("%s1", string(rune(i))), fmt.Sprintf("%s1", string(rune(i))), styleRotated)
			_ = xlsx.SetColWidth(SheetName3, string(rune(i)), string(rune(i)), 3)

			_ = xlsx.SetSheetCol(SheetName3, fmt.Sprintf("%s1", string(rune(i))), &[]interface{}{
				eduprogcomps.Mandatory[i-66].Type + " " + eduprogcomps.Mandatory[i-66].Code,
			})
			lastLetter = string(rune(i))
		} else if i > 90 && i <= 116 {
			bufLetter = string(rune(65))
			_ = xlsx.SetCellStyle(SheetName3, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), styleRotated)
			_ = xlsx.SetColWidth(SheetName3, fmt.Sprintf("%s%s", bufLetter, string(rune(i-26))), fmt.Sprintf("%s%s", bufLetter, string(rune(i-26))), 3)

			_ = xlsx.SetSheetCol(SheetName3, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), &[]interface{}{
				eduprogcomps.Mandatory[i-66].Type + " " + eduprogcomps.Mandatory[i-66].Code,
			})
			lastLetter = fmt.Sprintf("%s%s", bufLetter, string(rune(i-26)))
		} else if i > 116 && i <= 142 {
			bufLetter = string(rune(66))
			_ = xlsx.SetCellStyle(SheetName3, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), styleRotated)
			_ = xlsx.SetColWidth(SheetName3, fmt.Sprintf("%s%s", bufLetter, string(rune(i-52))), fmt.Sprintf("%s%s", bufLetter, string(rune(i-52))), 3)

			_ = xlsx.SetSheetCol(SheetName3, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), &[]interface{}{
				eduprogcomps.Mandatory[i-66].Type + " " + eduprogcomps.Mandatory[i-66].Code,
			})

			lastLetter = fmt.Sprintf("%s%s", bufLetter, string(rune(i-52)))
		}

	}

	for i := mandLen + 66; i < mandLen+selLen+66; i++ {

		if i <= 90 {
			_ = xlsx.SetCellStyle(SheetName3, fmt.Sprintf("%s1", string(rune(i))), fmt.Sprintf("%s1", string(rune(i))), styleRotated)
			_ = xlsx.SetColWidth(SheetName3, string(rune(i)), string(rune(i)), 3)

			_ = xlsx.SetSheetCol(SheetName3, fmt.Sprintf("%s1", string(rune(i))), &[]interface{}{
				selective[i-mandLen-66].Type + " " + selective[i-mandLen-66].BlockNum + "." + selective[i-mandLen-66].Code,
			})

			lastLetter = string(rune(i))
		} else if i > 90 && i <= 116 {
			bufLetter = string(rune(65))
			_ = xlsx.SetCellStyle(SheetName3, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), styleRotated)
			_ = xlsx.SetColWidth(SheetName3, fmt.Sprintf("%s%s", bufLetter, string(rune(i-26))), fmt.Sprintf("%s%s", bufLetter, string(rune(i-26))), 3)

			_ = xlsx.SetSheetCol(SheetName3, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), &[]interface{}{
				selective[i-mandLen-66].Type + " " + selective[i-mandLen-66].BlockNum + "." + selective[i-mandLen-66].Code,
			})

			lastLetter = fmt.Sprintf("%s%s", bufLetter, string(rune(i-26)))
		} else if i > 116 && i <= 142 {
			bufLetter = string(rune(66))
			_ = xlsx.SetCellStyle(SheetName3, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), styleRotated)
			_ = xlsx.SetColWidth(SheetName3, fmt.Sprintf("%s%s", bufLetter, string(rune(i-52))), fmt.Sprintf("%s%s", bufLetter, string(rune(i-52))), 3)

			_ = xlsx.SetSheetCol(SheetName3, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), &[]interface{}{
				selective[i-mandLen-66].Type + " " + selective[i-mandLen-66].BlockNum + "." + selective[i-mandLen-66].Code,
			})

			lastLetter = fmt.Sprintf("%s%s", bufLetter, string(rune(i-52)))
		}

	}

	competenicesPRLen := len(eduprogcompetenciesPR)
	if competenicesPRLen == 0 {
		_ = xlsx.MergeCell(SheetName3, "A2", "Z2")
		_ = xlsx.SetCellStyle(SheetName3, "A2", "Z2", styleError)
		_ = xlsx.SetCellValue(SheetName3, "A2", "Помилка: у освітньої програми відсутні ПР")
	} else if competenicesPRLen > 0 {
		for i := 2; i < competenicesPRLen+2; i++ {
			_ = xlsx.SetCellStyle(SheetName3, fmt.Sprintf("A%d", i), fmt.Sprintf("A%d", i), styleBold)
			_ = xlsx.SetRowHeight(SheetName3, i, 15)
			_ = xlsx.SetSheetRow(SheetName3, fmt.Sprintf("A%d", i), &[]interface{}{
				eduprogcompetenciesPR[i-2].Type + " " + strconv.FormatUint(eduprogcompetenciesPR[i-2].Code, 10),
			})

		}

		_ = xlsx.SetCellStyle(SheetName3, "B2", fmt.Sprintf("%s%d", lastLetter, competenicesPRLen+1), styleDot)

		resultsMatrix, _ := s.resultsMatrixService.ShowByEduprogId(eduprogId)
		if err != nil {
			log.Printf("EduprogService: %s", err)
			return "", nil, err
		}

		for i := 0; i < len(resultsMatrix); i++ {
			eduprogcomp, _ := s.eduprogcompService.FindById(resultsMatrix[i].ComponentId)
			if err != nil {
				log.Printf("EduprogService: %s", err)
				return "", nil, err
			}
			result, _ := s.eduprogcompetenciesService.FindById(resultsMatrix[i].EduprogresultId)
			if err != nil {
				log.Printf("EduprogService: %s", err)
				return "", nil, err
			}
			edcode, _ := strconv.ParseUint(eduprogcomp.Code, 10, 64)
			if eduprogcomp.Type == "ВБ" {
				edcode = edcode + uint64(len(eduprogcomps.Mandatory))
			}

			if edcode+65 <= 90 {
				_ = xlsx.SetCellValue(SheetName3, fmt.Sprintf("%s%d", string(rune(edcode+65)), result.Code+1), "·")
			} else if edcode+65 > 90 && edcode+65 <= 116 {
				bufLetter = string(rune(65))
				_ = xlsx.SetCellValue(SheetName3, fmt.Sprintf("%s%s%d", bufLetter, string(rune(edcode+65-26)), result.Code+1), "·")
			}

		}
	}

	_ = xlsx.SaveAs(fmt.Sprintf("./%s.xlsx", e.Name))
	if err != nil {
		log.Printf("EduprogService: %s", err)
		return "", nil, err
	}

	xlsx.SetActiveSheet(index)
	buf, _ := xlsx.WriteToBuffer()

	filename := fmt.Sprintf("%s.xlsx", e.Name)

	return filename, buf, nil

}
