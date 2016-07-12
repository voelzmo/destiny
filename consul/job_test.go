package consul_test

import (
	"github.com/pivotal-cf-experimental/destiny/consul"
	"github.com/pivotal-cf-experimental/destiny/iaas"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Job", func() {
	Describe("SetJobInstanceCount", func() {
		It("sets the correct values for instances and static_ips given a count", func() {
			manifest := consul.NewManifest(consul.Config{
				Networks: []consul.ConfigNetwork{
					{
						IPRange: "10.244.4.0/24",
						Nodes:   1,
					},
				},
			}, iaas.NewWardenConfig())
			job := manifest.Jobs[0]
			network := manifest.Networks[0]
			properties := manifest.Properties

			Expect(job.Instances).To(Equal(1))
			Expect(job.Networks[0].StaticIPs).To(HaveLen(1))
			Expect(job.Networks[0].Name).To(Equal(network.Name))
			Expect(job.Networks[0].StaticIPs).To(Equal([]string{"10.244.4.4"}))
			Expect(properties.Consul.Agent.Servers.Lan).To(Equal([]string{"10.244.4.4"}))

			job, properties = consul.SetJobInstanceCount(job, network, properties, 3)
			Expect(job.Instances).To(Equal(3))
			Expect(job.Networks[0].StaticIPs).To(HaveLen(3))
			Expect(job.Networks[0].StaticIPs).To(Equal([]string{"10.244.4.4", "10.244.4.5", "10.244.4.6"}))
			Expect(properties.Consul.Agent.Servers.Lan).To(Equal([]string{"10.244.4.4", "10.244.4.5", "10.244.4.6"}))
		})
	})
})
