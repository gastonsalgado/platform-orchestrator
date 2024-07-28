package secrets

import "fmt"

type GcpSecretManager struct{}

func (g GcpSecretManager) Get(id string) map[string]string {
	return map[string]string{}
}

func (g GcpSecretManager) Create(id string, data map[string]string) {

}

func (g GcpSecretManager) Update(id string, data map[string]string) {
	message := fmt.Sprintf("Storig %s in GSM", data)
	Logger.Info(message)
}

func (g GcpSecretManager) Delete(id string) {

}
