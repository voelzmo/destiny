package iaas_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
)

var _ = Describe("AWS Config", func() {
	var (
		awsConfig iaas.AWSConfig
	)

	BeforeEach(func() {
		awsConfig = iaas.AWSConfig{
			AccessKeyID:           "some-access-key-id",
			SecretAccessKey:       "some-secret-access-key",
			DefaultKeyName:        "some-default-key-name",
			DefaultSecurityGroups: []string{"some-default-security-group"},
			Region:                "some-region",
			Subnet:                "some-subnet",
			RegistryHost:          "some-host",
			RegistryPassword:      "some-password",
			RegistryPort:          1234,
			RegistryUsername:      "some-username",
		}
	})

	Describe("NetworkSubnet", func() {
		It("returns a network subnet specific to AWS", func() {
			subnetCloudProperties := awsConfig.NetworkSubnet()
			Expect(subnetCloudProperties).To(Equal(core.NetworkSubnetCloudProperties{
				Subnet: "some-subnet",
			}))
		})
	})

	Describe("Compilation", func() {
		It("returns a compilation specific to AWS", func() {
			compilationCloudProperties := awsConfig.Compilation()
			Expect(compilationCloudProperties).To(Equal(core.CompilationCloudProperties{
				InstanceType:     "m3.medium",
				AvailabilityZone: "us-east-1a",
				EphemeralDisk: &core.CompilationCloudPropertiesEphemeralDisk{
					Size: 1024,
					Type: "gp2",
				},
			}))
		})
	})

	Describe("ResourcePool", func() {
		It("returns a resource pool specific to AWS", func() {
			resourcePoolCloudProperties := awsConfig.ResourcePool()
			Expect(resourcePoolCloudProperties).To(Equal(core.ResourcePoolCloudProperties{
				InstanceType:     "m3.medium",
				AvailabilityZone: "us-east-1a",
				EphemeralDisk: &core.ResourcePoolCloudPropertiesEphemeralDisk{
					Size: 1024,
					Type: "gp2",
				},
			}))
		})
	})

	Describe("CPI", func() {
		It("returns the cpi specific to AWS", func() {
			cpi := awsConfig.CPI()
			Expect(cpi).To(Equal(iaas.CPI{
				JobName:     "aws_cpi",
				ReleaseName: "bosh-aws-cpi",
			}))
		})
	})

	Describe("Properties", func() {
		It("returns the properties specific to AWS", func() {
			properties := awsConfig.Properties("some-static-ip")

			Expect(properties).To(Equal(iaas.Properties{
				AWS: &iaas.PropertiesAWS{
					AccessKeyID:           "some-access-key-id",
					SecretAccessKey:       "some-secret-access-key",
					DefaultKeyName:        "some-default-key-name",
					DefaultSecurityGroups: []string{"some-default-security-group"},
					Region:                "some-region",
				},
				Registry: &core.PropertiesRegistry{
					Host:     "some-host",
					Password: "some-password",
					Port:     1234,
					Username: "some-username",
				},
				Blobstore: &core.PropertiesBlobstore{
					Address: "some-static-ip",
					Port:    2520,
					Agent: core.PropertiesBlobstoreAgent{
						User:     "agent",
						Password: "agent-password",
					},
				},
				Agent: &core.PropertiesAgent{
					Mbus: fmt.Sprintf("nats://nats:password@some-static-ip:4222"),
				},
			}))
		})
	})
})
