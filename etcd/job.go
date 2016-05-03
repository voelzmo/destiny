package etcd

import "github.com/pivotal-cf-experimental/destiny/core"

func SetJobInstanceCount(job core.Job, network core.Network, properties Properties, count int, ipOffset int) (core.Job, Properties) {
	job.Instances = count
	for i, net := range job.Networks {
		if net.Name == network.Name {
			net.StaticIPs = network.StaticIPs(count, ipOffset)

			if !properties.Etcd.RequireSSL {
				properties.Etcd.Machines = net.StaticIPs
			}

			properties.Etcd.Cluster[0].Instances = count
		}
		job.Networks[i] = net
	}

	return job, properties
}
