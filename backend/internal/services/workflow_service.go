package services

import (
	"fmt"

	v1 "github.com/gastonsalgado/platform-orchestrator/backend/internal/api/v1"
	"github.com/gastonsalgado/platform-orchestrator/backend/internal/workflows"
)

func apply(id string, infraTenant *v1.InfraTenant, template *v1.InfraTenantTemplate) (string, error) {
	workflowParameters := map[string]string{}
	for parameterName, parameterData := range template.Parameters {
		if parameterData.Sensitive {
			continue
		}

		parameterValue, ok := infraTenant.Parameters[parameterName]
		if ok {
			workflowParameters[parameterName] = parameterValue
		} else {
			workflowParameters[parameterName] = parameterData.Default
		}
	}

	factory := workflows.WorkflowFactory{}
	workflow := factory.CreateWorkflow(workflows.GCPCloudBuild)

	workflowSecret := fmt.Sprintf("%s-secrets", id)
	workflow.Apply(id, workflowParameters, workflowSecret)

	return "<WORKFLOW_MONITORING_URL>", nil
}

func Apply(infraTenant *v1.InfraTenant, template *v1.InfraTenantTemplate) (string, error) {
	if infraTenant.Replicas < 2 {
		return apply(infraTenant.Id, infraTenant, template)
	}

	for i := 1; i <= infraTenant.Replicas; i++ {
		id := fmt.Sprintf("%s-%d", infraTenant.Id, i)
		_, err := apply(id, infraTenant, template)
		if err != nil {
			return "", err
		}
	}

	return "<WORKFLOW_MONITORING_URL>", nil
}

func Destroy() (string, error) {
	return "", nil
}
