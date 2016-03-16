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
				Subnet: "random",
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
})
