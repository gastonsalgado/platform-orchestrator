package workflows

import "go.uber.org/zap"

var Logger *zap.Logger

type Workflow interface {
	Get(id string) map[string]string
	Apply(id string, parameters map[string]string, secretName string)
	Destroy(id string)
}
