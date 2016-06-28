package etcd_test

import (
	"io/ioutil"

	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/etcd"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/gomegamatchers"

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
		var (
			manifest etcd.Manifest
		)

		BeforeEach(func() {
			manifest = etcd.Manifest{
				Jobs: []core.Job{
					{
						Name: "some-other-job-1",
					},
					{
						Name:       "some-other-job-2",
						Properties: &core.JobProperties{},
					},
					{
						Name: "etcd_no_tls",
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
					{
						Name: "etcd_tls",
						Properties: &core.JobProperties{
							Etcd: &core.JobPropertiesEtcd{
								RequireSSL: true,
								CACert:     "some-ca-cert",
								ClientCert: "some-client-cert",
								ClientKey:  "some-client-key",
							},
						},
					},
				},
			}
		})

		It("returns a manifest with a proxy job", func() {
			manifest, err := manifest.ReplaceEtcdWithProxyJob("etcd_no_tls")
			Expect(err).NotTo(HaveOccurred())

			Expect(manifest.Jobs).To(Equal([]core.Job{
				{
					Name: "some-other-job-1",
				},
				{
					Name:       "some-other-job-2",
					Properties: &core.JobProperties{},
				},
				{
					Name: "etcd_no_tls",
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
				{
					Name: "etcd_tls",
					Properties: &core.JobProperties{
						Etcd: &core.JobPropertiesEtcd{
							RequireSSL: true,
							CACert:     "some-ca-cert",
							ClientCert: "some-client-cert",
							ClientKey:  "some-client-key",
						},
					},
				},
			}))

			Expect(manifest.Properties.EtcdProxy).To(Equal(&etcd.PropertiesEtcdProxy{
				Etcd: etcd.PropertiesEtcdProxyEtcd{
					URL:        "https://etcd.service.cf.internal",
					CACert:     "some-ca-cert",
					ClientCert: "some-client-cert",
					ClientKey:  "some-client-key",
					RequireSSL: true,
					Port:       4001,
				},
			}))
		})

		It("generates correct YAML", func() {
			manifest, err := manifest.ReplaceEtcdWithProxyJob("etcd_no_tls")
			Expect(err).NotTo(HaveOccurred())

			yaml, err := manifest.ToYAML()
			Expect(err).NotTo(HaveOccurred())

			expectedYAML, err := ioutil.ReadFile("fixtures/proxy_replacement.yml")
			Expect(err).NotTo(HaveOccurred())

			Expect(yaml).To(gomegamatchers.MatchYAML(expectedYAML))
		})

		Context("failure cases", func() {
			It("returns an error when job is not found", func() {
				manifest := etcd.Manifest{}
				_, err := manifest.ReplaceEtcdWithProxyJob("etcd_z1")
				Expect(err).To(MatchError("job not found"))
			})
		})
	})
})
