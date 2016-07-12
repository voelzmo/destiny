package consul

import (
	"fmt"

	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/destiny/network"
	"gopkg.in/yaml.v2"
)

type Manifest struct {
	DirectorUUID  string              `yaml:"director_uuid"`
	Name          string              `yaml:"name"`
	Jobs          []core.Job          `yaml:"jobs"`
	Properties    Properties          `yaml:"properties"`
	Compilation   core.Compilation    `yaml:"compilation"`
	Update        core.Update         `yaml:"update"`
	Networks      []core.Network      `yaml:"networks"`
	Releases      []core.Release      `yaml:"releases"`
	ResourcePools []core.ResourcePool `yaml:"resource_pools"`
}

type ConsulMember struct {
	Address string
}

func NewManifest(config Config, iaasConfig iaas.Config) Manifest {
	release := core.Release{
		Name:    "consul",
		Version: "latest",
	}

	config = populateDefaultConfigNodes(config)

	ipRanges := []network.IPRange{}

	for _, cfgNetwork := range config.Networks {
		ipRanges = append(ipRanges, network.IPRange(cfgNetwork.IPRange))
	}

	cloudProperties := iaasConfig.NetworkSubnet()

	consulNetworks := []core.Network{}
	for i, ipRange := range ipRanges {
		var staticIPs []string

		if i == 0 {
			staticIPs = []string{
				ipRange.IP(4),
				ipRange.IP(5),
				ipRange.IP(6),
				ipRange.IP(7),
				ipRange.IP(8),
				ipRange.IP(9),
			}
		} else {
			staticIPs = []string{
				ipRange.IP(4),
			}
		}

		consulNetwork := core.Network{
			Name: fmt.Sprintf("consul%d", i+1),
			Subnets: []core.NetworkSubnet{{
				CloudProperties: cloudProperties,
				Gateway:         ipRange.IP(1),
				Range:           string(ipRange),
				Reserved:        []string{ipRange.Range(2, 3), ipRange.Range(13, 254)},
				Static:          staticIPs,
			}},
			Type: "manual",
		}
		consulNetworks = append(consulNetworks, consulNetwork)
	}

	compilation := core.Compilation{
		Network:             consulNetworks[0].Name,
		ReuseCompilationVMs: true,
		Workers:             3,
		CloudProperties:     iaasConfig.Compilation(),
	}

	update := core.Update{
		Canaries:        1,
		CanaryWatchTime: "1000-180000",
		MaxInFlight:     50,
		Serial:          true,
		UpdateWatchTime: "1000-180000",
	}

	stemcell := core.ResourcePoolStemcell{
		Name:    iaasConfig.Stemcell(),
		Version: "latest",
	}

	resourcePools := []core.ResourcePool{}
	for i := range consulNetworks {
		resourcePool := core.ResourcePool{
			Name:            fmt.Sprintf("consul_z%d", i+1),
			Network:         consulNetworks[i].Name,
			Stemcell:        stemcell,
			CloudProperties: iaasConfig.ResourcePool(),
		}
		resourcePools = append(resourcePools, resourcePool)
	}

	jobs := []core.Job{}
	consulClusterStaticIPs := []string{}

	for i := range consulNetworks {
		instances := config.Networks[i].Nodes

		job := core.Job{
			Name:      fmt.Sprintf("consul_z%d", i+1),
			Instances: instances,
			Networks: []core.JobNetwork{{
				Name:      consulNetworks[i].Name,
				StaticIPs: consulNetworks[i].StaticIPs(instances),
			}},
			PersistentDisk: 1024,
			Properties: &core.JobProperties{
				Consul: core.JobPropertiesConsul{
					Agent: core.JobPropertiesConsulAgent{
						Mode:     "server",
						LogLevel: "info",
						Services: core.JobPropertiesConsulAgentServices{
							"router": core.JobPropertiesConsulAgentService{
								Name: "gorouter",
								Check: &core.JobPropertiesConsulAgentServiceCheck{
									Name:     "router-check",
									Script:   "/var/vcap/jobs/router/bin/script",
									Interval: "1m",
								},
								Tags: []string{"routing"},
							},
							"cloud_controller": core.JobPropertiesConsulAgentService{},
						},
					},
				},
			},
			ResourcePool: resourcePools[i].Name,
			Templates: []core.JobTemplate{{
				Name:    "consul_agent",
				Release: "consul",
			}},
			Update: &core.JobUpdate{
				MaxInFlight: 1,
			},
		}

		jobs = append(jobs, job)
		consulClusterStaticIPs = append(consulClusterStaticIPs, consulNetworks[i].StaticIPs(instances)...)
	}

	jobs = append(jobs, core.Job{
		Name:      "consul_test_consumer",
		Instances: 1,
		Networks: []core.JobNetwork{{
			Name:      consulNetworks[0].Name,
			StaticIPs: []string{consulNetworks[0].StaticIPs(6)[5]},
		}},
		PersistentDisk: 1024,
		ResourcePool:   resourcePools[0].Name,
		Templates: []core.JobTemplate{
			{
				Name:    "consul_agent",
				Release: "consul",
			},
			{
				Name:    "consul-test-consumer",
				Release: "consul",
			},
		},
	})

	properties := Properties{
		Consul: &PropertiesConsul{
			Agent: PropertiesConsulAgent{
				Domain: "cf.internal",
				Servers: PropertiesConsulAgentServers{
					Lan: consulClusterStaticIPs,
				},
			},
			EncryptKeys: []string{EncryptKey},
		},
	}

	overrideTLS(properties.Consul, config.DC)

	return Manifest{
		DirectorUUID:  config.DirectorUUID,
		Name:          config.Name,
		Releases:      []core.Release{release},
		Compilation:   compilation,
		Update:        update,
		ResourcePools: resourcePools,
		Jobs:          jobs,
		Networks:      consulNetworks,
		Properties:    properties,
	}
}

