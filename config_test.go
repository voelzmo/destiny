package destiny_test

import (
	"github.com/pivotal-cf-experimental/destiny"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("NewConfigWithDefaults", func() {

		It("applies the default values for fields that are missing", func() {
			config := destiny.NewConfigWithDefaults(destiny.Config{})
			Expect(config.Secrets.Consul.CACert).To(Equal(destiny.ConsulCACert))
			Expect(config.Secrets.Consul.EncryptKey).To(Equal(destiny.ConsulEncryptKey))
			Expect(config.Secrets.Consul.AgentKey).To(Equal(destiny.ConsulAgentKey))
			Expect(config.Secrets.Consul.AgentCert).To(Equal(destiny.ConsulAgentCert))
			Expect(config.Secrets.Consul.ServerKey).To(Equal(destiny.ConsulServerKey))
			Expect(config.Secrets.Consul.ServerCert).To(Equal(destiny.ConsulServerCert))

			Expect(config.Secrets.Etcd.CACert).To(Equal(destiny.EtcdCACert))
			Expect(config.Secrets.Etcd.ClientCert).To(Equal(destiny.EtcdClientCert))
			Expect(config.Secrets.Etcd.ClientKey).To(Equal(destiny.EtcdClientKey))
			Expect(config.Secrets.Etcd.PeerCACert).To(Equal(destiny.EtcdPeerCACert))
			Expect(config.Secrets.Etcd.PeerCert).To(Equal(destiny.EtcdPeerCert))
			Expect(config.Secrets.Etcd.PeerKey).To(Equal(destiny.EtcdPeerKey))
			Expect(config.Secrets.Etcd.ServerCert).To(Equal(destiny.EtcdServerCert))
			Expect(config.Secrets.Etcd.ServerKey).To(Equal(destiny.EtcdServerKey))
		})
	})
})
