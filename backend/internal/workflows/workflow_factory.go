package workflows

type WorkflowFactory struct{}

type WorkflowType string

const (
	GCPCloudBuild   WorkflowType = "GCPCloudBuild"
	AWSCodePipeline WorkflowType = "AWSCodePipeline"
	ArgoWorkflows   WorkflowType = "ArgoWorkflows"
	GitHubActions   WorkflowType = "GitHubActions"
	Tekton          WorkflowType = "Tekton"
	Jenkins         WorkflowType = "Jenkins"
)

func (w *WorkflowFactory) CreateWorkflow(workflowType WorkflowType) Workflow {
	switch workflowType {
	case GCPCloudBuild:
		return GCPCloudBuildWorkflow{}
	case AWSCodePipeline:
		return AWSCodePipelineWorkflow{}
	default:
		return nil
	}
}
