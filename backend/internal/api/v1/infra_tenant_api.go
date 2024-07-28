package v1

type InfraTenant struct {
	Id         string            `json:"id"`
	Template   string            `json:"template"`
	Replicas   int               `json:"replicas,omitempty"`
	Parameters map[string]string `json:"parameters"`
	AutoSync   bool              `json:"autoSync"`
	Expiry     string            `json:"expiry"`
}
