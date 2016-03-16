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
	release := Release{
		Name:    "etcd",
		Version: "latest",
	}

	ipRange := IPRange(config.IPRange)
	iaasConfig := IAASConfig(config)
	cloudProperties := iaasConfig.NetworkSubnet()

	etcdNetwork1 := Network{
		Name: "etcd1",
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

	z1Job := Job{
		Name:      "etcd_z1",
		Instances: 1,
		Networks: []JobNetwork{{
			Name:      etcdNetwork1.Name,
			StaticIPs: etcdNetwork1.StaticIPs(1),
		}},
		PersistentDisk: 1024,
		ResourcePool:   z1ResourcePool.Name,
		Templates: []JobTemplate{{
			Name:    "etcd",
			Release: "etcd",
		}},
	}

	return Manifest{
		DirectorUUID: config.DirectorUUID,
		Name:         config.Name,
		Compilation:  compilation,
		Jobs: []Job{
			z1Job,
		},
		Networks: []Network{
			etcdNetwork1,
		},
		Properties: Properties{
			Etcd: &PropertiesEtcd{
				Machines:                        etcdNetwork1.StaticIPs(1),
				PeerRequireSSL:                  false,
				RequireSSL:                      false,
				HeartbeatIntervalInMilliseconds: 50,
			},
		},
		Releases: []Release{
			release,
		},
		ResourcePools: []ResourcePool{
			z1ResourcePool,
		},
		Update: update,
	}
}
