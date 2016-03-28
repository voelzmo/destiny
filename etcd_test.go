package destiny_test

import (
	"github.com/pivotal-cf-experimental/destiny"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Etcd Manifest", func() {
	Describe("NewEtcd", func() {
		It("generates a valid Etcd BOSH-Lite manifest", func() {
			manifest := destiny.NewEtcd(destiny.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "etcd-some-random-guid",
				IPRange:      "10.244.4.0/24",
				IAAS:         destiny.Warden,
			})

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))
			Expect(manifest.Name).To(Equal("etcd-some-random-guid"))
			Expect(manifest.Releases).To(HaveLen(2))
			Expect(manifest.Releases).To(ContainElement(destiny.Release{
				Name:    "etcd",
				Version: "latest",
			}))
			Expect(manifest.Releases).To(ContainElement(destiny.Release{
				Name:    "consul",
				Version: "latest",
			}))
			Expect(manifest.Compilation).To(Equal(destiny.Compilation{
				Network:             "etcd1",
				ReuseCompilationVMs: true,
				Workers:             3,
			}))
			Expect(manifest.Update).To(Equal(destiny.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     1,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))
			Expect(manifest.ResourcePools).To(HaveLen(1))
			Expect(manifest.ResourcePools).To(ContainElement(destiny.ResourcePool{
				Name:    "etcd_z1",
				Network: "etcd1",
				Stemcell: destiny.ResourcePoolStemcell{
					Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
					Version: "latest",
				},
			}))
			Expect(manifest.Jobs).To(HaveLen(2))
			Expect(manifest.Jobs[0]).To(Equal(destiny.Job{
				Name:      "consul_z1",
				Instances: 1,
				Networks: []destiny.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.9"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []destiny.JobTemplate{
					{
						Name:    "consul_agent",
						Release: "consul",
					},
				},
				Properties: &destiny.JobProperties{
					Consul: destiny.JobPropertiesConsul{
						Agent: destiny.JobPropertiesConsulAgent{
							Mode: "server",
						},
					},
				},
			}))
			Expect(manifest.Jobs[1]).To(Equal(destiny.Job{
				Name:      "etcd_z1",
				Instances: 1,
				Networks: []destiny.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.4"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []destiny.JobTemplate{
					{
						Name:    "etcd",
						Release: "etcd",
					},
					{
						Name:    "consul_agent",
						Release: "consul",
					},
				},
			}))
			Expect(manifest.Networks).To(HaveLen(1))
			Expect(manifest.Networks).To(ContainElement(destiny.Network{
				Name: "etcd1",
				Subnets: []destiny.NetworkSubnet{
					{
						CloudProperties: destiny.NetworkSubnetCloudProperties{Name: "random"},
						Gateway:         "10.244.4.1",
						Range:           "10.244.4.0/24",
						Reserved: []string{
							"10.244.4.2-10.244.4.3",
							"10.244.4.13-10.244.4.254",
						},
						Static: []string{
							"10.244.4.4",
							"10.244.4.5",
							"10.244.4.6",
							"10.244.4.7",
							"10.244.4.8",
							"10.244.4.9",
						},
					},
				},
				Type: "manual",
			}))
			Expect(manifest.Properties.Etcd).To(Equal(&destiny.PropertiesEtcd{
				Machines: []string{
					"10.244.4.4",
				},
				PeerRequireSSL:                  true,
				RequireSSL:                      true,
				HeartbeatIntervalInMilliseconds: 50,
				AdvertiseURLsDNSSuffix:          "etcd.service.cf.internal",
				CACert:                          destiny.EtcdCACert,
				ClientCert:                      destiny.EtcdClientCert,
				ClientKey:                       destiny.EtcdClientKey,
				PeerCACert:                      destiny.EtcdPeerCACert,
				PeerCert:                        destiny.EtcdPeerCert,
				PeerKey:                         destiny.EtcdPeerKey,
				ServerCert:                      destiny.EtcdServerCert,
				ServerKey:                       destiny.EtcdServerKey,
			}))
			Expect(manifest.Properties.Consul).To(Equal(&destiny.PropertiesConsul{
				Agent: destiny.PropertiesConsulAgent{
					Domain: "cf.internal",
					Servers: destiny.PropertiesConsulAgentServers{
						Lan: []string{"10.244.4.9"},
					},
				},
				CACert:      destiny.ConsulCACert,
				AgentCert:   destiny.ConsulAgentCert,
				AgentKey:    destiny.ConsulAgentKey,
				ServerCert:  destiny.ConsulServerCert,
				ServerKey:   destiny.ConsulServerKey,
				EncryptKeys: []string{destiny.ConsulEncryptKey},
			}))
		})

		It("generates a valid Etcd AWS manifest", func() {
			manifest := destiny.NewEtcd(destiny.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "etcd-some-random-guid",
				IPRange:      "10.0.16.0/24",
				IAAS:         destiny.AWS,
				AWS: destiny.ConfigAWS{
					Subnet: "subnet-1234",
				},
			})

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))
			Expect(manifest.Name).To(Equal("etcd-some-random-guid"))

			Expect(manifest.Releases).To(HaveLen(2))
			Expect(manifest.Releases).To(ContainElement(destiny.Release{
				Name:    "etcd",
				Version: "latest",
			}))
			Expect(manifest.Releases).To(ContainElement(destiny.Release{
				Name:    "consul",
				Version: "latest",
			}))

			Expect(manifest.Compilation).To(Equal(destiny.Compilation{
				Network:             "etcd1",
				ReuseCompilationVMs: true,
				Workers:             3,
				CloudProperties: destiny.CompilationCloudProperties{
					InstanceType:     "m3.medium",
					AvailabilityZone: "us-east-1a",
					EphemeralDisk: &destiny.CompilationCloudPropertiesEphemeralDisk{
						Size: 1024,
						Type: "gp2",
					},
				},
			}))

			Expect(manifest.Update).To(Equal(destiny.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     1,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))

			Expect(manifest.ResourcePools).To(HaveLen(1))
			Expect(manifest.ResourcePools).To(ContainElement(destiny.ResourcePool{
				Name:    "etcd_z1",
				Network: "etcd1",
				Stemcell: destiny.ResourcePoolStemcell{
					Name:    "bosh-aws-xen-hvm-ubuntu-trusty-go_agent",
					Version: "latest",
				},
				CloudProperties: destiny.ResourcePoolCloudProperties{
					InstanceType:     "m3.medium",
					AvailabilityZone: "us-east-1a",
					EphemeralDisk: &destiny.ResourcePoolCloudPropertiesEphemeralDisk{
						Size: 1024,
						Type: "gp2",
					},
				},
			}))

			Expect(manifest.Jobs).To(HaveLen(2))
			Expect(manifest.Jobs[0]).To(Equal(destiny.Job{
				Name:      "consul_z1",
				Instances: 1,
				Networks: []destiny.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.0.16.9"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []destiny.JobTemplate{{
					Name:    "consul_agent",
					Release: "consul",
				}},
				Properties: &destiny.JobProperties{
					Consul: destiny.JobPropertiesConsul{
						Agent: destiny.JobPropertiesConsulAgent{
							Mode: "server",
						},
					},
				},
			}))
			Expect(manifest.Jobs[1]).To(Equal(destiny.Job{
				Name:      "etcd_z1",
				Instances: 1,
				Networks: []destiny.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.0.16.4"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []destiny.JobTemplate{
					{
						Name:    "etcd",
						Release: "etcd",
					},
					{
						Name:    "consul_agent",
						Release: "consul",
					},
				},
			}))

			Expect(manifest.Networks).To(HaveLen(1))
			Expect(manifest.Networks).To(ContainElement(destiny.Network{
				Name: "etcd1",
				Subnets: []destiny.NetworkSubnet{
					{
						CloudProperties: destiny.NetworkSubnetCloudProperties{Subnet: "subnet-1234"},
						Gateway:         "10.0.16.1",
						Range:           "10.0.16.0/24",
						Reserved: []string{
							"10.0.16.2-10.0.16.3",
							"10.0.16.13-10.0.16.254",
						},
						Static: []string{
							"10.0.16.4",
							"10.0.16.5",
							"10.0.16.6",
							"10.0.16.7",
							"10.0.16.8",
							"10.0.16.9",
						},
					},
				},
				Type: "manual",
			}))

			Expect(manifest.Properties.Etcd).To(Equal(&destiny.PropertiesEtcd{
				Machines: []string{
					"10.0.16.4",
				},
				PeerRequireSSL:                  true,
				RequireSSL:                      true,
				HeartbeatIntervalInMilliseconds: 50,
				AdvertiseURLsDNSSuffix:          "etcd.service.cf.internal",
				CACert:                          destiny.EtcdCACert,
				ClientCert:                      destiny.EtcdClientCert,
				ClientKey:                       destiny.EtcdClientKey,
				PeerCACert:                      destiny.EtcdPeerCACert,
				PeerCert:                        destiny.EtcdPeerCert,
				PeerKey:                         destiny.EtcdPeerKey,
				ServerCert:                      destiny.EtcdServerCert,
				ServerKey:                       destiny.EtcdServerKey,
			}))

			Expect(manifest.Properties.Consul).To(Equal(&destiny.PropertiesConsul{
				Agent: destiny.PropertiesConsulAgent{
					Domain: "cf.internal",
					Servers: destiny.PropertiesConsulAgentServers{
						Lan: []string{"10.0.16.9"},
					},
				},
				CACert:      destiny.ConsulCACert,
				AgentCert:   destiny.ConsulAgentCert,
				AgentKey:    destiny.ConsulAgentKey,
				ServerCert:  destiny.ConsulServerCert,
				ServerKey:   destiny.ConsulServerKey,
				EncryptKeys: []string{destiny.ConsulEncryptKey},
			}))
		})
	})

	Describe("EtcdMembers", func() {
		Context("when there is a single job with a single instance", func() {
			It("returns a list of members in the cluster", func() {
				manifest := destiny.Manifest{
					Jobs: []destiny.Job{
						{
							Instances: 1,
							Networks: []destiny.JobNetwork{{
								StaticIPs: []string{"10.244.4.2"},
							}},
						},
					},
				}

				members := manifest.EtcdMembers()
				Expect(members).To(Equal([]destiny.EtcdMember{{
					Address: "10.244.4.2",
				}}))
			})
		})

		Context("when there are multiple jobs with multiple instances", func() {
			It("returns a list of members in the cluster", func() {
				manifest := destiny.Manifest{
					Jobs: []destiny.Job{
						{
							Instances: 0,
						},
						{
							Instances: 1,
							Networks: []destiny.JobNetwork{{
								StaticIPs: []string{"10.244.4.2"},
							}},
						},
						{
							Instances: 2,
							Networks: []destiny.JobNetwork{{
								StaticIPs: []string{"10.244.5.2", "10.244.5.6"},
							}},
						},
					},
				}

				members := manifest.EtcdMembers()
				Expect(members).To(Equal([]destiny.EtcdMember{
					{
						Address: "10.244.4.2",
					},
					{
						Address: "10.244.5.2",
					},
					{
						Address: "10.244.5.6",
					},
				}))
			})
		})

		Context("when the job does not have a network", func() {
			It("returns an empty list", func() {
				manifest := destiny.Manifest{
					Jobs: []destiny.Job{
						{
							Instances: 1,
							Networks:  []destiny.JobNetwork{},
						},
					},
				}

				members := manifest.EtcdMembers()
				Expect(members).To(BeEmpty())
			})
		})

		Context("when the job network does not have enough static IPs", func() {
			It("returns as much about the list as possible", func() {
				manifest := destiny.Manifest{
					Jobs: []destiny.Job{
						{
							Instances: 2,
							Networks: []destiny.JobNetwork{{
								StaticIPs: []string{"10.244.5.2"},
							}},
						},
					},
				}

				members := manifest.EtcdMembers()
				Expect(members).To(Equal([]destiny.EtcdMember{
					{
						Address: "10.244.5.2",
					},
				}))
			})
		})
	})
})
