package core_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/destiny/core"
)

var _ = Describe("Network", func() {
	var network core.Network

	BeforeEach(func() {
		network = core.Network{
			Subnets: []core.NetworkSubnet{
				{Static: []string{"10.0.0.1", "10.0.0.2"}},
				{Static: []string{"10.0.0.3", "10.0.0.4"}},
			},
		}
	})

	Describe("StaticIPs", func() {
		It("returns the requested number of ips", func() {
			ips := network.StaticIPs(3)

			Expect(ips).To(HaveLen(3))
			Expect(ips).To(ConsistOf([]string{"10.0.0.1", "10.0.0.2", "10.0.0.3"}))
		})

		It("returns an empty slice when there are fewer ips available than requested", func() {
			ips := network.StaticIPs(5)
			Expect(ips).To(HaveLen(0))
		})
	})
})
