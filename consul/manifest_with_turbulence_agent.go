package consul

import (
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/destiny/turbulence"
)

func NewManifestWithTurbulenceAgent(config ConfigV2, iaasConfig iaas.Config) (ManifestV2, error) {
	manifest, err := NewManifestWithFakeDNSServer(config, iaasConfig)
	if err != nil {
		return ManifestV2{}, err
	}

	manifest.Releases = append(manifest.Releases, core.Release{
		Name:    "turbulence",
		Version: "latest",
	})

	manifest.InstanceGroups[2].Jobs = append(manifest.InstanceGroups[2].Jobs, core.InstanceGroupJob{
		Name:    "turbulence_agent",
		Release: "turbulence",
	})

	manifest.Properties.TurbulenceAgent = &core.PropertiesTurbulenceAgent{
		API: core.PropertiesTurbulenceAgentAPI{
			Host:     config.TurbulenceHost,
			Password: turbulence.DefaultPassword,
			CACert:   turbulence.APICACert,
		},
	}

	return manifest, nil
}
