package etcd

import (
	"github.com/cloudfoundry-incubator/candiedyaml"
	"github.com/pivotal-cf-experimental/destiny/core"
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
