package etcd

import (
	"github.com/pivotal-cf-experimental/destiny/bosh"
	yaml "gopkg.in/yaml.v2"
)

type Op struct {
	Type  string      `yaml:"type"`
	Path  string      `yaml:"path"`
	Value interface{} `yaml:"value"`
}

func NewManifestV2(config ConfigV2) (string, error) {
	ops := []Op{
		createOp("replace", "/director_uuid", config.DirectorUUID),
		createOp("replace", "/name", config.Name),
		createOp("replace", "/instance_groups/name=consul/azs", config.AZs),
		createOp("replace", "/instance_groups/name=etcd/azs", config.AZs),
		createOp("replace", "/instance_groups/name=testconsumer/azs", config.AZs),
	}

	manifestContents, err := applyOps(manifestV2, ops)
	if err != nil {
		panic(err)
	}

	return string(manifestContents), nil
}

func ApplyOp(manifest, opType, opPath string, opValue interface{}) (string, error) {
	ops := []Op{createOp(opType, opPath, opValue)}

	manifestContents, err := applyOps(manifest, ops)
	if err != nil {
		panic(err)
	}

	return manifestContents, nil
}

func createOp(opType, opPath string, opValue interface{}) Op {
	return Op{
		Type:  opType,
		Path:  opPath,
		Value: opValue,
	}
}

func applyOps(manifest string, ops []Op) (string, error) {
	opsYaml, err := yaml.Marshal(ops)
	if err != nil {
		panic(err)
	}

	boshCLI := bosh.BOSHCLI{}

	manifestContents, err := boshCLI.Interpolate(manifest, opsYaml)
	if err != nil {
		panic(err)
	}
	return string(manifestContents), nil
}
