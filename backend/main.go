package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gastonsalgado/platform-orchestrator/backend/internal/controllers"
	"github.com/gastonsalgado/platform-orchestrator/backend/internal/managers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func setupRoutes() *mux.Router {
	router := mux.NewRouter()

	infraTenantTemplatePath := "/api/v1/infraTenantTemplate"
	infraTenantTemplatePathById := fmt.Sprintf("%s/{id}", infraTenantTemplatePath)
	router.HandleFunc(infraTenantTemplatePath, controllers.GetInfraTenantTemplates).Methods("GET")
	router.HandleFunc(infraTenantTemplatePathById, controllers.GetInfraTenantTemplateByID).Methods("GET")
	router.HandleFunc(infraTenantTemplatePath, controllers.CreateInfraTenantTemplate).Methods("POST")
	router.HandleFunc(infraTenantTemplatePathById, controllers.UpdateInfraTenantTemplate).Methods("PUT")
	router.HandleFunc(infraTenantTemplatePathById, controllers.DeleteInfraTenantTemplate).Methods("DELETE")

	infraTenantPath := "/api/v1/infraTenant"
	infraTenantPathById := fmt.Sprintf("%s/{id}", infraTenantPath)
	router.HandleFunc(infraTenantPath, controllers.GetInfraTenants).Methods("GET")
	router.HandleFunc(infraTenantPathById, controllers.GetInfraTenantByID).Methods("GET")
	router.HandleFunc(infraTenantPath, controllers.CreateInfraTenant).Methods("POST")
	router.HandleFunc(infraTenantPathById, controllers.UpdateInfraTenant).Methods("PUT")
	router.HandleFunc(infraTenantPathById, controllers.DeleteInfraTenant).Methods("DELETE")

	return router
}

func main() {
	logger, _ := zap.NewProduction()
	controllers.Logger = logger
	managers.Logger = logger

	gitManagerInstance := managers.GetGitManagerInstance()
	err := gitManagerInstance.Init()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	http.Handle("/", setupRoutes())
	port := 8080
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	logger.Fatal(err.Error())
}
