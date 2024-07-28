package secrets

import "go.uber.org/zap"

var Logger *zap.Logger

type SecretManager interface {
	Get(id string) map[string]string
	Create(id string, data map[string]string)
	Update(id string, data map[string]string)
	Delete(id string)
}
