package destiny

type Config struct {
	DirectorUUID string
	Name         string
	IPRange      string
	IAAS         int
	AWS          ConfigAWS
	BOSH         ConfigBOSH
	Registry     ConfigRegistry
	Secrets      ConfigSecrets
}

type ConfigBOSH struct {
	Target         string
	Username       string
	Password       string
	DirectorCACert string
}

type ConfigAWS struct {
	AccessKeyID           string
	SecretAccessKey       string
	DefaultKeyName        string
	DefaultSecurityGroups []string
	Region                string
	Subnet                string
}

type ConfigRegistry struct {
	Host     string
	Password string
	Port     int
	Username string
}

type ConfigSecrets struct {
	Consul ConfigSecretsConsul
	Etcd   ConfigSecretsEtcd
}

type ConfigSecretsConsul struct {
	EncryptKey string
	AgentKey   string
	AgentCert  string
	ServerKey  string
	ServerCert string
	CACert     string
}

type ConfigSecretsEtcd struct {
	CACert     string
	ClientCert string
	ClientKey  string
	PeerCACert string
	PeerCert   string
	PeerKey    string
	ServerCert string
	ServerKey  string
}

func NewConfigWithDefaults(config Config) Config {
	if config.Secrets.Consul.CACert == "" {
		config.Secrets.Consul.CACert = ConsulCACert
	}

	if config.Secrets.Consul.EncryptKey == "" {
		config.Secrets.Consul.EncryptKey = ConsulEncryptKey
	}

	if config.Secrets.Consul.AgentKey == "" {
		config.Secrets.Consul.AgentKey = ConsulAgentKey
	}

	if config.Secrets.Consul.AgentCert == "" {
		config.Secrets.Consul.AgentCert = ConsulAgentCert
	}

	if config.Secrets.Consul.ServerKey == "" {
		config.Secrets.Consul.ServerKey = ConsulServerKey
	}

	if config.Secrets.Consul.ServerCert == "" {
		config.Secrets.Consul.ServerCert = ConsulServerCert
	}

	if config.Secrets.Etcd.CACert == "" {
		config.Secrets.Etcd.CACert = EtcdCACert
	}

	if config.Secrets.Etcd.ClientCert == "" {
		config.Secrets.Etcd.ClientCert = EtcdClientCert
	}

	if config.Secrets.Etcd.ClientKey == "" {
		config.Secrets.Etcd.ClientKey = EtcdClientKey
	}

	if config.Secrets.Etcd.PeerCACert == "" {
		config.Secrets.Etcd.PeerCACert = EtcdPeerCACert
	}

	if config.Secrets.Etcd.PeerCert == "" {
		config.Secrets.Etcd.PeerCert = EtcdPeerCert
	}

	if config.Secrets.Etcd.PeerKey == "" {
		config.Secrets.Etcd.PeerKey = EtcdPeerKey
	}

	if config.Secrets.Etcd.ServerCert == "" {
		config.Secrets.Etcd.ServerCert = EtcdServerCert
	}

	if config.Secrets.Etcd.ServerKey == "" {
		config.Secrets.Etcd.ServerKey = EtcdServerKey
	}

	return config
}
