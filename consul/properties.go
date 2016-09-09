package consul

import (
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
)

type Properties struct {
	Consul          *PropertiesConsul         `yaml:"consul,omitempty"`
	WardenCPI       *iaas.PropertiesWardenCPI `yaml:"warden_cpi,omitempty"`
	AWS             *iaas.PropertiesAWS       `yaml:"aws,omitempty"`
	Registry        *core.PropertiesRegistry  `yaml:"registry,omitempty"`
	Blobstore       *core.PropertiesBlobstore `yaml:"blobstore,omitempty"`
	Agent           *core.PropertiesAgent     `yaml:"agent,omitempty"`
	TurbulenceAgent *core.PropertiesTurbulenceAgent     `yaml:"turbulence_agent,omitempty"`
}

type PropertiesConsul struct {
	Agent               PropertiesConsulAgent `yaml:"agent"`
	CACert              string                `yaml:"ca_cert"`
	AgentCert           string                `yaml:"agent_cert"`
	AgentKey            string                `yaml:"agent_key"`
	ServerCert          string                `yaml:"server_cert"`
	ServerKey           string                `yaml:"server_key"`
	EncryptKeys         []string              `yaml:"encrypt_keys"`
	DNSTimeoutInSeconds int                   `yaml:"dns_timeout_in_seconds,omitempty"`
}

type PropertiesConsulAgent struct {
	Domain     string                       `yaml:"domain"`
	LogLevel   string                       `yaml:"log_level,omitempty"`
	Servers    PropertiesConsulAgentServers `yaml:"servers"`
	Mode       string                       `yaml:"mode,omitempty"`
	Datacenter string                       `yaml:"datacenter,omitempty"`
}

type PropertiesConsulAgentServers struct {
	Lan []string `yaml:"lan"`
}
