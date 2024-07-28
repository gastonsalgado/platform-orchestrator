package workflows

import "fmt"

type GCPCloudBuildWorkflow struct{}

func (g GCPCloudBuildWorkflow) Get(id string) map[string]string {
	return map[string]string{}
}

func (g GCPCloudBuildWorkflow) Create(id string, parameters map[string]string, secretName string) {
	message := fmt.Sprintf("Triggering GCPCloudBuild workflow to create InfraTenant %s, parameters %s and secret %s", id, parameters, secretName)
	Logger.Info(message)
}

func (g GCPCloudBuildWorkflow) Update(id string, parameters map[string]string, secretName string) {
	message := fmt.Sprintf("Triggering GCPCloudBuild workflow to update InfraTenant %s, parameters %s and secret %s", id, parameters, secretName)
	Logger.Info(message)
}

func (g GCPCloudBuildWorkflow) Delete(id string) {

}
