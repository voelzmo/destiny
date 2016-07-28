package consul_test

import (
	"github.com/pivotal-cf-experimental/destiny/consul"
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Job", func() {
	Describe("SetJobInstanceCount", func() {
		It("sets the correct values for instances and static_ips given a count", func() {
			manifest, err := consul.NewManifest(consul.Config{
				Networks: []consul.ConfigNetwork{
					{
						IPRange: "10.244.4.0/24",
						Nodes:   1,
					},
				},
			}, iaas.NewWardenConfig())
			Expect(err).NotTo(HaveOccurred())

			job := manifest.Jobs[0]
			network := manifest.Networks[0]
			properties := manifest.Properties

			Expect(job.Instances).To(Equal(1))
			Expect(job.Networks[0].StaticIPs).To(HaveLen(1))
			Expect(job.Networks[0].Name).To(Equal(network.Name))
			Expect(job.Networks[0].StaticIPs).To(Equal([]string{"10.244.4.4"}))
			Expect(properties.Consul.Agent.Servers.Lan).To(Equal([]string{"10.244.4.4"}))

			job, properties, err = consul.SetJobInstanceCount(job, network, properties, 3)
			Expect(err).NotTo(HaveOccurred())
			Expect(job.Instances).To(Equal(3))
			Expect(job.Networks[0].StaticIPs).To(HaveLen(3))
			Expect(job.Networks[0].StaticIPs).To(Equal([]string{"10.244.4.4", "10.244.4.5", "10.244.4.6"}))
			Expect(properties.Consul.Agent.Servers.Lan).To(Equal([]string{"10.244.4.4", "10.244.4.5", "10.244.4.6"}))
		})

		Context("failure cases", func() {
			It("returns an error when set job instance count fails to retrieve list of static ips", func() {
				manifest, err := consul.NewManifest(consul.Config{
					Networks: []consul.ConfigNetwork{
						{
							IPRange: "10.244.4.0/24",
							Nodes:   1,
						},
					},
				}, iaas.NewWardenConfig())
				Expect(err).NotTo(HaveOccurred())

				job := manifest.Jobs[0]
				network := manifest.Networks[0]
				properties := manifest.Properties

				network.Subnets = []core.NetworkSubnet{
					{
						Static: []string{"fake-ip"},
					},
				}

				_, _, err = consul.SetJobInstanceCount(job, network, properties, 100)
				Expect(err).To(MatchError("'fake' is not a valid ip address"))
			})
		})
	})
})
