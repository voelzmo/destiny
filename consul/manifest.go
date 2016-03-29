package consul

import (
	"github.com/cloudfoundry-incubator/candiedyaml"
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/destiny/network"
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

	ipRange := network.IPRange(config.IPRange)

	cloudProperties := iaasConfig.NetworkSubnet()

	consulNetwork1 := core.Network{
		Name: "consul1",
		Subnets: []core.NetworkSubnet{{
			CloudProperties: cloudProperties,
			Gateway:         ipRange.IP(1),
			Range:           string(ipRange),
			Reserved:        []string{ipRange.Range(2, 3), ipRange.Range(13, 254)},
			Static: []string{
				ipRange.IP(4),
				ipRange.IP(5),
				ipRange.IP(6),
				ipRange.IP(7),
				ipRange.IP(8),
				ipRange.IP(9),
			},
		}},
		Type: "manual",
	}

	compilation := core.Compilation{
		Network:             consulNetwork1.Name,
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

	z1ResourcePool := core.ResourcePool{
		Name:            "consul_z1",
		Network:         consulNetwork1.Name,
		Stemcell:        stemcell,
		CloudProperties: iaasConfig.ResourcePool(),
	}

	z1Job := core.Job{
		Name:      "consul_z1",
		Instances: 1,
		Networks: []core.JobNetwork{{
			Name:      consulNetwork1.Name,
			StaticIPs: consulNetwork1.StaticIPs(1),
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
		ResourcePool: z1ResourcePool.Name,
		Templates: []core.JobTemplate{{
			Name:    "consul_agent",
			Release: "consul",
		}},
		Update: &core.JobUpdate{
			MaxInFlight: 1,
		},
	}

	testConsumerJob := core.Job{
		Name:      "consul_test_consumer",
		Instances: 1,
		Networks: []core.JobNetwork{{
			Name:      consulNetwork1.Name,
			StaticIPs: []string{consulNetwork1.StaticIPs(6)[5]},
		}},
		PersistentDisk: 1024,
		ResourcePool:   z1ResourcePool.Name,
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
	}

	properties := Properties{
		Consul: &PropertiesConsul{
			Agent: PropertiesConsulAgent{
				Domain: "cf.internal",
				Servers: PropertiesConsulAgentServers{
					Lan: consulNetwork1.StaticIPs(1),
				},
			},
			CACert:      CACert,
			AgentCert:   AgentCert,
			AgentKey:    AgentKey,
			ServerCert:  ServerCert,
			ServerKey:   ServerKey,
			EncryptKeys: []string{EncryptKey},
		},
	}

	return Manifest{
		DirectorUUID:  config.DirectorUUID,
		Name:          config.Name,
		Releases:      []core.Release{release},
		Compilation:   compilation,
		Update:        update,
		ResourcePools: []core.ResourcePool{z1ResourcePool},
		Jobs:          []core.Job{z1Job, testConsumerJob},
		Networks:      []core.Network{consulNetwork1},
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
	return candiedyaml.Marshal(m)
}

func FromYAML(yaml []byte) (Manifest, error) {
	var m Manifest
	if err := candiedyaml.Unmarshal(yaml, &m); err != nil {
		return m, err
	}
	return m, nil
}
