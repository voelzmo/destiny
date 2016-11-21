package consul

import (
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
)

func NewManifestWithFakeDNSServer(config ConfigV2, iaasConfig iaas.Config) (ManifestV2, error) {
	manifest, err := NewManifestV2(config, iaasConfig)
	if err != nil {
		return ManifestV2{}, err
	}

	cidr, err := core.ParseCIDRBlock(config.AZs[0].IPRange)
	if err != nil {
		// not tested
		return ManifestV2{}, err
	}

	ip := cidr.GetFirstIP().Add(13).String()

	persistentDiskType := "default"
	if config.PersistentDiskType != "" {
		persistentDiskType = config.PersistentDiskType
	}

	vmType := "default"
	if config.VMType != "" {
		vmType = config.VMType
	}

	manifest.InstanceGroups = append(manifest.InstanceGroups, core.InstanceGroup{
		Name:               "fake-dns-server",
		Instances:          1,
		AZs:                []string{"z1"},
		VMType:             vmType,
		Stemcell:           "linux",
		PersistentDiskType: persistentDiskType,
		Networks: []core.InstanceGroupNetwork{
			{
				Name:      "private",
				StaticIPs: []string{ip},
			},
		},
		Jobs: []core.InstanceGroupJob{
			{
				Name:    "fake-dns-server",
				Release: "consul",
			},
		},
	})

	manifest.Properties.ConsulTestConsumer = &core.ConsulTestConsumer{
		NameServer: ip,
	}

	return manifest, nil
}
