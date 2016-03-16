package destiny

type AWSConfig struct {
	Subnet string
}

func NewAWSConfig(subnet string) AWSConfig {
	return AWSConfig{
		Subnet: subnet,
	}
}

func (a AWSConfig) NetworkSubnet() NetworkSubnetCloudProperties {
	return NetworkSubnetCloudProperties{
		Subnet: a.Subnet,
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
