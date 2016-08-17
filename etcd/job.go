package etcd

import "github.com/pivotal-cf-experimental/destiny/core"

func SetJobInstanceCount(job core.Job, network core.Network, count int, staticIPOffset int) (core.Job, error) {
	job.Instances = count
	for i, net := range job.Networks {
		if net.Name == network.Name {
			staticIps, err := network.StaticIPsFromRange(count + staticIPOffset)
			if err != nil {
				return core.Job{}, err
			}

			net.StaticIPs = staticIps[staticIPOffset:]
		}
		job.Networks[i] = net
	}

	return job, nil
}

func SetEtcdProperties(job core.Job, properties Properties) Properties {
	if !properties.Etcd.RequireSSL {
		properties.EtcdTestConsumer.Etcd.Machines = job.Networks[0].StaticIPs
		properties.Etcd.Machines = job.Networks[0].StaticIPs
	}

	properties.Etcd.Cluster[0].Instances = job.Instances

	return properties
}

func (m Manifest) RemoveJob(jobName string) Manifest {
	for i, job := range m.Jobs {
		if job.Name == jobName {
			m.Jobs = append(m.Jobs[:i], m.Jobs[i+1:]...)
		}
	}
	return m
}
