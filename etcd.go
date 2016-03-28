package destiny

func (m Manifest) EtcdMembers() []EtcdMember {
	members := []EtcdMember{}
	for _, job := range m.Jobs {
		if len(job.Networks) == 0 {
			continue
		}

		for i := 0; i < job.Instances; i++ {
			if len(job.Networks[0].StaticIPs) > i {
				members = append(members, EtcdMember{
					Address: job.Networks[0].StaticIPs[i],
				})
			}
		}
	}

	return members
}

type EtcdMember struct {
	Address string
}

func NewEtcd(config Config) Manifest {
	config = NewConfigWithDefaults(config)

	etcdRelease := Release{
		Name:    "etcd",
		Version: "latest",
	}

	consulRelease := Release{
		Name:    "consul",
		Version: "latest",
	}

	ipRange := IPRange(config.IPRange)
	iaasConfig := IAASConfig(config, "unused")
	cloudProperties := iaasConfig.NetworkSubnet()

	etcdNetwork1 := Network{
		Name: "etcd1",
		Subnets: []NetworkSubnet{{
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

	compilation := Compilation{
		Network:             etcdNetwork1.Name,
		ReuseCompilationVMs: true,
		Workers:             3,
		CloudProperties:     iaasConfig.Compilation(),
	}

	update := Update{
		Canaries:        1,
		CanaryWatchTime: "1000-180000",
		MaxInFlight:     1,
		Serial:          true,
		UpdateWatchTime: "1000-180000",
	}

	stemcell := ResourcePoolStemcell{
		Name:    StemcellForIAAS(config.IAAS),
		Version: "latest",
	}

	z1ResourcePool := ResourcePool{
		Name:            "etcd_z1",
		Network:         etcdNetwork1.Name,
		Stemcell:        stemcell,
		CloudProperties: iaasConfig.ResourcePool(),
	}

	consulZ1Job := Job{
		Name:      "consul_z1",
		Instances: 1,
		Networks: []JobNetwork{{
			Name:      etcdNetwork1.Name,
			StaticIPs: []string{etcdNetwork1.StaticIPs(6)[5]},
		}},
		PersistentDisk: 1024,
		ResourcePool:   z1ResourcePool.Name,
		Templates: []JobTemplate{
			{
				Name:    "consul_agent",
				Release: "consul",
			},
		},
		Properties: &JobProperties{
			Consul: JobPropertiesConsul{
				Agent: JobPropertiesConsulAgent{
					Mode: "server",
				},
			},
		},
	}

	etcdZ1Job := Job{
		Name:      "etcd_z1",
		Instances: 1,
		Networks: []JobNetwork{{
			Name:      etcdNetwork1.Name,
			StaticIPs: etcdNetwork1.StaticIPs(1),
		}},
		PersistentDisk: 1024,
		ResourcePool:   z1ResourcePool.Name,
		Templates: []JobTemplate{
			{
				Name:    "etcd",
				Release: "etcd",
			},
			{
				Name:    "consul_agent",
				Release: "consul",
			},
		},
	}

	return Manifest{
		DirectorUUID: config.DirectorUUID,
		Name:         config.Name,
		Compilation:  compilation,
		Jobs: []Job{
			consulZ1Job,
			etcdZ1Job,
		},
		Networks: []Network{
			etcdNetwork1,
		},
		Properties: Properties{
			Etcd: &PropertiesEtcd{
				Machines:                        etcdNetwork1.StaticIPs(1),
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
			Consul: &PropertiesConsul{
				Agent: PropertiesConsulAgent{
					Domain: "cf.internal",
					Servers: PropertiesConsulAgentServers{
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
		Releases: []Release{
			consulRelease,
			etcdRelease,
		},
		ResourcePools: []ResourcePool{
			z1ResourcePool,
		},
		Update: update,
	}
}
