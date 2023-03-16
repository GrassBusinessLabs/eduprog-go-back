package eduprog

import (
	"fmt"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers"
	_ "github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers"
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/resources"
	"github.com/go-chi/chi/v5"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const SheetName1 = "Перелік компонент"
const SheetName2 = "Матриця компетентностей"

const SheetName3 = "Матриця відповідності ПР"

func (c EduprogController) ExportEduprogToExcel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprog, _ := c.eduprogService.FindById(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		eduprogcomps, _ := c.eduprogcompService.SortComponentsByMnS(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		var creditsDto resources.CreditsDto

		for _, comp := range eduprogcomps.Selective {
			creditsDto.SelectiveCredits += comp.Credits
		}
		for _, comp := range eduprogcomps.Mandatory {
			creditsDto.MandatoryCredits += comp.Credits
		}
		creditsDto.TotalCredits = creditsDto.SelectiveCredits + creditsDto.MandatoryCredits
		creditsDto.TotalFreeCredits = 240 - creditsDto.TotalCredits
		creditsDto.MandatoryFreeCredits = 180 - creditsDto.MandatoryCredits
		creditsDto.SelectiveFreeCredits = 60 - creditsDto.SelectiveCredits

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
		_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("C%d", mandLen+4), fmt.Sprintf("%d кредитів", creditsDto.MandatoryCredits))

		selLen := len(eduprogcomps.Selective)

		for i := mandLen + 5; i < selLen+mandLen+5; i++ {

			_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
			_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("B%d", i), fmt.Sprintf("B%d", i), styleAlignLeft)

			_ = xlsx.SetSheetRow(SheetName1, fmt.Sprintf("A%d", i), &[]interface{}{
				eduprogcomps.Selective[i-mandLen-5].Type + " " + eduprogcomps.Selective[i-mandLen-5].Code + ".",
				eduprogcomps.Selective[i-mandLen-5].Name,
				eduprogcomps.Selective[i-mandLen-5].Credits,
				eduprogcomps.Selective[i-mandLen-5].ControlType,
			})

		}

		_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("A%d", selLen+mandLen+5), fmt.Sprintf("B%d", selLen+mandLen+5))
		_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("C%d", selLen+mandLen+5), fmt.Sprintf("D%d", selLen+mandLen+5))
		_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("A%d", selLen+mandLen+6), fmt.Sprintf("B%d", selLen+mandLen+6))
		_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("C%d", selLen+mandLen+6), fmt.Sprintf("D%d", selLen+mandLen+6))
		_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", selLen+mandLen+5), fmt.Sprintf("B%d", selLen+mandLen+5), styleBoldAlignLeft)
		_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("C%d", selLen+mandLen+5), fmt.Sprintf("D%d", selLen+mandLen+5), styleBold)
		_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("A%d", selLen+mandLen+5), "Загальний обсяг вибіркових компонент: ")
		_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("C%d", selLen+mandLen+5), fmt.Sprintf("%d кредитів", creditsDto.SelectiveCredits))
		_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", selLen+mandLen+6), fmt.Sprintf("B%d", selLen+mandLen+6), styleBoldAlignLeft)
		_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("C%d", selLen+mandLen+6), fmt.Sprintf("D%d", selLen+mandLen+6), styleBold)
		_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("A%d", selLen+mandLen+6), "ЗАГАЛЬНИЙ ОБСЯГ ОСВІТНЬОЇ ПРОГРАМИ: ")
		_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("C%d", selLen+mandLen+6), fmt.Sprintf("%d кредитів", creditsDto.TotalCredits))

		//----------------------------EXPORT COMPETENCIES MATRIX LOGIC----------------------------------//

		//eduprogcompetencies, _ := c.eduprogcompetenciesService.ShowCompetenciesByEduprogId(id)
		//if err != nil {
		//	log.Printf("EduprogController: %s", err)
		//	controllers.InternalServerError(w, err)
		//	return
		//}

		eduprogcompetenciesZK, _ := c.eduprogcompetenciesService.ShowCompetenciesByType(id, "ZK")
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		eduprogcompetenciesFK, _ := c.eduprogcompetenciesService.ShowCompetenciesByType(id, "FK")
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		eduprogcompetenciesPR, _ := c.eduprogcompetenciesService.ShowCompetenciesByType(id, "PR")
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
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
		selLen = len(eduprogcomps.Selective)
		lastLetter := ""
		bufLetter := ""
		_ = xlsx.SetRowHeight(SheetName2, 1, 40)
		for i := 66; i < mandLen+66; i++ {

			//	_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
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
				//	_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
				_ = xlsx.SetCellStyle(SheetName2, fmt.Sprintf("%s1", string(rune(i))), fmt.Sprintf("%s1", string(rune(i))), styleRotated)
				_ = xlsx.SetColWidth(SheetName2, string(rune(i)), string(rune(i)), 3)

				_ = xlsx.SetSheetCol(SheetName2, fmt.Sprintf("%s1", string(rune(i))), &[]interface{}{
					eduprogcomps.Selective[i-mandLen-66].Type + " " + eduprogcomps.Selective[i-mandLen-66].Code,
				})

				lastLetter = string(rune(i))
			} else if i > 90 && i <= 116 {
				bufLetter = string(rune(65))
				//	_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
				_ = xlsx.SetCellStyle(SheetName2, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), styleRotated)
				_ = xlsx.SetColWidth(SheetName2, fmt.Sprintf("%s%s", bufLetter, string(rune(i-26))), fmt.Sprintf("%s%s", bufLetter, string(rune(i-26))), 3)

				_ = xlsx.SetSheetCol(SheetName2, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), &[]interface{}{
					eduprogcomps.Selective[i-mandLen-66].Type + " " + eduprogcomps.Selective[i-mandLen-66].Code,
				})

				lastLetter = fmt.Sprintf("%s%s", bufLetter, string(rune(i-26)))
			} else if i > 116 && i <= 142 {
				bufLetter = string(rune(66))
				//	_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
				_ = xlsx.SetCellStyle(SheetName2, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), styleRotated)
				_ = xlsx.SetColWidth(SheetName2, fmt.Sprintf("%s%s", bufLetter, string(rune(i-52))), fmt.Sprintf("%s%s", bufLetter, string(rune(i-52))), 3)

				_ = xlsx.SetSheetCol(SheetName2, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), &[]interface{}{
					eduprogcomps.Selective[i-mandLen-66].Type + " " + eduprogcomps.Selective[i-mandLen-66].Code,
				})

				lastLetter = fmt.Sprintf("%s%s", bufLetter, string(rune(i-52)))
			}

		}

		competenicesZKLen := len(eduprogcompetenciesZK)

		for i := 2; i < competenicesZKLen+2; i++ {

			//_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
			_ = xlsx.SetCellStyle(SheetName2, fmt.Sprintf("A%d", i), fmt.Sprintf("A%d", i), style)
			_ = xlsx.SetRowHeight(SheetName2, i, 15)
			_ = xlsx.SetSheetRow(SheetName2, fmt.Sprintf("A%d", i), &[]interface{}{
				eduprogcompetenciesZK[i-2].Type + " " + strconv.FormatUint(eduprogcompetenciesZK[i-2].Code, 10),
			})

		}
		competenicesFKLen := len(eduprogcompetenciesFK)
		for i := competenicesZKLen + 2; i < competenicesZKLen+competenicesFKLen+2; i++ {

			//_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
			_ = xlsx.SetCellStyle(SheetName2, fmt.Sprintf("A%d", i), fmt.Sprintf("A%d", i), style)
			_ = xlsx.SetRowHeight(SheetName2, i, 15)
			_ = xlsx.SetSheetRow(SheetName2, fmt.Sprintf("A%d", i), &[]interface{}{
				eduprogcompetenciesFK[i-competenicesZKLen-2].Type + " " + strconv.FormatUint(eduprogcompetenciesFK[i-competenicesZKLen-2].Code, 10),
			})

		}

		_ = xlsx.SetCellStyle(SheetName2, "B2", fmt.Sprintf("%s%d", lastLetter, competenicesZKLen+1), styleDot)

		competenciesMatrix, _ := c.competenciesMatrixService.ShowByEduprogId(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for i := 0; i < len(competenciesMatrix); i++ {
			eduprogcomp, _ := c.eduprogcompService.FindById(competenciesMatrix[i].ComponentId)
			if err != nil {
				log.Printf("EduprogController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
			competency, _ := c.eduprogcompetenciesService.FindById(competenciesMatrix[i].CompetencyId)
			if err != nil {
				log.Printf("EduprogController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
			edcode, _ := strconv.ParseUint(eduprogcomp.Code, 10, 64)
			if eduprogcomp.Type == "ВБ" {
				edcode = edcode + uint64(len(eduprogcomps.Mandatory))
			}

			if edcode+65 <= 90 {
				_ = xlsx.SetCellValue(SheetName2, fmt.Sprintf("%s%d", string(rune(edcode+65)), competency.Code+1), "·")
			} else if edcode+65 > 90 && edcode+65 <= 116 {
				bufLetter = string(rune(65))
				_ = xlsx.SetCellValue(SheetName2, fmt.Sprintf("%s%s%d", bufLetter, string(rune(edcode+65-26)), competency.Code+1), "·")
			}

		}

		//----------------------------EXPORT EDUPROGRESULTS MATRIX LOGIC----------------------------------//

		xlsx.SetActiveSheet(index3)
		_ = xlsx.SetSheetName("Sheet3", SheetName3)

		mandLen = len(eduprogcomps.Mandatory)
		lastLetter = ""
		_ = xlsx.SetRowHeight(SheetName3, 1, 40)
		for i := 66; i < mandLen+66; i++ {

			//	_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
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
				//	_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
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
				//	_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
				_ = xlsx.SetCellStyle(SheetName3, fmt.Sprintf("%s1", string(rune(i))), fmt.Sprintf("%s1", string(rune(i))), styleRotated)
				_ = xlsx.SetColWidth(SheetName3, string(rune(i)), string(rune(i)), 3)

				_ = xlsx.SetSheetCol(SheetName3, fmt.Sprintf("%s1", string(rune(i))), &[]interface{}{
					eduprogcomps.Selective[i-mandLen-66].Type + " " + eduprogcomps.Selective[i-mandLen-66].Code,
				})

				lastLetter = string(rune(i))
			} else if i > 90 && i <= 116 {
				bufLetter = string(rune(65))
				//	_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
				_ = xlsx.SetCellStyle(SheetName3, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), styleRotated)
				_ = xlsx.SetColWidth(SheetName3, fmt.Sprintf("%s%s", bufLetter, string(rune(i-26))), fmt.Sprintf("%s%s", bufLetter, string(rune(i-26))), 3)

				_ = xlsx.SetSheetCol(SheetName3, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-26))), &[]interface{}{
					eduprogcomps.Selective[i-mandLen-66].Type + " " + eduprogcomps.Selective[i-mandLen-66].Code,
				})

				lastLetter = fmt.Sprintf("%s%s", bufLetter, string(rune(i-26)))
			} else if i > 116 && i <= 142 {
				bufLetter = string(rune(66))
				//	_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
				_ = xlsx.SetCellStyle(SheetName3, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), styleRotated)
				_ = xlsx.SetColWidth(SheetName3, fmt.Sprintf("%s%s", bufLetter, string(rune(i-52))), fmt.Sprintf("%s%s", bufLetter, string(rune(i-52))), 3)

				_ = xlsx.SetSheetCol(SheetName3, fmt.Sprintf("%s%s1", bufLetter, string(rune(i-52))), &[]interface{}{
					eduprogcomps.Selective[i-mandLen-66].Type + " " + eduprogcomps.Selective[i-mandLen-66].Code,
				})

				lastLetter = fmt.Sprintf("%s%s", bufLetter, string(rune(i-52)))
			}

		}

		competenicesPRLen := len(eduprogcompetenciesPR)
		for i := 2; i < competenicesPRLen+2; i++ {

			//_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
			_ = xlsx.SetCellStyle(SheetName3, fmt.Sprintf("A%d", i), fmt.Sprintf("A%d", i), style)
			_ = xlsx.SetRowHeight(SheetName3, i, 15)
			_ = xlsx.SetSheetRow(SheetName3, fmt.Sprintf("A%d", i), &[]interface{}{
				eduprogcompetenciesPR[i-2].Type + " " + strconv.FormatUint(eduprogcompetenciesPR[i-2].Code, 10),
			})

		}

		_ = xlsx.SetCellStyle(SheetName3, "B2", fmt.Sprintf("%s%d", lastLetter, competenicesPRLen+1), styleDot)

		resultsMatrix, _ := c.resultsMatrixService.ShowByEduprogId(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		for i := 0; i < len(resultsMatrix); i++ {
			eduprogcomp, _ := c.eduprogcompService.FindById(resultsMatrix[i].ComponentId)
			if err != nil {
				log.Printf("EduprogController: %s", err)
				controllers.InternalServerError(w, err)
				return
			}
			result, _ := c.eduprogresultsService.FindById(resultsMatrix[i].EduprogresultId)
			if err != nil {
				log.Printf("EduprogController: %s", err)
				controllers.InternalServerError(w, err)
				return
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

		_ = xlsx.SaveAs(fmt.Sprintf("./%s.xlsx", eduprog.Name))
		if err != nil {
			fmt.Println(err)
			return
		}
		xlsx.SetActiveSheet(index)
		w.Header().Set("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s.xlsx", eduprog.Name))
		buf, _ := xlsx.WriteToBuffer()
		http.ServeContent(w, r, fmt.Sprintf("%s.xlsx", eduprog.Name), time.Time{}, strings.NewReader(buf.String()))

	}
}

