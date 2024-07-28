package secrets

type SecretManagerFactory struct{}

type SecretManagerType string

const (
	GSM    SecretManagerType = "GcpSecretManager"
	AWSSSM SecretManagerType = "AwsSSM"
)

func (s *SecretManagerFactory) CreateSecretManager(secretManagerType SecretManagerType) SecretManager {
	switch secretManagerType {
	case GSM:
		return GcpSecretManager{}
	case AWSSSM:
		return AwsSSM{}
	default:
		return nil
	}
}
