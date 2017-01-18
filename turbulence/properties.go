package turbulence

import (
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/iaas"
)

type Properties struct {
	WardenCPI     *iaas.PropertiesWardenCPI     `yaml:"warden_cpi,omitempty"`
	AWS           *iaas.PropertiesAWS           `yaml:"aws,omitempty"`
	Registry      *core.PropertiesRegistry      `yaml:"registry,omitempty"`
	Blobstore     *core.PropertiesBlobstore     `yaml:"blobstore,omitempty"`
	Agent         *core.PropertiesAgent         `yaml:"agent,omitempty"`
	TurbulenceAPI *core.PropertiesTurbulenceAPI `yaml:"turbulence_api,omitempty"`
}
