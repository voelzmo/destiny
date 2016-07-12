package consul

import (
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"gopkg.in/yaml.v2"
)

type ManifestV2 struct {
	DirectorUUID   string               `yaml:"director_uuid"`
	Name           string               `yaml:"name"`
	Releases       []core.Release       `yaml:"releases"`
	Stemcells      []core.Stemcell      `yaml:"stemcells"`
	Update         core.Update          `yaml:"update"`
	InstanceGroups []core.InstanceGroup `yaml:"instance_groups"`
	Properties     PropertiesV2         `yaml:"properties"`
}

type PropertiesV2 struct {
	Consul ConsulProperties `yaml:"consul"`
}

type ConsulProperties struct {
	Agent       AgentProperties `yaml:"agent"`
	AgentCert   string          `yaml:"agent_cert"`
	AgentKey    string          `yaml:"agent_key"`
	CACert      string          `yaml:"ca_cert"`
	EncryptKeys []string        `yaml:"encrypt_keys"`
	ServerCert  string          `yaml:"server_cert"`
	ServerKey   string          `yaml:"server_key"`
}

type AgentProperties struct {
	Domain     string                `yaml:"domain"`
	Datacenter string                `yaml:"datacenter"`
	Servers    AgentServerProperties `yaml:"servers"`
}

type AgentServerProperties struct {
	Lan []string `yaml:"lan"`
}

func NewManifestV2(config Config, iaasConfig iaas.Config) ManifestV2 {
	return ManifestV2{
		DirectorUUID: config.DirectorUUID,
		Name:         config.Name,
		Releases:     releases(),
		Stemcells:    stemcells(),
		Update:       update(),
		InstanceGroups: []core.InstanceGroup{
			consulInstanceGroup(),
			consulTestConsumerInstanceGroup(),
		},
		Properties: properties(),
	}
}

func releases() []core.Release {
	return []core.Release{
		{
			Name:    "consul",
			Version: "latest",
		},
	}
}

func stemcells() []core.Stemcell {
	return []core.Stemcell{
		{
			Alias:   "default",
			OS:      "ubuntu-trusty",
			Version: "latest",
		},
	}
}

func update() core.Update {
	return core.Update{
		Canaries:        1,
		CanaryWatchTime: "1000-180000",
		MaxInFlight:     50,
		Serial:          true,
		UpdateWatchTime: "1000-180000",
	}
}

func consulInstanceGroup() core.InstanceGroup {
	return core.InstanceGroup{
		Instances: 1,
		Name:      "consul",
		AZs:       []string{"z1"},
		Networks: []core.InstanceGroupNetwork{
			{
				Name: "consul1",
				StaticIPs: []string{
					"10.244.4.4",
				},
			},
		},
		VMType:             "default",
		Stemcell:           "default",
		PersistentDiskType: "default",
		Update: core.Update{
			MaxInFlight: 1,
		},
		Jobs: []core.JobV2{
			{
				Name:    "consul_agent",
				Release: "consul",
			},
		},
		Properties: core.InstanceGroupProperties{
			Consul: core.ConsulInstanceGroupProperties{
				Agent: core.ConsulAgentProperties{
					Mode:     "server",
					LogLevel: "info",
					Services: map[string]core.ConsulAgentServiceProperties{
						"router": core.ConsulAgentServiceProperties{
							Name: "gorouter",
							Check: core.ConsulServiceCheckProperties{
								Name:     "router-check",
								Script:   "/var/vcap/jobs/router/bin/script",
								Interval: "1m",
							},
							Tags: []string{"routing"},
						},
						"cloud_controller": core.ConsulAgentServiceProperties{},
					},
				},
			},
		},
	}
}

func consulTestConsumerInstanceGroup() core.InstanceGroup {
	return core.InstanceGroup{
		Instances: 1,
		Name:      "consul_test_consumer",
		AZs:       []string{"z1"},
		Networks: []core.InstanceGroupNetwork{
			{
				Name: "consul1",
				StaticIPs: []string{
					"10.244.4.9",
				},
			},
		},
		VMType:             "default",
		Stemcell:           "default",
		PersistentDiskType: "default",
		Jobs: []core.JobV2{
			{
				Name:    "consul_agent",
				Release: "consul",
			},
			{
				Name:    "consul-test-consumer",
				Release: "consul",
			},
		},
	}
}

func properties() PropertiesV2 {
	return PropertiesV2{
		Consul: ConsulProperties{
			Agent: AgentProperties{
				Domain:     "cf.internal",
				Datacenter: "dc1",
				Servers: AgentServerProperties{
					Lan: []string{"10.244.4.4"},
				},
			},
			AgentCert: DC1AgentCert,
			AgentKey:  DC1AgentKey,
			CACert:    CACert,
			EncryptKeys: []string{
				EncryptKey,
			},
			ServerCert: DC1ServerCert,
			ServerKey:  DC1ServerKey,
		},
	}
}

func (m ManifestV2) ToYAML() ([]byte, error) {
	return yaml.Marshal(m)
}
