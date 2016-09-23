package core_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/destiny/core"
	. "github.com/pivotal-cf-experimental/gomegamatchers"
	"gopkg.in/yaml.v2"
)

var _ = Describe("Manifest", func() {
	Describe("PropertiesTurbulenceAgentAPI", func() {
		It("serializes the turbulence properties", func() {
			expectedYAML := `host: 1.2.3.4
password: secret
ca_cert: some-cert`
			actualYAML, err := yaml.Marshal(core.PropertiesTurbulenceAgentAPI{
				Host:     "1.2.3.4",
				Password: "secret",
				CACert:   "some-cert",
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(actualYAML).To(MatchYAML(expectedYAML))
		})
	})

})
