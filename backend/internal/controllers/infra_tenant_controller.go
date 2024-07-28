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

func GetInfraTenants(w http.ResponseWriter, r *http.Request) {
	// gitManager.Pull()

	infraTenantsBytes, err := services.GetAllInfraTenants()
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, infraTenantsBytes)
}

func GetInfraTenantByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	infraTenantId := vars["id"]

	// gitManager.Pull()

	infraTenantBytes, err := services.GetInfraTenantBytes(infraTenantId)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(infraTenantBytes) == 0 {
		http.Error(w, fmt.Sprintf("InfraTenant not found with Id %s", infraTenantId), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(infraTenantBytes))
}

func CreateInfraTenant(w http.ResponseWriter, r *http.Request) {
	var newInfraTenant *v1.InfraTenant
	err := json.NewDecoder(r.Body).Decode(&newInfraTenant)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// gitManager.Pull()

	infraTenantBytes, err := services.GetInfraTenantBytes(newInfraTenant.Id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(infraTenantBytes) > 0 {
		http.Error(w, fmt.Sprintf("InfraTenant already exists with Id %s", newInfraTenant.Id), http.StatusConflict)
		return
	}

	template, err := services.GetInfraTenantTemplate(newInfraTenant.Template)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if template == nil {
		http.Error(w, fmt.Sprintf("InfraTenantTemplate not found with Id %s", newInfraTenant.Template), http.StatusNotFound)
		return
	}

	err = services.AddInfraTenant(newInfraTenant, template)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// commitMessage := fmt.Sprintf("add InfraTenant with id %s", newInfraTenant.Id)
	// err = gitManager.Push(gitManager.InfraTenantsPath, commitMessage)
	// if err != nil {
	// 	Logger.Error(err.Error())
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func UpdateInfraTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	infraTenantId := vars["id"]

	// gitManager.Pull()

	infraTenantBytes, err := services.GetInfraTenantBytes(infraTenantId)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(infraTenantBytes) == 0 {
		http.Error(w, fmt.Sprintf("InfraTenant not found with Id %s", infraTenantId), http.StatusNotFound)
		return
	}

	var updatedInfraTenant *v1.InfraTenant
	err = json.NewDecoder(r.Body).Decode(&updatedInfraTenant)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedInfraTenant.Id = infraTenantId

	template, err := services.GetInfraTenantTemplate(updatedInfraTenant.Template)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if template == nil {
		http.Error(w, fmt.Sprintf("InfraTenantTemplate not found with Id %s", updatedInfraTenant.Template), http.StatusNotFound)
		return
	}

	updated, err := services.UpdateInfraTenant(infraTenantBytes, updatedInfraTenant, template)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updated {
		// commitMessage := fmt.Sprintf("update InfraTenant with id %s", infraTenantId)
		// err = gitManager.Push(gitManager.InfraTenantsPath, commitMessage)
		// if err != nil {
		// 	Logger.Error(err.Error())
		// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
		// 	return
		// }
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func DeleteInfraTenant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	infraTenantId := vars["id"]

	// gitManager.Pull()

	infraTenantBytes, err := services.GetInfraTenantBytes(infraTenantId)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if len(infraTenantBytes) == 0 {
		http.Error(w, fmt.Sprintf("InfraTenant not found with Id %s", infraTenantId), http.StatusNotFound)
		return
	}

	err = services.DeleteInfraTenant(infraTenantId)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// commitMessage := fmt.Sprintf("delete InfraTenant with id %s", infraTenantId)
	// err = gitManager.Push(gitManager.InfraTenantTemplatesPath, commitMessage)
	// if err != nil {
	// 	Logger.Error(err.Error())
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
