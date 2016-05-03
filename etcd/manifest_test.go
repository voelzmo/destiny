package etcd_test

import (
	"io/ioutil"

	"github.com/pivotal-cf-experimental/destiny/consul"
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/etcd"
	"github.com/pivotal-cf-experimental/destiny/iaas"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/gomegamatchers"
)

var _ = Describe("Manifest", func() {
	Describe("NewTLSManifest", func() {
		It("generates a valid Etcd BOSH-Lite manifest", func() {
			manifest := etcd.NewTLSManifest(etcd.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "etcd-some-random-guid",
				IPRange:      "10.244.4.0/24",
				IPOffset:     10,
			}, iaas.NewWardenConfig())

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))

			Expect(manifest.Name).To(Equal("etcd-some-random-guid"))

			Expect(manifest.Releases).To(HaveLen(2))

			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "etcd",
				Version: "latest",
			}))

			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "consul",
				Version: "latest",
			}))

			Expect(manifest.Compilation).To(Equal(core.Compilation{
				Network:             "etcd1",
				ReuseCompilationVMs: true,
				Workers:             3,
			}))

			Expect(manifest.Update).To(Equal(core.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     1,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))

			Expect(manifest.ResourcePools).To(HaveLen(1))

			Expect(manifest.ResourcePools).To(ContainElement(core.ResourcePool{
				Name:    "etcd_z1",
				Network: "etcd1",
				Stemcell: core.ResourcePoolStemcell{
					Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
					Version: "latest",
				},
			}))

			Expect(manifest.Jobs).To(HaveLen(3))

			Expect(manifest.Jobs[0]).To(Equal(core.Job{
				Name:      "consul_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.19"},
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
					Consul: core.JobPropertiesConsul{
						Agent: core.JobPropertiesConsulAgent{
							Mode: "server",
						},
					},
				},
			}))

			Expect(manifest.Jobs[1]).To(Equal(core.Job{
				Name:      "etcd_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.14"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
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
			}))

			Expect(manifest.Jobs[2]).To(Equal(core.Job{
				Name:      "testconsumer_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.20"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
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
			}))

			Expect(manifest.Networks).To(HaveLen(1))

			Expect(manifest.Networks).To(ContainElement(core.Network{
				Name: "etcd1",
				Subnets: []core.NetworkSubnet{
					{
						CloudProperties: core.NetworkSubnetCloudProperties{Name: "random"},
						Gateway:         "10.244.4.1",
						Range:           "10.244.4.0/24",
						Reserved: []string{
							"10.244.4.2-10.244.4.3",
							"10.244.4.104-10.244.4.254",
						},
						Static: []string{
							"10.244.4.4",
							"10.244.4.5",
							"10.244.4.6",
							"10.244.4.7",
							"10.244.4.8",
							"10.244.4.9",
							"10.244.4.10",
							"10.244.4.11",
							"10.244.4.12",
							"10.244.4.13",
							"10.244.4.14",
							"10.244.4.15",
							"10.244.4.16",
							"10.244.4.17",
							"10.244.4.18",
							"10.244.4.19",
							"10.244.4.20",
							"10.244.4.21",
							"10.244.4.22",
							"10.244.4.23",
							"10.244.4.24",
							"10.244.4.25",
							"10.244.4.26",
							"10.244.4.27",
							"10.244.4.28",
							"10.244.4.29",
							"10.244.4.30",
							"10.244.4.31",
							"10.244.4.32",
							"10.244.4.33",
							"10.244.4.34",
							"10.244.4.35",
							"10.244.4.36",
							"10.244.4.37",
							"10.244.4.38",
							"10.244.4.39",
							"10.244.4.40",
							"10.244.4.41",
							"10.244.4.42",
							"10.244.4.43",
							"10.244.4.44",
							"10.244.4.45",
							"10.244.4.46",
							"10.244.4.47",
							"10.244.4.48",
							"10.244.4.49",
							"10.244.4.50",
							"10.244.4.51",
							"10.244.4.52",
							"10.244.4.53",
							"10.244.4.54",
							"10.244.4.55",
							"10.244.4.56",
							"10.244.4.57",
							"10.244.4.58",
							"10.244.4.59",
							"10.244.4.60",
							"10.244.4.61",
							"10.244.4.62",
							"10.244.4.63",
							"10.244.4.64",
							"10.244.4.65",
							"10.244.4.66",
							"10.244.4.67",
							"10.244.4.68",
							"10.244.4.69",
							"10.244.4.70",
							"10.244.4.71",
							"10.244.4.72",
							"10.244.4.73",
							"10.244.4.74",
							"10.244.4.75",
							"10.244.4.76",
							"10.244.4.77",
							"10.244.4.78",
							"10.244.4.79",
							"10.244.4.80",
							"10.244.4.81",
							"10.244.4.82",
							"10.244.4.83",
							"10.244.4.84",
							"10.244.4.85",
							"10.244.4.86",
							"10.244.4.87",
							"10.244.4.88",
							"10.244.4.89",
							"10.244.4.90",
							"10.244.4.91",
							"10.244.4.92",
							"10.244.4.93",
							"10.244.4.94",
							"10.244.4.95",
							"10.244.4.96",
							"10.244.4.97",
							"10.244.4.98",
							"10.244.4.99",
							"10.244.4.100",
						},
					},
				},
				Type: "manual",
			}))

			Expect(manifest.Properties.Etcd).To(Equal(&etcd.PropertiesEtcd{
				Cluster: []etcd.PropertiesEtcdCluster{{
					Instances: 1,
					Name:      "etcd_z1",
				}},
				Machines: []string{
					"etcd.service.cf.internal",
				},
				PeerRequireSSL:                  true,
				RequireSSL:                      true,
				HeartbeatIntervalInMilliseconds: 50,
				AdvertiseURLsDNSSuffix:          "etcd.service.cf.internal",
				CACert:                          etcd.CACert,
				ClientCert:                      etcd.ClientCert,
				ClientKey:                       etcd.ClientKey,
				PeerCACert:                      etcd.PeerCACert,
				PeerCert:                        etcd.PeerCert,
				PeerKey:                         etcd.PeerKey,
				ServerCert:                      etcd.ServerCert,
				ServerKey:                       etcd.ServerKey,
			}))

			Expect(manifest.Properties.Consul).To(Equal(&consul.PropertiesConsul{
				Agent: consul.PropertiesConsulAgent{
					Domain: "cf.internal",
					Servers: consul.PropertiesConsulAgentServers{
						Lan: []string{"10.244.4.19"},
					},
				},
				CACert:      consul.CACert,
				AgentCert:   consul.AgentCert,
				AgentKey:    consul.AgentKey,
				ServerCert:  consul.ServerCert,
				ServerKey:   consul.ServerKey,
				EncryptKeys: []string{consul.EncryptKey},
			}))
		})

		It("generates a valid Etcd AWS manifest", func() {
			manifest := etcd.NewTLSManifest(etcd.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "etcd-some-random-guid",
				IPRange:      "10.0.16.0/24",
			}, iaas.AWSConfig{
				Subnet: "subnet-1234",
			})

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))
			Expect(manifest.Name).To(Equal("etcd-some-random-guid"))

			Expect(manifest.Releases).To(HaveLen(2))
			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "etcd",
				Version: "latest",
			}))
			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "consul",
				Version: "latest",
			}))

			Expect(manifest.Compilation).To(Equal(core.Compilation{
				Network:             "etcd1",
				ReuseCompilationVMs: true,
				Workers:             3,
				CloudProperties: core.CompilationCloudProperties{
					InstanceType:     "m3.medium",
					AvailabilityZone: "us-east-1a",
					EphemeralDisk: &core.CompilationCloudPropertiesEphemeralDisk{
						Size: 1024,
						Type: "gp2",
					},
				},
			}))

			Expect(manifest.Update).To(Equal(core.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     1,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))

			Expect(manifest.ResourcePools).To(HaveLen(1))
			Expect(manifest.ResourcePools).To(ContainElement(core.ResourcePool{
				Name:    "etcd_z1",
				Network: "etcd1",
				Stemcell: core.ResourcePoolStemcell{
					Name:    "bosh-aws-xen-hvm-ubuntu-trusty-go_agent",
					Version: "latest",
				},
				CloudProperties: core.ResourcePoolCloudProperties{
					InstanceType:     "m3.medium",
					AvailabilityZone: "us-east-1a",
					EphemeralDisk: &core.ResourcePoolCloudPropertiesEphemeralDisk{
						Size: 1024,
						Type: "gp2",
					},
				},
			}))

			Expect(manifest.Jobs).To(HaveLen(3))

			Expect(manifest.Jobs[0]).To(Equal(core.Job{
				Name:      "consul_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.0.16.9"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{{
					Name:    "consul_agent",
					Release: "consul",
				}},
				Properties: &core.JobProperties{
					Consul: core.JobPropertiesConsul{
						Agent: core.JobPropertiesConsulAgent{
							Mode: "server",
						},
					},
				},
			}))

			Expect(manifest.Jobs[1]).To(Equal(core.Job{
				Name:      "etcd_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.0.16.4"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
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
			}))

			Expect(manifest.Jobs[2]).To(Equal(core.Job{
				Name:      "testconsumer_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.0.16.10"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
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
			}))

			Expect(manifest.Networks).To(HaveLen(1))
			Expect(manifest.Networks).To(ContainElement(core.Network{
				Name: "etcd1",
				Subnets: []core.NetworkSubnet{
					{
						CloudProperties: core.NetworkSubnetCloudProperties{Subnet: "subnet-1234"},
						Gateway:         "10.0.16.1",
						Range:           "10.0.16.0/24",
						Reserved: []string{
							"10.0.16.2-10.0.16.3",
							"10.0.16.104-10.0.16.254",
						},
						Static: []string{
							"10.0.16.4",
							"10.0.16.5",
							"10.0.16.6",
							"10.0.16.7",
							"10.0.16.8",
							"10.0.16.9",
							"10.0.16.10",
							"10.0.16.11",
							"10.0.16.12",
							"10.0.16.13",
							"10.0.16.14",
							"10.0.16.15",
							"10.0.16.16",
							"10.0.16.17",
							"10.0.16.18",
							"10.0.16.19",
							"10.0.16.20",
							"10.0.16.21",
							"10.0.16.22",
							"10.0.16.23",
							"10.0.16.24",
							"10.0.16.25",
							"10.0.16.26",
							"10.0.16.27",
							"10.0.16.28",
							"10.0.16.29",
							"10.0.16.30",
							"10.0.16.31",
							"10.0.16.32",
							"10.0.16.33",
							"10.0.16.34",
							"10.0.16.35",
							"10.0.16.36",
							"10.0.16.37",
							"10.0.16.38",
							"10.0.16.39",
							"10.0.16.40",
							"10.0.16.41",
							"10.0.16.42",
							"10.0.16.43",
							"10.0.16.44",
							"10.0.16.45",
							"10.0.16.46",
							"10.0.16.47",
							"10.0.16.48",
							"10.0.16.49",
							"10.0.16.50",
							"10.0.16.51",
							"10.0.16.52",
							"10.0.16.53",
							"10.0.16.54",
							"10.0.16.55",
							"10.0.16.56",
							"10.0.16.57",
							"10.0.16.58",
							"10.0.16.59",
							"10.0.16.60",
							"10.0.16.61",
							"10.0.16.62",
							"10.0.16.63",
							"10.0.16.64",
							"10.0.16.65",
							"10.0.16.66",
							"10.0.16.67",
							"10.0.16.68",
							"10.0.16.69",
							"10.0.16.70",
							"10.0.16.71",
							"10.0.16.72",
							"10.0.16.73",
							"10.0.16.74",
							"10.0.16.75",
							"10.0.16.76",
							"10.0.16.77",
							"10.0.16.78",
							"10.0.16.79",
							"10.0.16.80",
							"10.0.16.81",
							"10.0.16.82",
							"10.0.16.83",
							"10.0.16.84",
							"10.0.16.85",
							"10.0.16.86",
							"10.0.16.87",
							"10.0.16.88",
							"10.0.16.89",
							"10.0.16.90",
							"10.0.16.91",
							"10.0.16.92",
							"10.0.16.93",
							"10.0.16.94",
							"10.0.16.95",
							"10.0.16.96",
							"10.0.16.97",
							"10.0.16.98",
							"10.0.16.99",
							"10.0.16.100",
						},
					},
				},
				Type: "manual",
			}))

			Expect(manifest.Properties.Etcd).To(Equal(&etcd.PropertiesEtcd{
				Cluster: []etcd.PropertiesEtcdCluster{{
					Instances: 1,
					Name:      "etcd_z1",
				}},
				Machines: []string{
					"etcd.service.cf.internal",
				},
				PeerRequireSSL:                  true,
				RequireSSL:                      true,
				HeartbeatIntervalInMilliseconds: 50,
				AdvertiseURLsDNSSuffix:          "etcd.service.cf.internal",
				CACert:                          etcd.CACert,
				ClientCert:                      etcd.ClientCert,
				ClientKey:                       etcd.ClientKey,
				PeerCACert:                      etcd.PeerCACert,
				PeerCert:                        etcd.PeerCert,
				PeerKey:                         etcd.PeerKey,
				ServerCert:                      etcd.ServerCert,
				ServerKey:                       etcd.ServerKey,
			}))

			Expect(manifest.Properties.Consul).To(Equal(&consul.PropertiesConsul{
				Agent: consul.PropertiesConsulAgent{
					Domain: "cf.internal",
					Servers: consul.PropertiesConsulAgentServers{
						Lan: []string{"10.0.16.9"},
					},
				},
				CACert:      consul.CACert,
				AgentCert:   consul.AgentCert,
				AgentKey:    consul.AgentKey,
				ServerCert:  consul.ServerCert,
				ServerKey:   consul.ServerKey,
				EncryptKeys: []string{consul.EncryptKey},
			}))
		})
	})

	Describe("NewManifest", func() {
		It("generates a valid Etcd BOSH-Lite manifest", func() {
			manifest := etcd.NewManifest(etcd.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "etcd-some-random-guid",
				IPRange:      "10.244.4.0/24",
				IPOffset:     5,
			}, iaas.NewWardenConfig())

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))

			Expect(manifest.Name).To(Equal("etcd-some-random-guid"))

			Expect(manifest.Releases).To(HaveLen(1))

			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "etcd",
				Version: "latest",
			}))

			Expect(manifest.Compilation).To(Equal(core.Compilation{
				Network:             "etcd1",
				ReuseCompilationVMs: true,
				Workers:             3,
			}))

			Expect(manifest.Update).To(Equal(core.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     1,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))

			Expect(manifest.ResourcePools).To(HaveLen(1))

			Expect(manifest.ResourcePools).To(ContainElement(core.ResourcePool{
				Name:    "etcd_z1",
				Network: "etcd1",
				Stemcell: core.ResourcePoolStemcell{
					Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
					Version: "latest",
				},
			}))

			Expect(manifest.Jobs).To(HaveLen(3))

			Expect(manifest.Jobs[0]).To(Equal(core.Job{
				Name:      "consul_z1",
				Instances: 0,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
			}))

			Expect(manifest.Jobs[1]).To(Equal(core.Job{
				Name:      "etcd_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.9"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "etcd",
						Release: "etcd",
					},
				},
			}))

			Expect(manifest.Jobs[2]).To(Equal(core.Job{
				Name:      "testconsumer_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.15"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "etcd-test-consumer",
						Release: "etcd",
					},
				},
			}))

			Expect(manifest.Networks).To(HaveLen(1))

			Expect(manifest.Networks).To(ContainElement(core.Network{
				Name: "etcd1",
				Subnets: []core.NetworkSubnet{
					{
						CloudProperties: core.NetworkSubnetCloudProperties{Name: "random"},
						Gateway:         "10.244.4.1",
						Range:           "10.244.4.0/24",
						Reserved: []string{
							"10.244.4.2-10.244.4.3",
							"10.244.4.104-10.244.4.254",
						},
						Static: []string{
							"10.244.4.4",
							"10.244.4.5",
							"10.244.4.6",
							"10.244.4.7",
							"10.244.4.8",
							"10.244.4.9",
							"10.244.4.10",
							"10.244.4.11",
							"10.244.4.12",
							"10.244.4.13",
							"10.244.4.14",
							"10.244.4.15",
							"10.244.4.16",
							"10.244.4.17",
							"10.244.4.18",
							"10.244.4.19",
							"10.244.4.20",
							"10.244.4.21",
							"10.244.4.22",
							"10.244.4.23",
							"10.244.4.24",
							"10.244.4.25",
							"10.244.4.26",
							"10.244.4.27",
							"10.244.4.28",
							"10.244.4.29",
							"10.244.4.30",
							"10.244.4.31",
							"10.244.4.32",
							"10.244.4.33",
							"10.244.4.34",
							"10.244.4.35",
							"10.244.4.36",
							"10.244.4.37",
							"10.244.4.38",
							"10.244.4.39",
							"10.244.4.40",
							"10.244.4.41",
							"10.244.4.42",
							"10.244.4.43",
							"10.244.4.44",
							"10.244.4.45",
							"10.244.4.46",
							"10.244.4.47",
							"10.244.4.48",
							"10.244.4.49",
							"10.244.4.50",
							"10.244.4.51",
							"10.244.4.52",
							"10.244.4.53",
							"10.244.4.54",
							"10.244.4.55",
							"10.244.4.56",
							"10.244.4.57",
							"10.244.4.58",
							"10.244.4.59",
							"10.244.4.60",
							"10.244.4.61",
							"10.244.4.62",
							"10.244.4.63",
							"10.244.4.64",
							"10.244.4.65",
							"10.244.4.66",
							"10.244.4.67",
							"10.244.4.68",
							"10.244.4.69",
							"10.244.4.70",
							"10.244.4.71",
							"10.244.4.72",
							"10.244.4.73",
							"10.244.4.74",
							"10.244.4.75",
							"10.244.4.76",
							"10.244.4.77",
							"10.244.4.78",
							"10.244.4.79",
							"10.244.4.80",
							"10.244.4.81",
							"10.244.4.82",
							"10.244.4.83",
							"10.244.4.84",
							"10.244.4.85",
							"10.244.4.86",
							"10.244.4.87",
							"10.244.4.88",
							"10.244.4.89",
							"10.244.4.90",
							"10.244.4.91",
							"10.244.4.92",
							"10.244.4.93",
							"10.244.4.94",
							"10.244.4.95",
							"10.244.4.96",
							"10.244.4.97",
							"10.244.4.98",
							"10.244.4.99",
							"10.244.4.100",
						},
					},
				},
				Type: "manual",
			}))

			Expect(manifest.Properties.Etcd).To(Equal(&etcd.PropertiesEtcd{
				Cluster: []etcd.PropertiesEtcdCluster{{
					Instances: 1,
					Name:      "etcd_z1",
				}},
				Machines: []string{
					"10.244.4.9",
				},
				PeerRequireSSL:                  false,
				RequireSSL:                      false,
				HeartbeatIntervalInMilliseconds: 50,
			}))
		})

		It("generates a valid Etcd AWS manifest", func() {
			manifest := etcd.NewManifest(etcd.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "etcd-some-random-guid",
				IPRange:      "10.0.16.0/24",
			}, iaas.AWSConfig{
				Subnet: "subnet-1234",
			})

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))
			Expect(manifest.Name).To(Equal("etcd-some-random-guid"))

			Expect(manifest.Releases).To(HaveLen(1))
			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "etcd",
				Version: "latest",
			}))

			Expect(manifest.Compilation).To(Equal(core.Compilation{
				Network:             "etcd1",
				ReuseCompilationVMs: true,
				Workers:             3,
				CloudProperties: core.CompilationCloudProperties{
					InstanceType:     "m3.medium",
					AvailabilityZone: "us-east-1a",
					EphemeralDisk: &core.CompilationCloudPropertiesEphemeralDisk{
						Size: 1024,
						Type: "gp2",
					},
				},
			}))

			Expect(manifest.Update).To(Equal(core.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     1,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))

			Expect(manifest.ResourcePools).To(HaveLen(1))
			Expect(manifest.ResourcePools).To(ContainElement(core.ResourcePool{
				Name:    "etcd_z1",
				Network: "etcd1",
				Stemcell: core.ResourcePoolStemcell{
					Name:    "bosh-aws-xen-hvm-ubuntu-trusty-go_agent",
					Version: "latest",
				},
				CloudProperties: core.ResourcePoolCloudProperties{
					InstanceType:     "m3.medium",
					AvailabilityZone: "us-east-1a",
					EphemeralDisk: &core.ResourcePoolCloudPropertiesEphemeralDisk{
						Size: 1024,
						Type: "gp2",
					},
				},
			}))

			Expect(manifest.Jobs).To(HaveLen(3))

			Expect(manifest.Jobs[0]).To(Equal(core.Job{
				Name:      "consul_z1",
				Instances: 0,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
			}))

			Expect(manifest.Jobs[1]).To(Equal(core.Job{
				Name:      "etcd_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.0.16.4"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "etcd",
						Release: "etcd",
					},
				},
			}))

			Expect(manifest.Jobs[2]).To(Equal(core.Job{
				Name:      "testconsumer_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.0.16.10"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "etcd-test-consumer",
						Release: "etcd",
					},
				},
			}))

			Expect(manifest.Networks).To(HaveLen(1))
			Expect(manifest.Networks).To(ContainElement(core.Network{
				Name: "etcd1",
				Subnets: []core.NetworkSubnet{
					{
						CloudProperties: core.NetworkSubnetCloudProperties{Subnet: "subnet-1234"},
						Gateway:         "10.0.16.1",
						Range:           "10.0.16.0/24",
						Reserved: []string{
							"10.0.16.2-10.0.16.3",
							"10.0.16.104-10.0.16.254",
						},
						Static: []string{
							"10.0.16.4",
							"10.0.16.5",
							"10.0.16.6",
							"10.0.16.7",
							"10.0.16.8",
							"10.0.16.9",
							"10.0.16.10",
							"10.0.16.11",
							"10.0.16.12",
							"10.0.16.13",
							"10.0.16.14",
							"10.0.16.15",
							"10.0.16.16",
							"10.0.16.17",
							"10.0.16.18",
							"10.0.16.19",
							"10.0.16.20",
							"10.0.16.21",
							"10.0.16.22",
							"10.0.16.23",
							"10.0.16.24",
							"10.0.16.25",
							"10.0.16.26",
							"10.0.16.27",
							"10.0.16.28",
							"10.0.16.29",
							"10.0.16.30",
							"10.0.16.31",
							"10.0.16.32",
							"10.0.16.33",
							"10.0.16.34",
							"10.0.16.35",
							"10.0.16.36",
							"10.0.16.37",
							"10.0.16.38",
							"10.0.16.39",
							"10.0.16.40",
							"10.0.16.41",
							"10.0.16.42",
							"10.0.16.43",
							"10.0.16.44",
							"10.0.16.45",
							"10.0.16.46",
							"10.0.16.47",
							"10.0.16.48",
							"10.0.16.49",
							"10.0.16.50",
							"10.0.16.51",
							"10.0.16.52",
							"10.0.16.53",
							"10.0.16.54",
							"10.0.16.55",
							"10.0.16.56",
							"10.0.16.57",
							"10.0.16.58",
							"10.0.16.59",
							"10.0.16.60",
							"10.0.16.61",
							"10.0.16.62",
							"10.0.16.63",
							"10.0.16.64",
							"10.0.16.65",
							"10.0.16.66",
							"10.0.16.67",
							"10.0.16.68",
							"10.0.16.69",
							"10.0.16.70",
							"10.0.16.71",
							"10.0.16.72",
							"10.0.16.73",
							"10.0.16.74",
							"10.0.16.75",
							"10.0.16.76",
							"10.0.16.77",
							"10.0.16.78",
							"10.0.16.79",
							"10.0.16.80",
							"10.0.16.81",
							"10.0.16.82",
							"10.0.16.83",
							"10.0.16.84",
							"10.0.16.85",
							"10.0.16.86",
							"10.0.16.87",
							"10.0.16.88",
							"10.0.16.89",
							"10.0.16.90",
							"10.0.16.91",
							"10.0.16.92",
							"10.0.16.93",
							"10.0.16.94",
							"10.0.16.95",
							"10.0.16.96",
							"10.0.16.97",
							"10.0.16.98",
							"10.0.16.99",
							"10.0.16.100",
						},
					},
				},
				Type: "manual",
			}))

			Expect(manifest.Properties.Etcd).To(Equal(&etcd.PropertiesEtcd{
				Cluster: []etcd.PropertiesEtcdCluster{{
					Instances: 1,
					Name:      "etcd_z1",
				}},
				Machines: []string{
					"10.0.16.4",
				},
				PeerRequireSSL:                  false,
				RequireSSL:                      false,
				HeartbeatIntervalInMilliseconds: 50,
			}))
		})
	})

	Describe("EtcdMembers", func() {
		Context("when there is a single job with a single instance", func() {
			It("returns a list of members in the cluster", func() {
				manifest := etcd.Manifest{
					Jobs: []core.Job{
						{
							Instances: 1,
							Networks: []core.JobNetwork{{
								StaticIPs: []string{"10.244.4.2"},
							}},
							Templates: []core.JobTemplate{
								{
									Name:    "etcd",
									Release: "etcd",
								},
							},
						},
					},
				}

				members := manifest.EtcdMembers()
				Expect(members).To(Equal([]etcd.EtcdMember{{
					Address: "10.244.4.2",
				}}))
			})
		})

		Context("when there are multiple jobs with multiple instances", func() {
			It("returns a list of etcd members in the cluster", func() {
				manifest := etcd.Manifest{
					Jobs: []core.Job{
						{
							Instances: 0,
						},
						{
							Instances: 1,
							Networks: []core.JobNetwork{{
								StaticIPs: []string{"10.244.4.2"},
							}},
							Templates: []core.JobTemplate{
								{
									Name:    "consul_agent",
									Release: "consul",
								},
							},
						},
						{
							Instances: 2,
							Networks: []core.JobNetwork{{
								StaticIPs: []string{"10.244.5.2", "10.244.5.6"},
							}},
							Templates: []core.JobTemplate{
								{
									Name:    "etcd",
									Release: "etcd",
								},
							},
						},
						{
							Instances: 2,
							Networks: []core.JobNetwork{{
								StaticIPs: []string{"10.244.6.2", "10.244.6.6"},
							}},
							Templates: []core.JobTemplate{
								{
									Name:    "etcd",
									Release: "etcd",
								},
							},
						},
					},
				}

				members := manifest.EtcdMembers()
				Expect(members).To(Equal([]etcd.EtcdMember{
					{
						Address: "10.244.5.2",
					},
					{
						Address: "10.244.5.6",
					},
					{
						Address: "10.244.6.2",
					},
					{
						Address: "10.244.6.6",
					},
				}))
			})
		})

		Context("when the job does not have a network", func() {
			It("returns an empty list", func() {
				manifest := etcd.Manifest{
					Jobs: []core.Job{
						{
							Instances: 1,
							Networks:  []core.JobNetwork{},
						},
					},
				}

				members := manifest.EtcdMembers()
				Expect(members).To(BeEmpty())
			})
		})

		Context("when the job network does not have enough static IPs", func() {
			It("returns as much about the list as possible", func() {
				manifest := etcd.Manifest{
					Jobs: []core.Job{
						{
							Instances: 2,
							Networks: []core.JobNetwork{{
								StaticIPs: []string{"10.244.5.2"},
							}},
							Templates: []core.JobTemplate{
								{
									Name:    "etcd",
									Release: "etcd",
								},
							},
						},
					},
				}

				members := manifest.EtcdMembers()
				Expect(members).To(Equal([]etcd.EtcdMember{
					{
						Address: "10.244.5.2",
					},
				}))
			})
		})
	})

	Describe("ToYAML", func() {
		It("returns a YAML representation of the etcd manifest", func() {
			etcdManifest, err := ioutil.ReadFile("fixtures/etcd_manifest.yml")
			Expect(err).NotTo(HaveOccurred())

			manifest := etcd.NewTLSManifest(etcd.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "etcd",
				IPRange:      "10.244.4.0/24",
				IPOffset:     30,
				Secrets: etcd.ConfigSecrets{
					Consul: consul.ConfigSecretsConsul{
						EncryptKey: "consul-encrypt-key",
						AgentKey:   "consul-agent-key",
						AgentCert:  "consul-agent-cert",
						ServerKey:  "consul-server-key",
						ServerCert: "consul-server-cert",
						CACert:     "consul-ca-cert",
					},
					Etcd: etcd.ConfigSecretsEtcd{
						CACert:     "etcd-ca-cert",
						ClientCert: "etcd-client-cert",
						ClientKey:  "etcd-client-key",
						PeerCACert: "etcd-peer-ca-cert",
						PeerCert:   "etcd-peer-cert",
						PeerKey:    "etcd-peer-key",
						ServerCert: "etcd-server-cert",
						ServerKey:  "etcd-server-key",
					},
				},
			}, iaas.NewWardenConfig())

			yaml, err := manifest.ToYAML()
			Expect(err).NotTo(HaveOccurred())
			Expect(yaml).To(MatchYAML(etcdManifest))
		})
	})

	Describe("FromYAML", func() {
		It("returns a Manifest matching the given YAML", func() {
			etcdManifest, err := ioutil.ReadFile("fixtures/etcd_manifest.yml")
			Expect(err).NotTo(HaveOccurred())

			manifest, err := etcd.FromYAML(etcdManifest)
			Expect(err).NotTo(HaveOccurred())

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))
			Expect(manifest.Name).To(Equal("etcd"))
			Expect(manifest.Releases).To(HaveLen(2))
			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "etcd",
				Version: "latest",
			}))
			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "consul",
				Version: "latest",
			}))
			Expect(manifest.Compilation).To(Equal(core.Compilation{
				Network:             "etcd1",
				ReuseCompilationVMs: true,
				Workers:             3,
			}))
			Expect(manifest.Update).To(Equal(core.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     1,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))
			Expect(manifest.ResourcePools).To(HaveLen(1))
			Expect(manifest.ResourcePools).To(ContainElement(core.ResourcePool{
				Name:    "etcd_z1",
				Network: "etcd1",
				Stemcell: core.ResourcePoolStemcell{
					Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
					Version: "latest",
				},
			}))

			Expect(manifest.Jobs).To(HaveLen(3))

			Expect(manifest.Jobs[0]).To(Equal(core.Job{
				Name:      "consul_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.39"},
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
					Consul: core.JobPropertiesConsul{
						Agent: core.JobPropertiesConsulAgent{
							Mode: "server",
						},
					},
				},
			}))

			Expect(manifest.Jobs[1]).To(Equal(core.Job{
				Name:      "etcd_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.34"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
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
			}))

			Expect(manifest.Jobs[2]).To(Equal(core.Job{
				Name:      "testconsumer_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.40"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
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
			}))

			Expect(manifest.Networks).To(HaveLen(1))
			Expect(manifest.Networks).To(ContainElement(core.Network{
				Name: "etcd1",
				Subnets: []core.NetworkSubnet{
					{
						CloudProperties: core.NetworkSubnetCloudProperties{Name: "random"},
						Gateway:         "10.244.4.1",
						Range:           "10.244.4.0/24",
						Reserved: []string{
							"10.244.4.2-10.244.4.3",
							"10.244.4.104-10.244.4.254",
						},
						Static: []string{
							"10.244.4.4",
							"10.244.4.5",
							"10.244.4.6",
							"10.244.4.7",
							"10.244.4.8",
							"10.244.4.9",
							"10.244.4.10",
							"10.244.4.11",
							"10.244.4.12",
							"10.244.4.13",
							"10.244.4.14",
							"10.244.4.15",
							"10.244.4.16",
							"10.244.4.17",
							"10.244.4.18",
							"10.244.4.19",
							"10.244.4.20",
							"10.244.4.21",
							"10.244.4.22",
							"10.244.4.23",
							"10.244.4.24",
							"10.244.4.25",
							"10.244.4.26",
							"10.244.4.27",
							"10.244.4.28",
							"10.244.4.29",
							"10.244.4.30",
							"10.244.4.31",
							"10.244.4.32",
							"10.244.4.33",
							"10.244.4.34",
							"10.244.4.35",
							"10.244.4.36",
							"10.244.4.37",
							"10.244.4.38",
							"10.244.4.39",
							"10.244.4.40",
							"10.244.4.41",
							"10.244.4.42",
							"10.244.4.43",
							"10.244.4.44",
							"10.244.4.45",
							"10.244.4.46",
							"10.244.4.47",
							"10.244.4.48",
							"10.244.4.49",
							"10.244.4.50",
							"10.244.4.51",
							"10.244.4.52",
							"10.244.4.53",
							"10.244.4.54",
							"10.244.4.55",
							"10.244.4.56",
							"10.244.4.57",
							"10.244.4.58",
							"10.244.4.59",
							"10.244.4.60",
							"10.244.4.61",
							"10.244.4.62",
							"10.244.4.63",
							"10.244.4.64",
							"10.244.4.65",
							"10.244.4.66",
							"10.244.4.67",
							"10.244.4.68",
							"10.244.4.69",
							"10.244.4.70",
							"10.244.4.71",
							"10.244.4.72",
							"10.244.4.73",
							"10.244.4.74",
							"10.244.4.75",
							"10.244.4.76",
							"10.244.4.77",
							"10.244.4.78",
							"10.244.4.79",
							"10.244.4.80",
							"10.244.4.81",
							"10.244.4.82",
							"10.244.4.83",
							"10.244.4.84",
							"10.244.4.85",
							"10.244.4.86",
							"10.244.4.87",
							"10.244.4.88",
							"10.244.4.89",
							"10.244.4.90",
							"10.244.4.91",
							"10.244.4.92",
							"10.244.4.93",
							"10.244.4.94",
							"10.244.4.95",
							"10.244.4.96",
							"10.244.4.97",
							"10.244.4.98",
							"10.244.4.99",
							"10.244.4.100",
						},
					},
				},
				Type: "manual",
			}))
			Expect(manifest.Properties.Etcd).To(Equal(&etcd.PropertiesEtcd{
				Cluster: []etcd.PropertiesEtcdCluster{{
					Instances: 1,
					Name:      "etcd_z1",
				}},
				Machines: []string{
					"etcd.service.cf.internal",
				},
				PeerRequireSSL:                  true,
				RequireSSL:                      true,
				HeartbeatIntervalInMilliseconds: 50,
				AdvertiseURLsDNSSuffix:          "etcd.service.cf.internal",
				CACert:                          "etcd-ca-cert",
				ClientCert:                      "etcd-client-cert",
				ClientKey:                       "etcd-client-key",
				PeerCACert:                      "etcd-peer-ca-cert",
				PeerCert:                        "etcd-peer-cert",
				PeerKey:                         "etcd-peer-key",
				ServerCert:                      "etcd-server-cert",
				ServerKey:                       "etcd-server-key",
			}))
		})

		Context("failure cases", func() {
			It("should error on malformed YAML", func() {
				_, err := etcd.FromYAML([]byte("%%%%%%%%%%"))
				Expect(err).To(MatchError(ContainSubstring("yaml: ")))
			})
		})
	})
})
