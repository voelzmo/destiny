package destiny

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

type ConsulMember struct {
	Address string
}

func NewConsul(config Config) Manifest {
	release := Release{
		Name:    "consul",
		Version: "latest",
	}

	ipRange := IPRange(config.IPRange)
	iaasConfig := IAASConfig(config)

	cloudProperties := iaasConfig.NetworkSubnet()

	consulNetwork1 := Network{
		Name: "consul1",
		Subnets: []NetworkSubnet{{
			CloudProperties: cloudProperties,
			Gateway:         ipRange.IP(1),
			Range:           string(ipRange),
			Reserved:        []string{ipRange.Range(2, 3), ipRange.Range(12, 254)},
			Static: []string{
				ipRange.IP(4),
				ipRange.IP(5),
				ipRange.IP(6),
				ipRange.IP(7),
				ipRange.IP(8),
			},
		}},
		Type: "manual",
	}

	compilation := Compilation{
		Network:             consulNetwork1.Name,
		ReuseCompilationVMs: true,
		Workers:             3,
		CloudProperties:     iaasConfig.Compilation(),
	}

	update := Update{
		Canaries:        1,
		CanaryWatchTime: "1000-180000",
		MaxInFlight:     50,
		Serial:          true,
		UpdateWatchTime: "1000-180000",
	}

	stemcell := ResourcePoolStemcell{
		Name:    StemcellForIAAS(config.IAAS),
		Version: "latest",
	}

	z1ResourcePool := ResourcePool{
		Name:            "consul_z1",
		Network:         consulNetwork1.Name,
		Stemcell:        stemcell,
		CloudProperties: iaasConfig.ResourcePool(),
	}

	z1Job := Job{
		Name:      "consul_z1",
		Instances: 1,
		Networks: []JobNetwork{{
			Name:      consulNetwork1.Name,
			StaticIPs: consulNetwork1.StaticIPs(1),
		}},
		PersistentDisk: 1024,
		Properties: &JobProperties{
			Consul: JobPropertiesConsul{
				Agent: JobPropertiesConsulAgent{
					Mode: "server",
					Services: JobPropertiesConsulAgentServices{
						"router": JobPropertiesConsulAgentService{
							Name: "gorouter",
							Check: &JobPropertiesConsulAgentServiceCheck{
								Name:     "router-check",
								Script:   "/var/vcap/jobs/router/bin/script",
								Interval: "1m",
							},
							Tags: []string{"routing"},
						},
						"cloud_controller": JobPropertiesConsulAgentService{},
					},
				},
			},
		},
		ResourcePool: z1ResourcePool.Name,
		Templates: []JobTemplate{{
			Name:    "consul_agent",
			Release: "consul",
		}},
		Update: &JobUpdate{
			MaxInFlight: 1,
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
			RequireSSL:  true,
		},
	}

	return Manifest{
		DirectorUUID:  config.DirectorUUID,
		Name:          config.Name,
		Releases:      []Release{release},
		Compilation:   compilation,
		Update:        update,
		ResourcePools: []ResourcePool{z1ResourcePool},
		Jobs:          []Job{z1Job},
		Networks:      []Network{consulNetwork1},
		Properties:    properties,
	}
}
