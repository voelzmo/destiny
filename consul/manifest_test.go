package consul_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/destiny/consul"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	. "github.com/pivotal-cf-experimental/gomegamatchers"
)

var _ = Describe("Manifest", func() {
	Describe("ToYAML", func() {
		It("returns a YAML representation of the consul manifest", func() {
			consulManifest, err := ioutil.ReadFile("fixtures/consul_manifest.yml")
			Expect(err).NotTo(HaveOccurred())

			manifest := consul.NewConsul(consul.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "consul",
				IPRange:      "10.244.4.0/24",
			}, iaas.NewWardenConfig())

			yaml, err := manifest.ToYAML()
			Expect(err).NotTo(HaveOccurred())
			Expect(yaml).To(MatchYAML(consulManifest))
		})
	})
})
