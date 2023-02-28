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

const SheetName = "Перелік компонент"

func (c EduprogController) ExportEduprogListToExcel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
		if err != nil {
			log.Printf("EduprogschemeController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		eduprogcomps, _ := c.eduprogcompService.SortComponentsByMnS(id)
		if err != nil {
			log.Printf("EduprogcompController: %s", err)
			//InternalServerError(w, err)
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

		xlsx := excelize.NewFile()
		index, _ := xlsx.NewSheet("Sheet1")
		xlsx.SetActiveSheet(index)
		err = xlsx.SetSheetName("Sheet1", SheetName)

		style, _ := xlsx.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Size: 12},
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
			Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
				{Type: "top", Color: "#000000", Style: 1},
				{Type: "bottom", Color: "#000000", Style: 1},
				{Type: "right", Color: "#000000", Style: 1},
				{Type: "left", Color: "#000000", Style: 1},
			},
		})
		styleAlignLeft, _ := xlsx.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Size: 12},
			Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
			Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
				{Type: "top", Color: "#000000", Style: 1},
				{Type: "bottom", Color: "#000000", Style: 1},
				{Type: "right", Color: "#000000", Style: 1},
				{Type: "left", Color: "#000000", Style: 1},
			},
		})
		styleBold, _ := xlsx.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Size: 12, Bold: true},
			Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
			Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
				{Type: "top", Color: "#000000", Style: 1},
				{Type: "bottom", Color: "#000000", Style: 1},
				{Type: "right", Color: "#000000", Style: 1},
				{Type: "left", Color: "#000000", Style: 1},
			},
		})
		styleBoldAlignLeft, _ := xlsx.NewStyle(&excelize.Style{
			Font:      &excelize.Font{Size: 12, Bold: true},
			Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center", WrapText: true},
			Border: []excelize.Border{{Type: "left", Color: "#000000", Style: 2},
				{Type: "top", Color: "#000000", Style: 1},
				{Type: "bottom", Color: "#000000", Style: 1},
				{Type: "right", Color: "#000000", Style: 1},
				{Type: "left", Color: "#000000", Style: 1},
			},
		})
		xlsx.SetCellStyle(SheetName, "A1", "D3", style)
		xlsx.MergeCell(SheetName, "A3", "D3")
		xlsx.SetColWidth(SheetName, "A", "A", 10)
		xlsx.SetColWidth(SheetName, "B", "B", 50)
		xlsx.SetColWidth(SheetName, "C", "C", 15)
		xlsx.SetColWidth(SheetName, "D", "D", 20)

		data := [][]interface{}{
			{"Код н/д", "Компоненти освітньої програми (навчальні дисципліни, курсові проекти (роботи), практики, кваліфікаційна робота)", "Кількість кредитів", "Форма підсумкового контролю"},
			{1, 2, 3, 4},
			{"Обов'язкові компоненти ОП"},
		}

		xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", 3), fmt.Sprintf("D%d", 3), styleBold)
		startRow := 1

		for i := startRow; i < len(data)+startRow; i++ {

			xlsx.SetSheetRow(SheetName, fmt.Sprintf("A%d", i), &data[i-1])

		}

		mandLen := len(eduprogcomps.Mandatory)

		for i := 4; i < mandLen+4; i++ {

			xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
			xlsx.SetCellStyle(SheetName, fmt.Sprintf("B%d", i), fmt.Sprintf("B%d", i), styleAlignLeft)

			xlsx.SetSheetRow(SheetName, fmt.Sprintf("A%d", i), &[]interface{}{
				eduprogcomps.Mandatory[i-4].Type + " " + eduprogcomps.Mandatory[i-4].Code,
				eduprogcomps.Mandatory[i-4].Name,
				eduprogcomps.Mandatory[i-4].Credits,
				eduprogcomps.Mandatory[i-4].ControlType,
			})

		}

		xlsx.MergeCell(SheetName, fmt.Sprintf("A%d", mandLen+4), fmt.Sprintf("B%d", mandLen+4))
		xlsx.MergeCell(SheetName, fmt.Sprintf("C%d", mandLen+4), fmt.Sprintf("D%d", mandLen+4))
		xlsx.SetCellStyle(SheetName, fmt.Sprintf("C%d", mandLen+4), fmt.Sprintf("D%d", mandLen+4), styleBold)
		xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", mandLen+4), fmt.Sprintf("B%d", mandLen+4), styleBoldAlignLeft)
		xlsx.SetCellValue(SheetName, fmt.Sprintf("A%d", mandLen+4), "Загальний обсяг обов'язкових компонент: ")
		xlsx.SetCellValue(SheetName, fmt.Sprintf("C%d", mandLen+4), fmt.Sprintf("%d кредитів", creditsDto.MandatoryCredits))

		selLen := len(eduprogcomps.Selective)

		for i := mandLen + 5; i < selLen+mandLen+5; i++ {

			xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", i), fmt.Sprintf("D%d", i), style)
			xlsx.SetCellStyle(SheetName, fmt.Sprintf("B%d", i), fmt.Sprintf("B%d", i), styleAlignLeft)

			xlsx.SetSheetRow(SheetName, fmt.Sprintf("A%d", i), &[]interface{}{
				eduprogcomps.Selective[i-mandLen-5].Type + " " + eduprogcomps.Selective[i-mandLen-5].Code,
				eduprogcomps.Selective[i-mandLen-5].Name,
				eduprogcomps.Selective[i-mandLen-5].Credits,
				eduprogcomps.Selective[i-mandLen-5].ControlType,
			})

		}

		xlsx.MergeCell(SheetName, fmt.Sprintf("A%d", selLen+mandLen+5), fmt.Sprintf("B%d", selLen+mandLen+5))
		xlsx.MergeCell(SheetName, fmt.Sprintf("C%d", selLen+mandLen+5), fmt.Sprintf("D%d", selLen+mandLen+5))
		xlsx.MergeCell(SheetName, fmt.Sprintf("A%d", selLen+mandLen+6), fmt.Sprintf("B%d", selLen+mandLen+6))
		xlsx.MergeCell(SheetName, fmt.Sprintf("C%d", selLen+mandLen+6), fmt.Sprintf("D%d", selLen+mandLen+6))
		xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", selLen+mandLen+5), fmt.Sprintf("B%d", selLen+mandLen+5), styleBoldAlignLeft)
		xlsx.SetCellStyle(SheetName, fmt.Sprintf("C%d", selLen+mandLen+5), fmt.Sprintf("D%d", selLen+mandLen+5), styleBold)
		xlsx.SetCellValue(SheetName, fmt.Sprintf("A%d", selLen+mandLen+5), "Загальний обсяг вибіркових компонент: ")
		xlsx.SetCellValue(SheetName, fmt.Sprintf("C%d", selLen+mandLen+5), fmt.Sprintf("%d кредитів", creditsDto.SelectiveCredits))
		xlsx.SetCellStyle(SheetName, fmt.Sprintf("A%d", selLen+mandLen+6), fmt.Sprintf("B%d", selLen+mandLen+6), styleBoldAlignLeft)
		xlsx.SetCellStyle(SheetName, fmt.Sprintf("C%d", selLen+mandLen+6), fmt.Sprintf("D%d", selLen+mandLen+6), styleBold)
		xlsx.SetCellValue(SheetName, fmt.Sprintf("A%d", selLen+mandLen+6), "ЗАГАЛЬНИЙ ОБСЯГ ОСВІТНЬОЇ ПРОГРАМИ: ")
		xlsx.SetCellValue(SheetName, fmt.Sprintf("C%d", selLen+mandLen+6), fmt.Sprintf("%d кредитів", creditsDto.TotalCredits))

		err = xlsx.SaveAs("./ComponentsCollection.xlsx")
		if err != nil {
			fmt.Println(err)
			return
		}
		//w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+"ComponentsCollection.xlsx")
		//w.Header().Set("Content-Transfer-Encoding", "binary")
		//w.Header().Set("Expires", "0")
		//xlsx.Write(w)
		//
		//buff, err := xlsx.WriteToBuffer()
		//if err != nil {
		//    fmt.Println(err)
		//    return
		//}
		//SuccessExport(w, buff.Bytes())

		buf, _ := xlsx.WriteToBuffer()
		http.ServeContent(w, r, "ComponentsCollection.xlsx", time.Time{}, strings.NewReader(buf.String()))
	}
}

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
