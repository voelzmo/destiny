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

	consulRelease := core.Release{
		Name:    "consul",
		Version: "latest",
	}

	etcdNetwork1 := etcdNetwork1(iaasConfig.NetworkSubnet(), network.IPRange(config.IPRange))

	compilation := compilation(etcdNetwork1.Name, iaasConfig.Compilation())

	z1ResourcePool := z1ResourcePool(etcdNetwork1.Name, iaasConfig.ResourcePool(), iaasConfig.Stemcell())

	consulZ1Job := core.Job{
		Name:      "consul_z1",
		Instances: 1,
		Networks: []core.JobNetwork{{
			Name:      etcdNetwork1.Name,
			StaticIPs: []string{etcdNetwork1.StaticIPs(6, config.IPOffset)[5]},
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
			Consul: core.JobPropertiesConsul{
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
			StaticIPs: etcdNetwork1.StaticIPs(1, config.IPOffset),
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
			Consul: core.JobPropertiesConsul{
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
			StaticIPs: []string{etcdNetwork1.StaticIPs(7, config.IPOffset)[6]},
		}},
		PersistentDisk: 1024,
		ResourcePool:   z1ResourcePool.Name,
		Templates: []core.JobTemplate{
			{
				Name:    "consul_agent",
				Release: "consul",
			},
			{
				Name:    "etcd-test-consumer",
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
						Lan: []string{etcdNetwork1.StaticIPs(6, config.IPOffset)[5]},
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
			etcdRelease(),
		},
		ResourcePools: []core.ResourcePool{
			z1ResourcePool,
		},
		Update: update(),
	}
}

func NewManifest(config Config, iaasConfig iaas.Config) Manifest {
	config = NewConfigWithDefaults(config)

	etcdNetwork1 := etcdNetwork1(iaasConfig.NetworkSubnet(), network.IPRange(config.IPRange))

	compilation := compilation(etcdNetwork1.Name, iaasConfig.Compilation())

	z1ResourcePool := z1ResourcePool(etcdNetwork1.Name, iaasConfig.ResourcePool(), iaasConfig.Stemcell())

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
			StaticIPs: etcdNetwork1.StaticIPs(1, config.IPOffset),
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
			StaticIPs: []string{etcdNetwork1.StaticIPs(7, config.IPOffset)[6]},
		}},
		PersistentDisk: 1024,
		ResourcePool:   z1ResourcePool.Name,
		Templates: []core.JobTemplate{
			{
				Name:    "etcd-test-consumer",
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
				Machines:                        etcdNetwork1.StaticIPs(1, config.IPOffset),
				PeerRequireSSL:                  false,
				RequireSSL:                      false,
				HeartbeatIntervalInMilliseconds: 50,
			},
		},
		Releases: []core.Release{
			etcdRelease(),
		},
		ResourcePools: []core.ResourcePool{
			z1ResourcePool,
		},
		Update: update(),
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

func etcdRelease() core.Release {
	return core.Release{
		Name:    "etcd",
		Version: "latest",
	}
}

func etcdNetwork1(cloudProperties core.NetworkSubnetCloudProperties, ipRange network.IPRange) core.Network {
	return core.Network{
		Name: "etcd1",
		Subnets: []core.NetworkSubnet{{
			CloudProperties: cloudProperties,
			Gateway:         ipRange.IP(1),
			Range:           string(ipRange),
			Reserved:        []string{ipRange.Range(2, 3), ipRange.Range(104, 254)},
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
				ipRange.IP(41),
				ipRange.IP(42),
				ipRange.IP(43),
				ipRange.IP(44),
				ipRange.IP(45),
				ipRange.IP(46),
				ipRange.IP(47),
				ipRange.IP(48),
				ipRange.IP(49),
				ipRange.IP(50),
				ipRange.IP(51),
				ipRange.IP(52),
				ipRange.IP(53),
				ipRange.IP(54),
				ipRange.IP(55),
				ipRange.IP(56),
				ipRange.IP(57),
				ipRange.IP(58),
				ipRange.IP(59),
				ipRange.IP(60),
				ipRange.IP(61),
				ipRange.IP(62),
				ipRange.IP(63),
				ipRange.IP(64),
				ipRange.IP(65),
				ipRange.IP(66),
				ipRange.IP(67),
				ipRange.IP(68),
				ipRange.IP(69),
				ipRange.IP(70),
				ipRange.IP(71),
				ipRange.IP(72),
				ipRange.IP(73),
				ipRange.IP(74),
				ipRange.IP(75),
				ipRange.IP(76),
				ipRange.IP(77),
				ipRange.IP(78),
				ipRange.IP(79),
				ipRange.IP(80),
				ipRange.IP(81),
				ipRange.IP(82),
				ipRange.IP(83),
				ipRange.IP(84),
				ipRange.IP(85),
				ipRange.IP(86),
				ipRange.IP(87),
				ipRange.IP(88),
				ipRange.IP(89),
				ipRange.IP(90),
				ipRange.IP(91),
				ipRange.IP(92),
				ipRange.IP(93),
				ipRange.IP(94),
				ipRange.IP(95),
				ipRange.IP(96),
				ipRange.IP(97),
				ipRange.IP(98),
				ipRange.IP(99),
				ipRange.IP(100),
			},
		}},
		Type: "manual",
	}
}

func update() core.Update {
	return core.Update{
		Canaries:        1,
		CanaryWatchTime: "1000-180000",
		MaxInFlight:     1,
		Serial:          true,
		UpdateWatchTime: "1000-180000",
	}
}

func compilation(networkName string, cloudProperties core.CompilationCloudProperties) core.Compilation {
	return core.Compilation{
		Network:             networkName,
		ReuseCompilationVMs: true,
		Workers:             3,
		CloudProperties:     cloudProperties,
	}
}

func z1ResourcePool(networkName string, resourcePool core.ResourcePoolCloudProperties, stemcellName string) core.ResourcePool {
	stemcell := core.ResourcePoolStemcell{
		Name:    stemcellName,
		Version: "latest",
	}

	return core.ResourcePool{
		Name:            "etcd_z1",
		Network:         networkName,
		Stemcell:        stemcell,
		CloudProperties: resourcePool,
	}
}
