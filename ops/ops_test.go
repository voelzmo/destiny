package ops_test

import (
	"errors"

	"github.com/pivotal-cf-experimental/destiny/ops"
	"github.com/pivotal-cf-experimental/gomegamatchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ops", func() {
	Describe("ApplyOp", func() {
		It("returns a manifest with a ops change", func() {
			manifest := "name: some-name"
			modifiedManifest, err := ops.ApplyOp(manifest, ops.Op{
				Type:  "replace",
				Path:  "/name",
				Value: "some-changed-name",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(modifiedManifest).To(Equal("name: some-changed-name"))
		})
	})

	Describe("ApplyOps", func() {
		It("returns a manifest with a set of ops changes", func() {
			manifest := `
---
name: some-name
favorite_color: blue`
			modifiedManifest, err := ops.ApplyOps(manifest, []ops.Op{
				{
					Type:  "replace",
					Path:  "/name",
					Value: "some-changed-name",
				},
				{
					Type:  "replace",
					Path:  "/favorite_color",
					Value: "red",
				},
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(modifiedManifest).To(gomegamatchers.MatchYAML(`
---
name: some-changed-name
favorite_color: red`))
		})

		Context("failure cases", func() {
			Context("when apply ops fails to marshal", func() {
				BeforeEach(func() {
					ops.SetMarshal(func(interface{}) ([]byte, error) {
						return []byte{}, errors.New("failed to marshal")
					})
				})

				AfterEach(func() {
					ops.ResetMarshal()
				})

				It("returns an error", func() {
					_, err := ops.ApplyOps("some-manifest", []ops.Op{})
					Expect(err).To(MatchError("failed to marshal"))
				})
			})
		})
	})
})
