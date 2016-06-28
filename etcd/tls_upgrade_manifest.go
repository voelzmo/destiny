package etcd

import (
	"github.com/pivotal-cf-experimental/destiny/consul"
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/destiny/network"
)

func NewTLSUpgradeManifest(config Config, iaasConfig iaas.Config) Manifest {
	config = NewConfigWithDefaults(config)

	ipRange := network.IPRange(config.IPRange)
	cloudProperties := iaasConfig.NetworkSubnet()

	releases := []core.Release{
		{
			Name:    "etcd",
			Version: "latest",
		},
		{
			Name:    "consul",
			Version: "latest",
		},
	}

	compilation := core.Compilation{
		Network:             "etcd1",
		ReuseCompilationVMs: true,
		Workers:             3,
	}

	update := core.Update{
		Canaries:        1,
		CanaryWatchTime: "1000-180000",
		MaxInFlight:     1,
		Serial:          true,
		UpdateWatchTime: "1000-180000",
	}

	resourcePools := []core.ResourcePool{
		{
			Name:    "etcd_z1",
			Network: "etcd1",
			Stemcell: core.ResourcePoolStemcell{
				Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
				Version: "latest",
			},
		},
	}

	networks := []core.Network{
		{
			Name: "etcd1",
			Subnets: []core.NetworkSubnet{{
				CloudProperties: cloudProperties,
				Gateway:         ipRange.IP(1),
				Range:           string(ipRange),
				Reserved:        []string{ipRange.Range(2, 3), ipRange.Range(22, 254)},
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
				},
			}},
			Type: "manual",
		},
	}

	consulJob := core.Job{
		Name:      "consul_z1",
		Instances: 3,
		Networks: []core.JobNetwork{{
			Name: "etcd1",
			StaticIPs: []string{
				ipRange.IP(9),
				ipRange.IP(10),
				ipRange.IP(11),
			},
		}},
		PersistentDisk: 1024,
		ResourcePool:   "etcd_z1",
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

	etcdTLS := core.Job{
		Name:      "etcd_tls_z1",
		Instances: 3,
		Networks: []core.JobNetwork{{
			Name: "etcd1",
			StaticIPs: []string{
				ipRange.IP(15),
				ipRange.IP(16),
				ipRange.IP(17),
			},
		}},
		PersistentDisk: 1024,
		ResourcePool:   "etcd_z1",
		Templates: []core.JobTemplate{
			{
				Name:    "consul_agent",
				Release: "consul",
			},
			{
				Name:    "etcd",
				Release: "etcd",
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
			Etcd: &core.JobPropertiesEtcd{
				Cluster: []core.JobPropertiesEtcdCluster{{
					Instances: 3,
					Name:      "etcd_tls_z1",
				}},
				Machines: []string{
					"etcd.service.cf.internal",
				},
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
		},
	}

	etcdNoTLS := core.Job{
		Name:      "etcd_z1",
		Instances: 1,
		Networks: []core.JobNetwork{{
			Name: "etcd1",
			StaticIPs: []string{
				ipRange.IP(4),
			},
		}},
		PersistentDisk: 1024,
		ResourcePool:   "etcd_z1",
		Templates: []core.JobTemplate{
			{
				Name:    "etcd_proxy",
				Release: "etcd",
			},
		},
	}

	testconsumerJob := core.Job{
		Name:      "testconsumer_z1",
		Instances: 3,
		Networks: []core.JobNetwork{{
			Name: "etcd1",
			StaticIPs: []string{
				ipRange.IP(12),
				ipRange.IP(13),
				ipRange.IP(14),
			},
		}},
		PersistentDisk: 1024,
		ResourcePool:   "etcd_z1",
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
		Properties: &core.JobProperties{
			EtcdTestConsumer: &core.JobPropertiesEtcdTestConsumer{
				Etcd: core.JobPropertiesEtcdTestConsumerEtcd{
					Machines: []string{
						"etcd.service.cf.internal",
					},
					RequireSSL: true,
					CACert:     config.Secrets.Etcd.CACert,
					ClientCert: config.Secrets.Etcd.ClientCert,
					ClientKey:  config.Secrets.Etcd.ClientKey,
				},
			},
		},
	}

	properties := Properties{
		Consul: &consul.PropertiesConsul{
			Agent: consul.PropertiesConsulAgent{
				Domain: "cf.internal",
				Servers: consul.PropertiesConsulAgentServers{
					Lan: []string{
						ipRange.IP(9),
						ipRange.IP(10),
						ipRange.IP(11),
					},
				},
			},
			CACert:      config.Secrets.Consul.CACert,
			AgentCert:   config.Secrets.Consul.AgentCert,
			AgentKey:    config.Secrets.Consul.AgentKey,
			ServerCert:  config.Secrets.Consul.ServerCert,
			ServerKey:   config.Secrets.Consul.ServerKey,
			EncryptKeys: []string{config.Secrets.Consul.EncryptKey},
		},
		EtcdProxy: &PropertiesEtcdProxy{
			Etcd: PropertiesEtcdProxyEtcd{
				URL:        "https://etcd.service.cf.internal",
				CACert:     config.Secrets.Etcd.CACert,
				ClientCert: config.Secrets.Etcd.ClientCert,
				ClientKey:  config.Secrets.Etcd.ClientKey,
				RequireSSL: true,
				Port:       4001,
			},
		},
	}

	return Manifest{
		DirectorUUID:  config.DirectorUUID,
		Name:          config.Name,
		Releases:      releases,
		Compilation:   compilation,
		Update:        update,
		ResourcePools: resourcePools,
		Networks:      networks,
		Jobs: []core.Job{
			consulJob,
			etcdTLS,
			etcdNoTLS,
			testconsumerJob,
		},
		Properties: properties,
	}
}
