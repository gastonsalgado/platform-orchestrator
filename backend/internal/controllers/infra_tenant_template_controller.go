package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	v1 "github.com/gastonsalgado/platform-orchestrator/backend/internal/api/v1"
	"github.com/gastonsalgado/platform-orchestrator/backend/internal/managers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var Logger *zap.Logger
var gitManager = managers.GetGitManagerInstance()

func GetInfraTenantTemplates(w http.ResponseWriter, r *http.Request) {
	// gitManager.Pull()
	templatesPath := fmt.Sprintf("%s/%s", gitManager.BasePath, gitManager.InfraTenantTemplatesPath)

	templates, err := os.ReadDir(templatesPath)
	if err != nil {
		fmt.Println(err)
	}

	result := []v1.InfraTenantTemplate{}
	for _, template := range templates {
		templateBytes, err := os.ReadFile(fmt.Sprintf("%s/%s", templatesPath, template.Name()))
		if err != nil {
			Logger.Error(err.Error())
		}

		var template v1.InfraTenantTemplate
		err = json.Unmarshal(templateBytes, &template)
		if err != nil {
			Logger.Error(err.Error())
		}

		result = append(result, template)
	}

	jsonBytes, err := json.Marshal(result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(jsonBytes))
}

func GetInfraTenantTemplateByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateId := vars["id"]

	templateBytes, err := infraTenantTemplateExist(templateId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(templateBytes) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(templateBytes))
}

func CreateInfraTenantTemplate(w http.ResponseWriter, r *http.Request) {
	var newTemplate v1.InfraTenantTemplate
	err := json.NewDecoder(r.Body).Decode(&newTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	templateBytes, err := infraTenantTemplateExist(newTemplate.Id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(templateBytes) > 0 {
		http.Error(w, fmt.Sprintf("There is already an InfraTenantTemplate with with Id %s", newTemplate.Id), http.StatusConflict)
		return
	}

	newTemplateBytes, err := json.MarshalIndent(newTemplate, "", "    ")
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	newTemplatePath := getInfraTenantTemplatePath(newTemplate.Id)
	err = os.WriteFile(newTemplatePath, newTemplateBytes, 0644)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = gitManager.Push(gitManager.InfraTenantTemplatesPath, "add InfraTenantTemplate")
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func UpdateInfraTenantTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateId := vars["id"]

	templateBytes, err := infraTenantTemplateExist(templateId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(templateBytes) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var updatedTemplate v1.InfraTenantTemplate
	err = json.NewDecoder(r.Body).Decode(&updatedTemplate)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTemplateBytes, err := json.MarshalIndent(updatedTemplate, "", "    ")
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if string(templateBytes) != string(updatedTemplateBytes) {
		err = os.WriteFile(getInfraTenantTemplatePath(templateId), updatedTemplateBytes, 0644)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// err = gitManager.Push(gitManager.InfraTenantTemplatesPath, "update InfraTenantTemplate")
		// if err != nil {
		// 	Logger.Error(err.Error())
		// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
		// 	return
		// }
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func DeleteInfraTenantTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateId := vars["id"]

	templateBytes, err := infraTenantTemplateExist(templateId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(templateBytes) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = os.Remove(getInfraTenantTemplatePath(templateId))
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// err = gitManager.Push(gitManager.InfraTenantTemplatesPath, "delete InfraTenantTemplate")
	// if err != nil {
	// 	Logger.Error(err.Error())
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func infraTenantTemplateExist(id string) ([]byte, error) {
	// gitManager.Pull()
	templatePath := getInfraTenantTemplatePath(id)
	templateBytes, err := os.ReadFile(templatePath)
	if os.IsNotExist(err) {
		return templateBytes, nil
	}
	if err != nil {
		Logger.Error(err.Error())
		return templateBytes, err
	}

	return templateBytes, nil
}

func getInfraTenantTemplatePath(id string) string {
	templateFilename := fmt.Sprintf("%s.json", id)
	return fmt.Sprintf("%s/%s/%s", gitManager.BasePath, gitManager.InfraTenantTemplatesPath, templateFilename)
}
