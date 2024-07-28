package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	v1 "github.com/gastonsalgado/platform-orchestrator/backend/internal/api/v1"
	"github.com/gastonsalgado/platform-orchestrator/backend/internal/services"
	"github.com/gorilla/mux"
)

func GetInfraTenantTemplates(w http.ResponseWriter, r *http.Request) {
	// gitManager.Pull()

	result, err := services.GetAllInfraTenantTemplates()
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, result)
}

func GetInfraTenantTemplateByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateId := vars["id"]

	// gitManager.Pull()

	result, err := services.GetInfraTenantTemplateBytes(templateId)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(result) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(result))
}

func CreateInfraTenantTemplate(w http.ResponseWriter, r *http.Request) {
	var newTemplate *v1.InfraTenantTemplate
	err := json.NewDecoder(r.Body).Decode(&newTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// gitManager.Pull()

	templateBytes, err := services.GetInfraTenantTemplateBytes(newTemplate.Id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(templateBytes) > 0 {
		http.Error(w, fmt.Sprintf("InfraTenantTemplate already exists with Id %s", newTemplate.Id), http.StatusConflict)
		return
	}

	err = services.AddInfraTenantTemplate(newTemplate)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// commitMessage := fmt.Sprintf("add InfraTenantTemplate with id %s", newTemplate.Id)
	// err = gitManager.Push(gitManager.InfraTenantTemplatesPath, commitMessage)
	// if err != nil {
	// 	Logger.Error(err.Error())
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func UpdateInfraTenantTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateId := vars["id"]

	// gitManager.Pull()

	templateBytes, err := services.GetInfraTenantTemplateBytes(templateId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(templateBytes) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var updatedTemplate *v1.InfraTenantTemplate
	err = json.NewDecoder(r.Body).Decode(&updatedTemplate)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedTemplate.Id = templateId

	updated, err := services.UpdateInfraTenantTemplate(templateBytes, updatedTemplate)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updated {
		// commitMessage := fmt.Sprintf("update InfraTenantTemplate with id %s", templateId)
		// err = gitManager.Push(gitManager.InfraTenantTemplatesPath, commitMessage)
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

	// gitManager.Pull()

	templateBytes, err := services.GetInfraTenantTemplateBytes(templateId)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(templateBytes) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = services.DeleteInfraTenantTemplate(templateId)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	commitMessage := fmt.Sprintf("delete InfraTenantTemplate with id %s", templateId)
	err = gitManager.Push(gitManager.InfraTenantTemplatesPath, commitMessage)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