//func (c EduprogController) ExportEduprogListToExcel() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//
//		id, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
//		if err != nil {
//			log.Printf("EduprogschemeController: %s", err)
//			controllers.BadRequest(w, err)
//			return
//		}
//
//		eduprogcomps, _ := c.eduprogcompService.SortComponentsByMnS(id)
//		if err != nil {
//			log.Printf("EduprogcompController: %s", err)
//			//InternalServerError(w, err)
//			return
//		}
//
//		var creditsDto resources.CreditsDto
//
//		for _, comp := range eduprogcomps.Selective {
//			creditsDto.SelectiveCredits += comp.Credits
//		}
//		for _, comp := range eduprogcomps.Mandatory {
//			creditsDto.MandatoryCredits += comp.Credits
//		}
//		creditsDto.TotalCredits = creditsDto.SelectiveCredits + creditsDto.MandatoryCredits
//		creditsDto.TotalFreeCredits = 240 - creditsDto.TotalCredits
//		creditsDto.MandatoryFreeCredits = 180 - creditsDto.MandatoryCredits
//		creditsDto.SelectiveFreeCredits = 60 - creditsDto.SelectiveCredits
//
//		xlsx := excelize.NewFile()
//		index, _ := xlsx.NewSheet("Sheet1")
//		xlsx.SetActiveSheet(index)
//		err = xlsx.SetSheetName("Sheet1", SheetName1)
//
//		style, _ := xlsx.NewStyle(&excelize.Style{
//			Font:      &excelize.Font{Size: 12, Family: "Times New Roman"},
//			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
//			Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
//				{Type: "top", Color: "#000000", Style: 1},
//				{Type: "bottom", Color: "#000000", Style: 1},
//				{Type: "right", Color: "#000000", Style: 1},
//				{Type: "left", Color: "#000000", Style: 1},
//			},
//		})
//		styleAlignLeft, _ := xlsx.NewStyle(&excelize.Style{
//			Font:      &excelize.Font{Size: 12, Family: "Times New Roman"},
//			Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
//			Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
//				{Type: "top", Color: "#000000", Style: 1},
//				{Type: "bottom", Color: "#000000", Style: 1},
//				{Type: "right", Color: "#000000", Style: 1},
//				{Type: "left", Color: "#000000", Style: 1},
//			},
//		})
//		styleBold, _ := xlsx.NewStyle(&excelize.Style{
//			Font:      &excelize.Font{Size: 12, Bold: true, Family: "Times New Roman"},
//			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
//			Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
//				{Type: "top", Color: "#000000", Style: 1},
//				{Type: "bottom", Color: "#000000", Style: 1},
//				{Type: "right", Color: "#000000", Style: 1},
//				{Type: "left", Color: "#000000", Style: 1},
//			},
//		})
//		styleBoldAlignLeft, _ := xlsx.NewStyle(&excelize.Style{
//			Font:      &excelize.Font{Size: 12, Bold: true, Family: "Times New Roman"},
//			Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
//			Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
//				{Type: "top", Color: "#000000", Style: 1},
//				{Type: "bottom", Color: "#000000", Style: 1},
//				{Type: "right", Color: "#000000", Style: 1},
//				{Type: "left", Color: "#000000", Style: 1},
//			},
//		})
//		_ = xlsx.SetCellStyle(SheetName1, "A1", "D3", style)
//		_ = xlsx.MergeCell(SheetName1, "A3", "D3")
//		_ = xlsx.SetColWidth(SheetName1, "A", "A", 10)
//		_ = xlsx.SetColWidth(SheetName1, "B", "B", 50)
//		_ = xlsx.SetColWidth(SheetName1, "C", "C", 15)
//		_ = xlsx.SetColWidth(SheetName1, "D", "D", 20)
//
//		data := [][]interface{}{
//			{"Код н/д", "Компоненти освітньої програми (навчальні дисципліни, курсові проекти (роботи), практики, кваліфікаційна робота)", "Кількість кредитів", "Форма підсумкового контролю"},
//			{1, 2, 3, 4},
//			{"Обов'язкові компоненти ОП"},
//		}
//
//		_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", 3), fmt.Sprintf("D%d", 3), styleBold)
//		startRow := 1
//
//		for i := startRow; i < len(data)+startRow; i++ {
//
//			_ = xlsx.SetSheetRow(SheetName1, fmt.Sprintf("A%d", i), &data[i-1])
//
//		}
//
//		mandLen := len(eduprogcomps.Mandatory)
//
//		for i := 4; i < mandLen+4; i++ {
//
//			_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
//			_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("B%d", i), fmt.Sprintf("B%d", i), styleAlignLeft)
//
//			_ = xlsx.SetSheetRow(SheetName1, fmt.Sprintf("A%d", i), &[]interface{}{
//				eduprogcomps.Mandatory[i-4].Type + " " + eduprogcomps.Mandatory[i-4].Code + ".",
//				eduprogcomps.Mandatory[i-4].Name,
//				eduprogcomps.Mandatory[i-4].Credits,
//				eduprogcomps.Mandatory[i-4].ControlType,
//			})
//
//		}
//
//		_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("A%d", mandLen+4), fmt.Sprintf("B%d", mandLen+4))
//		_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("C%d", mandLen+4), fmt.Sprintf("D%d", mandLen+4))
//		_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("C%d", mandLen+4), fmt.Sprintf("D%d", mandLen+4), styleBold)
//		_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", mandLen+4), fmt.Sprintf("B%d", mandLen+4), styleBoldAlignLeft)
//		_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("A%d", mandLen+4), "Загальний обсяг обов'язкових компонент: ")
//		_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("C%d", mandLen+4), fmt.Sprintf("%d кредитів", creditsDto.MandatoryCredits))
//
//		selLen := len(eduprogcomps.Selective)
//
//		for i := mandLen + 5; i < selLen+mandLen+5; i++ {
//
//			_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
//			_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("B%d", i), fmt.Sprintf("B%d", i), styleAlignLeft)
//
//			_ = xlsx.SetSheetRow(SheetName1, fmt.Sprintf("A%d", i), &[]interface{}{
//				eduprogcomps.Selective[i-mandLen-5].Type + " " + eduprogcomps.Selective[i-mandLen-5].Code + ".",
//				eduprogcomps.Selective[i-mandLen-5].Name,
//				eduprogcomps.Selective[i-mandLen-5].Credits,
//				eduprogcomps.Selective[i-mandLen-5].ControlType,
//			})
//
//		}
//
//		_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("A%d", selLen+mandLen+5), fmt.Sprintf("B%d", selLen+mandLen+5))
//		_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("C%d", selLen+mandLen+5), fmt.Sprintf("D%d", selLen+mandLen+5))
//		_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("A%d", selLen+mandLen+6), fmt.Sprintf("B%d", selLen+mandLen+6))
//		_ = xlsx.MergeCell(SheetName1, fmt.Sprintf("C%d", selLen+mandLen+6), fmt.Sprintf("D%d", selLen+mandLen+6))
//		_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", selLen+mandLen+5), fmt.Sprintf("B%d", selLen+mandLen+5), styleBoldAlignLeft)
//		_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("C%d", selLen+mandLen+5), fmt.Sprintf("D%d", selLen+mandLen+5), styleBold)
//		_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("A%d", selLen+mandLen+5), "Загальний обсяг вибіркових компонент: ")
//		_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("C%d", selLen+mandLen+5), fmt.Sprintf("%d кредитів", creditsDto.SelectiveCredits))
//		_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("A%d", selLen+mandLen+6), fmt.Sprintf("B%d", selLen+mandLen+6), styleBoldAlignLeft)
//		_ = xlsx.SetCellStyle(SheetName1, fmt.Sprintf("C%d", selLen+mandLen+6), fmt.Sprintf("D%d", selLen+mandLen+6), styleBold)
//		_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("A%d", selLen+mandLen+6), "ЗАГАЛЬНИЙ ОБСЯГ ОСВІТНЬОЇ ПРОГРАМИ: ")
//		_ = xlsx.SetCellValue(SheetName1, fmt.Sprintf("C%d", selLen+mandLen+6), fmt.Sprintf("%d кредитів", creditsDto.TotalCredits))
//
//		_ = xlsx.SaveAs("./ComponentsCollection.xlsx")
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		//w.Header().Set("Content-Type", "application/octet-stream")
//		w.Header().Set("Content-Disposition", "attachment; filename="+"ComponentsCollection.xlsx")
//		//w.Header().Set("Content-Transfer-Encoding", "binary")
//		//w.Header().Set("Expires", "0")
//		//xlsx.Write(w)
//		//
//		//buff, err := xlsx.WriteToBuffer()
//		//if err != nil {
//		//    fmt.Println(err)
//		//    return
//		//}
//		//SuccessExport(w, buff.Bytes())
//
//		buf, _ := xlsx.WriteToBuffer()
//		http.ServeContent(w, r, "ComponentsCollection.xlsx", time.Time{}, strings.NewReader(buf.String()))
//	}
//}

