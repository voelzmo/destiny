package destiny_test

import (
	"github.com/pivotal-cf-experimental/destiny"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Properties", func() {
	Describe("Merge", func() {
		var (
			properties destiny.Properties
		)

		BeforeEach(func() {
			properties = destiny.Properties{
				Etcd: &destiny.PropertiesEtcd{
					RequireSSL: false,
				},
				Consul: &destiny.PropertiesConsul{
					CACert: "some-ca-cert",
				},
				TurbulenceAPI: &destiny.PropertiesTurbulenceAPI{
					Certificate: "some-certificate",
				},
				WardenCPI: &destiny.PropertiesWardenCPI{
					Warden: destiny.PropertiesWardenCPIWarden{
						ConnectAddress: "some-connect-address",
					},
				},
				AWS: &destiny.PropertiesAWS{
					AccessKeyID: "some-access-key-id",
				},
				Registry: &destiny.PropertiesRegistry{
					Host: "some-host",
				},
				Blobstore: &destiny.PropertiesBlobstore{
					Address: "some-address",
				},
				Agent: &destiny.PropertiesAgent{
					Mbus: "some-mbus",
				},
			}
		})

		It("returns a merged properties with updated etcd fields", func() {
			etcdProperties := destiny.Properties{
				Etcd: &destiny.PropertiesEtcd{
					RequireSSL: true,
				},
			}
			newProperties := properties.Merge(etcdProperties)

			Expect(newProperties.Etcd).To(Equal(&destiny.PropertiesEtcd{
				RequireSSL: true,
			}))
		})

		It("returns a merged properties with updated consul fields", func() {
			consulProperties := destiny.Properties{
				Consul: &destiny.PropertiesConsul{
					CACert: "some-new-ca-cert",
				},
			}
			newProperties := properties.Merge(consulProperties)

			Expect(newProperties.Consul).To(Equal(&destiny.PropertiesConsul{
				CACert: "some-new-ca-cert",
			}))
		})

		It("returns a merged properties with updated turbulence api", func() {
			turbulenceProperties := destiny.Properties{
				TurbulenceAPI: &destiny.PropertiesTurbulenceAPI{
					Certificate: "some-new-certificate",
				},
			}
			newProperties := properties.Merge(turbulenceProperties)

			Expect(newProperties.TurbulenceAPI).To(Equal(&destiny.PropertiesTurbulenceAPI{
				Certificate: "some-new-certificate",
			}))
		})

		It("returns a merged properties with updated warden cpi properties", func() {
			wardenProperties := destiny.Properties{
				WardenCPI: &destiny.PropertiesWardenCPI{
					Warden: destiny.PropertiesWardenCPIWarden{
						ConnectAddress: "some-new-connect-address",
					},
				},
			}
			newProperties := properties.Merge(wardenProperties)

			Expect(newProperties.WardenCPI).To(Equal(&destiny.PropertiesWardenCPI{
				Warden: destiny.PropertiesWardenCPIWarden{
					ConnectAddress: "some-new-connect-address",
				},
			}))
		})

		It("returns a merged properties with updated AWS properties", func() {
			awsProperties := destiny.Properties{
				AWS: &destiny.PropertiesAWS{
					AccessKeyID: "some-new-access-key-id",
				},
			}
			newProperties := properties.Merge(awsProperties)

			Expect(newProperties.AWS).To(Equal(&destiny.PropertiesAWS{
				AccessKeyID: "some-new-access-key-id",
			}))
		})

		It("returns a merged properties with updated Registry properties", func() {
			registryProperties := destiny.Properties{
				Registry: &destiny.PropertiesRegistry{
					Host: "some-new-host",
				},
			}
			newProperties := properties.Merge(registryProperties)

			Expect(newProperties.Registry).To(Equal(&destiny.PropertiesRegistry{
				Host: "some-new-host",
			}))
		})

		It("returns a merged properties with updated blobstore properties", func() {
			blobstoreProperties := destiny.Properties{
				Blobstore: &destiny.PropertiesBlobstore{
					Address: "some-new-address",
				},
			}
			newProperties := properties.Merge(blobstoreProperties)

			Expect(newProperties.Blobstore).To(Equal(&destiny.PropertiesBlobstore{
				Address: "some-new-address",
			}))
		})

		It("returns a merged properties with updated agent properties", func() {
			agentProperties := destiny.Properties{
				Agent: &destiny.PropertiesAgent{
					Mbus: "some-new-mbus",
				},
			}
			newProperties := properties.Merge(agentProperties)

			Expect(newProperties.Agent).To(Equal(&destiny.PropertiesAgent{
				Mbus: "some-new-mbus",
			}))
		})
	})
})
