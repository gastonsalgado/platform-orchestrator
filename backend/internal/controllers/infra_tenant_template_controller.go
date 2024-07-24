package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	v1 "github.com/gastonsalgado/platform-orchestrator/backend/internal/api/v1"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func GetInfraTenantTemplates(w http.ResponseWriter, r *http.Request) {
	path := ""

	templates, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}

	result := []v1.InfraTenantTemplate{}
	for _, template := range templates {
		templateBytes, err := os.ReadFile(fmt.Sprintf("%s/%s", path, template.Name()))
		if err != nil {
			fmt.Println(err)
		}

		var template v1.InfraTenantTemplate
		err = json.Unmarshal(templateBytes, &template)
		if err != nil {
			fmt.Println(err)
		}

		result = append(result, template)
	}

	jsonBytes, err := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(jsonBytes))
}

func GetInfraTenantTemplateByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"id": %v}`, vars["id"])
}

func CreateInfraTenantTemplate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `{"alive": true}`)
}

func UpdateInfraTenantTemplate(w http.ResponseWriter, r *http.Request) {

}

func DeleteInfraTenantTemplate(w http.ResponseWriter, r *http.Request) {

}
