package core

type Compilation struct {
	CloudProperties     CompilationCloudProperties `yaml:"cloud_properties"`
	Network             string                     `yaml:"network"`
	ReuseCompilationVMs bool                       `yaml:"reuse_compilation_vms"`
	Workers             int                        `yaml:"workers"`
}

type CompilationCloudProperties struct {
	InstanceType     string                                   `yaml:"instance_type,omitempty"`
	AvailabilityZone string                                   `yaml:"availability_zone,omitempty"`
	EphemeralDisk    *CompilationCloudPropertiesEphemeralDisk `yaml:"ephemeral_disk,omitempty"`
}

type CompilationCloudPropertiesEphemeralDisk struct {
	Size int    `yaml:"size"`
	Type string `yaml:"type"`
}

type Release struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type ResourcePool struct {
	CloudProperties ResourcePoolCloudProperties `yaml:"cloud_properties"`
	Name            string                      `yaml:"name"`
	Network         string                      `yaml:"network"`
	Stemcell        ResourcePoolStemcell        `yaml:"stemcell"`
}

type ResourcePoolCloudProperties struct {
	InstanceType     string                                    `yaml:"instance_type,omitempty"`
	AvailabilityZone string                                    `yaml:"availability_zone,omitempty"`
	EphemeralDisk    *ResourcePoolCloudPropertiesEphemeralDisk `yaml:"ephemeral_disk,omitempty"`
}

type ResourcePoolCloudPropertiesEphemeralDisk struct {
	Size int    `yaml:"size"`
	Type string `yaml:"type"`
}

type ResourcePoolStemcell struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Update struct {
	Canaries        int    `yaml:"canaries"`
	CanaryWatchTime string `yaml:"canary_watch_time"`
	MaxInFlight     int    `yaml:"max_in_flight"`
	Serial          bool   `yaml:"serial"`
	UpdateWatchTime string `yaml:"update_watch_time"`
}

type PropertiesBlobstore struct {
	Address string                   `yaml:"address"`
	Port    int                      `yaml:"port"`
	Agent   PropertiesBlobstoreAgent `yaml:"agent"`
}

type PropertiesBlobstoreAgent struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type PropertiesAgent struct {
	Mbus string `yaml:"mbus"`
}

type PropertiesRegistry struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
}

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
