package etcd

import (
	"github.com/cloudfoundry-incubator/candiedyaml"
	"github.com/pivotal-cf-experimental/destiny/consul"
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/destiny/network"
)

type Manifest struct {
	DirectorUUID  string              `yaml:"director_uuid"`
	Name          string              `yaml:"name"`
	Jobs          []core.Job          `yaml:"jobs"`
	Properties    Properties          `yaml:"properties"`
	Update        core.Update         `yaml:"update"`
	Compilation   core.Compilation    `yaml:"compilation"`
	Networks      []core.Network      `yaml:"networks"`
	Releases      []core.Release      `yaml:"releases"`
	ResourcePools []core.ResourcePool `yaml:"resource_pools"`
}

type EtcdMember struct {
	Address string
}

func NewTLSManifest(config Config, iaasConfig iaas.Config) Manifest {
	config = NewConfigWithDefaults(config)

	etcdRelease := core.Release{
		Name:    "etcd",
		Version: "latest",
	}

	consulRelease := core.Release{
		Name:    "consul",
		Version: "latest",
	}

	ipRange := network.IPRange(config.IPRange)
	cloudProperties := iaasConfig.NetworkSubnet()

	etcdNetwork1 := core.Network{
		Name: "etcd1",
		Subnets: []core.NetworkSubnet{{
			CloudProperties: cloudProperties,
			Gateway:         ipRange.IP(1),
			Range:           string(ipRange),
			Reserved:        []string{ipRange.Range(2, 3), ipRange.Range(50, 254)},
			Static: []string{
				ipRange.IP(4),
				ipRange.IP(5),
				ipRange.IP(6),
				ipRange.IP(7),
				ipRange.IP(8),
				ipRange.IP(9),
				ipRange.IP(10),
				ipRange.IP(11),
				ipRange.IP(12),
				ipRange.IP(13),
				ipRange.IP(14),
				ipRange.IP(15),
				ipRange.IP(16),
				ipRange.IP(17),
				ipRange.IP(18),
				ipRange.IP(19),
				ipRange.IP(20),
				ipRange.IP(21),
				ipRange.IP(22),
				ipRange.IP(23),
				ipRange.IP(24),
				ipRange.IP(25),
				ipRange.IP(26),
				ipRange.IP(27),
				ipRange.IP(28),
				ipRange.IP(29),
				ipRange.IP(30),
				ipRange.IP(31),
				ipRange.IP(32),
				ipRange.IP(33),
				ipRange.IP(34),
				ipRange.IP(35),
				ipRange.IP(36),
				ipRange.IP(37),
				ipRange.IP(38),
				ipRange.IP(39),
				ipRange.IP(40),
			},
		}},
		Type: "manual",
	}

	compilation := core.Compilation{
		Network:             etcdNetwork1.Name,
		ReuseCompilationVMs: true,
		Workers:             3,
		CloudProperties:     iaasConfig.Compilation(),
	}

	update := core.Update{
		Canaries:        1,
		CanaryWatchTime: "1000-180000",
		MaxInFlight:     1,
		Serial:          true,
		UpdateWatchTime: "1000-180000",
	}

	stemcell := core.ResourcePoolStemcell{
		Name:    iaasConfig.Stemcell(),
		Version: "latest",
	}

	z1ResourcePool := core.ResourcePool{
		Name:            "etcd_z1",
		Network:         etcdNetwork1.Name,
		Stemcell:        stemcell,
		CloudProperties: iaasConfig.ResourcePool(),
	}

	consulZ1Job := core.Job{
		Name:      "consul_z1",
		Instances: 1,
		Networks: []core.JobNetwork{{
			Name:      etcdNetwork1.Name,
			StaticIPs: []string{etcdNetwork1.StaticIPs(6)[5]},
		}},
		PersistentDisk: 1024,
		ResourcePool:   z1ResourcePool.Name,
		Templates: []core.JobTemplate{
			{
				Name:    "consul_agent",
				Release: "consul",
			},
		},
		Properties: &core.JobProperties{
			Consul: &core.JobPropertiesConsul{
				Agent: core.JobPropertiesConsulAgent{
					Mode: "server",
				},
			},
		},
	}

	etcdZ1Job := core.Job{
		Name:      "etcd_z1",
		Instances: 1,
		Networks: []core.JobNetwork{{
			Name:      etcdNetwork1.Name,
			StaticIPs: etcdNetwork1.StaticIPs(1),
		}},
		PersistentDisk: 1024,
		ResourcePool:   z1ResourcePool.Name,
		Templates: []core.JobTemplate{
			{
				Name:    "etcd",
				Release: "etcd",
			},
			{
				Name:    "consul_agent",
				Release: "consul",
			},
		},
		Properties: &core.JobProperties{
			Consul: &core.JobPropertiesConsul{
				Agent: core.JobPropertiesConsulAgent{
					Services: core.JobPropertiesConsulAgentServices{
						"etcd": core.JobPropertiesConsulAgentService{},
					},
				},
			},
		},
	}

	testconsumerZ1Job := core.Job{
		Name:      "testconsumer_z1",
		Instances: 1,
		Networks: []core.JobNetwork{{
			Name:      etcdNetwork1.Name,
			StaticIPs: []string{etcdNetwork1.StaticIPs(9)[8]},
		}},
		PersistentDisk: 1024,
		ResourcePool:   z1ResourcePool.Name,
		Templates: []core.JobTemplate{
			{
				Name:    "consul_agent",
				Release: "consul",
			},
			{
				Name:    "etcd_testconsumer",
				Release: "etcd",
			},
		},
	}

	return Manifest{
		DirectorUUID: config.DirectorUUID,
		Name:         config.Name,
		Compilation:  compilation,
		Jobs: []core.Job{
			consulZ1Job,
			etcdZ1Job,
			testconsumerZ1Job,
		},
		Networks: []core.Network{
			etcdNetwork1,
		},
		Properties: Properties{
			EtcdTestConsumer: &PropertiesEtcdTestConsumer{
				Etcd: PropertiesEtcdTestConsumerEtcd{
					RequireSSL: true,
					Machines:   []string{"etcd.service.cf.internal"},
					CACert:     config.Secrets.Etcd.CACert,
					ClientCert: config.Secrets.Etcd.ClientCert,
					ClientKey:  config.Secrets.Etcd.ClientKey,
				},
			},
			Etcd: &PropertiesEtcd{
				Cluster: []PropertiesEtcdCluster{{
					Instances: 1,
					Name:      "etcd_z1",
				}},
				Machines:                        []string{"etcd.service.cf.internal"},
				PeerRequireSSL:                  true,
				RequireSSL:                      true,
				HeartbeatIntervalInMilliseconds: 50,
				AdvertiseURLsDNSSuffix:          "etcd.service.cf.internal",
				CACert:                          config.Secrets.Etcd.CACert,
				ClientCert:                      config.Secrets.Etcd.ClientCert,
				ClientKey:                       config.Secrets.Etcd.ClientKey,
				PeerCACert:                      config.Secrets.Etcd.PeerCACert,
				PeerCert:                        config.Secrets.Etcd.PeerCert,
				PeerKey:                         config.Secrets.Etcd.PeerKey,
				ServerCert:                      config.Secrets.Etcd.ServerCert,
				ServerKey:                       config.Secrets.Etcd.ServerKey,
			},
			Consul: &consul.PropertiesConsul{
				Agent: consul.PropertiesConsulAgent{
					Domain: "cf.internal",
					Servers: consul.PropertiesConsulAgentServers{
						Lan: []string{etcdNetwork1.StaticIPs(6)[5]},
					},
				},
				CACert:      config.Secrets.Consul.CACert,
				AgentCert:   config.Secrets.Consul.AgentCert,
				AgentKey:    config.Secrets.Consul.AgentKey,
				ServerCert:  config.Secrets.Consul.ServerCert,
				ServerKey:   config.Secrets.Consul.ServerKey,
				EncryptKeys: []string{config.Secrets.Consul.EncryptKey},
			},
		},
		Releases: []core.Release{
			consulRelease,
			etcdRelease,
		},
		ResourcePools: []core.ResourcePool{
			z1ResourcePool,
		},
		Update: update,
	}
}

