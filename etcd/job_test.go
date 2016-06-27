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
				Expect(properties.EtcdTestConsumer.Etcd.Machines).To(Equal(job.Networks[0].StaticIPs))
				Expect(properties.Etcd.Machines).To(Equal(job.Networks[0].StaticIPs))
				Expect(properties.Etcd.Cluster[0].Instances).To(Equal(1))

				job, properties = etcd.SetJobInstanceCount(job, network, properties, 3)
				Expect(job.Instances).To(Equal(3))
				Expect(job.Networks[0].StaticIPs).To(HaveLen(3))
				Expect(properties.EtcdTestConsumer.Etcd.Machines).To(Equal(job.Networks[0].StaticIPs))
				Expect(properties.Etcd.Machines).To(Equal(job.Networks[0].StaticIPs))
				Expect(properties.Etcd.Cluster[0].Instances).To(Equal(3))
			})

			It("does not override the machines property if ssl is enabled", func() {
				manifest := etcd.NewTLSManifest(etcd.Config{
					IPRange: "10.244.4.0/24",
				}, iaas.NewWardenConfig())
				job := manifest.Jobs[1]
				network := manifest.Networks[0]
				properties := manifest.Properties

				Expect(job.Instances).To(Equal(1))
				Expect(job.Networks[0].StaticIPs).To(HaveLen(1))
				Expect(job.Networks[0].Name).To(Equal(network.Name))
				Expect(properties.EtcdTestConsumer.Etcd.Machines).To(Equal([]string{"etcd.service.cf.internal"}))
				Expect(properties.Etcd.Machines).To(Equal([]string{"etcd.service.cf.internal"}))
				Expect(properties.Etcd.Cluster[0].Instances).To(Equal(1))

				job, properties = etcd.SetJobInstanceCount(job, network, properties, 3)
				Expect(job.Instances).To(Equal(3))
				Expect(job.Networks[0].StaticIPs).To(HaveLen(3))
				Expect(properties.EtcdTestConsumer.Etcd.Machines).To(Equal([]string{"etcd.service.cf.internal"}))
				Expect(properties.Etcd.Machines).To(Equal([]string{"etcd.service.cf.internal"}))
				Expect(properties.Etcd.Cluster[0].Instances).To(Equal(3))
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

	Describe("ReplaceEtcdWithProxyJob", func() {
		It("returns a manifest with a proxy job", func() {
			manifest := etcd.Manifest{
				Jobs: []core.Job{
					{
						Name: "some-other-job",
					},
					{

						Name:      "etcd_z1",
						Instances: 1,
						Networks: []core.JobNetwork{{
							Name:      "etcd1",
							StaticIPs: []string{"10.244.4.4"},
						}},
						PersistentDisk: 1024,
						ResourcePool:   "etcd_z1",
						Templates: []core.JobTemplate{
							{
								Name:    "consul",
								Release: "consul",
							},
							{
								Name:    "etcd",
								Release: "etcd",
							},
						},
					},
				},
			}

			manifest = manifest.ReplaceEtcdWithProxyJob("etcd_z1")

			Expect(manifest.Jobs).To(Equal([]core.Job{
				{
					Name: "some-other-job",
				},
				{

					Name:      "etcd_z1",
					Instances: 1,
					Networks: []core.JobNetwork{{
						Name:      "etcd1",
						StaticIPs: []string{"10.244.4.4"},
					}},
					PersistentDisk: 1024,
					ResourcePool:   "etcd_z1",
					Templates: []core.JobTemplate{
						{
							Name:    "consul",
							Release: "consul",
						},
						{
							Name:    "etcd_proxy",
							Release: "etcd",
						},
					},
				},
			}))
		})
	})
})
