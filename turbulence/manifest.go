package turbulence

import (
	"github.com/cloudfoundry-incubator/candiedyaml"
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
	"github.com/pivotal-cf-experimental/destiny/network"
)

type Manifest struct {
	DirectorUUID  string              `yaml:"director_uuid"`
	Name          string              `yaml:"name"`
	Jobs          []core.Job          `yaml:"jobs"`
	Properties    Properties          `yaml:"properties"`
	Update        core.Update         `yaml:"update"`
	Compilation   core.Compilation    `yaml:"compilation"`
	Networks      []core.Network      `yaml:"networks"`
	Releases      []core.Release      `yaml:"releases"`
	ResourcePools []core.ResourcePool `yaml:"resource_pools"`
}

func NewManifest(config Config, iaasConfig iaas.Config) Manifest {
	turbulenceRelease := core.Release{
		Name:    "turbulence",
		Version: "latest",
	}

	ipRange := network.IPRange(config.IPRange)

	cloudProperties := iaasConfig.NetworkSubnet()
	cpi := iaasConfig.CPI()

	cpiRelease := core.Release{
		Name:    cpi.ReleaseName,
		Version: "latest",
	}

	turbulenceNetwork := core.Network{
		Name: "turbulence",
		Subnets: []core.NetworkSubnet{{
			CloudProperties: cloudProperties,
			Gateway:         ipRange.IP(1),
			Range:           string(ipRange),
			Reserved:        []string{ipRange.Range(2, 11), ipRange.Range(17, 254)},
			Static: []string{
				ipRange.IP(12),
				ipRange.IP(13),
			},
		}},
		Type: "manual",
	}

	compilation := core.Compilation{
		Network:             turbulenceNetwork.Name,
		ReuseCompilationVMs: true,
		Workers:             3,
		CloudProperties:     iaasConfig.Compilation(),
	}

	turbulenceResourcePool := core.ResourcePool{
		Name:    "turbulence",
		Network: turbulenceNetwork.Name,
		Stemcell: core.ResourcePoolStemcell{
			Name:    iaasConfig.Stemcell(),
			Version: "latest",
		},
		CloudProperties: iaasConfig.ResourcePool(),
	}

	update := core.Update{
		Canaries:        1,
		CanaryWatchTime: "1000-180000",
		MaxInFlight:     1,
		Serial:          true,
		UpdateWatchTime: "1000-180000",
	}

	apiJob := core.Job{
		Instances: 1,
		Name:      "api",
		Networks: []core.JobNetwork{{
			Name:      turbulenceNetwork.Name,
			StaticIPs: turbulenceNetwork.StaticIPs(1),
		}},
		PersistentDisk: 1024,
		ResourcePool:   turbulenceResourcePool.Name,
		Templates: []core.JobTemplate{
			{
				Name:    "turbulence_api",
				Release: turbulenceRelease.Name,
			},
			{
				Name:    cpi.JobName,
				Release: cpiRelease.Name,
			},
		},
	}

	directorCACert := APIDirectorCACert
	if config.BOSH.DirectorCACert != "" {
		directorCACert = config.BOSH.DirectorCACert
	}

	iaasProperties := iaasConfig.Properties(turbulenceNetwork.Subnets[0].Static[0])
	turbulenceProperties := Properties{
		WardenCPI: iaasProperties.WardenCPI,
		AWS:       iaasProperties.AWS,
		Registry:  iaasProperties.Registry,
		Blobstore: iaasProperties.Blobstore,
		Agent:     iaasProperties.Agent,
		TurbulenceAPI: &PropertiesTurbulenceAPI{
			Certificate: APICertificate,
			CPIJobName:  cpi.JobName,
			Director: PropertiesTurbulenceAPIDirector{
				CACert:   directorCACert,
				Host:     config.BOSH.Target,
				Password: config.BOSH.Password,
				Username: config.BOSH.Username,
			},
			Password:   "turbulence-password",
			PrivateKey: APIPrivateKey,
		},
	}

	return Manifest{
		DirectorUUID:  config.DirectorUUID,
		Name:          config.Name,
		Releases:      []core.Release{turbulenceRelease, cpiRelease},
		ResourcePools: []core.ResourcePool{turbulenceResourcePool},
		Compilation:   compilation,
		Update:        update,
		Jobs:          []core.Job{apiJob},
		Networks:      []core.Network{turbulenceNetwork},
		Properties:    turbulenceProperties,
	}
}

func (m Manifest) ToYAML() ([]byte, error) {
	return candiedyaml.Marshal(m)
}

func FromYAML(yaml []byte) (Manifest, error) {
	var m Manifest
	if err := candiedyaml.Unmarshal(yaml, &m); err != nil {
		return m, err
	}
	return m, nil
}
