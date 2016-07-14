package cloudconfig

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/destiny/network"
)

type Config struct {
	AZs []ConfigAZ
}

type ConfigAZ struct {
	IPRange   string
	StaticIPs int
}

type CloudConfig struct {
	AZs         []AZ        `yaml:"azs"`
	VMTypes     []VMType    `yaml:"vm_types"`
	DiskTypes   []DiskType  `yaml:"disk_types"`
	Compilation Compilation `yaml:"compilation"`
	Networks    []Network   `yaml:"networks"`
}

type AZ struct {
	Name string `yaml:"name"`
}

type VMType struct {
	Name string `yaml:"name"`
}

type DiskType struct {
	Name     string `yaml:"name"`
	DiskSize int    `yaml:"disk_size"`
}

type Compilation struct {
	Workers             int    `yaml:"workers"`
	ReuseCompilationVMs bool   `yaml:"reuse_compilation_vms"`
	AZ                  string `yaml:"az"`
	VMType              string `yaml:"vm_type"`
	Network             string `yaml:"network"`
}

type Network struct {
	Name    string   `yaml:"name"`
	Subnets []Subnet `yaml:"subnets"`
	Type    string   `yaml:"type"`
}

type Subnet struct {
	CloudProperties SubnetCloudProperties `yaml:"cloud_properties"`
	Range           string                `yaml:"range"`
	Gateway         string                `yaml:"gateway"`
	AZ              string                `yaml:"az"`
	Reserved        []string              `yaml:"reserved"`
	Static          []string              `yaml:"static"`
}

type SubnetCloudProperties struct {
	Name string
}

func NewCloudConfig(config Config, iaasConfig iaas.Config) CloudConfig {
	vmTypes := []VMType{
		{
			Name: "default",
		},
	}

	diskTypes := []DiskType{
		{
			Name:     "default",
			DiskSize: 1024,
		},
	}

	compilation := Compilation{
		Workers:             5,
		ReuseCompilationVMs: true,
		AZ:                  "z1",
		VMType:              "default",
		Network:             "private",
	}

	azs := []AZ{}
	subnets := []Subnet{}

	for i, cfgAZ := range config.AZs {
		azName := fmt.Sprintf("z%d", i+1)
		azs = append(azs, AZ{
			Name: azName,
		})

		ipRange := network.IPRange(cfgAZ.IPRange)

		staticIPs := []string{}
		for j := 0; j < cfgAZ.StaticIPs; j++ {
			staticIPs = append(staticIPs, ipRange.IP(4+j))
		}

		subnets = append(subnets, Subnet{
			CloudProperties: SubnetCloudProperties{
				Name: "random",
			},
			Range:   string(ipRange),
			Gateway: ipRange.IP(1),
			AZ:      azName,
			Reserved: []string{
				ipRange.Range(2, 3),
				ipRange.Range(13, 254),
			},
			Static: staticIPs,
		})
	}

	networks := []Network{}
	networks = append(networks, Network{
		Name:    "private",
		Subnets: subnets,
		Type:    "manual",
	})

	return CloudConfig{
		AZs:         azs,
		VMTypes:     vmTypes,
		DiskTypes:   diskTypes,
		Compilation: compilation,
		Networks:    networks,
	}
}

func (c CloudConfig) ToYAML() ([]byte, error) {
	return yaml.Marshal(c)
}
