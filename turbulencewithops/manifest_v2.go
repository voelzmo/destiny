package turbulencewithops

import "github.com/pivotal-cf-experimental/destiny/ops"

type Op struct {
	Type  string      `yaml:"type"`
	Path  string      `yaml:"path"`
	Value interface{} `yaml:"value"`
}

func NewManifestV2(config ConfigV2) (string, error) {
	return ops.ApplyOps(manifestV2, []ops.Op{
		{"replace", "/director_uuid", config.DirectorUUID},
		{"replace", "/name", config.Name},
		{"replace", "/instance_groups/name=api/azs", config.AZs},
		{"replace", "/instance_groups/name=api/properties/turbulence_api/director/host", config.DirectorHost},
		{"replace", "/instance_groups/name=api/properties/turbulence_api/director/username", config.DirectorUsername},
		{"replace", "/instance_groups/name=api/properties/turbulence_api/director/password", config.DirectorPassword},
	})
}
