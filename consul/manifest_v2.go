package consul

import "github.com/pivotal-cf-experimental/destiny/ops"

type Op struct {
	Type  string      `yaml:"type"`
	Path  string      `yaml:"path"`
	Value interface{} `yaml:"value"`
}

func NewManifestV2(config ConfigV2) (string, error) {
	return ops.ApplyOps(manifestV2, []ops.Op{
		{"replace", "/name", config.Name},
		{"replace", "/instance_groups/name=consul/azs", config.AZs},
		{"replace", "/instance_groups/name=testconsumer/azs", config.AZs},
	})
}

func NewManifestV2Windows(config ConfigV2) (string, error) {
	return ops.ApplyOps(manifestV2, []ops.Op{
		{"replace", "/name", config.Name},
		{"replace", "/instance_groups/name=consul/azs", config.AZs},
		{"replace", "/instance_groups/name=testconsumer/azs", config.AZs},
		{"replace", "/stemcells/-", map[string]string{
			"alias":   "windows",
			"os":      "windows2012R2",
			"version": "latest",
		}},
		{"replace", "/instance_groups/name=testconsumer/jobs/name=consul_agent/name", "consul_agent_windows"},
		{"replace", "/instance_groups/name=testconsumer/jobs/name=consul-test-consumer/name", "consul-test-consumer-windows"},
		{"replace", "/instance_groups/name=testconsumer/vm_extensions?", []string{"50GB_ephemeral_disk"}},
		{"replace", "/instance_groups/name=testconsumer/stemcell", "windows"},
	})
}
