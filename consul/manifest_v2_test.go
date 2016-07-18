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
			manifest := consul.NewManifestV2(consul.ConfigV2{
				DirectorUUID: "some-director-uuid",
				Name:         "consul-some-random-guid",
				AZs: []consul.ConfigAZ{
					{
						Name:    "z1",
						IPRange: "10.244.4.0/24",
						Nodes:   2,
					},
					{
						Name:    "z2",
						IPRange: "10.244.5.0/24",
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

			Expect(manifest.Stemcells).To(Equal([]consul.Stemcell{
				{
					Alias:   "default",
					Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
					Version: "latest",
				},
			}))

			Expect(manifest.Update).To(Equal(core.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     1,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))

			Expect(manifest.InstanceGroups).To(HaveLen(2))
			Expect(manifest.InstanceGroups[0]).To(Equal(core.InstanceGroup{
				Instances: 3,
				Name:      "consul",
				AZs:       []string{"z1", "z2"},
				Networks: []core.InstanceGroupNetwork{
					{
						Name: "private",
						StaticIPs: []string{
							"10.244.4.4",
							"10.244.4.5",
							"10.244.5.4",
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
				},
				MigratedFrom: []core.InstanceGroupMigratedFrom{
					{
						Name: "consul_z1",
						AZ:   "z1",
					},
					{
						Name: "consul_z2",
						AZ:   "z2",
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
			}))

			Expect(manifest.InstanceGroups[1]).To(Equal(core.InstanceGroup{
				Instances: 3,
				Name:      "consul_test_consumer",
				AZs:       []string{"z1"},
				Networks: []core.InstanceGroupNetwork{
					{
						Name: "private",
						StaticIPs: []string{
							"10.244.4.10",
							"10.244.4.11",
							"10.244.4.12",
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
				MigratedFrom: []core.InstanceGroupMigratedFrom{
					{
						Name: "consul_test_consumer",
						AZ:   "z1",
					},
				},
			}))

			Expect(manifest.Properties).To(Equal(consul.Properties{
				Consul: &consul.PropertiesConsul{
					Agent: consul.PropertiesConsulAgent{
						Domain:     "cf.internal",
						Datacenter: "dc1",
						Servers: consul.PropertiesConsulAgentServers{
							Lan: []string{
								"10.244.4.4",
								"10.244.4.5",
								"10.244.5.4",
							},
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
			manifest := consul.NewManifestV2(consul.ConfigV2{
				DirectorUUID: "some-director-uuid",
				Name:         "consul-some-random-guid",
				AZs: []consul.ConfigAZ{
					{
						Name:    "us-east-1",
						IPRange: "10.0.4.0/24",
						Nodes:   1,
					},
				},
			}, iaas.AWSConfig{})

			Expect(manifest.Stemcells).To(Equal([]consul.Stemcell{
				{
					Alias:   "default",
					Name:    "bosh-aws-xen-hvm-ubuntu-trusty-go_agent",
					Version: "latest",
				},
			}))

		})
	})

	Describe("ToYAML", func() {
		It("returns a YAML representation of the consul manifest", func() {
			consulManifest, err := ioutil.ReadFile("fixtures/consul_manifest_v2.yml")
			Expect(err).NotTo(HaveOccurred())

			manifest := consul.NewManifestV2(consul.ConfigV2{
				DirectorUUID: "some-director-uuid",
				Name:         "consul-some-random-guid",
				AZs: []consul.ConfigAZ{
					{
						Name:    "z1",
						IPRange: "10.244.4.0/24",
						Nodes:   2,
					},
					{
						Name:    "z2",
						IPRange: "10.244.5.0/24",
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
