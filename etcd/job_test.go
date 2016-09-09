package etcd_test

import (
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/etcd"
	"github.com/pivotal-cf-experimental/destiny/iaas"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manifest", func() {
	Describe("Job", func() {
		Describe("SetEtcdProperties", func() {
			It("updates the etcd and testconsumer properties to match the current job configuration", func() {
				manifest, err := etcd.NewManifest(etcd.Config{
					IPRange: "10.244.4.0/24",
				}, iaas.NewWardenConfig())
				Expect(err).NotTo(HaveOccurred())

				job := manifest.Jobs[1]
				properties := manifest.Properties

				Expect(properties.EtcdTestConsumer.Etcd.Machines).To(Equal(job.Networks[0].StaticIPs))
				Expect(properties.Etcd.Machines).To(Equal(job.Networks[0].StaticIPs))
				Expect(properties.Etcd.Cluster[0].Instances).To(Equal(1))

				job.Instances = 3
				job.Networks[0].StaticIPs = []string{"ip1", "ip2", "ip3"}

				properties = etcd.SetEtcdProperties(job, properties)
				Expect(properties.EtcdTestConsumer.Etcd.Machines).To(Equal(job.Networks[0].StaticIPs))
				Expect(properties.Etcd.Machines).To(Equal(job.Networks[0].StaticIPs))
				Expect(properties.Etcd.Cluster[0].Instances).To(Equal(3))
			})

			It("does not override the machines property if ssl is enabled", func() {
				manifest, err := etcd.NewTLSManifest(etcd.Config{
					IPRange: "10.244.4.0/24",
				}, iaas.NewWardenConfig())
				Expect(err).NotTo(HaveOccurred())

				job := manifest.Jobs[1]
				properties := manifest.Properties

				Expect(properties.EtcdTestConsumer.Etcd.Machines).To(Equal([]string{"etcd.service.cf.internal"}))
				Expect(properties.Etcd.Machines).To(Equal([]string{"etcd.service.cf.internal"}))
				Expect(properties.Etcd.Cluster[0].Instances).To(Equal(1))

				job.Instances = 3
				job.Networks[0].StaticIPs = []string{"ip1", "ip2", "ip3"}

				properties = etcd.SetEtcdProperties(job, properties)
				Expect(properties.EtcdTestConsumer.Etcd.Machines).To(Equal([]string{"etcd.service.cf.internal"}))
				Expect(properties.Etcd.Machines).To(Equal([]string{"etcd.service.cf.internal"}))
				Expect(properties.Etcd.Cluster[0].Instances).To(Equal(3))
			})
		})

		Describe("SetJobInstanceCount", func() {
			It("sets the correct values for instances and static_ips given a count", func() {
				manifest, err := etcd.NewManifest(etcd.Config{
					IPRange: "10.244.4.0/24",
				}, iaas.NewWardenConfig())
				Expect(err).NotTo(HaveOccurred())

				job := manifest.Jobs[1]
				network := manifest.Networks[0]

				Expect(job.Instances).To(Equal(1))
				Expect(job.Networks[0].StaticIPs).To(HaveLen(1))
				Expect(job.Networks[0].Name).To(Equal(network.Name))

				job, err = etcd.SetJobInstanceCount(job, network, 3, 0)

				Expect(err).NotTo(HaveOccurred())
				Expect(job.Instances).To(Equal(3))
				Expect(job.Networks[0].StaticIPs).To(HaveLen(3))
			})

			It("returns a job with an offsetted static ip list", func() {
				manifest, err := etcd.NewManifest(etcd.Config{
					IPRange: "10.244.4.0/24",
				}, iaas.NewWardenConfig())
				Expect(err).NotTo(HaveOccurred())

				job := manifest.Jobs[1]
				network := manifest.Networks[0]

				Expect(job.Instances).To(Equal(1))
				Expect(job.Networks[0].StaticIPs).To(HaveLen(1))
				Expect(job.Networks[0].Name).To(Equal(network.Name))

				job, err = etcd.SetJobInstanceCount(job, network, 3, 3)

				Expect(err).NotTo(HaveOccurred())
				Expect(job.Instances).To(Equal(3))
				Expect(job.Networks[0].StaticIPs).To(HaveLen(3))
				Expect(job.Networks[0].StaticIPs).To(ConsistOf([]string{
					"10.244.4.7",
					"10.244.4.8",
					"10.244.4.9",
				}))
			})

			Context("failure cases", func() {
				It("returns an error when there aren't enough static ips", func() {
					manifest, err := etcd.NewManifest(etcd.Config{
						IPRange: "10.244.4.0/24",
					}, iaas.NewWardenConfig())
					Expect(err).NotTo(HaveOccurred())

					job := manifest.Jobs[1]
					network := manifest.Networks[0]

					job, err = etcd.SetJobInstanceCount(job, network, 1000, 0)
					Expect(err).To(MatchError("can't allocate 1000 ips from 248 available ips"))
				})
			})
		})
	})

	Describe("RemoveJob", func() {
		It("returns a manifest with specified job removed", func() {
			manifest := etcd.Manifest{
				Jobs: []core.Job{
					{Name: "some-job-name-1"},
					{Name: "some-job-name-2"},
					{Name: "some-job-name-3"},
				},
			}

			manifest = manifest.RemoveJob("some-job-name-2")

			Expect(manifest.Jobs).To(Equal([]core.Job{
				{Name: "some-job-name-1"},
				{Name: "some-job-name-3"},
			}))
		})
	})
})
