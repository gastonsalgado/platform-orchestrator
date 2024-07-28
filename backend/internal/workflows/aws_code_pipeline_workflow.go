package workflows

type AWSCodePipelineWorkflow struct{}

func (a AWSCodePipelineWorkflow) Get(id string) map[string]string {
	return map[string]string{}
}

func (a AWSCodePipelineWorkflow) Apply(id string, parameters map[string]string, secretName string) {

}

func (a AWSCodePipelineWorkflow) Destroy(id string) {

}
