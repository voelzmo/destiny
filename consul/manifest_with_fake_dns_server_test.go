package consul_test

import (
	"github.com/pivotal-cf-experimental/destiny/consul"
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/gomegamatchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("manifest with fake dns server", func() {
	Context("NewManifestWithFakeDNSServer", func() {
		It("creates a deployment manifest with fake-dns-server", func() {
			manifest, err := consul.NewManifestWithFakeDNSServer(consul.ConfigV2{
				DirectorUUID:   "some-director-uuid",
				Name:           "consul-some-random-guid",
				TurbulenceHost: "10.244.4.32",
				AZs: []consul.ConfigAZ{
					{
						IPRange: "10.244.4.0/24",
						Nodes:   2,
						Name:    "z1",
					},
					{
						IPRange: "10.244.5.0/24",
						Nodes:   1,
						Name:    "z2",
					},
				},
			}, iaas.NewWardenConfig())
			Expect(err).NotTo(HaveOccurred())

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))
			Expect(manifest.Name).To(Equal("consul-some-random-guid"))

			Expect(manifest.InstanceGroups[2].Name).To(Equal("fake-dns-server"))
			Expect(manifest.InstanceGroups[2].Instances).To(Equal(1))
			Expect(manifest.InstanceGroups[2].VMType).To(Equal("default"))
			Expect(manifest.InstanceGroups[2].Stemcell).To(Equal("linux"))
			Expect(manifest.InstanceGroups[2].PersistentDiskType).To(Equal("default"))
			Expect(manifest.InstanceGroups[2].Jobs).To(gomegamatchers.ContainSequence([]core.InstanceGroupJob{
				{
					Name:    "fake-dns-server",
					Release: "consul",
				},
			}))
			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "consul",
				Version: "latest",
			}))
			Expect(manifest.Properties.ConsulTestConsumer.NameServer).To(Equal("10.244.4.13"))
		})
	})
})
