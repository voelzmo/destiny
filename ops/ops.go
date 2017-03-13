package ops

import (
	"github.com/pivotal-cf-experimental/destiny/bosh"
	yaml "gopkg.in/yaml.v2"
)

type Op struct {
	Type  string      `yaml:"type"`
	Path  string      `yaml:"path"`
	Value interface{} `yaml:"value"`
}

var marshal = yaml.Marshal

func ApplyOp(manifest string, op Op) (string, error) {
	return ApplyOps(manifest, []Op{op})
}

func ApplyOps(manifest string, ops []Op) (string, error) {
	opsYaml, err := marshal(ops)
	if err != nil {
		return "", err
	}

	boshCLI := bosh.CLI{}

	manifestContents, err := boshCLI.Interpolate(manifest, string(opsYaml))
	if err != nil {
		// not tested
		return "", err
	}

	return manifestContents, nil
}
