package core

type InstanceGroup struct {
	Instances          int                         `yaml:"instances"`
	Name               string                      `yaml:"name"`
	AZs                []string                    `yaml:"azs"`
	Networks           []InstanceGroupNetwork      `yaml:"networks"`
	VMType             string                      `yaml:"vm_type"`
	Stemcell           string                      `yaml:"stemcell"`
	PersistentDiskType string                      `yaml:"persistent_disk_type,omitempty"`
	Update             Update                      `yaml:"update,omitempty"`
	Jobs               []InstanceGroupJob          `yaml:"jobs"`
	MigratedFrom       []InstanceGroupMigratedFrom `yaml:"migrated_from,omitempty"`
	Properties         *JobProperties              `yaml:"properties,omitempty"`
}

type InstanceGroupMigratedFrom struct {
	Name string `yaml:"name"`
	AZ   string `yaml:"az"`
}

type InstanceGroupNetwork JobNetwork

type InstanceGroupJob struct {
	Name       string        `yaml:"name"`
	Release    string        `yaml:"release"`
	Properties APIProperties `yaml:"properties,omitempty"`
	Provides   JobProvides   `yaml:"provides,omitempty"`
	Consumes   JobConsumes   `yaml:"consumes,omitempty"`
}

type JobProvides struct {
	API Provider `yaml:"api"`
}

type Provider struct {
	Shared bool `yaml:"shared"`
}

type JobConsumes struct {
	API Consumer `yaml:"api"`
}

type Consumer struct {
	From       string `yaml:"from"`
	Deployment string `yaml:"deployment"`
}

type PropertiesTurbulenceAPI struct {
	Certificate string                          `yaml:"certificate,omitempty"`
	CPIJobName  string                          `yaml:"cpi_job_name"`
	Director    PropertiesTurbulenceAPIDirector `yaml:"director,omitempty"`
	Password    string                          `yaml:"password,omitempty"`
	PrivateKey  string                          `yaml:"private_key,omitempty"`
}

type PropertiesTurbulenceAPIDirector struct {
	CACert       string `yaml:"ca_cert"`
	Host         string `yaml:"host"`
	ClientSecret string `yaml:"client_secret"`
	Client       string `yaml:"client"`
}

type APIProperties struct {
	Cert     APIPropertiesCert               `yaml:"cert"`
	Password string                          `yaml:"password"`
	Director PropertiesTurbulenceAPIDirector `yaml:"director"`
}

type APIPropertiesCert struct {
	Certificate string `yaml:"certificate"`
	PrivateKey  string `yaml:"private_key"`
	CA          string `yaml:"ca"`
}
