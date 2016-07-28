package consul

import "github.com/pivotal-cf-experimental/destiny/core"

func SetJobInstanceCount(job core.Job, network core.Network, properties Properties, count int) (core.Job, Properties, error) {
	job.Instances = count
	for i, net := range job.Networks {
		if net.Name == network.Name {
			var err error
			net.StaticIPs, err = network.StaticIPsFromRange(count)
			if err != nil {
				return job, properties, err
			}
			properties.Consul.Agent.Servers.Lan = net.StaticIPs
		}
		job.Networks[i] = net
	}

	return job, properties, nil
}
