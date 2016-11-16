package consul

import "github.com/pivotal-cf-experimental/destiny/core"

func releases() []core.Release {
	return []core.Release{
		{
			Name:    "consul",
			Version: "latest",
		},
	}
}

func update() core.Update {
	return core.Update{
		Canaries:        1,
		CanaryWatchTime: "1000-180000",
		MaxInFlight:     1,
		Serial:          true,
		UpdateWatchTime: "1000-180000",
	}
}
