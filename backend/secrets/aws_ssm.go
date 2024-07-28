package secrets

type AwsSSM struct{}

func (a AwsSSM) Get(id string) map[string]string {
	return map[string]string{}
}

func (a AwsSSM) Create(id string, data map[string]string) {

}

func (a AwsSSM) Update(id string, data map[string]string) {

}

func (a AwsSSM) Delete(id string) {

}
