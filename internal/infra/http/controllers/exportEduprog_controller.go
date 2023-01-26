package controllers

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"net/http"
)

func (c EduprogController) ExportEduprogListToExcel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		defer file.Close()
		f, err := excelize.OpenReader(file)
		if err != nil {
			fmt.Fprint(w, err.Error())
			return
		}
		f.Path = "Book1.xlsx"
		f.NewSheet("NewSheet")

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", f.Path))
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

		if err := f.Write(w); err != nil {
			fmt.Fprint(w, err.Error())
		}
	}
}
