package consul_test

import (
	"io/ioutil"

	"github.com/pivotal-cf-experimental/destiny/consul"
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/gomegamatchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ManifestV2", func() {
	Describe("NewManifestV2", func() {
		It("returns a BOSH 2.0 manifest for bosh-lite", func() {
			manifest := consul.NewManifestV2(consul.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "consul-some-random-guid",
				Networks: []consul.ConfigNetwork{
					{
						IPRange: "10.244.4.0/24",
						Nodes:   1,
					},
				},
			}, iaas.NewWardenConfig())

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))
			Expect(manifest.Name).To(Equal("consul-some-random-guid"))

			Expect(manifest.Releases).To(Equal([]core.Release{
				{
					Name:    "consul",
					Version: "latest",
				},
			}))

			Expect(manifest.Stemcells).To(Equal([]core.Stemcell{
				{
					Alias:   "default",
					OS:      "ubuntu-trusty",
					Version: "latest",
				},
			}))

			Expect(manifest.Update).To(Equal(core.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     50,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))

			Expect(manifest.InstanceGroups).To(Equal([]core.InstanceGroup{
				core.InstanceGroup{
					Instances: 1,
					Name:      "consul",
					AZs:       []string{"z1"},
					Networks: []core.InstanceGroupNetwork{
						{
							Name: "private",
							StaticIPs: []string{
								"10.244.4.4",
							},
						},
					},
					VMType:             "default",
					Stemcell:           "default",
					PersistentDiskType: "default",
					Update: core.Update{
						MaxInFlight: 1,
					},
					Jobs: []core.InstanceGroupJob{
						{
							Name:    "consul_agent",
							Release: "consul",
						},
					},
					Properties: core.InstanceGroupProperties{
						Consul: core.InstanceGroupPropertiesConsul{
							Agent: core.InstanceGroupPropertiesConsulAgent{
								Mode:     "server",
								LogLevel: "info",
								Services: map[string]core.InstanceGroupPropertiesConsulAgentService{
									"router": core.InstanceGroupPropertiesConsulAgentService{
										Name: "gorouter",
										Check: core.InstanceGroupPropertiesConsulAgentServiceCheck{
											Name:     "router-check",
											Script:   "/var/vcap/jobs/router/bin/script",
											Interval: "1m",
										},
										Tags: []string{"routing"},
									},
									"cloud_controller": core.InstanceGroupPropertiesConsulAgentService{},
								},
							},
						},
					},
				},
				core.InstanceGroup{
					Instances: 1,
					Name:      "consul_test_consumer",
					AZs:       []string{"z1"},
					Networks: []core.InstanceGroupNetwork{
						{
							Name: "private",
							StaticIPs: []string{
								"10.244.4.9",
							},
						},
					},
					VMType:             "default",
					Stemcell:           "default",
					PersistentDiskType: "default",
					Jobs: []core.InstanceGroupJob{
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
			}))

			Expect(manifest.Properties).To(Equal(consul.Properties{
				Consul: &consul.PropertiesConsul{
					Agent: consul.PropertiesConsulAgent{
						Domain:     "cf.internal",
						Datacenter: "dc1",
						Servers: consul.PropertiesConsulAgentServers{
							Lan: []string{"10.244.4.4"},
						},
					},
					AgentCert: consul.DC1AgentCert,
					AgentKey:  consul.DC1AgentKey,
					CACert:    consul.CACert,
					EncryptKeys: []string{
						consul.EncryptKey,
					},
					ServerCert: consul.DC1ServerCert,
					ServerKey:  consul.DC1ServerKey,
				},
			}))
		})

		It("returns a bosh 2.0 consul manifest for aws", func() {
			manifest := consul.NewManifestV2(consul.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "consul-some-random-guid",
				Networks: []consul.ConfigNetwork{
					{
						IPRange: "10.0.4.0/24",
						Nodes:   1,
					},
				},
			}, iaas.AWSConfig{})

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))
			Expect(manifest.Name).To(Equal("consul-some-random-guid"))

			Expect(manifest.Releases).To(Equal([]core.Release{
				{
					Name:    "consul",
					Version: "latest",
				},
			}))

			Expect(manifest.Stemcells).To(Equal([]core.Stemcell{
				{
					Alias:   "default",
					OS:      "ubuntu-trusty",
					Version: "latest",
				},
			}))

			Expect(manifest.Update).To(Equal(core.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     50,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))

			Expect(manifest.InstanceGroups).To(Equal([]core.InstanceGroup{
				core.InstanceGroup{
					Instances: 1,
					Name:      "consul",
					AZs:       []string{"z1"},
					Networks: []core.InstanceGroupNetwork{
						{
							Name: "private",
							StaticIPs: []string{
								"10.0.4.4",
							},
						},
					},
					VMType:             "default",
					Stemcell:           "default",
					PersistentDiskType: "default",
					Update: core.Update{
						MaxInFlight: 1,
					},
					Jobs: []core.InstanceGroupJob{
						{
							Name:    "consul_agent",
							Release: "consul",
						},
					},
					Properties: core.InstanceGroupProperties{
						Consul: core.InstanceGroupPropertiesConsul{
							Agent: core.InstanceGroupPropertiesConsulAgent{
								Mode:     "server",
								LogLevel: "info",
								Services: map[string]core.InstanceGroupPropertiesConsulAgentService{
									"router": core.InstanceGroupPropertiesConsulAgentService{
										Name: "gorouter",
										Check: core.InstanceGroupPropertiesConsulAgentServiceCheck{
											Name:     "router-check",
											Script:   "/var/vcap/jobs/router/bin/script",
											Interval: "1m",
										},
										Tags: []string{"routing"},
									},
									"cloud_controller": core.InstanceGroupPropertiesConsulAgentService{},
								},
							},
						},
					},
				},
				core.InstanceGroup{
					Instances: 1,
					Name:      "consul_test_consumer",
					AZs:       []string{"z1"},
					Networks: []core.InstanceGroupNetwork{
						{
							Name: "private",
							StaticIPs: []string{
								"10.0.4.9",
							},
						},
					},
					VMType:             "default",
					Stemcell:           "default",
					PersistentDiskType: "default",
					Jobs: []core.InstanceGroupJob{
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
			}))

			Expect(manifest.Properties).To(Equal(consul.Properties{
				Consul: &consul.PropertiesConsul{
					Agent: consul.PropertiesConsulAgent{
						Domain:     "cf.internal",
						Datacenter: "dc1",
						Servers: consul.PropertiesConsulAgentServers{
							Lan: []string{"10.0.4.4"},
						},
					},
					AgentCert: consul.DC1AgentCert,
					AgentKey:  consul.DC1AgentKey,
					CACert:    consul.CACert,
					EncryptKeys: []string{
						consul.EncryptKey,
					},
					ServerCert: consul.DC1ServerCert,
					ServerKey:  consul.DC1ServerKey,
				},
			}))
		})
	})

	Describe("ToYAML", func() {
		It("returns a YAML representation of the consul manifest", func() {
			consulManifest, err := ioutil.ReadFile("fixtures/consul_manifest_v2.yml")
			Expect(err).NotTo(HaveOccurred())

			manifest := consul.NewManifestV2(consul.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "consul",
				Networks: []consul.ConfigNetwork{
					{
						IPRange: "10.244.4.0/24",
						Nodes:   1,
					},
				},
			}, iaas.NewWardenConfig())

			yaml, err := manifest.ToYAML()
			Expect(err).NotTo(HaveOccurred())
			Expect(yaml).To(gomegamatchers.MatchYAML(consulManifest))
		})
	})
})
