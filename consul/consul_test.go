package consul_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/destiny/consul"
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
)

var _ = Describe("Consul Manifest", func() {
	Describe("NewConsul", func() {
		It("generates a valid Consul BOSH-Lite manifest", func() {
			manifest := consul.NewConsul(consul.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "consul-some-random-guid",
				IPRange:      "10.244.4.0/24",
			}, iaas.NewWardenConfig())

			Expect(manifest).To(Equal(consul.Manifest{
				DirectorUUID: "some-director-uuid",
				Name:         "consul-some-random-guid",
				Releases: []core.Release{{
					Name:    "consul",
					Version: "latest",
				}},
				Compilation: core.Compilation{
					Network:             "consul1",
					ReuseCompilationVMs: true,
					Workers:             3,
				},
				Update: core.Update{
					Canaries:        1,
					CanaryWatchTime: "1000-180000",
					MaxInFlight:     50,
					Serial:          true,
					UpdateWatchTime: "1000-180000",
				},
				ResourcePools: []core.ResourcePool{
					{
						Name:    "consul_z1",
						Network: "consul1",
						Stemcell: core.ResourcePoolStemcell{
							Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
							Version: "latest",
						},
					},
				},
				Jobs: []core.Job{
					{
						Name:      "consul_z1",
						Instances: 1,
						Networks: []core.JobNetwork{{
							Name:      "consul1",
							StaticIPs: []string{"10.244.4.4"},
						}},
						PersistentDisk: 1024,
						Properties: &core.JobProperties{
							Consul: core.JobPropertiesConsul{
								Agent: core.JobPropertiesConsulAgent{
									Mode:     "server",
									LogLevel: "info",
									Services: core.JobPropertiesConsulAgentServices{
										"router": core.JobPropertiesConsulAgentService{
											Name: "gorouter",
											Check: &core.JobPropertiesConsulAgentServiceCheck{
												Name:     "router-check",
												Script:   "/var/vcap/jobs/router/bin/script",
												Interval: "1m",
											},
											Tags: []string{"routing"},
										},
										"cloud_controller": core.JobPropertiesConsulAgentService{},
									},
								},
							},
						},
						ResourcePool: "consul_z1",
						Templates: []core.JobTemplate{{
							Name:    "consul_agent",
							Release: "consul",
						}},
						Update: &core.JobUpdate{
							MaxInFlight: 1,
						},
					},
					{
						Name:      "consul_test_consumer",
						Instances: 1,
						Networks: []core.JobNetwork{{
							Name:      "consul1",
							StaticIPs: []string{"10.244.4.9"},
						}},
						PersistentDisk: 1024,
						ResourcePool:   "consul_z1",
						Templates: []core.JobTemplate{
							{
								Name:    "consul_agent",
								Release: "consul",
							},
							{
								Name:    "consul-test-consumer",
								Release: "consul",
							},
						},
					},
				},
				Networks: []core.Network{
					{
						Name: "consul1",
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
					},
				},
				Properties: consul.Properties{
					Consul: &consul.PropertiesConsul{
						Agent: consul.PropertiesConsulAgent{
							Domain:   "cf.internal",
							LogLevel: "",
							Servers: consul.PropertiesConsulAgentServers{
								Lan: []string{"10.244.4.4"},
							},
						},
						CACert:      consul.CACert,
						AgentCert:   consul.AgentCert,
						AgentKey:    consul.AgentKey,
						ServerCert:  consul.ServerCert,
						ServerKey:   consul.ServerKey,
						EncryptKeys: []string{consul.EncryptKey},
					},
				},
			}))
		})

		It("generates a valid Consul AWS manifest", func() {
			manifest := consul.NewConsul(consul.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "consul-some-random-guid",
				IPRange:      "10.0.4.0/24",
			}, iaas.AWSConfig{
				Subnet: "subnet-1234",
			})

			Expect(manifest).To(Equal(consul.Manifest{
				DirectorUUID: "some-director-uuid",
				Name:         "consul-some-random-guid",
				Releases: []core.Release{{
					Name:    "consul",
					Version: "latest",
				}},
				Compilation: core.Compilation{
					Network:             "consul1",
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
				},
				Update: core.Update{
					Canaries:        1,
					CanaryWatchTime: "1000-180000",
					MaxInFlight:     50,
					Serial:          true,
					UpdateWatchTime: "1000-180000",
				},
				ResourcePools: []core.ResourcePool{
					{
						Name:    "consul_z1",
						Network: "consul1",
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
					},
				},
				Jobs: []core.Job{
					{
						Name:      "consul_z1",
						Instances: 1,
						Networks: []core.JobNetwork{{
							Name:      "consul1",
							StaticIPs: []string{"10.0.4.4"},
						}},
						PersistentDisk: 1024,
						Properties: &core.JobProperties{
							Consul: core.JobPropertiesConsul{
								Agent: core.JobPropertiesConsulAgent{
									Mode:     "server",
									LogLevel: "info",
									Services: core.JobPropertiesConsulAgentServices{
										"router": core.JobPropertiesConsulAgentService{
											Name: "gorouter",
											Check: &core.JobPropertiesConsulAgentServiceCheck{
												Name:     "router-check",
												Script:   "/var/vcap/jobs/router/bin/script",
												Interval: "1m",
											},
											Tags: []string{"routing"},
										},
										"cloud_controller": core.JobPropertiesConsulAgentService{},
									},
								},
							},
						},
						ResourcePool: "consul_z1",
						Templates: []core.JobTemplate{{
							Name:    "consul_agent",
							Release: "consul",
						}},
						Update: &core.JobUpdate{
							MaxInFlight: 1,
						},
					},
					{
						Name:      "consul_test_consumer",
						Instances: 1,
						Networks: []core.JobNetwork{{
							Name:      "consul1",
							StaticIPs: []string{"10.0.4.9"},
						}},
						PersistentDisk: 1024,
						ResourcePool:   "consul_z1",
						Templates: []core.JobTemplate{
							{
								Name:    "consul_agent",
								Release: "consul",
							},
							{
								Name:    "consul-test-consumer",
								Release: "consul",
							},
						},
					},
				},
				Networks: []core.Network{
					{
						Name: "consul1",
						Subnets: []core.NetworkSubnet{
							{
								CloudProperties: core.NetworkSubnetCloudProperties{Subnet: "subnet-1234"},
								Gateway:         "10.0.4.1",
								Range:           "10.0.4.0/24",
								Reserved: []string{
									"10.0.4.2-10.0.4.3",
									"10.0.4.13-10.0.4.254",
								},
								Static: []string{
									"10.0.4.4",
									"10.0.4.5",
									"10.0.4.6",
									"10.0.4.7",
									"10.0.4.8",
									"10.0.4.9",
								},
							},
						},
						Type: "manual",
					},
				},
				Properties: consul.Properties{
					Consul: &consul.PropertiesConsul{
						Agent: consul.PropertiesConsulAgent{
							Domain:   "cf.internal",
							LogLevel: "",
							Servers: consul.PropertiesConsulAgentServers{
								Lan: []string{"10.0.4.4"},
							},
						},
						CACert:      consul.CACert,
						AgentCert:   consul.AgentCert,
						AgentKey:    consul.AgentKey,
						ServerCert:  consul.ServerCert,
						ServerKey:   consul.ServerKey,
						EncryptKeys: []string{consul.EncryptKey},
					},
				},
			}))
		})
	})

	Describe("ConsulMembers", func() {
		Context("when there is a single job with a single instance", func() {
			It("returns a list of members in the cluster", func() {
				manifest := consul.Manifest{
					Jobs: []core.Job{
						{
							Instances: 1,
							Networks: []core.JobNetwork{{
								StaticIPs: []string{"10.244.4.2"},
							}},
						},
					},
				}

				members := manifest.ConsulMembers()
				Expect(members).To(Equal([]consul.ConsulMember{{
					Address: "10.244.4.2",
				}}))
			})
		})

		Context("when there are multiple jobs with multiple instances", func() {
			It("returns a list of members in the cluster", func() {
				manifest := consul.Manifest{
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

				members := manifest.ConsulMembers()
				Expect(members).To(Equal([]consul.ConsulMember{
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
				manifest := consul.Manifest{
					Jobs: []core.Job{
						{
							Instances: 1,
							Networks:  []core.JobNetwork{},
						},
					},
				}

				members := manifest.ConsulMembers()
				Expect(members).To(BeEmpty())
			})
		})

		Context("when the job network does not have enough static IPs", func() {
			It("returns as much about the list as possible", func() {
				manifest := consul.Manifest{
					Jobs: []core.Job{
						{
							Instances: 2,
							Networks: []core.JobNetwork{{
								StaticIPs: []string{"10.244.5.2"},
							}},
						},
					},
				}

				members := manifest.ConsulMembers()
				Expect(members).To(Equal([]consul.ConsulMember{
					{
						Address: "10.244.5.2",
					},
				}))
			})
		})
	})
})
