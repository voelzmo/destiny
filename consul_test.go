package destiny_test

import (
	"github.com/pivotal-cf-experimental/destiny"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Consul Manifest", func() {
	Describe("NewConsul", func() {
		It("generates a valid Consul BOSH-Lite manifest", func() {
			manifest := destiny.NewConsul(destiny.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "consul-some-random-guid",
				IPRange:      "10.244.4.0/24",
				IAAS:         destiny.Warden,
			})

			Expect(manifest).To(Equal(destiny.Manifest{
				DirectorUUID: "some-director-uuid",
				Name:         "consul-some-random-guid",
				Releases: []destiny.Release{{
					Name:    "consul",
					Version: "latest",
				}},
				Compilation: destiny.Compilation{
					Network:             "consul1",
					ReuseCompilationVMs: true,
					Workers:             3,
				},
				Update: destiny.Update{
					Canaries:        1,
					CanaryWatchTime: "1000-180000",
					MaxInFlight:     50,
					Serial:          true,
					UpdateWatchTime: "1000-180000",
				},
				ResourcePools: []destiny.ResourcePool{
					{
						Name:    "consul_z1",
						Network: "consul1",
						Stemcell: destiny.ResourcePoolStemcell{
							Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
							Version: "latest",
						},
					},
				},
				Jobs: []destiny.Job{
					{
						Name:      "consul_z1",
						Instances: 1,
						Networks: []destiny.JobNetwork{{
							Name:      "consul1",
							StaticIPs: []string{"10.244.4.4"},
						}},
						PersistentDisk: 1024,
						Properties: &destiny.JobProperties{
							Consul: destiny.JobPropertiesConsul{
								Agent: destiny.JobPropertiesConsulAgent{
									Mode:     "server",
									LogLevel: "info",
									Services: destiny.JobPropertiesConsulAgentServices{
										"router": destiny.JobPropertiesConsulAgentService{
											Name: "gorouter",
											Check: &destiny.JobPropertiesConsulAgentServiceCheck{
												Name:     "router-check",
												Script:   "/var/vcap/jobs/router/bin/script",
												Interval: "1m",
											},
											Tags: []string{"routing"},
										},
										"cloud_controller": destiny.JobPropertiesConsulAgentService{},
									},
								},
							},
						},
						ResourcePool: "consul_z1",
						Templates: []destiny.JobTemplate{{
							Name:    "consul_agent",
							Release: "consul",
						}},
						Update: &destiny.JobUpdate{
							MaxInFlight: 1,
						},
					},
					{
						Name:      "consul_test_consumer",
						Instances: 1,
						Networks: []destiny.JobNetwork{{
							Name:      "consul1",
							StaticIPs: []string{"10.244.4.9"},
						}},
						PersistentDisk: 1024,
						ResourcePool:   "consul_z1",
						Templates: []destiny.JobTemplate{
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
				Networks: []destiny.Network{
					{
						Name: "consul1",
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
					},
				},
				Properties: destiny.Properties{
					Consul: &destiny.PropertiesConsul{
						Agent: destiny.PropertiesConsulAgent{
							Domain:   "cf.internal",
							LogLevel: "",
							Servers: destiny.PropertiesConsulAgentServers{
								Lan: []string{"10.244.4.4"},
							},
						},
						CACert:      destiny.CACert,
						AgentCert:   destiny.AgentCert,
						AgentKey:    destiny.AgentKey,
						ServerCert:  destiny.ServerCert,
						ServerKey:   destiny.ServerKey,
						EncryptKeys: []string{destiny.EncryptKey},
					},
				},
			}))
		})

		It("generates a valid Consul AWS manifest", func() {
			manifest := destiny.NewConsul(destiny.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "consul-some-random-guid",
				IPRange:      "10.0.4.0/24",
				IAAS:         destiny.AWS,
				AWS: destiny.ConfigAWS{
					Subnet: "subnet-1234",
				},
			})

			Expect(manifest).To(Equal(destiny.Manifest{
				DirectorUUID: "some-director-uuid",
				Name:         "consul-some-random-guid",
				Releases: []destiny.Release{{
					Name:    "consul",
					Version: "latest",
				}},
				Compilation: destiny.Compilation{
					Network:             "consul1",
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
				},
				Update: destiny.Update{
					Canaries:        1,
					CanaryWatchTime: "1000-180000",
					MaxInFlight:     50,
					Serial:          true,
					UpdateWatchTime: "1000-180000",
				},
				ResourcePools: []destiny.ResourcePool{
					{
						Name:    "consul_z1",
						Network: "consul1",
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
					},
				},
				Jobs: []destiny.Job{
					{
						Name:      "consul_z1",
						Instances: 1,
						Networks: []destiny.JobNetwork{{
							Name:      "consul1",
							StaticIPs: []string{"10.0.4.4"},
						}},
						PersistentDisk: 1024,
						Properties: &destiny.JobProperties{
							Consul: destiny.JobPropertiesConsul{
								Agent: destiny.JobPropertiesConsulAgent{
									Mode:     "server",
									LogLevel: "info",
									Services: destiny.JobPropertiesConsulAgentServices{
										"router": destiny.JobPropertiesConsulAgentService{
											Name: "gorouter",
											Check: &destiny.JobPropertiesConsulAgentServiceCheck{
												Name:     "router-check",
												Script:   "/var/vcap/jobs/router/bin/script",
												Interval: "1m",
											},
											Tags: []string{"routing"},
										},
										"cloud_controller": destiny.JobPropertiesConsulAgentService{},
									},
								},
							},
						},
						ResourcePool: "consul_z1",
						Templates: []destiny.JobTemplate{{
							Name:    "consul_agent",
							Release: "consul",
						}},
						Update: &destiny.JobUpdate{
							MaxInFlight: 1,
						},
					},
					{
						Name:      "consul_test_consumer",
						Instances: 1,
						Networks: []destiny.JobNetwork{{
							Name:      "consul1",
							StaticIPs: []string{"10.0.4.9"},
						}},
						PersistentDisk: 1024,
						ResourcePool:   "consul_z1",
						Templates: []destiny.JobTemplate{
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
				Networks: []destiny.Network{
					{
						Name: "consul1",
						Subnets: []destiny.NetworkSubnet{
							{
								CloudProperties: destiny.NetworkSubnetCloudProperties{Subnet: "subnet-1234"},
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
				Properties: destiny.Properties{
					Consul: &destiny.PropertiesConsul{
						Agent: destiny.PropertiesConsulAgent{
							Domain:   "cf.internal",
							LogLevel: "",
							Servers: destiny.PropertiesConsulAgentServers{
								Lan: []string{"10.0.4.4"},
							},
						},
						CACert:      destiny.CACert,
						AgentCert:   destiny.AgentCert,
						AgentKey:    destiny.AgentKey,
						ServerCert:  destiny.ServerCert,
						ServerKey:   destiny.ServerKey,
						EncryptKeys: []string{destiny.EncryptKey},
					},
				},
			}))
		})
	})

	Describe("ConsulMembers", func() {
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

				members := manifest.ConsulMembers()
				Expect(members).To(Equal([]destiny.ConsulMember{{
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

				members := manifest.ConsulMembers()
				Expect(members).To(Equal([]destiny.ConsulMember{
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

				members := manifest.ConsulMembers()
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

				members := manifest.ConsulMembers()
				Expect(members).To(Equal([]destiny.ConsulMember{
					{
						Address: "10.244.5.2",
					},
				}))
			})
		})
	})
})
