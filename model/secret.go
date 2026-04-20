package model

type SecretProvider string

const (
	SecretProviderAnthropic SecretProvider = "anthropic"
)

type SecretType string

const (
	SecretTypeAnthropicAPIKey   SecretType = "anthropic_api_key"
	SecretTypeAnthropicAdminKey SecretType = "anthropic_admin_key"
)

type SecretStatus string

const (
	SecretStatusValid     SecretStatus = "valid"
	SecretStatusInvalid   SecretStatus = "invalid"
	SecretStatusError     SecretStatus = "error"
	SecretStatusUntested  SecretStatus = "untested"
)

type Secret struct {
	Provider SecretProvider `json:"provider"`
	Type     SecretType     `json:"type"`
	Value    string         `json:"value"`
	Status   SecretStatus   `json:"status,omitempty"`
}
