package etcd

import (
	"errors"

	"github.com/pivotal-cf-experimental/destiny/core"
)

func SetJobInstanceCount(job core.Job, network core.Network, properties Properties, count int) (core.Job, Properties) {
	job.Instances = count
	for i, net := range job.Networks {
		if net.Name == network.Name {
			net.StaticIPs = network.StaticIPs(count)

			if !properties.Etcd.RequireSSL {
				properties.EtcdTestConsumer.Etcd.Machines = net.StaticIPs
				properties.Etcd.Machines = net.StaticIPs
			}

			properties.Etcd.Cluster[0].Instances = count
		}
		job.Networks[i] = net
	}

	return job, properties
}

func (m Manifest) RemoveJob(jobName string) Manifest {
	for i, job := range m.Jobs {
		if job.Name == jobName {
			m.Jobs = append(m.Jobs[:i], m.Jobs[i+1:]...)
		}
	}
	return m
}

func (m Manifest) ReplaceEtcdWithProxyJob(jobToReplace string) (Manifest, error) {
	m, err := m.replaceEtcdTemplateWithProxy(jobToReplace)
	if err != nil {
		return m, err
	}

	for _, job := range m.Jobs {
		if job.Properties != nil {
			if job.Properties.Etcd.RequireSSL {
				m.Properties.EtcdProxy = &PropertiesEtcdProxy{
					Etcd: PropertiesEtcdProxyEtcd{
						URL:        "https://etcd.service.cf.internal",
						CACert:     job.Properties.Etcd.CACert,
						ClientCert: job.Properties.Etcd.ClientCert,
						ClientKey:  job.Properties.Etcd.ClientKey,
						RequireSSL: true,
						Port:       4001,
					},
				}
			}
		}
	}

	return m, nil
}

func (m Manifest) replaceEtcdTemplateWithProxy(jobName string) (Manifest, error) {
	for _, job := range m.Jobs {
		if job.Name == jobName {
			for i, template := range job.Templates {
				if template.Name == "etcd" {
					job.Templates[i] = core.JobTemplate{
						Name:    "etcd_proxy",
						Release: "etcd",
					}
					return m, nil
				}
			}
		}
	}
	return m, errors.New("job not found")
}
