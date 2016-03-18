package destiny

import "fmt"

type AWSConfig struct {
	config   Config
	subnet   string
	staticIP string
}

func NewAWSConfig(config Config, subnet string, staticIP string) AWSConfig {
	return AWSConfig{
		config:   config,
		subnet:   subnet,
		staticIP: staticIP,
	}
}

func (a AWSConfig) NetworkSubnet() NetworkSubnetCloudProperties {
	return NetworkSubnetCloudProperties{
		Subnet: a.subnet,
	}
}

func (a AWSConfig) Compilation() CompilationCloudProperties {
	return CompilationCloudProperties{
		InstanceType:     "m3.medium",
		AvailabilityZone: "us-east-1a",
		EphemeralDisk: &CompilationCloudPropertiesEphemeralDisk{
			Size: 1024,
			Type: "gp2",
		},
	}
}

func (a AWSConfig) ResourcePool() ResourcePoolCloudProperties {
	return ResourcePoolCloudProperties{
		InstanceType:     "m3.medium",
		AvailabilityZone: "us-east-1a",
		EphemeralDisk: &ResourcePoolCloudPropertiesEphemeralDisk{
			Size: 1024,
			Type: "gp2",
		},
	}
}

func (a AWSConfig) CPI() CPI {
	return CPI{
		JobName:     "aws_cpi",
		ReleaseName: "bosh-aws-cpi",
	}
}

func (a AWSConfig) Properties() Properties {
	return Properties{
		AWS: &PropertiesAWS{
			AccessKeyID:           a.config.AWS.AccessKeyID,
			SecretAccessKey:       a.config.AWS.SecretAccessKey,
			DefaultKeyName:        a.config.AWS.DefaultKeyName,
			DefaultSecurityGroups: a.config.AWS.DefaultSecurityGroups,
			Region:                a.config.AWS.Region,
		},
		Registry: &PropertiesRegistry{
			Host:     a.config.Registry.Host,
			Password: a.config.Registry.Password,
			Port:     a.config.Registry.Port,
			Username: a.config.Registry.Username,
		},
		Blobstore: &PropertiesBlobstore{
			Address: a.staticIP,
			Port:    2520,
			Agent: PropertiesBlobstoreAgent{
				User:     "agent",
				Password: "agent-password",
			},
		},
		Agent: &PropertiesAgent{
			Mbus: fmt.Sprintf("nats://nats:password@%s:4222", a.staticIP),
		},
	}
}
