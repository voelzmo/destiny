package consul

import (
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
)

func NewManifestWithJobLevelProperties(config Config, iaasConfig iaas.Config) (Manifest, error) {
	manifest, err := NewManifest(config, iaasConfig)
	if err != nil {
		return Manifest{}, err
	}

	for _, consulJob := range manifest.Jobs {
		consulJob.Properties.Consul = &core.JobPropertiesConsul{
			Agent: core.JobPropertiesConsulAgent{
				Domain:     manifest.Properties.Consul.Agent.Domain,
				Datacenter: manifest.Properties.Consul.Agent.Datacenter,
				Servers: core.JobPropertiesConsulAgentServers{
					Lan: manifest.Properties.Consul.Agent.Servers.Lan,
				},
				Mode:     consulJob.Properties.Consul.Agent.Mode,
				LogLevel: consulJob.Properties.Consul.Agent.LogLevel,
				Services: consulJob.Properties.Consul.Agent.Services,
			},
			CACert:      manifest.Properties.Consul.CACert,
			ServerCert:  manifest.Properties.Consul.ServerCert,
			ServerKey:   manifest.Properties.Consul.ServerKey,
			EncryptKeys: manifest.Properties.Consul.EncryptKeys,
			AgentCert:   manifest.Properties.Consul.AgentCert,
			AgentKey:    manifest.Properties.Consul.AgentKey,
		}
	}

	consulTestConsumerJob := findJob(manifest, "consul_test_consumer")

	consulTestConsumerJob.Properties.Consul = &core.JobPropertiesConsul{
		Agent: core.JobPropertiesConsulAgent{
			Domain:     manifest.Properties.Consul.Agent.Domain,
			Datacenter: manifest.Properties.Consul.Agent.Datacenter,
			Servers: core.JobPropertiesConsulAgentServers{
				Lan: manifest.Properties.Consul.Agent.Servers.Lan,
			},
			Mode:     consulTestConsumerJob.Properties.Consul.Agent.Mode,
			LogLevel: consulTestConsumerJob.Properties.Consul.Agent.LogLevel,
			Services: consulTestConsumerJob.Properties.Consul.Agent.Services,
		},
		CACert:      manifest.Properties.Consul.CACert,
		AgentCert:   manifest.Properties.Consul.AgentCert,
		AgentKey:    manifest.Properties.Consul.AgentKey,
		EncryptKeys: manifest.Properties.Consul.EncryptKeys,
	}

	manifest.Properties.Consul = nil

	return manifest, nil
}

func findJob(manifest Manifest, name string) *core.Job {
	for index := range manifest.Jobs {
		if manifest.Jobs[index].Name == name {
			return &manifest.Jobs[index]
		}
	}
	return &core.Job{}
}