func (m Manifest) ConsulMembers() []ConsulMember {
	members := []ConsulMember{}
	for _, job := range m.Jobs {
		if len(job.Networks) == 0 {
			continue
		}

		for i := 0; i < job.Instances; i++ {
			if len(job.Networks[0].StaticIPs) > i {
				members = append(members, ConsulMember{
					Address: job.Networks[0].StaticIPs[i],
				})
			}
		}
	}

	return members
}

func (m Manifest) ToYAML() ([]byte, error) {
	return yaml.Marshal(m)
}

func FromYAML(manifestYAML []byte) (Manifest, error) {
	var m Manifest
	if err := yaml.Unmarshal(manifestYAML, &m); err != nil {
		return m, err
	}
	return m, nil
}

func overrideTLS(properties *PropertiesConsul, dc string) {
	switch dc {
	case "dc1":
		properties.Agent.Datacenter = dc
		properties.AgentCert = DC1AgentCert
		properties.AgentKey = DC1AgentKey
		properties.ServerCert = DC1ServerCert
		properties.ServerKey = DC1ServerKey
	case "dc2":
		properties.Agent.Datacenter = dc
		properties.AgentCert = DC2AgentCert
		properties.AgentKey = DC2AgentKey
		properties.ServerCert = DC2ServerCert
		properties.ServerKey = DC2ServerKey
	case "dc3":
		properties.Agent.Datacenter = dc
		properties.AgentCert = DC3AgentCert
		properties.AgentKey = DC3AgentKey
		properties.ServerCert = DC3ServerCert
		properties.ServerKey = DC3ServerKey
	default:
		properties.Agent.Datacenter = "dc1"
		properties.AgentCert = DC1AgentCert
		properties.AgentKey = DC1AgentKey
		properties.ServerCert = DC1ServerCert
		properties.ServerKey = DC1ServerKey
	}

	properties.CACert = CACert
}

func populateDefaultConfigNodes(config Config) Config {
	for i, cfgNetworks := range config.Networks {
		if cfgNetworks.Nodes == 0 {
			config.Networks[i].Nodes = 1
		}
	}

	return config
}
