package consul

type Config struct {
	DirectorUUID  string
	Name          string
	IPRange       string
	ConsulSecrets ConfigSecretsConsul
	DC            string
}

type ConfigSecretsConsul struct {
	EncryptKey string
	AgentKey   string
	AgentCert  string
	ServerKey  string
	ServerCert string
	CACert     string
}
