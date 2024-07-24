package v1

type InfraTenant struct {
	Template   string            `json:"template"`
	Parameters map[string]string `json:"parameters"`
	AutoSync   bool              `json:"autoSync"`
	Expiry     string            `json:"expiry"`
}
