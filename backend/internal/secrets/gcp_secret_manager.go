package secrets

import "fmt"

type GcpSecretManager struct{}

func (g GcpSecretManager) Get(id string) map[string]string {
	return map[string]string{}
}

func (g GcpSecretManager) Create(id string, data map[string]string) {

}

func (g GcpSecretManager) Update(id string, data map[string]string) {
	message := fmt.Sprintf("Storig %s for InfraTenant %s in Google Secret Manager", data, id)
	Logger.Info(message)
}

func (g GcpSecretManager) Delete(id string) {

}
