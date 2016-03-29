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
	Describe("ToYAML", func() {
		It("returns a YAML representation of the etcd manifest", func() {
			etcdManifest, err := ioutil.ReadFile("fixtures/etcd_manifest.yml")
			Expect(err).NotTo(HaveOccurred())

			manifest := etcd.NewEtcd(etcd.Config{
				DirectorUUID: "some-director-uuid",
				Name:         "etcd",
				IPRange:      "10.244.4.0/24",
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
			Expect(manifest.Jobs).To(HaveLen(2))
			Expect(manifest.Jobs[0]).To(Equal(core.Job{
				Name:      "consul_z1",
				Instances: 1,
				Networks: []core.JobNetwork{{
					Name:      "etcd1",
					StaticIPs: []string{"10.244.4.9"},
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
					StaticIPs: []string{"10.244.4.4"},
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
							"10.244.4.13-10.244.4.254",
						},
						Static: []string{
							"10.244.4.4",
							"10.244.4.5",
							"10.244.4.6",
							"10.244.4.7",
							"10.244.4.8",
							"10.244.4.9",
						},
					},
				},
				Type: "manual",
			}))
			Expect(manifest.Properties.Etcd).To(Equal(&etcd.PropertiesEtcd{
				Machines: []string{
					"10.244.4.4",
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
