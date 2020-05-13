package google

type CredentialsFile struct {
	Type          string `json:"type"`
	ProjectId     string `json:"project_id"`
	PrivateKeyId  string `json:"private_key_id"`
	PrivateKey    string `json:"private_key"`
	ClientEmail   string `json:"client_email"`
	ClientId      string `json:"client_id"`
	AuthUrl       string `json:"auth_uri"`
	TokenUrl      string `json:"token_uri"`
	AuthCertUrl   string `json:"auth_provider_x509_cert_url"`
	ClientCertUrl string `json:"client_x509_cert_url"`
}
