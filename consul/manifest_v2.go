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
	Stemcells      []Stemcell           `yaml:"stemcells"`
	Update         core.Update          `yaml:"update"`
	InstanceGroups []core.InstanceGroup `yaml:"instance_groups"`
	Properties     Properties           `yaml:"properties"`
}

type Stemcell struct {
	Alias   string
	Name    string
	Version string
}

func NewManifestV2(config ConfigV2, iaasConfig iaas.Config) ManifestV2 {
	return ManifestV2{
		DirectorUUID: config.DirectorUUID,
		Name:         config.Name,
		Releases:     releases(),
		Stemcells: []Stemcell{
			{
				Alias:   "default",
				Version: "latest",
				Name:    iaasConfig.Stemcell(),
			},
		},
		Update: update(),
		InstanceGroups: []core.InstanceGroup{
			consulInstanceGroup(config.AZs, config.PersistentDiskType, config.VMType),
			consulTestConsumerInstanceGroup(config.AZs, config.VMType),
		},
		Properties: properties(config.AZs),
	}
}

func consulInstanceGroup(azs []ConfigAZ, persistentDiskType string, vmType string) core.InstanceGroup {
	totalNodes := 0
	for _, az := range azs {
		totalNodes += az.Nodes
	}

	if persistentDiskType == "" {
		persistentDiskType = "default"
	}

	if vmType == "" {
		vmType = "default"
	}

	return core.InstanceGroup{
		Instances: totalNodes,
		Name:      "consul",
		AZs:       core.AZs(len(azs)),
		Networks: []core.InstanceGroupNetwork{
			{
				Name:      "private",
				StaticIPs: consulInstanceGroupStaticIPs(azs),
			},
		},
		VMType:             vmType,
		Stemcell:           "default",
		PersistentDiskType: persistentDiskType,
		Jobs: []core.InstanceGroupJob{
			{
				Name:    "consul_agent",
				Release: "consul",
			},
		},
		MigratedFrom: []core.InstanceGroupMigratedFrom{
			{
				Name: "consul_z1",
				AZ:   "z1",
			},
			{
				Name: "consul_z2",
				AZ:   "z2",
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

func consulTestConsumerInstanceGroup(azs []ConfigAZ, vmType string) core.InstanceGroup {
	ipRange := network.IPRange(azs[0].IPRange)
	totalNodesInFirstAZ := azs[0].Nodes

	if vmType == "" {
		vmType = "default"
	}

	return core.InstanceGroup{
		Instances: 3,
		Name:      "test_consumer",
		AZs:       []string{azs[0].Name},
		Networks: []core.InstanceGroupNetwork{
			{
				Name: "private",
				StaticIPs: []string{
					ipRange.IP(totalNodesInFirstAZ + 8),
					ipRange.IP(totalNodesInFirstAZ + 9),
					ipRange.IP(totalNodesInFirstAZ + 10),
				},
			},
		},
		VMType:   vmType,
		Stemcell: "default",
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
		MigratedFrom: []core.InstanceGroupMigratedFrom{
			{
				Name: "consul_test_consumer",
				AZ:   "z1",
			},
		},
	}
}

func properties(azs []ConfigAZ) Properties {
	return Properties{
		Consul: &PropertiesConsul{
			Agent: PropertiesConsulAgent{
				Domain:     "cf.internal",
				Datacenter: "dc1",
				Servers: PropertiesConsulAgentServers{
					Lan: consulInstanceGroupStaticIPs(azs),
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

func consulInstanceGroupStaticIPs(azs []ConfigAZ) []string {
	staticIPs := []string{}
	for _, cfgAZs := range azs {
		ipRange := network.IPRange(cfgAZs.IPRange)
		for n := 0; n < cfgAZs.Nodes; n++ {
			staticIPs = append(staticIPs, ipRange.IP(n+4))
		}
	}

	return staticIPs
}

func (m ManifestV2) ToYAML() ([]byte, error) {
	return yaml.Marshal(m)
}
