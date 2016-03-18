package destiny

type WardenConfig struct {
}

func NewWardenConfig() WardenConfig {
	return WardenConfig{}
}

func (WardenConfig) NetworkSubnet() NetworkSubnetCloudProperties {
	return NetworkSubnetCloudProperties{Name: "random"}
}

func (WardenConfig) Compilation() CompilationCloudProperties {
	return CompilationCloudProperties{}
}

func (WardenConfig) ResourcePool() ResourcePoolCloudProperties {
	return ResourcePoolCloudProperties{}
}

func (WardenConfig) CPI() CPI {
	return CPI{
		JobName:     "warden_cpi",
		ReleaseName: "bosh-warden-cpi",
	}
}

func (WardenConfig) Properties() Properties {
	return Properties{
		WardenCPI: &PropertiesWardenCPI{
			Agent: PropertiesWardenCPIAgent{
				Blobstore: PropertiesWardenCPIAgentBlobstore{
					Options: PropertiesWardenCPIAgentBlobstoreOptions{
						Endpoint: "http://10.254.50.4:25251",
						Password: "agent-password",
						User:     "agent",
					},
					Provider: "dav",
				},
				Mbus: "nats://nats:nats-password@10.254.50.4:4222",
			},
			Warden: PropertiesWardenCPIWarden{
				ConnectAddress: "10.254.50.4:7777",
				ConnectNetwork: "tcp",
			},
		},
	}
}
