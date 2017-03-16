package turbulencewithops_test

import (
	"io/ioutil"

	"github.com/pivotal-cf-experimental/destiny/turbulencewithops"
	"github.com/pivotal-cf-experimental/gomegamatchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ManifestV2", func() {
	Describe("NewManifestV2", func() {
		It("returns a YAML representation of the turbulence manifest", func() {
			turbulenceManifest, err := ioutil.ReadFile("fixtures/turbulence_manifest_v2.yml")
			Expect(err).NotTo(HaveOccurred())

			manifest, err := turbulencewithops.NewManifestV2(turbulencewithops.ConfigV2{
				DirectorUUID:     "some-director-uuid",
				Name:             "turbulence",
				AZs:              []string{"z1"},
				DirectorHost:     "some-director-host",
				DirectorUsername: "some-director-user",
				DirectorPassword: "some-director-password",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(manifest).To(gomegamatchers.MatchYAML(turbulenceManifest))
		})
	})
})
