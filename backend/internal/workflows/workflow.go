package workflows

import "go.uber.org/zap"

var Logger *zap.Logger

type Workflow interface {
	Get(id string) map[string]string
	Create(id string, parameters map[string]string, secretName string)
	Update(id string, parameters map[string]string, secretName string)
	Delete(id string)
}
