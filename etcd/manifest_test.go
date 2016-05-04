package etcd_test

import (
	"io/ioutil"

	"github.com/pivotal-cf-experimental/destiny/consul"
	"github.com/pivotal-cf-experimental/destiny/core"
	"github.com/pivotal-cf-experimental/destiny/etcd"
	"github.com/pivotal-cf-experimental/destiny/iaas"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotal-cf-experimental/gomegamatchers"
)

var _ = Describe("Manifest", func() {
	Describe("NewTLSManifest", func() {
		It("generates a valid Etcd BOSH-Lite manifest", func() {
			manifest := etcd.NewTLSManifest(etcd.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "etcd-some-random-guid",
				IPRange:      "10.244.4.0/24",
				IPOffset:     10,
			}, iaas.NewWardenConfig())

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))

			Expect(manifest.Name).To(Equal("etcd-some-random-guid"))

			Expect(manifest.Releases).To(HaveLen(2))

			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "etcd",
				Version: "latest",
			}))

			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "consul",
				Version: "latest",
			}))

			Expect(manifest.Compilation).To(Equal(core.Compilation{
				Network:             "etcd1",
				ReuseCompilationVMs: true,
				Workers:             3,
			}))

			Expect(manifest.Update).To(Equal(core.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     1,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))

			Expect(manifest.ResourcePools).To(HaveLen(1))

			Expect(manifest.ResourcePools).To(ContainElement(core.ResourcePool{
				Name:    "etcd_z1",
				Network: "etcd1",
				Stemcell: core.ResourcePoolStemcell{
					Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
					Version: "latest",
				},
			}))

			Expect(manifest.Jobs).To(HaveLen(3))

			Expect(manifest.Jobs[0]).To(Equal(core.Job{
				Name:      "consul_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.25"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "consul_agent",
						Release: "consul",
					},
				},
				Properties: &core.JobProperties{
					Consul: core.JobPropertiesConsul{
						Agent: core.JobPropertiesConsulAgent{
							Mode: "server",
						},
					},
				},
			}))

			Expect(manifest.Jobs[1]).To(Equal(core.Job{
				Name:      "etcd_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.20"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "etcd",
						Release: "etcd",
					},
					{
						Name:    "consul_agent",
						Release: "consul",
					},
				},
				Properties: &core.JobProperties{
					Consul: core.JobPropertiesConsul{
						Agent: core.JobPropertiesConsulAgent{
							Services: core.JobPropertiesConsulAgentServices{
								"etcd": core.JobPropertiesConsulAgentService{},
							},
						},
					},
				},
			}))

			Expect(manifest.Jobs[2]).To(Equal(core.Job{
				Name:      "testconsumer_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.26"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "consul_agent",
						Release: "consul",
					},
					{
						Name:    "etcd-test-consumer",
						Release: "etcd",
					},
				},
			}))

			Expect(manifest.Networks).To(HaveLen(1))

			Expect(manifest.Networks).To(ContainElement(core.Network{
				Name: "etcd1",
				Subnets: []core.NetworkSubnet{
					{
						CloudProperties: core.NetworkSubnetCloudProperties{Name: "random"},
						Gateway:         "10.244.4.1",
						Range:           "10.244.4.0/24",
						Reserved: []string{
							"10.244.4.2-10.244.4.3",
							"10.244.4.255",
						},
						Static: []string{
							"10.244.4.10",
							"10.244.4.11",
							"10.244.4.12",
							"10.244.4.13",
							"10.244.4.14",
							"10.244.4.15",
							"10.244.4.16",
							"10.244.4.17",
							"10.244.4.18",
							"10.244.4.19",
							"10.244.4.20",
							"10.244.4.21",
							"10.244.4.22",
							"10.244.4.23",
							"10.244.4.24",
							"10.244.4.25",
							"10.244.4.26",
							"10.244.4.27",
							"10.244.4.28",
							"10.244.4.29",
							"10.244.4.30",
							"10.244.4.31",
							"10.244.4.32",
							"10.244.4.33",
							"10.244.4.34",
							"10.244.4.35",
							"10.244.4.36",
							"10.244.4.37",
							"10.244.4.38",
							"10.244.4.39",
							"10.244.4.40",
							"10.244.4.41",
							"10.244.4.42",
							"10.244.4.43",
							"10.244.4.44",
							"10.244.4.45",
							"10.244.4.46",
							"10.244.4.47",
							"10.244.4.48",
							"10.244.4.49",
							"10.244.4.50",
							"10.244.4.51",
							"10.244.4.52",
							"10.244.4.53",
							"10.244.4.54",
							"10.244.4.55",
							"10.244.4.56",
							"10.244.4.57",
							"10.244.4.58",
							"10.244.4.59",
							"10.244.4.60",
							"10.244.4.61",
							"10.244.4.62",
							"10.244.4.63",
							"10.244.4.64",
							"10.244.4.65",
							"10.244.4.66",
							"10.244.4.67",
							"10.244.4.68",
							"10.244.4.69",
							"10.244.4.70",
							"10.244.4.71",
							"10.244.4.72",
							"10.244.4.73",
							"10.244.4.74",
							"10.244.4.75",
							"10.244.4.76",
							"10.244.4.77",
							"10.244.4.78",
							"10.244.4.79",
							"10.244.4.80",
							"10.244.4.81",
							"10.244.4.82",
							"10.244.4.83",
							"10.244.4.84",
							"10.244.4.85",
							"10.244.4.86",
							"10.244.4.87",
							"10.244.4.88",
							"10.244.4.89",
							"10.244.4.90",
							"10.244.4.91",
							"10.244.4.92",
							"10.244.4.93",
							"10.244.4.94",
							"10.244.4.95",
							"10.244.4.96",
							"10.244.4.97",
							"10.244.4.98",
							"10.244.4.99",
							"10.244.4.100",
							"10.244.4.101",
							"10.244.4.102",
							"10.244.4.103",
							"10.244.4.104",
							"10.244.4.105",
							"10.244.4.106",
							"10.244.4.107",
							"10.244.4.108",
							"10.244.4.109",
							"10.244.4.110",
							"10.244.4.111",
							"10.244.4.112",
							"10.244.4.113",
							"10.244.4.114",
							"10.244.4.115",
							"10.244.4.116",
							"10.244.4.117",
							"10.244.4.118",
							"10.244.4.119",
							"10.244.4.120",
							"10.244.4.121",
							"10.244.4.122",
							"10.244.4.123",
							"10.244.4.124",
							"10.244.4.125",
							"10.244.4.126",
							"10.244.4.127",
							"10.244.4.128",
							"10.244.4.129",
							"10.244.4.130",
							"10.244.4.131",
							"10.244.4.132",
							"10.244.4.133",
							"10.244.4.134",
							"10.244.4.135",
							"10.244.4.136",
							"10.244.4.137",
							"10.244.4.138",
							"10.244.4.139",
							"10.244.4.140",
							"10.244.4.141",
							"10.244.4.142",
							"10.244.4.143",
							"10.244.4.144",
							"10.244.4.145",
							"10.244.4.146",
							"10.244.4.147",
							"10.244.4.148",
							"10.244.4.149",
							"10.244.4.150",
							"10.244.4.151",
							"10.244.4.152",
							"10.244.4.153",
							"10.244.4.154",
							"10.244.4.155",
							"10.244.4.156",
							"10.244.4.157",
							"10.244.4.158",
							"10.244.4.159",
							"10.244.4.160",
							"10.244.4.161",
							"10.244.4.162",
							"10.244.4.163",
							"10.244.4.164",
							"10.244.4.165",
							"10.244.4.166",
							"10.244.4.167",
							"10.244.4.168",
							"10.244.4.169",
							"10.244.4.170",
							"10.244.4.171",
							"10.244.4.172",
							"10.244.4.173",
							"10.244.4.174",
							"10.244.4.175",
							"10.244.4.176",
							"10.244.4.177",
							"10.244.4.178",
							"10.244.4.179",
							"10.244.4.180",
							"10.244.4.181",
							"10.244.4.182",
							"10.244.4.183",
							"10.244.4.184",
							"10.244.4.185",
							"10.244.4.186",
							"10.244.4.187",
							"10.244.4.188",
							"10.244.4.189",
							"10.244.4.190",
							"10.244.4.191",
							"10.244.4.192",
							"10.244.4.193",
							"10.244.4.194",
							"10.244.4.195",
							"10.244.4.196",
							"10.244.4.197",
							"10.244.4.198",
							"10.244.4.199",
							"10.244.4.200",
							"10.244.4.201",
							"10.244.4.202",
							"10.244.4.203",
							"10.244.4.204",
							"10.244.4.205",
							"10.244.4.206",
							"10.244.4.207",
							"10.244.4.208",
							"10.244.4.209",
							"10.244.4.210",
							"10.244.4.211",
							"10.244.4.212",
							"10.244.4.213",
							"10.244.4.214",
							"10.244.4.215",
							"10.244.4.216",
							"10.244.4.217",
							"10.244.4.218",
							"10.244.4.219",
							"10.244.4.220",
							"10.244.4.221",
							"10.244.4.222",
							"10.244.4.223",
							"10.244.4.224",
							"10.244.4.225",
							"10.244.4.226",
							"10.244.4.227",
							"10.244.4.228",
							"10.244.4.229",
							"10.244.4.230",
							"10.244.4.231",
							"10.244.4.232",
							"10.244.4.233",
							"10.244.4.234",
							"10.244.4.235",
							"10.244.4.236",
							"10.244.4.237",
							"10.244.4.238",
							"10.244.4.239",
							"10.244.4.240",
							"10.244.4.241",
							"10.244.4.242",
							"10.244.4.243",
							"10.244.4.244",
							"10.244.4.245",
							"10.244.4.246",
							"10.244.4.247",
							"10.244.4.248",
							"10.244.4.249",
							"10.244.4.250",
							"10.244.4.251",
							"10.244.4.252",
							"10.244.4.253",
							"10.244.4.254",
						},
					},
				},
				Type: "manual",
			}))

			Expect(manifest.Properties.Etcd).To(Equal(&etcd.PropertiesEtcd{
				Cluster: []etcd.PropertiesEtcdCluster{{
					Instances: 1,
					Name:      "etcd_z1",
				}},
				Machines: []string{
					"etcd.service.cf.internal",
				},
				PeerRequireSSL:                  true,
				RequireSSL:                      true,
				HeartbeatIntervalInMilliseconds: 50,
				AdvertiseURLsDNSSuffix:          "etcd.service.cf.internal",
				CACert:                          etcd.CACert,
				ClientCert:                      etcd.ClientCert,
				ClientKey:                       etcd.ClientKey,
				PeerCACert:                      etcd.PeerCACert,
				PeerCert:                        etcd.PeerCert,
				PeerKey:                         etcd.PeerKey,
				ServerCert:                      etcd.ServerCert,
				ServerKey:                       etcd.ServerKey,
			}))

			Expect(manifest.Properties.Consul).To(Equal(&consul.PropertiesConsul{
				Agent: consul.PropertiesConsulAgent{
					Domain: "cf.internal",
					Servers: consul.PropertiesConsulAgentServers{
						Lan: []string{"10.244.4.25"},
					},
				},
				CACert:      consul.CACert,
				AgentCert:   consul.AgentCert,
				AgentKey:    consul.AgentKey,
				ServerCert:  consul.ServerCert,
				ServerKey:   consul.ServerKey,
				EncryptKeys: []string{consul.EncryptKey},
			}))
		})

		It("generates a valid Etcd AWS manifest", func() {
			manifest := etcd.NewTLSManifest(etcd.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "etcd-some-random-guid",
				IPRange:      "10.0.16.0/24",
			}, iaas.AWSConfig{
				Subnet: "subnet-1234",
			})

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))
			Expect(manifest.Name).To(Equal("etcd-some-random-guid"))

			Expect(manifest.Releases).To(HaveLen(2))
			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "etcd",
				Version: "latest",
			}))
			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "consul",
				Version: "latest",
			}))

			Expect(manifest.Compilation).To(Equal(core.Compilation{
				Network:             "etcd1",
				ReuseCompilationVMs: true,
				Workers:             3,
				CloudProperties: core.CompilationCloudProperties{
					InstanceType:     "m3.medium",
					AvailabilityZone: "us-east-1a",
					EphemeralDisk: &core.CompilationCloudPropertiesEphemeralDisk{
						Size: 1024,
						Type: "gp2",
					},
				},
			}))

			Expect(manifest.Update).To(Equal(core.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     1,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))

			Expect(manifest.ResourcePools).To(HaveLen(1))
			Expect(manifest.ResourcePools).To(ContainElement(core.ResourcePool{
				Name:    "etcd_z1",
				Network: "etcd1",
				Stemcell: core.ResourcePoolStemcell{
					Name:    "bosh-aws-xen-hvm-ubuntu-trusty-go_agent",
					Version: "latest",
				},
				CloudProperties: core.ResourcePoolCloudProperties{
					InstanceType:     "m3.medium",
					AvailabilityZone: "us-east-1a",
					EphemeralDisk: &core.ResourcePoolCloudPropertiesEphemeralDisk{
						Size: 1024,
						Type: "gp2",
					},
				},
			}))

			Expect(manifest.Jobs).To(HaveLen(3))

			Expect(manifest.Jobs[0]).To(Equal(core.Job{
				Name:      "consul_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.0.16.15"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{{
					Name:    "consul_agent",
					Release: "consul",
				}},
				Properties: &core.JobProperties{
					Consul: core.JobPropertiesConsul{
						Agent: core.JobPropertiesConsulAgent{
							Mode: "server",
						},
					},
				},
			}))

			Expect(manifest.Jobs[1]).To(Equal(core.Job{
				Name:      "etcd_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.0.16.10"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "etcd",
						Release: "etcd",
					},
					{
						Name:    "consul_agent",
						Release: "consul",
					},
				},
				Properties: &core.JobProperties{
					Consul: core.JobPropertiesConsul{
						Agent: core.JobPropertiesConsulAgent{
							Services: core.JobPropertiesConsulAgentServices{
								"etcd": core.JobPropertiesConsulAgentService{},
							},
						},
					},
				},
			}))

			Expect(manifest.Jobs[2]).To(Equal(core.Job{
				Name:      "testconsumer_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.0.16.16"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "consul_agent",
						Release: "consul",
					},
					{
						Name:    "etcd-test-consumer",
						Release: "etcd",
					},
				},
			}))

			Expect(manifest.Networks).To(HaveLen(1))
			Expect(manifest.Networks).To(ContainElement(core.Network{
				Name: "etcd1",
				Subnets: []core.NetworkSubnet{
					{
						CloudProperties: core.NetworkSubnetCloudProperties{Subnet: "subnet-1234"},
						Gateway:         "10.0.16.1",
						Range:           "10.0.16.0/24",
						Reserved: []string{
							"10.0.16.2-10.0.16.3",
							"10.0.16.255",
						},
						Static: []string{
							"10.0.16.10",
							"10.0.16.11",
							"10.0.16.12",
							"10.0.16.13",
							"10.0.16.14",
							"10.0.16.15",
							"10.0.16.16",
							"10.0.16.17",
							"10.0.16.18",
							"10.0.16.19",
							"10.0.16.20",
							"10.0.16.21",
							"10.0.16.22",
							"10.0.16.23",
							"10.0.16.24",
							"10.0.16.25",
							"10.0.16.26",
							"10.0.16.27",
							"10.0.16.28",
							"10.0.16.29",
							"10.0.16.30",
							"10.0.16.31",
							"10.0.16.32",
							"10.0.16.33",
							"10.0.16.34",
							"10.0.16.35",
							"10.0.16.36",
							"10.0.16.37",
							"10.0.16.38",
							"10.0.16.39",
							"10.0.16.40",
							"10.0.16.41",
							"10.0.16.42",
							"10.0.16.43",
							"10.0.16.44",
							"10.0.16.45",
							"10.0.16.46",
							"10.0.16.47",
							"10.0.16.48",
							"10.0.16.49",
							"10.0.16.50",
							"10.0.16.51",
							"10.0.16.52",
							"10.0.16.53",
							"10.0.16.54",
							"10.0.16.55",
							"10.0.16.56",
							"10.0.16.57",
							"10.0.16.58",
							"10.0.16.59",
							"10.0.16.60",
							"10.0.16.61",
							"10.0.16.62",
							"10.0.16.63",
							"10.0.16.64",
							"10.0.16.65",
							"10.0.16.66",
							"10.0.16.67",
							"10.0.16.68",
							"10.0.16.69",
							"10.0.16.70",
							"10.0.16.71",
							"10.0.16.72",
							"10.0.16.73",
							"10.0.16.74",
							"10.0.16.75",
							"10.0.16.76",
							"10.0.16.77",
							"10.0.16.78",
							"10.0.16.79",
							"10.0.16.80",
							"10.0.16.81",
							"10.0.16.82",
							"10.0.16.83",
							"10.0.16.84",
							"10.0.16.85",
							"10.0.16.86",
							"10.0.16.87",
							"10.0.16.88",
							"10.0.16.89",
							"10.0.16.90",
							"10.0.16.91",
							"10.0.16.92",
							"10.0.16.93",
							"10.0.16.94",
							"10.0.16.95",
							"10.0.16.96",
							"10.0.16.97",
							"10.0.16.98",
							"10.0.16.99",
							"10.0.16.100",
							"10.0.16.101",
							"10.0.16.102",
							"10.0.16.103",
							"10.0.16.104",
							"10.0.16.105",
							"10.0.16.106",
							"10.0.16.107",
							"10.0.16.108",
							"10.0.16.109",
							"10.0.16.110",
							"10.0.16.111",
							"10.0.16.112",
							"10.0.16.113",
							"10.0.16.114",
							"10.0.16.115",
							"10.0.16.116",
							"10.0.16.117",
							"10.0.16.118",
							"10.0.16.119",
							"10.0.16.120",
							"10.0.16.121",
							"10.0.16.122",
							"10.0.16.123",
							"10.0.16.124",
							"10.0.16.125",
							"10.0.16.126",
							"10.0.16.127",
							"10.0.16.128",
							"10.0.16.129",
							"10.0.16.130",
							"10.0.16.131",
							"10.0.16.132",
							"10.0.16.133",
							"10.0.16.134",
							"10.0.16.135",
							"10.0.16.136",
							"10.0.16.137",
							"10.0.16.138",
							"10.0.16.139",
							"10.0.16.140",
							"10.0.16.141",
							"10.0.16.142",
							"10.0.16.143",
							"10.0.16.144",
							"10.0.16.145",
							"10.0.16.146",
							"10.0.16.147",
							"10.0.16.148",
							"10.0.16.149",
							"10.0.16.150",
							"10.0.16.151",
							"10.0.16.152",
							"10.0.16.153",
							"10.0.16.154",
							"10.0.16.155",
							"10.0.16.156",
							"10.0.16.157",
							"10.0.16.158",
							"10.0.16.159",
							"10.0.16.160",
							"10.0.16.161",
							"10.0.16.162",
							"10.0.16.163",
							"10.0.16.164",
							"10.0.16.165",
							"10.0.16.166",
							"10.0.16.167",
							"10.0.16.168",
							"10.0.16.169",
							"10.0.16.170",
							"10.0.16.171",
							"10.0.16.172",
							"10.0.16.173",
							"10.0.16.174",
							"10.0.16.175",
							"10.0.16.176",
							"10.0.16.177",
							"10.0.16.178",
							"10.0.16.179",
							"10.0.16.180",
							"10.0.16.181",
							"10.0.16.182",
							"10.0.16.183",
							"10.0.16.184",
							"10.0.16.185",
							"10.0.16.186",
							"10.0.16.187",
							"10.0.16.188",
							"10.0.16.189",
							"10.0.16.190",
							"10.0.16.191",
							"10.0.16.192",
							"10.0.16.193",
							"10.0.16.194",
							"10.0.16.195",
							"10.0.16.196",
							"10.0.16.197",
							"10.0.16.198",
							"10.0.16.199",
							"10.0.16.200",
							"10.0.16.201",
							"10.0.16.202",
							"10.0.16.203",
							"10.0.16.204",
							"10.0.16.205",
							"10.0.16.206",
							"10.0.16.207",
							"10.0.16.208",
							"10.0.16.209",
							"10.0.16.210",
							"10.0.16.211",
							"10.0.16.212",
							"10.0.16.213",
							"10.0.16.214",
							"10.0.16.215",
							"10.0.16.216",
							"10.0.16.217",
							"10.0.16.218",
							"10.0.16.219",
							"10.0.16.220",
							"10.0.16.221",
							"10.0.16.222",
							"10.0.16.223",
							"10.0.16.224",
							"10.0.16.225",
							"10.0.16.226",
							"10.0.16.227",
							"10.0.16.228",
							"10.0.16.229",
							"10.0.16.230",
							"10.0.16.231",
							"10.0.16.232",
							"10.0.16.233",
							"10.0.16.234",
							"10.0.16.235",
							"10.0.16.236",
							"10.0.16.237",
							"10.0.16.238",
							"10.0.16.239",
							"10.0.16.240",
							"10.0.16.241",
							"10.0.16.242",
							"10.0.16.243",
							"10.0.16.244",
							"10.0.16.245",
							"10.0.16.246",
							"10.0.16.247",
							"10.0.16.248",
							"10.0.16.249",
							"10.0.16.250",
							"10.0.16.251",
							"10.0.16.252",
							"10.0.16.253",
							"10.0.16.254",
						},
					},
				},
				Type: "manual",
			}))

			Expect(manifest.Properties.Etcd).To(Equal(&etcd.PropertiesEtcd{
				Cluster: []etcd.PropertiesEtcdCluster{{
					Instances: 1,
					Name:      "etcd_z1",
				}},
				Machines: []string{
					"etcd.service.cf.internal",
				},
				PeerRequireSSL:                  true,
				RequireSSL:                      true,
				HeartbeatIntervalInMilliseconds: 50,
				AdvertiseURLsDNSSuffix:          "etcd.service.cf.internal",
				CACert:                          etcd.CACert,
				ClientCert:                      etcd.ClientCert,
				ClientKey:                       etcd.ClientKey,
				PeerCACert:                      etcd.PeerCACert,
				PeerCert:                        etcd.PeerCert,
				PeerKey:                         etcd.PeerKey,
				ServerCert:                      etcd.ServerCert,
				ServerKey:                       etcd.ServerKey,
			}))

			Expect(manifest.Properties.Consul).To(Equal(&consul.PropertiesConsul{
				Agent: consul.PropertiesConsulAgent{
					Domain: "cf.internal",
					Servers: consul.PropertiesConsulAgentServers{
						Lan: []string{"10.0.16.15"},
					},
				},
				CACert:      consul.CACert,
				AgentCert:   consul.AgentCert,
				AgentKey:    consul.AgentKey,
				ServerCert:  consul.ServerCert,
				ServerKey:   consul.ServerKey,
				EncryptKeys: []string{consul.EncryptKey},
			}))
		})
	})

	Describe("NewManifest", func() {
		It("generates a valid Etcd BOSH-Lite manifest", func() {
			manifest := etcd.NewManifest(etcd.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "etcd-some-random-guid",
				IPRange:      "10.244.4.0/24",
				IPOffset:     5,
			}, iaas.NewWardenConfig())

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))

			Expect(manifest.Name).To(Equal("etcd-some-random-guid"))

			Expect(manifest.Releases).To(HaveLen(1))

			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "etcd",
				Version: "latest",
			}))

			Expect(manifest.Compilation).To(Equal(core.Compilation{
				Network:             "etcd1",
				ReuseCompilationVMs: true,
				Workers:             3,
			}))

			Expect(manifest.Update).To(Equal(core.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     1,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))

			Expect(manifest.ResourcePools).To(HaveLen(1))

			Expect(manifest.ResourcePools).To(ContainElement(core.ResourcePool{
				Name:    "etcd_z1",
				Network: "etcd1",
				Stemcell: core.ResourcePoolStemcell{
					Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
					Version: "latest",
				},
			}))

			Expect(manifest.Jobs).To(HaveLen(3))

			Expect(manifest.Jobs[0]).To(Equal(core.Job{
				Name:      "consul_z1",
				Instances: 0,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
			}))

			Expect(manifest.Jobs[1]).To(Equal(core.Job{
				Name:      "etcd_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.15"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "etcd",
						Release: "etcd",
					},
				},
			}))

			Expect(manifest.Jobs[2]).To(Equal(core.Job{
				Name:      "testconsumer_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.21"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "etcd-test-consumer",
						Release: "etcd",
					},
				},
			}))

			Expect(manifest.Networks).To(HaveLen(1))

			Expect(manifest.Networks).To(ContainElement(core.Network{
				Name: "etcd1",
				Subnets: []core.NetworkSubnet{
					{
						CloudProperties: core.NetworkSubnetCloudProperties{Name: "random"},
						Gateway:         "10.244.4.1",
						Range:           "10.244.4.0/24",
						Reserved: []string{
							"10.244.4.2-10.244.4.3",
							"10.244.4.255",
						},
						Static: []string{
							"10.244.4.10",
							"10.244.4.11",
							"10.244.4.12",
							"10.244.4.13",
							"10.244.4.14",
							"10.244.4.15",
							"10.244.4.16",
							"10.244.4.17",
							"10.244.4.18",
							"10.244.4.19",
							"10.244.4.20",
							"10.244.4.21",
							"10.244.4.22",
							"10.244.4.23",
							"10.244.4.24",
							"10.244.4.25",
							"10.244.4.26",
							"10.244.4.27",
							"10.244.4.28",
							"10.244.4.29",
							"10.244.4.30",
							"10.244.4.31",
							"10.244.4.32",
							"10.244.4.33",
							"10.244.4.34",
							"10.244.4.35",
							"10.244.4.36",
							"10.244.4.37",
							"10.244.4.38",
							"10.244.4.39",
							"10.244.4.40",
							"10.244.4.41",
							"10.244.4.42",
							"10.244.4.43",
							"10.244.4.44",
							"10.244.4.45",
							"10.244.4.46",
							"10.244.4.47",
							"10.244.4.48",
							"10.244.4.49",
							"10.244.4.50",
							"10.244.4.51",
							"10.244.4.52",
							"10.244.4.53",
							"10.244.4.54",
							"10.244.4.55",
							"10.244.4.56",
							"10.244.4.57",
							"10.244.4.58",
							"10.244.4.59",
							"10.244.4.60",
							"10.244.4.61",
							"10.244.4.62",
							"10.244.4.63",
							"10.244.4.64",
							"10.244.4.65",
							"10.244.4.66",
							"10.244.4.67",
							"10.244.4.68",
							"10.244.4.69",
							"10.244.4.70",
							"10.244.4.71",
							"10.244.4.72",
							"10.244.4.73",
							"10.244.4.74",
							"10.244.4.75",
							"10.244.4.76",
							"10.244.4.77",
							"10.244.4.78",
							"10.244.4.79",
							"10.244.4.80",
							"10.244.4.81",
							"10.244.4.82",
							"10.244.4.83",
							"10.244.4.84",
							"10.244.4.85",
							"10.244.4.86",
							"10.244.4.87",
							"10.244.4.88",
							"10.244.4.89",
							"10.244.4.90",
							"10.244.4.91",
							"10.244.4.92",
							"10.244.4.93",
							"10.244.4.94",
							"10.244.4.95",
							"10.244.4.96",
							"10.244.4.97",
							"10.244.4.98",
							"10.244.4.99",
							"10.244.4.100",
							"10.244.4.101",
							"10.244.4.102",
							"10.244.4.103",
							"10.244.4.104",
							"10.244.4.105",
							"10.244.4.106",
							"10.244.4.107",
							"10.244.4.108",
							"10.244.4.109",
							"10.244.4.110",
							"10.244.4.111",
							"10.244.4.112",
							"10.244.4.113",
							"10.244.4.114",
							"10.244.4.115",
							"10.244.4.116",
							"10.244.4.117",
							"10.244.4.118",
							"10.244.4.119",
							"10.244.4.120",
							"10.244.4.121",
							"10.244.4.122",
							"10.244.4.123",
							"10.244.4.124",
							"10.244.4.125",
							"10.244.4.126",
							"10.244.4.127",
							"10.244.4.128",
							"10.244.4.129",
							"10.244.4.130",
							"10.244.4.131",
							"10.244.4.132",
							"10.244.4.133",
							"10.244.4.134",
							"10.244.4.135",
							"10.244.4.136",
							"10.244.4.137",
							"10.244.4.138",
							"10.244.4.139",
							"10.244.4.140",
							"10.244.4.141",
							"10.244.4.142",
							"10.244.4.143",
							"10.244.4.144",
							"10.244.4.145",
							"10.244.4.146",
							"10.244.4.147",
							"10.244.4.148",
							"10.244.4.149",
							"10.244.4.150",
							"10.244.4.151",
							"10.244.4.152",
							"10.244.4.153",
							"10.244.4.154",
							"10.244.4.155",
							"10.244.4.156",
							"10.244.4.157",
							"10.244.4.158",
							"10.244.4.159",
							"10.244.4.160",
							"10.244.4.161",
							"10.244.4.162",
							"10.244.4.163",
							"10.244.4.164",
							"10.244.4.165",
							"10.244.4.166",
							"10.244.4.167",
							"10.244.4.168",
							"10.244.4.169",
							"10.244.4.170",
							"10.244.4.171",
							"10.244.4.172",
							"10.244.4.173",
							"10.244.4.174",
							"10.244.4.175",
							"10.244.4.176",
							"10.244.4.177",
							"10.244.4.178",
							"10.244.4.179",
							"10.244.4.180",
							"10.244.4.181",
							"10.244.4.182",
							"10.244.4.183",
							"10.244.4.184",
							"10.244.4.185",
							"10.244.4.186",
							"10.244.4.187",
							"10.244.4.188",
							"10.244.4.189",
							"10.244.4.190",
							"10.244.4.191",
							"10.244.4.192",
							"10.244.4.193",
							"10.244.4.194",
							"10.244.4.195",
							"10.244.4.196",
							"10.244.4.197",
							"10.244.4.198",
							"10.244.4.199",
							"10.244.4.200",
							"10.244.4.201",
							"10.244.4.202",
							"10.244.4.203",
							"10.244.4.204",
							"10.244.4.205",
							"10.244.4.206",
							"10.244.4.207",
							"10.244.4.208",
							"10.244.4.209",
							"10.244.4.210",
							"10.244.4.211",
							"10.244.4.212",
							"10.244.4.213",
							"10.244.4.214",
							"10.244.4.215",
							"10.244.4.216",
							"10.244.4.217",
							"10.244.4.218",
							"10.244.4.219",
							"10.244.4.220",
							"10.244.4.221",
							"10.244.4.222",
							"10.244.4.223",
							"10.244.4.224",
							"10.244.4.225",
							"10.244.4.226",
							"10.244.4.227",
							"10.244.4.228",
							"10.244.4.229",
							"10.244.4.230",
							"10.244.4.231",
							"10.244.4.232",
							"10.244.4.233",
							"10.244.4.234",
							"10.244.4.235",
							"10.244.4.236",
							"10.244.4.237",
							"10.244.4.238",
							"10.244.4.239",
							"10.244.4.240",
							"10.244.4.241",
							"10.244.4.242",
							"10.244.4.243",
							"10.244.4.244",
							"10.244.4.245",
							"10.244.4.246",
							"10.244.4.247",
							"10.244.4.248",
							"10.244.4.249",
							"10.244.4.250",
							"10.244.4.251",
							"10.244.4.252",
							"10.244.4.253",
							"10.244.4.254",
						},
					},
				},
				Type: "manual",
			}))

			Expect(manifest.Properties.Etcd).To(Equal(&etcd.PropertiesEtcd{
				Cluster: []etcd.PropertiesEtcdCluster{{
					Instances: 1,
					Name:      "etcd_z1",
				}},
				Machines: []string{
					"10.244.4.15",
				},
				PeerRequireSSL:                  false,
				RequireSSL:                      false,
				HeartbeatIntervalInMilliseconds: 50,
			}))
		})

		It("generates a valid Etcd AWS manifest", func() {
			manifest := etcd.NewManifest(etcd.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "etcd-some-random-guid",
				IPRange:      "10.0.16.0/24",
			}, iaas.AWSConfig{
				Subnet: "subnet-1234",
			})

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))
			Expect(manifest.Name).To(Equal("etcd-some-random-guid"))

			Expect(manifest.Releases).To(HaveLen(1))
			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "etcd",
				Version: "latest",
			}))

			Expect(manifest.Compilation).To(Equal(core.Compilation{
				Network:             "etcd1",
				ReuseCompilationVMs: true,
				Workers:             3,
				CloudProperties: core.CompilationCloudProperties{
					InstanceType:     "m3.medium",
					AvailabilityZone: "us-east-1a",
					EphemeralDisk: &core.CompilationCloudPropertiesEphemeralDisk{
						Size: 1024,
						Type: "gp2",
					},
				},
			}))

			Expect(manifest.Update).To(Equal(core.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     1,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))

			Expect(manifest.ResourcePools).To(HaveLen(1))
			Expect(manifest.ResourcePools).To(ContainElement(core.ResourcePool{
				Name:    "etcd_z1",
				Network: "etcd1",
				Stemcell: core.ResourcePoolStemcell{
					Name:    "bosh-aws-xen-hvm-ubuntu-trusty-go_agent",
					Version: "latest",
				},
				CloudProperties: core.ResourcePoolCloudProperties{
					InstanceType:     "m3.medium",
					AvailabilityZone: "us-east-1a",
					EphemeralDisk: &core.ResourcePoolCloudPropertiesEphemeralDisk{
						Size: 1024,
						Type: "gp2",
					},
				},
			}))

			Expect(manifest.Jobs).To(HaveLen(3))

			Expect(manifest.Jobs[0]).To(Equal(core.Job{
				Name:      "consul_z1",
				Instances: 0,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
			}))

			Expect(manifest.Jobs[1]).To(Equal(core.Job{
				Name:      "etcd_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.0.16.10"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "etcd",
						Release: "etcd",
					},
				},
			}))

			Expect(manifest.Jobs[2]).To(Equal(core.Job{
				Name:      "testconsumer_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.0.16.16"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "etcd-test-consumer",
						Release: "etcd",
					},
				},
			}))

			Expect(manifest.Networks).To(HaveLen(1))
			Expect(manifest.Networks).To(ContainElement(core.Network{
				Name: "etcd1",
				Subnets: []core.NetworkSubnet{
					{
						CloudProperties: core.NetworkSubnetCloudProperties{Subnet: "subnet-1234"},
						Gateway:         "10.0.16.1",
						Range:           "10.0.16.0/24",
						Reserved: []string{
							"10.0.16.2-10.0.16.3",
							"10.0.16.255",
						},
						Static: []string{
							"10.0.16.10",
							"10.0.16.11",
							"10.0.16.12",
							"10.0.16.13",
							"10.0.16.14",
							"10.0.16.15",
							"10.0.16.16",
							"10.0.16.17",
							"10.0.16.18",
							"10.0.16.19",
							"10.0.16.20",
							"10.0.16.21",
							"10.0.16.22",
							"10.0.16.23",
							"10.0.16.24",
							"10.0.16.25",
							"10.0.16.26",
							"10.0.16.27",
							"10.0.16.28",
							"10.0.16.29",
							"10.0.16.30",
							"10.0.16.31",
							"10.0.16.32",
							"10.0.16.33",
							"10.0.16.34",
							"10.0.16.35",
							"10.0.16.36",
							"10.0.16.37",
							"10.0.16.38",
							"10.0.16.39",
							"10.0.16.40",
							"10.0.16.41",
							"10.0.16.42",
							"10.0.16.43",
							"10.0.16.44",
							"10.0.16.45",
							"10.0.16.46",
							"10.0.16.47",
							"10.0.16.48",
							"10.0.16.49",
							"10.0.16.50",
							"10.0.16.51",
							"10.0.16.52",
							"10.0.16.53",
							"10.0.16.54",
							"10.0.16.55",
							"10.0.16.56",
							"10.0.16.57",
							"10.0.16.58",
							"10.0.16.59",
							"10.0.16.60",
							"10.0.16.61",
							"10.0.16.62",
							"10.0.16.63",
							"10.0.16.64",
							"10.0.16.65",
							"10.0.16.66",
							"10.0.16.67",
							"10.0.16.68",
							"10.0.16.69",
							"10.0.16.70",
							"10.0.16.71",
							"10.0.16.72",
							"10.0.16.73",
							"10.0.16.74",
							"10.0.16.75",
							"10.0.16.76",
							"10.0.16.77",
							"10.0.16.78",
							"10.0.16.79",
							"10.0.16.80",
							"10.0.16.81",
							"10.0.16.82",
							"10.0.16.83",
							"10.0.16.84",
							"10.0.16.85",
							"10.0.16.86",
							"10.0.16.87",
							"10.0.16.88",
							"10.0.16.89",
							"10.0.16.90",
							"10.0.16.91",
							"10.0.16.92",
							"10.0.16.93",
							"10.0.16.94",
							"10.0.16.95",
							"10.0.16.96",
							"10.0.16.97",
							"10.0.16.98",
							"10.0.16.99",
							"10.0.16.100",
							"10.0.16.101",
							"10.0.16.102",
							"10.0.16.103",
							"10.0.16.104",
							"10.0.16.105",
							"10.0.16.106",
							"10.0.16.107",
							"10.0.16.108",
							"10.0.16.109",
							"10.0.16.110",
							"10.0.16.111",
							"10.0.16.112",
							"10.0.16.113",
							"10.0.16.114",
							"10.0.16.115",
							"10.0.16.116",
							"10.0.16.117",
							"10.0.16.118",
							"10.0.16.119",
							"10.0.16.120",
							"10.0.16.121",
							"10.0.16.122",
							"10.0.16.123",
							"10.0.16.124",
							"10.0.16.125",
							"10.0.16.126",
							"10.0.16.127",
							"10.0.16.128",
							"10.0.16.129",
							"10.0.16.130",
							"10.0.16.131",
							"10.0.16.132",
							"10.0.16.133",
							"10.0.16.134",
							"10.0.16.135",
							"10.0.16.136",
							"10.0.16.137",
							"10.0.16.138",
							"10.0.16.139",
							"10.0.16.140",
							"10.0.16.141",
							"10.0.16.142",
							"10.0.16.143",
							"10.0.16.144",
							"10.0.16.145",
							"10.0.16.146",
							"10.0.16.147",
							"10.0.16.148",
							"10.0.16.149",
							"10.0.16.150",
							"10.0.16.151",
							"10.0.16.152",
							"10.0.16.153",
							"10.0.16.154",
							"10.0.16.155",
							"10.0.16.156",
							"10.0.16.157",
							"10.0.16.158",
							"10.0.16.159",
							"10.0.16.160",
							"10.0.16.161",
							"10.0.16.162",
							"10.0.16.163",
							"10.0.16.164",
							"10.0.16.165",
							"10.0.16.166",
							"10.0.16.167",
							"10.0.16.168",
							"10.0.16.169",
							"10.0.16.170",
							"10.0.16.171",
							"10.0.16.172",
							"10.0.16.173",
							"10.0.16.174",
							"10.0.16.175",
							"10.0.16.176",
							"10.0.16.177",
							"10.0.16.178",
							"10.0.16.179",
							"10.0.16.180",
							"10.0.16.181",
							"10.0.16.182",
							"10.0.16.183",
							"10.0.16.184",
							"10.0.16.185",
							"10.0.16.186",
							"10.0.16.187",
							"10.0.16.188",
							"10.0.16.189",
							"10.0.16.190",
							"10.0.16.191",
							"10.0.16.192",
							"10.0.16.193",
							"10.0.16.194",
							"10.0.16.195",
							"10.0.16.196",
							"10.0.16.197",
							"10.0.16.198",
							"10.0.16.199",
							"10.0.16.200",
							"10.0.16.201",
							"10.0.16.202",
							"10.0.16.203",
							"10.0.16.204",
							"10.0.16.205",
							"10.0.16.206",
							"10.0.16.207",
							"10.0.16.208",
							"10.0.16.209",
							"10.0.16.210",
							"10.0.16.211",
							"10.0.16.212",
							"10.0.16.213",
							"10.0.16.214",
							"10.0.16.215",
							"10.0.16.216",
							"10.0.16.217",
							"10.0.16.218",
							"10.0.16.219",
							"10.0.16.220",
							"10.0.16.221",
							"10.0.16.222",
							"10.0.16.223",
							"10.0.16.224",
							"10.0.16.225",
							"10.0.16.226",
							"10.0.16.227",
							"10.0.16.228",
							"10.0.16.229",
							"10.0.16.230",
							"10.0.16.231",
							"10.0.16.232",
							"10.0.16.233",
							"10.0.16.234",
							"10.0.16.235",
							"10.0.16.236",
							"10.0.16.237",
							"10.0.16.238",
							"10.0.16.239",
							"10.0.16.240",
							"10.0.16.241",
							"10.0.16.242",
							"10.0.16.243",
							"10.0.16.244",
							"10.0.16.245",
							"10.0.16.246",
							"10.0.16.247",
							"10.0.16.248",
							"10.0.16.249",
							"10.0.16.250",
							"10.0.16.251",
							"10.0.16.252",
							"10.0.16.253",
							"10.0.16.254",
						},
					},
				},
				Type: "manual",
			}))

			Expect(manifest.Properties.Etcd).To(Equal(&etcd.PropertiesEtcd{
				Cluster: []etcd.PropertiesEtcdCluster{{
					Instances: 1,
					Name:      "etcd_z1",
				}},
				Machines: []string{
					"10.0.16.10",
				},
				PeerRequireSSL:                  false,
				RequireSSL:                      false,
				HeartbeatIntervalInMilliseconds: 50,
			}))
		})
	})

	Describe("EtcdMembers", func() {
		Context("when there is a single job with a single instance", func() {
			It("returns a list of members in the cluster", func() {
				manifest := etcd.Manifest{
					Jobs: []core.Job{
						{
							Instances: 1,
							Networks: []core.JobNetwork{{
								StaticIPs: []string{"10.244.4.2"},
							}},
							Templates: []core.JobTemplate{
								{
									Name:    "etcd",
									Release: "etcd",
								},
							},
						},
					},
				}

				members := manifest.EtcdMembers()
				Expect(members).To(Equal([]etcd.EtcdMember{{
					Address: "10.244.4.2",
				}}))
			})
		})

		Context("when there are multiple jobs with multiple instances", func() {
			It("returns a list of etcd members in the cluster", func() {
				manifest := etcd.Manifest{
					Jobs: []core.Job{
						{
							Instances: 0,
						},
						{
							Instances: 1,
							Networks: []core.JobNetwork{{
								StaticIPs: []string{"10.244.4.2"},
							}},
							Templates: []core.JobTemplate{
								{
									Name:    "consul_agent",
									Release: "consul",
								},
							},
						},
						{
							Instances: 2,
							Networks: []core.JobNetwork{{
								StaticIPs: []string{"10.244.5.2", "10.244.5.6"},
							}},
							Templates: []core.JobTemplate{
								{
									Name:    "etcd",
									Release: "etcd",
								},
							},
						},
						{
							Instances: 2,
							Networks: []core.JobNetwork{{
								StaticIPs: []string{"10.244.6.2", "10.244.6.6"},
							}},
							Templates: []core.JobTemplate{
								{
									Name:    "etcd",
									Release: "etcd",
								},
							},
						},
					},
				}

				members := manifest.EtcdMembers()
				Expect(members).To(Equal([]etcd.EtcdMember{
					{
						Address: "10.244.5.2",
					},
					{
						Address: "10.244.5.6",
					},
					{
						Address: "10.244.6.2",
					},
					{
						Address: "10.244.6.6",
					},
				}))
			})
		})

		Context("when the job does not have a network", func() {
			It("returns an empty list", func() {
				manifest := etcd.Manifest{
					Jobs: []core.Job{
						{
							Instances: 1,
							Networks:  []core.JobNetwork{},
						},
					},
				}

				members := manifest.EtcdMembers()
				Expect(members).To(BeEmpty())
			})
		})

		Context("when the job network does not have enough static IPs", func() {
			It("returns as much about the list as possible", func() {
				manifest := etcd.Manifest{
					Jobs: []core.Job{
						{
							Instances: 2,
							Networks: []core.JobNetwork{{
								StaticIPs: []string{"10.244.5.2"},
							}},
							Templates: []core.JobTemplate{
								{
									Name:    "etcd",
									Release: "etcd",
								},
							},
						},
					},
				}

				members := manifest.EtcdMembers()
				Expect(members).To(Equal([]etcd.EtcdMember{
					{
						Address: "10.244.5.2",
					},
				}))
			})
		})
	})

	Describe("ToYAML", func() {
		It("returns a YAML representation of the etcd manifest", func() {
			etcdManifest, err := ioutil.ReadFile("fixtures/etcd_manifest.yml")
			Expect(err).NotTo(HaveOccurred())

			manifest := etcd.NewTLSManifest(etcd.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "etcd",
				IPRange:      "10.244.4.0/24",
				IPOffset:     30,
				Secrets: etcd.ConfigSecrets{
					Consul: consul.ConfigSecretsConsul{
						EncryptKey: "consul-encrypt-key",
						AgentKey:   "consul-agent-key",
						AgentCert:  "consul-agent-cert",
						ServerKey:  "consul-server-key",
						ServerCert: "consul-server-cert",
						CACert:     "consul-ca-cert",
					},
					Etcd: etcd.ConfigSecretsEtcd{
						CACert:     "etcd-ca-cert",
						ClientCert: "etcd-client-cert",
						ClientKey:  "etcd-client-key",
						PeerCACert: "etcd-peer-ca-cert",
						PeerCert:   "etcd-peer-cert",
						PeerKey:    "etcd-peer-key",
						ServerCert: "etcd-server-cert",
						ServerKey:  "etcd-server-key",
					},
				},
			}, iaas.NewWardenConfig())

			yaml, err := manifest.ToYAML()
			Expect(err).NotTo(HaveOccurred())
			Expect(yaml).To(MatchYAML(etcdManifest))
		})
	})

	Describe("FromYAML", func() {
		It("returns a Manifest matching the given YAML", func() {
			etcdManifest, err := ioutil.ReadFile("fixtures/etcd_manifest.yml")
			Expect(err).NotTo(HaveOccurred())

			manifest, err := etcd.FromYAML(etcdManifest)
			Expect(err).NotTo(HaveOccurred())

			Expect(manifest.DirectorUUID).To(Equal("some-director-uuid"))
			Expect(manifest.Name).To(Equal("etcd"))
			Expect(manifest.Releases).To(HaveLen(2))
			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "etcd",
				Version: "latest",
			}))
			Expect(manifest.Releases).To(ContainElement(core.Release{
				Name:    "consul",
				Version: "latest",
			}))
			Expect(manifest.Compilation).To(Equal(core.Compilation{
				Network:             "etcd1",
				ReuseCompilationVMs: true,
				Workers:             3,
			}))
			Expect(manifest.Update).To(Equal(core.Update{
				Canaries:        1,
				CanaryWatchTime: "1000-180000",
				MaxInFlight:     1,
				Serial:          true,
				UpdateWatchTime: "1000-180000",
			}))
			Expect(manifest.ResourcePools).To(HaveLen(1))
			Expect(manifest.ResourcePools).To(ContainElement(core.ResourcePool{
				Name:    "etcd_z1",
				Network: "etcd1",
				Stemcell: core.ResourcePoolStemcell{
					Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
					Version: "latest",
				},
			}))

			Expect(manifest.Jobs).To(HaveLen(3))

			Expect(manifest.Jobs[0]).To(Equal(core.Job{
				Name:      "consul_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.45"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "consul_agent",
						Release: "consul",
					},
				},
				Properties: &core.JobProperties{
					Consul: core.JobPropertiesConsul{
						Agent: core.JobPropertiesConsulAgent{
							Mode: "server",
						},
					},
				},
			}))

			Expect(manifest.Jobs[1]).To(Equal(core.Job{
				Name:      "etcd_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.40"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "etcd",
						Release: "etcd",
					},
					{
						Name:    "consul_agent",
						Release: "consul",
					},
				},
				Properties: &core.JobProperties{
					Consul: core.JobPropertiesConsul{
						Agent: core.JobPropertiesConsulAgent{
							Services: core.JobPropertiesConsulAgentServices{
								"etcd": core.JobPropertiesConsulAgentService{},
							},
						},
					},
				},
			}))

			Expect(manifest.Jobs[2]).To(Equal(core.Job{
				Name:      "testconsumer_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.46"},
				}},
				PersistentDisk: 1024,
				ResourcePool:   "etcd_z1",
				Templates: []core.JobTemplate{
					{
						Name:    "consul_agent",
						Release: "consul",
					},
					{
						Name:    "etcd-test-consumer",
						Release: "etcd",
					},
				},
			}))

			Expect(manifest.Networks).To(HaveLen(1))
			Expect(manifest.Networks).To(ContainElement(core.Network{
				Name: "etcd1",
				Subnets: []core.NetworkSubnet{
					{
						CloudProperties: core.NetworkSubnetCloudProperties{Name: "random"},
						Gateway:         "10.244.4.1",
						Range:           "10.244.4.0/24",
						Reserved: []string{
							"10.244.4.2-10.244.4.3",
							"10.244.4.255",
						},
						Static: []string{
							"10.244.4.10",
							"10.244.4.11",
							"10.244.4.12",
							"10.244.4.13",
							"10.244.4.14",
							"10.244.4.15",
							"10.244.4.16",
							"10.244.4.17",
							"10.244.4.18",
							"10.244.4.19",
							"10.244.4.20",
							"10.244.4.21",
							"10.244.4.22",
							"10.244.4.23",
							"10.244.4.24",
							"10.244.4.25",
							"10.244.4.26",
							"10.244.4.27",
							"10.244.4.28",
							"10.244.4.29",
							"10.244.4.30",
							"10.244.4.31",
							"10.244.4.32",
							"10.244.4.33",
							"10.244.4.34",
							"10.244.4.35",
							"10.244.4.36",
							"10.244.4.37",
							"10.244.4.38",
							"10.244.4.39",
							"10.244.4.40",
							"10.244.4.41",
							"10.244.4.42",
							"10.244.4.43",
							"10.244.4.44",
							"10.244.4.45",
							"10.244.4.46",
							"10.244.4.47",
							"10.244.4.48",
							"10.244.4.49",
							"10.244.4.50",
							"10.244.4.51",
							"10.244.4.52",
							"10.244.4.53",
							"10.244.4.54",
							"10.244.4.55",
							"10.244.4.56",
							"10.244.4.57",
							"10.244.4.58",
							"10.244.4.59",
							"10.244.4.60",
							"10.244.4.61",
							"10.244.4.62",
							"10.244.4.63",
							"10.244.4.64",
							"10.244.4.65",
							"10.244.4.66",
							"10.244.4.67",
							"10.244.4.68",
							"10.244.4.69",
							"10.244.4.70",
							"10.244.4.71",
							"10.244.4.72",
							"10.244.4.73",
							"10.244.4.74",
							"10.244.4.75",
							"10.244.4.76",
							"10.244.4.77",
							"10.244.4.78",
							"10.244.4.79",
							"10.244.4.80",
							"10.244.4.81",
							"10.244.4.82",
							"10.244.4.83",
							"10.244.4.84",
							"10.244.4.85",
							"10.244.4.86",
							"10.244.4.87",
							"10.244.4.88",
							"10.244.4.89",
							"10.244.4.90",
							"10.244.4.91",
							"10.244.4.92",
							"10.244.4.93",
							"10.244.4.94",
							"10.244.4.95",
							"10.244.4.96",
							"10.244.4.97",
							"10.244.4.98",
							"10.244.4.99",
							"10.244.4.100",
							"10.244.4.101",
							"10.244.4.102",
							"10.244.4.103",
							"10.244.4.104",
							"10.244.4.105",
							"10.244.4.106",
							"10.244.4.107",
							"10.244.4.108",
							"10.244.4.109",
							"10.244.4.110",
							"10.244.4.111",
							"10.244.4.112",
							"10.244.4.113",
							"10.244.4.114",
							"10.244.4.115",
							"10.244.4.116",
							"10.244.4.117",
							"10.244.4.118",
							"10.244.4.119",
							"10.244.4.120",
							"10.244.4.121",
							"10.244.4.122",
							"10.244.4.123",
							"10.244.4.124",
							"10.244.4.125",
							"10.244.4.126",
							"10.244.4.127",
							"10.244.4.128",
							"10.244.4.129",
							"10.244.4.130",
							"10.244.4.131",
							"10.244.4.132",
							"10.244.4.133",
							"10.244.4.134",
							"10.244.4.135",
							"10.244.4.136",
							"10.244.4.137",
							"10.244.4.138",
							"10.244.4.139",
							"10.244.4.140",
							"10.244.4.141",
							"10.244.4.142",
							"10.244.4.143",
							"10.244.4.144",
							"10.244.4.145",
							"10.244.4.146",
							"10.244.4.147",
							"10.244.4.148",
							"10.244.4.149",
							"10.244.4.150",
							"10.244.4.151",
							"10.244.4.152",
							"10.244.4.153",
							"10.244.4.154",
							"10.244.4.155",
							"10.244.4.156",
							"10.244.4.157",
							"10.244.4.158",
							"10.244.4.159",
							"10.244.4.160",
							"10.244.4.161",
							"10.244.4.162",
							"10.244.4.163",
							"10.244.4.164",
							"10.244.4.165",
							"10.244.4.166",
							"10.244.4.167",
							"10.244.4.168",
							"10.244.4.169",
							"10.244.4.170",
							"10.244.4.171",
							"10.244.4.172",
							"10.244.4.173",
							"10.244.4.174",
							"10.244.4.175",
							"10.244.4.176",
							"10.244.4.177",
							"10.244.4.178",
							"10.244.4.179",
							"10.244.4.180",
							"10.244.4.181",
							"10.244.4.182",
							"10.244.4.183",
							"10.244.4.184",
							"10.244.4.185",
							"10.244.4.186",
							"10.244.4.187",
							"10.244.4.188",
							"10.244.4.189",
							"10.244.4.190",
							"10.244.4.191",
							"10.244.4.192",
							"10.244.4.193",
							"10.244.4.194",
							"10.244.4.195",
							"10.244.4.196",
							"10.244.4.197",
							"10.244.4.198",
							"10.244.4.199",
							"10.244.4.200",
							"10.244.4.201",
							"10.244.4.202",
							"10.244.4.203",
							"10.244.4.204",
							"10.244.4.205",
							"10.244.4.206",
							"10.244.4.207",
							"10.244.4.208",
							"10.244.4.209",
							"10.244.4.210",
							"10.244.4.211",
							"10.244.4.212",
							"10.244.4.213",
							"10.244.4.214",
							"10.244.4.215",
							"10.244.4.216",
							"10.244.4.217",
							"10.244.4.218",
							"10.244.4.219",
							"10.244.4.220",
							"10.244.4.221",
							"10.244.4.222",
							"10.244.4.223",
							"10.244.4.224",
							"10.244.4.225",
							"10.244.4.226",
							"10.244.4.227",
							"10.244.4.228",
							"10.244.4.229",
							"10.244.4.230",
							"10.244.4.231",
							"10.244.4.232",
							"10.244.4.233",
							"10.244.4.234",
							"10.244.4.235",
							"10.244.4.236",
							"10.244.4.237",
							"10.244.4.238",
							"10.244.4.239",
							"10.244.4.240",
							"10.244.4.241",
							"10.244.4.242",
							"10.244.4.243",
							"10.244.4.244",
							"10.244.4.245",
							"10.244.4.246",
							"10.244.4.247",
							"10.244.4.248",
							"10.244.4.249",
							"10.244.4.250",
							"10.244.4.251",
							"10.244.4.252",
							"10.244.4.253",
							"10.244.4.254",
						},
					},
				},
				Type: "manual",
			}))
			Expect(manifest.Properties.Etcd).To(Equal(&etcd.PropertiesEtcd{
				Cluster: []etcd.PropertiesEtcdCluster{{
					Instances: 1,
					Name:      "etcd_z1",
				}},
				Machines: []string{
					"etcd.service.cf.internal",
				},
				PeerRequireSSL:                  true,
				RequireSSL:                      true,
				HeartbeatIntervalInMilliseconds: 50,
				AdvertiseURLsDNSSuffix:          "etcd.service.cf.internal",
				CACert:                          "etcd-ca-cert",
				ClientCert:                      "etcd-client-cert",
				ClientKey:                       "etcd-client-key",
				PeerCACert:                      "etcd-peer-ca-cert",
				PeerCert:                        "etcd-peer-cert",
				PeerKey:                         "etcd-peer-key",
				ServerCert:                      "etcd-server-cert",
				ServerKey:                       "etcd-server-key",
			}))
		})

		Context("failure cases", func() {
			It("should error on malformed YAML", func() {
				_, err := etcd.FromYAML([]byte("%%%%%%%%%%"))
				Expect(err).To(MatchError(ContainSubstring("yaml: ")))
			})
		})
	})
})
