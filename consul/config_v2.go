package consul

type ConfigV2 struct {
	DirectorUUID string
	Name         string
	AZs          []ConfigAZ
}

type ConfigAZ struct {
	Name    string
	IPRange string
	Nodes   int
}
