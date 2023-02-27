package eduprog

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strings"
	"time"
)

func (c EduprogController) ExportEduprogListToExcel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		xlsx := excelize.NewFile()
		_, err := xlsx.NewSheet("Sheet2")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = xlsx.SetCellValue("Sheet2", "A2", "Hello world.")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = xlsx.SetCellValue("Sheet1", "B2", 100)
		if err != nil {
			fmt.Println(err)
			return
		}
		xlsx.SetActiveSheet(2)
		err = xlsx.SaveAs("./Workbook2.xlsx")
		if err != nil {
			fmt.Println(err)
			return
		}
		//w.Header().Set("Content-Type", "application/octet-stream")
		//w.Header().Set("Content-Disposition", "attachment; filename="+"Workbook.xlsx")
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
		http.ServeContent(w, r, "test.xlsx", time.Time{}, strings.NewReader(buf.String()))
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
