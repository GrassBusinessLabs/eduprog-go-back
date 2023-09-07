package eduprog

import (
	"github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers"
	_ "github.com/GrassBusinessLabs/eduprog-go-back/internal/infra/http/controllers"
	"github.com/go-chi/chi/v5"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func (c EduprogController) ExportEduprogToWord() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		err = c.eduprogService.ExportEduprogToWord(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		filename := "OPP.docx"
		header := make(http.Header)
		header.Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": filename}))
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		for k, v := range header {
			w.Header()[k] = v
		}
		f, err := os.Open(filename)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		http.ServeContent(w, r, "OPP.docx", time.Time{}, f)
	}
}

func (c EduprogController) ExportEduprogToExcel() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		filename, buffer, err := c.eduprogService.ExportEduprogToExcel(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		header := make(http.Header)
		header.Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": filename}))
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		for k, v := range header {
			w.Header()[k] = v
		}

		http.ServeContent(w, r, filename, time.Time{}, strings.NewReader(buffer.String()))

	}
}

func (c EduprogController) ExportEducompRelationsToJpg() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "edId"), 10, 64)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.BadRequest(w, err)
			return
		}

		filename, err := c.eduprogService.ExportEducompRealtionsToJpg(id)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}

		header := make(http.Header)
		header.Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": filename}))
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		for k, v := range header {
			w.Header()[k] = v
		}
		f, err := os.Open(filename)
		if err != nil {
			log.Printf("EduprogController: %s", err)
			controllers.InternalServerError(w, err)
			return
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		http.ServeContent(w, r, filename, time.Time{}, f)
	}
}
