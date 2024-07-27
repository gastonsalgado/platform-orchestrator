package v1

type InfraTenantTemplateParameters struct {
	Format    string `json:"format"`
	Default   string `json:"default"`
	Sensitive bool   `json:"sensitive"`
}

type InfraTenantTemplate struct {
	Id         string                                   `json:"id"`
	Repository string                                   `json:"repository"`
	Reference  string                                   `json:"reference"`
	Parameters map[string]InfraTenantTemplateParameters `json:"parameters"`
}
