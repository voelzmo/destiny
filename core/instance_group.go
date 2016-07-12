package core

type InstanceGroup struct {
	Instances          int                     `yaml:"instances"`
	Name               string                  `yaml:"name"`
	AZs                []string                `yaml:"azs"`
	Networks           []InstanceGroupNetwork  `yaml:"networks"`
	VMType             string                  `yaml:"vm_type"`
	Stemcell           string                  `yaml:"stemcell"`
	PersistentDiskType string                  `yaml:"persistent_disk_type"`
	Update             Update                  `yaml:"update,omitempty"`
	Jobs               []JobV2                 `yaml:"jobs"`
	Properties         InstanceGroupProperties `yaml:"properties,omitempty"`
}

type InstanceGroupProperties struct {
	Consul ConsulInstanceGroupProperties `yaml:"consul"`
}

type ConsulInstanceGroupProperties struct {
	Agent ConsulAgentProperties `yaml:"agent"`
}

type ConsulAgentProperties struct {
	Mode     string                                  `yaml:"mode"`
	LogLevel string                                  `yaml:"log_level"`
	Services map[string]ConsulAgentServiceProperties `yaml:"services"`
}

type ConsulAgentServiceProperties struct {
	Name  string                       `yaml:"name,omitempty"`
	Check ConsulServiceCheckProperties `yaml:"check,omitempty"`
	Tags  []string                     `yaml:"tags,omitempty"`
}

type ConsulServiceCheckProperties struct {
	Name     string `yaml:"name,omitempty"`
	Script   string `yaml:"script,omitempty"`
	Interval string `yaml:"interval,omitempty"`
}

type InstanceGroupNetwork struct {
	Name      string   `yaml:"name"`
	StaticIPs []string `yaml:"static_ips"`
}

type JobV2 struct {
	Name    string `yaml:"name"`
	Release string `yaml:"release"`
}
