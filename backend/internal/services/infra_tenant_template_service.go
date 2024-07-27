package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	v1 "github.com/gastonsalgado/platform-orchestrator/backend/internal/api/v1"
	"github.com/gastonsalgado/platform-orchestrator/backend/internal/managers"
	"go.uber.org/zap"
)

var Logger *zap.Logger
var gitManager = managers.GetGitManagerInstance()

func getInfraTenantTemplatePath(id string) string {
	templateFilename := id
	if !strings.HasSuffix(templateFilename, ".json") {
		templateFilename = fmt.Sprintf("%s.json", templateFilename)
	}
	return fmt.Sprintf("%s/%s/%s", gitManager.BasePath, gitManager.InfraTenantTemplatesPath, templateFilename)
}

func getInfraTenantTemplatesPath() string {
	return fmt.Sprintf("%s/%s", gitManager.BasePath, gitManager.InfraTenantTemplatesPath)
}

func GetAllInfraTenantTemplate() (string, error) {
	templatesPath := getInfraTenantTemplatesPath()
	templateFiles, err := os.ReadDir(templatesPath)
	if err != nil {
		return "", err
	}

	templates := []*v1.InfraTenantTemplate{}
	for _, templateFile := range templateFiles {
		template, err := GetInfraTenantTemplate(templateFile.Name())
		if err != nil {
			return "", err
		}
		templates = append(templates, template)
	}

	templatesBytes, err := json.Marshal(templates)
	if err != nil {
		return "", err
	}

	return string(templatesBytes), nil
}

func GetInfraTenantTemplateBytes(id string) ([]byte, error) {
	templatePath := getInfraTenantTemplatePath(id)
	templateBytes, err := os.ReadFile(templatePath)
	if err != nil && !os.IsNotExist(err) {
		return templateBytes, err
	}

	return templateBytes, nil
}

func GetInfraTenantTemplate(id string) (*v1.InfraTenantTemplate, error) {
	templateBytes, err := GetInfraTenantTemplateBytes(id)
	if err != nil {
		return nil, err
	}

	if len(templateBytes) == 0 {
		return nil, nil
	}

	template := &v1.InfraTenantTemplate{}
	err = json.Unmarshal(templateBytes, &template)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func AddInfraTenantTemplate(newTemplate *v1.InfraTenantTemplate) error {
	newTemplateBytes, err := json.MarshalIndent(newTemplate, "", "    ")
	if err != nil {
		return err
	}

	newTemplatePath := getInfraTenantTemplatePath(newTemplate.Id)
	return os.WriteFile(newTemplatePath, newTemplateBytes, 0644)
}

func UpdateInfraTenantTemplate(templateBytes []byte, updatedTemplate *v1.InfraTenantTemplate) (bool, error) {
	updatedTemplateBytes, err := json.MarshalIndent(updatedTemplate, "", "    ")
	if err != nil {
		return false, err
	}

	if string(templateBytes) == string(updatedTemplateBytes) {
		return false, nil
	}

	return true, os.WriteFile(getInfraTenantTemplatePath(updatedTemplate.Id), updatedTemplateBytes, 0644)
}

func DeleteInfraTenantTemplate(id string) error {
	return os.Remove(getInfraTenantTemplatePath(id))
}
