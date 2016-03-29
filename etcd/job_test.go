package etcd_test

import (
	"github.com/pivotal-cf-experimental/destiny/etcd"
	"github.com/pivotal-cf-experimental/destiny/iaas"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Job", func() {
	Describe("SetJobInstanceCount", func() {
		It("sets the correct values for instances and static_ips given a count", func() {
			manifest := etcd.NewManifest(etcd.Config{
				IPRange: "10.244.4.0/24",
			}, iaas.NewWardenConfig())
			job := manifest.Jobs[1]
			network := manifest.Networks[0]
			properties := manifest.Properties

			Expect(job.Instances).To(Equal(1))
			Expect(job.Networks[0].StaticIPs).To(HaveLen(1))
			Expect(job.Networks[0].Name).To(Equal(network.Name))
			Expect(properties.Etcd.Machines).To(Equal(job.Networks[0].StaticIPs))

			job, properties = etcd.SetJobInstanceCount(job, network, properties, 3)
			Expect(job.Instances).To(Equal(3))
			Expect(job.Networks[0].StaticIPs).To(HaveLen(3))
			Expect(properties.Etcd.Machines).To(Equal(job.Networks[0].StaticIPs))
		})
	})
})
