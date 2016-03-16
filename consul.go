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
	//iaasConfig := IAASConfig(config)

	//cloudProperties := iaasConfig.NetworkSubnet()
	cloudProperties := NetworkSubnetCloudProperties{Name: "random"}
	if config.IAAS == AWS {
		cloudProperties = NetworkSubnetCloudProperties{Subnet: config.AWS.Subnet}
	}

	consulNetwork1 := Network{
		Name: "consul1",
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
		Network:             consulNetwork1.Name,
		ReuseCompilationVMs: true,
		Workers:             3,
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
		Name:     "consul_z1",
		Network:  consulNetwork1.Name,
		Stemcell: stemcell,
	}

	if config.IAAS == AWS {
		compilation.CloudProperties = CompilationCloudProperties{
			InstanceType:     "m3.medium",
			AvailabilityZone: "us-east-1a",
			EphemeralDisk: &CompilationCloudPropertiesEphemeralDisk{
				Size: 1024,
				Type: "gp2",
			},
		}

		z1ResourcePool.CloudProperties = ResourcePoolCloudProperties{
			InstanceType:     "m3.medium",
			AvailabilityZone: "us-east-1a",
			EphemeralDisk: &ResourcePoolCloudPropertiesEphemeralDisk{
				Size: 1024,
				Type: "gp2",
			},
		}
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
					Mode:     "server",
					LogLevel: "info",
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

	testConsumerJob := Job{
		Name:      "consul_test_consumer",
		Instances: 1,
		Networks: []JobNetwork{{
			Name:      consulNetwork1.Name,
			StaticIPs: []string{consulNetwork1.StaticIPs(6)[5]},
		}},
		PersistentDisk: 1024,
		ResourcePool:   z1ResourcePool.Name,
		Templates: []JobTemplate{
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
		Jobs:          []Job{z1Job, testConsumerJob},
		Networks:      []Network{consulNetwork1},
		Properties:    properties,
	}
}
