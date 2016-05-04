package turbulence

type Config struct {
	DirectorUUID string
	Name         string
	IPRange      string
	IPOffset     int
	BOSH         ConfigBOSH
}

type ConfigBOSH struct {
	Target         string
	Username       string
	Password       string
	DirectorCACert string
}
