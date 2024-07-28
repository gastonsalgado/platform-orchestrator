package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	v1 "github.com/gastonsalgado/platform-orchestrator/backend/internal/api/v1"
	"github.com/gastonsalgado/platform-orchestrator/backend/internal/secrets"
)

func getInfraTenantPath(id string) string {
	infraTenantFilename := id
	if !strings.HasSuffix(infraTenantFilename, ".json") {
		infraTenantFilename = fmt.Sprintf("%s.json", infraTenantFilename)
	}
	return fmt.Sprintf("%s/%s/%s", gitManager.BasePath, gitManager.InfraTenantsPath, infraTenantFilename)
}

func getInfraTenantsPath() string {
	return fmt.Sprintf("%s/%s", gitManager.BasePath, gitManager.InfraTenantsPath)
}

func GetAllInfraTenants() (string, error) {
	infraTenantsPath := getInfraTenantsPath()
	infraTenantFiles, err := os.ReadDir(infraTenantsPath)
	if err != nil {
		return "", err
	}

	infraTenants := []*v1.InfraTenant{}
	for _, infraTenantFile := range infraTenantFiles {
		infraTenant, err := GetInfraTenant(infraTenantFile.Name())
		if err != nil {
			return "", err
		}
		infraTenants = append(infraTenants, infraTenant)
	}

	infraTenantsBytes, err := json.Marshal(infraTenants)
	if err != nil {
		return "", err
	}

	return string(infraTenantsBytes), nil
}

func GetInfraTenantBytes(id string) ([]byte, error) {
	infraTenantPath := getInfraTenantPath(id)
	infraTenantBytes, err := os.ReadFile(infraTenantPath)
	if err != nil && !os.IsNotExist(err) {
		return infraTenantBytes, err
	}

	return infraTenantBytes, nil
}

func GetInfraTenant(id string) (*v1.InfraTenant, error) {
	infraTenantBytes, err := GetInfraTenantBytes(id)
	if err != nil {
		return nil, err
	}

	if len(infraTenantBytes) == 0 {
		return nil, nil
	}

	infraTenant := &v1.InfraTenant{}
	err = json.Unmarshal(infraTenantBytes, &infraTenant)
	if err != nil {
		return nil, err
	}

	return infraTenant, nil
}

func addOrUpdateInfraTenantNonSensitiveData(id string, infraTenant *v1.InfraTenant, template *v1.InfraTenantTemplate) error {
	newInfraTenant := &v1.InfraTenant{
		Id:         id,
		Template:   infraTenant.Template,
		AutoSync:   infraTenant.AutoSync,
		Parameters: map[string]string{},
	}

	for parameterName, parameterData := range template.Parameters {
		if parameterData.Sensitive {
			continue
		}

		parameterValue, ok := infraTenant.Parameters[parameterName]
		if ok {
			newInfraTenant.Parameters[parameterName] = parameterValue
		} else {
			newInfraTenant.Parameters[parameterName] = parameterData.Default
		}
	}

	newInfraTenantBytes, err := json.MarshalIndent(newInfraTenant, "", "    ")
	if err != nil {
		return err
	}

	newInfraTenantPath := getInfraTenantPath(id)
	return os.WriteFile(newInfraTenantPath, newInfraTenantBytes, 0644)
}

func addOrUpdateInfraTenantSensitiveData(id string, infraTenant *v1.InfraTenant, template *v1.InfraTenantTemplate) error {
	factory := secrets.SecretManagerFactory{}
	secretManager := factory.CreateSecretManager(secrets.GSM)

	infraTenantSensitiveParameters := map[string]string{}
	templateSensitiveParameters := secretManager.Get(template.Id)

	for parameterName, parameterData := range template.Parameters {
		if parameterData.Sensitive {
			parameterValue, ok := infraTenant.Parameters[parameterName]
			if ok {
				infraTenantSensitiveParameters[parameterName] = parameterValue
			} else {
				infraTenantSensitiveParameters[parameterName] = templateSensitiveParameters[parameterName]
			}

		}
	}

	secretManager.Update(id, infraTenantSensitiveParameters)

	return nil
}

func addOrUpdateInfraTenant(id string, infraTenant *v1.InfraTenant, template *v1.InfraTenantTemplate) error {
	err := addOrUpdateInfraTenantNonSensitiveData(id, infraTenant, template)
	if err != nil {
		return nil
	}

	err = addOrUpdateInfraTenantSensitiveData(id, infraTenant, template)
	if err != nil {
		return nil
	}

	return nil
}

func AddInfraTenant(infraTenant *v1.InfraTenant, template *v1.InfraTenantTemplate) error {
	if infraTenant.Replicas < 2 {
		return addOrUpdateInfraTenant(infraTenant.Id, infraTenant, template)
	}

	for i := 1; i <= infraTenant.Replicas; i++ {
		id := fmt.Sprintf("%s-%d", infraTenant.Id, i)
		err := addOrUpdateInfraTenant(id, infraTenant, template)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateInfraTenant(infraTenantBytes []byte, updatedInfraTenant *v1.InfraTenant, template *v1.InfraTenantTemplate) (bool, error) {
	updatedInfraTenantBytes, err := json.MarshalIndent(updatedInfraTenant, "", "    ")
	if err != nil {
		return false, err
	}

	if string(infraTenantBytes) == string(updatedInfraTenantBytes) {
		return false, nil
	}

	return true, addOrUpdateInfraTenant(updatedInfraTenant.Id, updatedInfraTenant, template)
}

func DeleteInfraTenant(id string) error {
	return os.Remove(getInfraTenantPath(id))
}
