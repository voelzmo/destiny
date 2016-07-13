package cloudconfig_test

import (
	"io/ioutil"

	"github.com/pivotal-cf-experimental/destiny/cloudconfig"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/gomegamatchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudConfig", func() {
	Describe("NewCloudConfig", func() {
		It("generates a valid cloud config for bosh lite", func() {
			cc := cloudconfig.NewCloudConfig(cloudconfig.Config{
				Networks: []cloudconfig.ConfigNetwork{
					{
						IPRange: "10.244.4.0/24",
						Nodes:   1,
					},
				},
			}, iaas.NewWardenConfig())

			Expect(cc.AZs).To(Equal([]cloudconfig.AZ{
				{
					Name: "z1",
				},
			}))

			Expect(cc.VMTypes).To(Equal([]cloudconfig.VMType{
				{
					Name: "default",
				},
			}))

			Expect(cc.DiskTypes).To(Equal([]cloudconfig.DiskType{
				{
					Name:     "default",
					DiskSize: 1024,
				},
			}))

			Expect(cc.Compilation).To(Equal(cloudconfig.Compilation{
				Workers:             5,
				ReuseCompilationVMs: true,
				AZ:                  "z1",
				VMType:              "default",
				Network:             "default",
			}))

			Expect(cc.Networks).To(Equal([]cloudconfig.Network{
				{
					Name: "default",
					Subnets: []cloudconfig.Subnet{
						{
							CloudProperties: cloudconfig.SubnetCloudProperties{
								Name: "random",
							},
							Range:   "10.244.4.0/24",
							Gateway: "10.244.4.1",
							AZ:      "z1",
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
			}))
		})

		It("generates a valid cloud config for aws", func() {
			cc := cloudconfig.NewCloudConfig(cloudconfig.Config{
				Networks: []cloudconfig.ConfigNetwork{
					{
						IPRange: "10.0.4.0/24",
						Nodes:   1,
					},
				},
			}, iaas.AWSConfig{})

			Expect(cc.AZs).To(Equal([]cloudconfig.AZ{
				{
					Name: "z1",
				},
			}))

			Expect(cc.VMTypes).To(Equal([]cloudconfig.VMType{
				{
					Name: "default",
				},
			}))

			Expect(cc.DiskTypes).To(Equal([]cloudconfig.DiskType{
				{
					Name:     "default",
					DiskSize: 1024,
				},
			}))

			Expect(cc.Compilation).To(Equal(cloudconfig.Compilation{
				Workers:             5,
				ReuseCompilationVMs: true,
				AZ:                  "z1",
				VMType:              "default",
				Network:             "default",
			}))

			Expect(cc.Networks).To(Equal([]cloudconfig.Network{
				{
					Name: "default",
					Subnets: []cloudconfig.Subnet{
						{
							CloudProperties: cloudconfig.SubnetCloudProperties{
								Name: "random",
							},
							Range:   "10.0.4.0/24",
							Gateway: "10.0.4.1",
							AZ:      "z1",
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
			}))
		})
	})

	Describe("ToYAML", func() {
		It("returns a YAML representation of the cloud config", func() {
			fixture, err := ioutil.ReadFile("fixtures/cloud_config.yml")
			Expect(err).NotTo(HaveOccurred())

			cc := cloudconfig.NewCloudConfig(cloudconfig.Config{
				Networks: []cloudconfig.ConfigNetwork{
					{
						IPRange: "10.244.4.0/24",
						Nodes:   1,
					},
				},
			}, iaas.NewWardenConfig())

			yaml, err := cc.ToYAML()
			Expect(err).NotTo(HaveOccurred())
			Expect(yaml).To(gomegamatchers.MatchYAML(fixture))
		})
	})
})
