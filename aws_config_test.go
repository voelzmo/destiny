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
		//destiny.Config{
		//AWS: destiny.ConfigAWS{
		//AccessKeyID:           "some-access-key-id",
		//SecretAccessKey:       "some-secret-access-key",
		//DefaultKeyName:        "some-default-key-name",
		//DefaultSecurityGroups: []string{"some-default-security-group"},
		//Region:                "some-region",
		//Subnet:                "some-subnet",
		//},
		//Registry: destiny.ConfigRegistry{
		//Host:     "some-host",
		//Password: "some-password",
		//Port:     "some-port",
		//Username: "some-username",
		//},
		//})
	})

	Describe("NetworkSubnet", func() {
		It("returns a network subnet specific to AWS", func() {
			subnetCloudProperties := awsConfig.NetworkSubnet()
			Expect(subnetCloudProperties).To(Equal(destiny.NetworkSubnetCloudProperties{
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

	Describe("Properties", func() {
		PIt("returns the properties specific to AWS", func() {
			//properties = awsConfig.Properties()
			//Expect(properties).To(Equal(destiny.Properties{
			//AWS: &destiny.PropertiesAWS{
			//AccessKeyID:           config.AWS.AccessKeyID,
			//SecretAccessKey:       config.AWS.SecretAccessKey,
			//DefaultKeyName:        config.AWS.DefaultKeyName,
			//DefaultSecurityGroups: config.AWS.DefaultSecurityGroups,
			//Region:                config.AWS.Region,
			//},
			//Registry: &destiny.PropertiesRegistry{
			//Host:     config.Registry.Host,
			//Password: config.Registry.Password,
			//Port:     config.Registry.Port,
			//Username: config.Registry.Username,
			//},
			//Blobstore: &destiny.PropertiesBlobstore{
			//Address: turbulenceNetwork.StaticIPs(1)[0],
			//Port:    2520,
			//Agent: destiny.PropertiesBlobstoreAgent{
			//User:     "agent",
			//Password: "agent-password",
			//},
			//},
			//Agent: &destiny.PropertiesAgent{
			//Mbus: fmt.Sprintf("nats://nats:password@%s:4222", turbulenceNetwork.StaticIPs(1)[0]),
			//},
			//}))
		})
	})
})
