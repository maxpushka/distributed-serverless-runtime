package routes_management

import (
	"log"
	"net/http"
	"serverless/cdn"
	"strconv"

	"github.com/gorilla/mux"
)

func ParseRouteID(r *http.Request) (string, int, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Print(err)
		return "", 0, err
	}
	return idStr, id, nil
}

func SaveFile(w http.ResponseWriter, r *http.Request, idStr string, command *cdn.CommandCDN) error {
	r.Body = http.MaxBytesReader(w, r.Body, 20<<20) // 20 MB max
	err := r.ParseMultipartForm(10 << 20)           // 10 MB max
	if err != nil {
		log.Print(err)
		return err
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		log.Print(err)
		return err
	}
	defer file.Close()

	err = command.Upload(idStr, file)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
