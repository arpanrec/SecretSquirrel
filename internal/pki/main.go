package pki

type ConfigJsonPki struct {
	// The path to the CA certificate file.
	CaCertFile string `json:"openssl_root_ca_cert_path,omitempty"`
	// The path to the CA private key file.
	CaPrivateKeyFile string `json:"openssl_root_ca_key_path,omitempty"`
	// The path to the CA private key password file.
	CaPrivateKeyPasswordFile string `json:"openssl_root_ca_key_password_path,omitempty"`
	// The path to the server certificate file.
	CaPrivateKeyNopassFile string `json:"openssl_root_ca_key_nopasswd_path,omitempty"`
}
