package destiny_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/destiny"
)

var _ = Describe("Warden Config", func() {

	var (
		wardenConfig destiny.WardenConfig
	)

	BeforeEach(func() {
		wardenConfig = destiny.NewWardenConfig()
	})

	Describe("NetworkSubnet", func() {
		It("returns a network subnet specific to Warden", func() {
			subnetCloudProperties := wardenConfig.NetworkSubnet()
			Expect(subnetCloudProperties).To(Equal(destiny.NetworkSubnetCloudProperties{
				Name: "random",
			}))
		})
	})

	Describe("Compilation", func() {
		It("returns an empty compilation cloud properties for Warden", func() {
			compilationCloudProperties := wardenConfig.Compilation()
			Expect(compilationCloudProperties).To(Equal(destiny.CompilationCloudProperties{}))
		})
	})

	Describe("ResourcePool", func() {
		It("returns an empty resource pool cloud properties for Warden", func() {
			resourcePoolCloudProperties := wardenConfig.ResourcePool()
			Expect(resourcePoolCloudProperties).To(Equal(destiny.ResourcePoolCloudProperties{}))
		})
	})

	Describe("CPI", func() {
		It("returns the cpi specific to Warden", func() {
			cpi := wardenConfig.CPI()
			Expect(cpi).To(Equal(destiny.CPI{
				JobName:     "warden_cpi",
				ReleaseName: "bosh-warden-cpi",
			}))
		})
	})

	Describe("Properties", func() {
		It("returns the properties specific to Warden", func() {
			properties := wardenConfig.Properties()
			Expect(properties).To(Equal(destiny.Properties{
				WardenCPI: &destiny.PropertiesWardenCPI{
					Agent: destiny.PropertiesWardenCPIAgent{
						Blobstore: destiny.PropertiesWardenCPIAgentBlobstore{
							Options: destiny.PropertiesWardenCPIAgentBlobstoreOptions{
								Endpoint: "http://10.254.50.4:25251",
								Password: "agent-password",
								User:     "agent",
							},
							Provider: "dav",
						},
						Mbus: "nats://nats:nats-password@10.254.50.4:4222",
					},
					Warden: destiny.PropertiesWardenCPIWarden{
						ConnectAddress: "10.254.50.4:7777",
						ConnectNetwork: "tcp",
					},
				},
			}))
		})
	})
})
