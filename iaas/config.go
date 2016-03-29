package iaas

import "github.com/pivotal-cf-experimental/destiny/core"

type Config interface {
	NetworkSubnet() core.NetworkSubnetCloudProperties
	Compilation() core.CompilationCloudProperties
	ResourcePool() core.ResourcePoolCloudProperties
	CPI() CPI
	Properties(staticIP string) Properties
	Stemcell() string
}
