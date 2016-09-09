package consul

import (
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/destiny/turbulence"
)

func NewManifestWithTurbulenceAgent(config Config, iaasConfig iaas.Config) (Manifest, error) {
	manifest, err := NewManifest(config, iaasConfig)
	if err != nil {
		return Manifest{}, err
	}
	consulTestConsumerJob := findJob(manifest, "consul_test_consumer")
	consulTestConsumerJob.Templates = append(consulTestConsumerJob.Templates, core.JobTemplate{
		Name:    "turbulence_agent",
		Release: "turbulence",
	})

	manifest.Releases = append(manifest.Releases, core.Release{
		Name:    "turbulence",
		Version: "latest",
	})

	manifest.Properties.TurbulenceAgent = &core.PropertiesTurbulenceAgent{
		API: core.PropertiesTurbulenceAgentAPI{
			Host:     config.TurbulenceHost,
			Password: turbulence.DEFAULT_PASSWORD,
			CACert:   turbulence.APICACert,
		},
	}

	return manifest, nil
}