//func (c EduprogController) ExportCompetenciesMatrixToExcel() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		id, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
//		if err != nil {
//			log.Printf("EduprogschemeController: %s", err)
//			controllers.BadRequest(w, err)
//			return
//		}
//
//		eduprogcomps, _ := c.eduprogcompService.SortComponentsByMnS(id)
//		if err != nil {
//			log.Printf("EduprogcompController: %s", err)
//			controllers.InternalServerError(w, err)
//			return
//		}
//
//		eduprogcompetencies, _ := c.eduprogcompetenciesService.ShowCompetenciesByEduprogId(id)
//		if err != nil {
//			log.Printf("EduprogcompController: %s", err)
//			controllers.InternalServerError(w, err)
//			return
//		}
//
//		xlsx := excelize.NewFile()
//		index, _ := xlsx.NewSheet("Sheet1")
//		xlsx.SetActiveSheet(index)
//		err = xlsx.SetSheetName("Sheet1", SheetName1)
//
//		style, _ := xlsx.NewStyle(&excelize.Style{
//			Font:      &excelize.Font{Size: 12, Family: "Times New Roman", Bold: true},
//			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
//			Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
//				{Type: "top", Color: "#000000", Style: 1},
//				{Type: "bottom", Color: "#000000", Style: 1},
//				{Type: "right", Color: "#000000", Style: 1},
//				{Type: "left", Color: "#000000", Style: 1},
//			},
//		})
//
//		styleRotated, _ := xlsx.NewStyle(&excelize.Style{
//			Font:      &excelize.Font{Size: 12, Family: "Times New Roman", Bold: true},
//			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true, TextRotation: 90},
//			Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
//				{Type: "top", Color: "#000000", Style: 1},
//				{Type: "bottom", Color: "#000000", Style: 1},
//				{Type: "right", Color: "#000000", Style: 1},
//				{Type: "left", Color: "#000000", Style: 1},
//			},
//		})
//
//		styleDot, _ := xlsx.NewStyle(&excelize.Style{
//			Font:      &excelize.Font{Size: 24, Family: "Times New Roman", Bold: true},
//			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
//			Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
//				{Type: "top", Color: "#000000", Style: 1},
//				{Type: "bottom", Color: "#000000", Style: 1},
//				{Type: "right", Color: "#000000", Style: 1},
//				{Type: "left", Color: "#000000", Style: 1},
//			},
//		})
//
//		mandLen := len(eduprogcomps.Mandatory)
//		lastLetter := ""
//		_ = xlsx.SetRowHeight(SheetName1, 1, 40)
//		for i := 66; i < mandLen+66; i++ {
//
//			//	_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
//			_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("%s1", string(rune(i))), fmt.Sprintf("%s1", string(rune(i))), styleRotated)
//			_ = xlsx.SetColWidth(SheetName, string(rune(i)), string(rune(i)), 3)
//
//			_ = xlsx.SetSheetCol(SheetName, fmt.Sprintf("%s1", string(rune(i))), &[]interface{}{
//				eduprogcomps.Mandatory[i-66].Type + " " + eduprogcomps.Mandatory[i-66].Code,
//			})
//			lastLetter = string(rune(i))
//		}
//
//		competenicesLen := len(eduprogcompetencies)
//
//		for i := 2; i < competenicesLen+2; i++ {
//
//			//_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
//			_ = xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("A%d", i), style)
//			_ = xlsx.SetRowHeight(SheetName, i, 15)
//			_ = xlsx.SetSheetRow(SheetName, fmt.Sprintf("A%d", i), &[]interface{}{
//				eduprogcompetencies[i-2].Type + " " + strconv.FormatUint(eduprogcompetencies[i-2].Code, 10),
//			})
//
//		}
//		//styleAlignLeft, _ := xlsx.NewStyle(&excelize.Style{
//		//	Font:      &excelize.Font{Size: 12, Family: "Times New Roman"},
//		//	Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
//		//	Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
//		//		{Type: "top", Color: "#000000", Style: 1},
//		//		{Type: "bottom", Color: "#000000", Style: 1},
//		//		{Type: "right", Color: "#000000", Style: 1},
//		//		{Type: "left", Color: "#000000", Style: 1},
//		//	},
//		//})
//		//styleBold, _ := xlsx.NewStyle(&excelize.Style{
//		//	Font:      &excelize.Font{Size: 12, Bold: true, Family: "Times New Roman"},
//		//	Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
//		//	Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
//		//		{Type: "top", Color: "#000000", Style: 1},
//		//		{Type: "bottom", Color: "#000000", Style: 1},
//		//		{Type: "right", Color: "#000000", Style: 1},
//		//		{Type: "left", Color: "#000000", Style: 1},
//		//	},
//		//})
//		//styleBoldAlignLeft, _ := xlsx.NewStyle(&excelize.Style{
//		//	Font:      &excelize.Font{Size: 12, Bold: true, Family: "Times New Roman"},
//		//	Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
//		//	Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
//		//		{Type: "top", Color: "#000000", Style: 1},
//		//		{Type: "bottom", Color: "#000000", Style: 1},
//		//		{Type: "right", Color: "#000000", Style: 1},
//		//		{Type: "left", Color: "#000000", Style: 1},
//		//	},
//		//})
//
//		_ = xlsx.SetCellStyle(SheetName, "B2", fmt.Sprintf("%s%d", lastLetter, competenicesLen+1), styleDot)
//
//		competenciesMatrix, _ := c.competenciesMatrixService.ShowByEduprogId(id)
//		if err != nil {
//			log.Printf("EduprogcompController: %s", err)
//			controllers.InternalServerError(w, err)
//			return
//		}
//
//		for i := 0; i < len(competenciesMatrix); i++ {
//			eduprogcomp, _ := c.eduprogcompService.FindById(competenciesMatrix[i].ComponentId)
//			if err != nil {
//				log.Printf("EduprogcompController: %s", err)
//				controllers.InternalServerError(w, err)
//				return
//			}
//			competency, _ := c.eduprogcompetenciesService.FindById(competenciesMatrix[i].CompetencyId)
//			if err != nil {
//				log.Printf("EduprogcompController: %s", err)
//				controllers.InternalServerError(w, err)
//				return
//			}
//			edcode, _ := strconv.ParseUint(eduprogcomp.Code, 10, 64)
//			_ = xlsx.SetCellValue(SheetName, fmt.Sprintf("%s%d", string(rune(edcode+65)), competency.Code+1), "·")
//
//		}
//
//		_ = xlsx.SaveAs("./CompetenciesMatrix.xlsx")
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		//w.Header().Set("Content-Type", "application/octet-stream")
//		w.Header().Set("Content-Disposition", "attachment; filename="+"CompetenciesMatrix.xlsx")
//		//w.Header().Set("Content-Transfer-Encoding", "binary")
//		//w.Header().Set("Expires", "0")
//		//xlsx.Write(w)
//		//
//		//buff, err := xlsx.WriteToBuffer()
//		//if err != nil {
//		//    fmt.Println(err)
//		//    return
//		//}
//		//SuccessExport(w, buff.Bytes())
//
//		buf, _ := xlsx.WriteToBuffer()
//		http.ServeContent(w, r, "CompetenciesMatrix.xlsx", time.Time{}, strings.NewReader(buf.String()))
//	}
//
//}

