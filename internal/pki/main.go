package pki

type ConfigPki struct {
	CaCertFile                 string `json:"openssl_root_ca_cert_file,omitempty"`
	CaPrivateKeyFile           string `json:"openssl_root_ca_key_file,omitempty"`
	CaPrivateKeyPasswordFile   string `json:"openssl_root_ca_key_password_file,omitempty"`
	CaPrivateKeyNoPasswordFile string `json:"openssl_root_ca_key_NoPassword_file,omitempty"`
	CaCert                     string `json:"openssl_root_ca_cert,omitempty"`
	CaPrivateKey               string `json:"openssl_root_ca_key,omitempty"`
	CaPrivateKeyPassword       string `json:"openssl_root_ca_key_password,omitempty"`
	CaPrivateKeyNoPassword     string `json:"openssl_root_ca_key_NoPassword,omitempty"`
}
