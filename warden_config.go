package destiny

type WardenConfig struct {
}

func NewWardenConfig() WardenConfig {
	return WardenConfig{}
}

func (a WardenConfig) NetworkSubnet() NetworkSubnetCloudProperties {
	return NetworkSubnetCloudProperties{Subnet: "random"}
}

func (a WardenConfig) Compilation() CompilationCloudProperties {
	return CompilationCloudProperties{}
}

func (a WardenConfig) ResourcePool() ResourcePoolCloudProperties {
	return ResourcePoolCloudProperties{}
}

func (a WardenConfig) CPI() CPI {
	return CPI{
		JobName:     "warden_cpi",
		ReleaseName: "bosh-warden-cpi",
	}
}
