package consul_test

import (
	"github.com/pivotal-cf-experimental/destiny/consul"
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/destiny/turbulence"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manifest", func() {
	Describe("NewManifestWithTurbulenceAgent", func() {
		It("generates a valid Consul BOSH-lite manifest with additional turbulence agent on test consumer", func() {
			manifest, err := consul.NewManifestWithTurbulenceAgent(consul.Config{
				DirectorUUID:   "some-director-uuid",
				Name:           "consul-some-random-guid",
				TurbulenceHost: "10.244.4.32",
				Networks: []consul.ConfigNetwork{
					{
						IPRange: "10.244.4.0/24",
						Nodes:   2,
					},
					{
						IPRange: "10.244.5.0/24",
						Nodes:   1,
					},
				},
			}, iaas.NewWardenConfig())
			Expect(err).NotTo(HaveOccurred())

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))
			Expect(manifest.Name).To(Equal("consul-some-random-guid"))

			consulTestConsumerJob := findJob(manifest, "consul_test_consumer")
			Expect(consulTestConsumerJob.Templates).To(ContainElement(core.JobTemplate{
				Name:    "turbulence_agent",
				Release: "turbulence",
			}))

			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "turbulence",
				Version: "latest",
			}))

			Expect(manifest.Properties.TurbulenceAgent.API).To(Equal(core.PropertiesTurbulenceAgentAPI{
				Host:     "10.244.4.32",
				Password: turbulence.DEFAULT_PASSWORD,
				CACert:   turbulence.APICACert,
			}))
		})
	})

	Context("failure cases", func() {
		It("returns an error when the manifest creation fails", func() {
			_, err := consul.NewManifestWithTurbulenceAgent(consul.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "consul-some-random-guid",
				Networks: []consul.ConfigNetwork{
					{
						IPRange: "fake-cidr-block",
						Nodes:   1,
					},
				},
			}, iaas.NewWardenConfig())
			Expect(err).To(MatchError(`"fake-cidr-block" cannot parse CIDR block`))
		})
	})
})

func findJob(manifest consul.Manifest, name string) *core.Job {
	for _, job := range manifest.Jobs {
		if job.Name == name {
			return &job
		}
	}
	return &core.Job{}
}