func NewManifest(config Config, iaasConfig iaas.Config) Manifest {
	config = NewConfigWithDefaults(config)

	etcdRelease := core.Release{
		Name:    "etcd",
		Version: "latest",
	}

	ipRange := network.IPRange(config.IPRange)
	cloudProperties := iaasConfig.NetworkSubnet()

	etcdNetwork1 := core.Network{
		Name: "etcd1",
		Subnets: []core.NetworkSubnet{{
			CloudProperties: cloudProperties,
			Gateway:         ipRange.IP(1),
			Range:           string(ipRange),
			Reserved:        []string{ipRange.Range(2, 3), ipRange.Range(50, 254)},
			Static: []string{
				ipRange.IP(4),
				ipRange.IP(5),
				ipRange.IP(6),
				ipRange.IP(7),
				ipRange.IP(8),
				ipRange.IP(9),
				ipRange.IP(10),
				ipRange.IP(11),
				ipRange.IP(12),
				ipRange.IP(13),
				ipRange.IP(14),
				ipRange.IP(15),
				ipRange.IP(16),
				ipRange.IP(17),
				ipRange.IP(18),
				ipRange.IP(19),
				ipRange.IP(20),
				ipRange.IP(21),
				ipRange.IP(22),
				ipRange.IP(23),
				ipRange.IP(24),
				ipRange.IP(25),
				ipRange.IP(26),
				ipRange.IP(27),
				ipRange.IP(28),
				ipRange.IP(29),
				ipRange.IP(30),
				ipRange.IP(31),
				ipRange.IP(32),
				ipRange.IP(33),
				ipRange.IP(34),
				ipRange.IP(35),
				ipRange.IP(36),
				ipRange.IP(37),
				ipRange.IP(38),
				ipRange.IP(39),
				ipRange.IP(40),
			},
		}},
		Type: "manual",
	}

	compilation := core.Compilation{
		Network:             etcdNetwork1.Name,
		ReuseCompilationVMs: true,
		Workers:             3,
		CloudProperties:     iaasConfig.Compilation(),
	}

	update := core.Update{
		Canaries:        1,
		CanaryWatchTime: "1000-180000",
		MaxInFlight:     1,
		Serial:          true,
		UpdateWatchTime: "1000-180000",
	}

	stemcell := core.ResourcePoolStemcell{
		Name:    iaasConfig.Stemcell(),
		Version: "latest",
	}

	z1ResourcePool := core.ResourcePool{
		Name:            "etcd_z1",
		Network:         etcdNetwork1.Name,
		Stemcell:        stemcell,
		CloudProperties: iaasConfig.ResourcePool(),
	}

	consulZ1Job := core.Job{
		Name:      "consul_z1",
		Instances: 0,
		Networks: []core.JobNetwork{{
			Name:      etcdNetwork1.Name,
			StaticIPs: []string{},
		}},
		PersistentDisk: 1024,
		ResourcePool:   z1ResourcePool.Name,
	}

	etcdZ1Job := core.Job{
		Name:      "etcd_z1",
		Instances: 1,
		Networks: []core.JobNetwork{{
			Name:      etcdNetwork1.Name,
			StaticIPs: etcdNetwork1.StaticIPs(1),
		}},
		PersistentDisk: 1024,
		ResourcePool:   z1ResourcePool.Name,
		Templates: []core.JobTemplate{
			{
				Name:    "etcd",
				Release: "etcd",
			},
		},
	}

	testconsumerZ1Job := core.Job{
		Name:      "testconsumer_z1",
		Instances: 1,
		Networks: []core.JobNetwork{{
			Name:      etcdNetwork1.Name,
			StaticIPs: []string{etcdNetwork1.StaticIPs(9)[8]},
		}},
		PersistentDisk: 1024,
		ResourcePool:   z1ResourcePool.Name,
		Templates: []core.JobTemplate{
			{
				Name:    "etcd_testconsumer",
				Release: "etcd",
			},
		},
	}

	return Manifest{
		DirectorUUID: config.DirectorUUID,
		Name:         config.Name,
		Compilation:  compilation,
		Jobs: []core.Job{
			consulZ1Job,
			etcdZ1Job,
			testconsumerZ1Job,
		},
		Networks: []core.Network{
			etcdNetwork1,
		},
		Properties: Properties{
			Etcd: &PropertiesEtcd{
				Cluster: []PropertiesEtcdCluster{{
					Instances: 1,
					Name:      "etcd_z1",
				}},
				Machines:                        etcdNetwork1.StaticIPs(1),
				PeerRequireSSL:                  false,
				RequireSSL:                      false,
				HeartbeatIntervalInMilliseconds: 50,
			},
			EtcdTestConsumer: &PropertiesEtcdTestConsumer{
				Etcd: PropertiesEtcdTestConsumerEtcd{
					Machines: etcdNetwork1.StaticIPs(1),
				},
			},
		},
		Releases: []core.Release{
			etcdRelease,
		},
		ResourcePools: []core.ResourcePool{
			z1ResourcePool,
		},
		Update: update,
	}
}

func (m Manifest) EtcdMembers() []EtcdMember {
	members := []EtcdMember{}
	for _, job := range m.Jobs {
		if len(job.Networks) == 0 {
			continue
		}

		if job.HasTemplate("etcd", "etcd") {
			for i := 0; i < job.Instances; i++ {
				if len(job.Networks[0].StaticIPs) > i {
					members = append(members, EtcdMember{
						Address: job.Networks[0].StaticIPs[i],
					})
				}
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
