package iaas

import "github.com/pivotal-cf-experimental/destiny/core"

type GCPConfig struct {
	SubnetCloudProperties GCPSubnetCloudProperties
}

type GCPSubnetCloudProperties struct {
	NetworkName    string
	SubnetworkName string
	Tags           []string
}

func (g GCPConfig) NetworkSubnet(ipRange string) core.NetworkSubnetCloudProperties {
	return core.NetworkSubnetCloudProperties{
		NetworkName:         g.SubnetCloudProperties.NetworkName,
		SubnetworkName:      g.SubnetCloudProperties.SubnetworkName,
		EphemeralExternalIP: true,
		Tags:                g.SubnetCloudProperties.Tags,
	}
}

func (g GCPConfig) Compilation(zone string) core.CompilationCloudProperties {
	return core.CompilationCloudProperties{
		MachineType:    "n1-standard-2",
		Zone:           zone,
		RootDiskSizeGB: "100",
		RootDiskType:   "pd-ssd",
	}
}

func (g GCPConfig) ResourcePool(ipRange string) core.ResourcePoolCloudProperties {
	return core.ResourcePoolCloudProperties{}
}

func (g GCPConfig) CPI() CPI {
	return CPI{}
}

func (g GCPConfig) Properties(staticIP string) Properties {
	return Properties{}
}

func (g GCPConfig) Stemcell() string {
	return ""
}

func (g GCPConfig) WindowsStemcell() string {
	return ""
}
