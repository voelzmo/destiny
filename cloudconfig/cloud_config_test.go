package cloudconfig_test

import (
	"io/ioutil"

	"github.com/pivotal-cf-experimental/destiny/cloudconfig"
	"github.com/pivotal-cf-experimental/gomegamatchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudConfig", func() {
	Describe("NewWardenCloudConfig", func() {
		It("generates a valid cloud config", func() {
			cc := cloudconfig.NewWardenCloudConfig(cloudconfig.Config{
				AZs: []cloudconfig.ConfigAZ{
					{
						IPRange:   "10.244.4.0/24",
						StaticIPs: 6,
					},
					{
						IPRange:   "10.244.5.0/24",
						StaticIPs: 3,
					},
					{
						IPRange:   "10.244.6.0/24",
						StaticIPs: 2,
					},
				},
			})

			Expect(cc.AZs).To(HaveLen(3))
			Expect(cc.AZs[0]).To(Equal(cloudconfig.AZ{
				Name: "z1",
			}))
			Expect(cc.AZs[1]).To(Equal(cloudconfig.AZ{
				Name: "z2",
			}))
			Expect(cc.AZs[2]).To(Equal(cloudconfig.AZ{
				Name: "z3",
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
				Workers:             3,
				ReuseCompilationVMs: true,
				AZ:                  "z1",
				Network:             "private",
				VMType:              "default",
			}))

			Expect(cc.Networks).To(HaveLen(1))
			Expect(cc.Networks[0].Name).To(Equal("private"))
			Expect(cc.Networks[0].Type).To(Equal("manual"))

			Expect(cc.Networks[0].Subnets).To(HaveLen(3))
			Expect(cc.Networks[0].Subnets[0]).To(Equal(cloudconfig.Subnet{
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
				CloudProperties: cloudconfig.SubnetCloudProperties{
					Name: "random",
				},
			}))
			Expect(cc.Networks[0].Subnets[1]).To(Equal(cloudconfig.Subnet{
				Range:   "10.244.5.0/24",
				Gateway: "10.244.5.1",
				AZ:      "z2",
				Reserved: []string{
					"10.244.5.2-10.244.5.3",
					"10.244.5.7-10.244.5.254",
				},
				Static: []string{
					"10.244.5.4",
					"10.244.5.5",
					"10.244.5.6",
				},
				CloudProperties: cloudconfig.SubnetCloudProperties{
					Name: "random",
				},
			}))
			Expect(cc.Networks[0].Subnets[2]).To(Equal(cloudconfig.Subnet{
				Range:   "10.244.6.0/24",
				Gateway: "10.244.6.1",
				AZ:      "z3",
				Reserved: []string{
					"10.244.6.2-10.244.6.3",
					"10.244.6.6-10.244.6.254",
				},
				Static: []string{
					"10.244.6.4",
					"10.244.6.5",
				},
				CloudProperties: cloudconfig.SubnetCloudProperties{
					Name: "random",
				},
			}))
		})
	})

	Describe("ToYAML", func() {
		It("returns a YAML representation of the cloud config", func() {
			fixture, err := ioutil.ReadFile("fixtures/cloud_config.yml")
			Expect(err).NotTo(HaveOccurred())

			cc := cloudconfig.NewWardenCloudConfig(cloudconfig.Config{
				AZs: []cloudconfig.ConfigAZ{
					{
						IPRange:   "10.244.4.0/24",
						StaticIPs: 11,
					},
					{
						IPRange:   "10.244.5.0/24",
						StaticIPs: 5,
					},
				},
			})

			yaml, err := cc.ToYAML()
			Expect(err).NotTo(HaveOccurred())
			Expect(yaml).To(gomegamatchers.MatchYAML(fixture))
		})
	})
})
