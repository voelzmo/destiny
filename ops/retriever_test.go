package ops_test

import (
	"github.com/pivotal-cf-experimental/destiny/ops"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Retrievers", func() {
	Describe("ManifestName", func() {
		It("returns the deployment name given a manifest", func() {
			name, err := ops.ManifestName("name: some-name")
			Expect(err).NotTo(HaveOccurred())

			Expect(name).To(Equal("some-name"))
		})

		Context("failure cases", func() {
			Context("when the manifest yaml is invalid", func() {
				It("returns an error", func() {
					_, err := ops.ManifestName("%%%")
					Expect(err).To(MatchError("yaml: could not find expected directive name"))
				})
			})

			Context("when the manifest name is empty", func() {
				It("returns an error", func() {
					_, err := ops.ManifestName("hello: world")
					Expect(err).To(MatchError("could not find name in manifest"))
				})
			})
		})
	})

	Describe("InstanceGroups", func() {
		It("returns the instance groups given a manifest", func() {
			instanceGroups, err := ops.InstanceGroups(`
instance_groups:
- name: consul
  instances: 1
- name: etcd
  instances: 3
- name: testconsumer
  instances: 1`)
			Expect(err).NotTo(HaveOccurred())
			Expect(instanceGroups).To(Equal([]ops.InstanceGroup{
				{
					Name:      "consul",
					Instances: 1,
				},
				{
					Name:      "etcd",
					Instances: 3,
				},
				{
					Name:      "testconsumer",
					Instances: 1,
				},
			}))
		})

		Context("failure cases", func() {
			Context("when the manifest yaml is invalid", func() {
				It("returns an error", func() {
					_, err := ops.InstanceGroups("%%%")
					Expect(err).To(MatchError("yaml: could not find expected directive name"))
				})
			})
		})
	})
})