//
//func (c EduprogcompController) ExportEduprogcompListToExcel() http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		eduprogcomps, err := c.eduprogcompService.ShowList()
//		if err != nil {
//			log.Printf("EduprogcompController: %s", err)
//			InternalServerError(w, err)
//			return
//		}
//		xlsx := excelize.NewFile()
//
//		data := [][]interface{}{
//			{"Код н/д", "Компоненти освітньої програми (навчальні дисципліни, курсові проекти (роботи), практики, кваліфікаційна робота)", "Кількість кредитів", "Форма підсумкового контролю"},
//			{1, 2, 3, 4},
//			{"Обов'язкові компоненти ОП"},
//		}
//
//		//for i := range eduprogcomps.Items {
//		//
//		//}
//
//		for i, row := range data {
//			startCell, err := excelize.JoinCellName("A", i+1)
//			if err != nil {
//				log.Printf("EduprogcompController: %s", err)
//				return
//			}
//			err = xlsx.SetSheetRow("Sheet1", startCell, &row)
//			if err != nil {
//				log.Printf("EduprogcompController: %s", err)
//				return
//			}
//		}
//
//		err = xlsx.SaveAs("./Workbook2.xlsx")
//		if err != nil {
//			log.Printf("EduprogcompController: %s", err)
//			InternalServerError(w, err)
//			return
//		}
//
//	}
//}
