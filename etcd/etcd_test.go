package etcd_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/destiny/consul"
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/etcd"
	"github.com/pivotal-cf-experimental/destiny/iaas"
)

var _ = Describe("Etcd Manifest", func() {
	Describe("NewEtcd", func() {
		It("generates a valid Etcd BOSH-Lite manifest", func() {
			manifest := etcd.NewEtcd(etcd.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "etcd-some-random-guid",
				IPRange:      "10.244.4.0/24",
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
			Expect(manifest.Jobs).To(HaveLen(2))
			Expect(manifest.Jobs[0]).To(Equal(core.Job{
				Name:      "consul_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.9"},
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
					StaticIPs: []string{"10.244.4.4"},
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
			Expect(manifest.Properties.Etcd).To(Equal(&etcd.PropertiesEtcd{
				Machines: []string{
					"10.244.4.4",
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
						Lan: []string{"10.244.4.9"},
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
			manifest := etcd.NewEtcd(etcd.Config{
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

			Expect(manifest.Jobs).To(HaveLen(2))
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

			Expect(manifest.Properties.Etcd).To(Equal(&etcd.PropertiesEtcd{
				Machines: []string{
					"10.0.16.4",
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
			It("returns a list of members in the cluster", func() {
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
						},
						{
							Instances: 2,
							Networks: []core.JobNetwork{{
								StaticIPs: []string{"10.244.5.2", "10.244.5.6"},
							}},
						},
					},
				}

				members := manifest.EtcdMembers()
				Expect(members).To(Equal([]etcd.EtcdMember{
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
})
