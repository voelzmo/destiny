package turbulence_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/destiny/turbulence"
	. "github.com/pivotal-cf-experimental/gomegamatchers"
)

var _ = Describe("Manifest", func() {
	Describe("ToYAML", func() {
		It("returns a YAML representation of the turbulence manifest", func() {
			turbulenceManifest, err := ioutil.ReadFile("fixtures/turbulence_manifest.yml")
			Expect(err).NotTo(HaveOccurred())

			manifest := turbulence.NewTurbulence(turbulence.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "turbulence",
				IPRange:      "10.244.4.0/24",
				BOSH: turbulence.ConfigBOSH{
					Target:   "some-bosh-target",
					Username: "some-bosh-username",
					Password: "some-bosh-password",
				},
			}, iaas.NewWardenConfig())

			yaml, err := manifest.ToYAML()
			Expect(err).NotTo(HaveOccurred())
			Expect(yaml).To(MatchYAML(turbulenceManifest))
		})
	})
})
