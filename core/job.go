package core

type Job struct {
	Instances      int            `yaml:"instances"`
	Lifecycle      string         `yaml:"lifecycle,omitempty"`
	Name           string         `yaml:"name"`
	Networks       []JobNetwork   `yaml:"networks"`
	ResourcePool   string         `yaml:"resource_pool"`
	Templates      []JobTemplate  `yaml:"templates"`
	PersistentDisk int            `yaml:"persistent_disk"`
	Properties     *JobProperties `yaml:"properties,omitempty"`
	Update         *JobUpdate     `yaml:"update,omitempty"`
}

type JobUpdate struct {
	MaxInFlight int `yaml:"max_in_flight"`
}

type JobNetwork struct {
	Name      string   `yaml:"name"`
	StaticIPs []string `yaml:"static_ips"`
}

type JobTemplate struct {
	Name    string `yaml:"name"`
	Release string `yaml:"release"`
}

type JobProperties struct {
	Consul JobPropertiesConsul `yaml:"consul"`
}

type JobPropertiesConsul struct {
	Agent JobPropertiesConsulAgent `yaml:"agent"`
}

type JobPropertiesConsulAgent struct {
	Mode     string                           `yaml:"mode,omitempty"`
	LogLevel string                           `yaml:"log_level,omitempty"`
	Services JobPropertiesConsulAgentServices `yaml:"services,omitempty"`
}

type JobPropertiesConsulAgentServices map[string]JobPropertiesConsulAgentService

type JobPropertiesConsulAgentService struct {
	Name  string                                `yaml:"name,omitempty"`
	Check *JobPropertiesConsulAgentServiceCheck `yaml:"check,omitempty"`
	Tags  []string                              `yaml:"tags,omitempty"`
}

type JobPropertiesConsulAgentServiceCheck struct {
	Name     string `yaml:"name"`
	Script   string `yaml:"script,omitempty"`
	Interval string `yaml:"interval,omitempty"`
}

func (j Job) HasTemplate(name, release string) bool {
	for _, template := range j.Templates {
		if template.Name == name && template.Release == release {
			return true
		}
	}
	return false
}
