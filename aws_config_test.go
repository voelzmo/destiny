package destiny_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/destiny"
)

var _ = Describe("AWS Config", func() {

	var (
		subnet    string
		awsConfig destiny.AWSConfig
	)

	BeforeEach(func() {
		subnet = "some-subnet"
		awsConfig = destiny.NewAWSConfig(subnet)
	})

	Describe("NetworkSubnet", func() {
		It("returns a network subnet specific to AWS", func() {
			subnetCloudProperties := awsConfig.NetworkSubnet()
			Expect(subnetCloudProperties).To(Equal(destiny.NetworkSubnetCloudProperties{
				Name:   "random",
				Subnet: subnet,
			}))
		})
	})

	Describe("Compilation", func() {
		It("returns a compilation specific to AWS", func() {
			compilationCloudProperties := awsConfig.Compilation()
			Expect(compilationCloudProperties).To(Equal(destiny.CompilationCloudProperties{
				InstanceType:     "m3.medium",
				AvailabilityZone: "us-east-1a",
				EphemeralDisk: &destiny.CompilationCloudPropertiesEphemeralDisk{
					Size: 1024,
					Type: "gp2",
				},
			}))
		})
	})

	Describe("ResourcePool", func() {
		It("returns a resource pool specific to AWS", func() {
			resourcePoolCloudProperties := awsConfig.ResourcePool()
			Expect(resourcePoolCloudProperties).To(Equal(destiny.ResourcePoolCloudProperties{
				InstanceType:     "m3.medium",
				AvailabilityZone: "us-east-1a",
				EphemeralDisk: &destiny.ResourcePoolCloudPropertiesEphemeralDisk{
					Size: 1024,
					Type: "gp2",
				},
			}))
		})
	})

	Describe("CPI", func() {
		It("returns the cpi specific to AWS", func() {
			cpi := awsConfig.CPI()
			Expect(cpi).To(Equal(destiny.CPI{
				JobName:     "aws_cpi",
				ReleaseName: "bosh-aws-cpi",
			}))
		})
	})
})
