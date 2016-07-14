package consul

import (
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/destiny/network"
	"gopkg.in/yaml.v2"
)

type ManifestV2 struct {
	DirectorUUID   string               `yaml:"director_uuid"`
	Name           string               `yaml:"name"`
	Releases       []core.Release       `yaml:"releases"`
	Stemcells      []core.Stemcell      `yaml:"stemcells"`
	Update         core.Update          `yaml:"update"`
	InstanceGroups []core.InstanceGroup `yaml:"instance_groups"`
	Properties     Properties           `yaml:"properties"`
}

func NewManifestV2(config Config, iaasConfig iaas.Config) ManifestV2 {
	return ManifestV2{
		DirectorUUID: config.DirectorUUID,
		Name:         config.Name,
		Releases:     releases(),
		Stemcells:    stemcells(),
		Update:       update(),
		InstanceGroups: []core.InstanceGroup{
			consulInstanceGroup(config.Networks),
			consulTestConsumerInstanceGroup(config.Networks),
		},
		Properties: properties(config.Networks),
	}
}

func consulInstanceGroup(networks []ConfigNetwork) core.InstanceGroup {
	return core.InstanceGroup{
		Instances: 1,
		Name:      "consul",
		AZs:       core.AZs(len(networks)),
		Networks: []core.InstanceGroupNetwork{
			{
				Name:      "private",
				StaticIPs: consulInstanceGroupStaticIPs(networks),
			},
		},
		VMType:             "default",
		Stemcell:           "default",
		PersistentDiskType: "default",
		Update: core.Update{
			MaxInFlight: 1,
		},
		Jobs: []core.InstanceGroupJob{
			{
				Name:    "consul_agent",
				Release: "consul",
			},
		},
		Properties: core.InstanceGroupProperties{
			Consul: core.InstanceGroupPropertiesConsul{
				Agent: core.InstanceGroupPropertiesConsulAgent{
					Mode:     "server",
					LogLevel: "info",
					Services: map[string]core.InstanceGroupPropertiesConsulAgentService{
						"router": core.InstanceGroupPropertiesConsulAgentService{
							Name: "gorouter",
							Check: core.InstanceGroupPropertiesConsulAgentServiceCheck{
								Name:     "router-check",
								Script:   "/var/vcap/jobs/router/bin/script",
								Interval: "1m",
							},
							Tags: []string{"routing"},
						},
						"cloud_controller": core.InstanceGroupPropertiesConsulAgentService{},
					},
				},
			},
		},
	}
}

func consulTestConsumerInstanceGroup(networks []ConfigNetwork) core.InstanceGroup {
	ipRange := network.IPRange(networks[0].IPRange)
	return core.InstanceGroup{
		Instances: 1,
		Name:      "consul_test_consumer",
		AZs:       []string{"z1"},
		Networks: []core.InstanceGroupNetwork{
			{
				Name: "private",
				StaticIPs: []string{
					ipRange.IP(9),
				},
			},
		},
		VMType:             "default",
		Stemcell:           "default",
		PersistentDiskType: "default",
		Jobs: []core.InstanceGroupJob{
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

func properties(networks []ConfigNetwork) Properties {
	return Properties{
		Consul: &PropertiesConsul{
			Agent: PropertiesConsulAgent{
				Domain:     "cf.internal",
				Datacenter: "dc1",
				Servers: PropertiesConsulAgentServers{
					Lan: consulInstanceGroupStaticIPs(networks),
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

func consulInstanceGroupStaticIPs(networks []ConfigNetwork) []string {
	staticIPs := []string{}
	for _, cfgNetwork := range networks {
		ipRange := network.IPRange(cfgNetwork.IPRange)
		for n := 0; n < cfgNetwork.Nodes; n++ {
			staticIPs = append(staticIPs, ipRange.IP(n+4))
		}
	}

	return staticIPs
}

func (m ManifestV2) ToYAML() ([]byte, error) {
	return yaml.Marshal(m)
}
