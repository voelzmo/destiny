package turbulencewithops

const manifestV2 = `
director_uuid: REPLACE_ME
name: REPLACE_ME

stemcells:
- alias: default
  os: ubuntu-trusty
  version: latest

releases:
- name: turbulence
  version: latest
- name: bosh-warden-cpi
  version: latest

instance_groups:
- name: api
  instances: 1
  azs: []
  jobs:
  - name: turbulence_api
    release: turbulence
  - name: warden_cpi
    release: bosh-warden-cpi
  vm_type: default
  stemcell: default
  persistent_disk_type: 1GB
  networks:
  - name: private
  properties:
    turbulence_api:
      certificate: |+
        -----BEGIN CERTIFICATE-----
        MIIEYDCCAkigAwIBAgIQZemY/h/wilbjexxXyuR9DTANBgkqhkiG9w0BAQsFADAa
        MRgwFgYDVQQDEw9CT1NIIEJvb3Rsb2FkZXIwHhcNMTYwOTI3MjEyNzQzWhcNMTgw
        OTI3MjEyNzQzWjAXMRUwEwYDVQQDEwwxMC4yNDQuNC4yNDQwggEiMA0GCSqGSIb3
        DQEBAQUAA4IBDwAwggEKAoIBAQCeSQemenvCEgbq+UP3Hdigzh1+KfRG+ytV0Zii
        IQ5xghxXxd2HqumXQeyyLC6usMdnGMu8DjmihNFBzrlnC6b44S2jPu0nttNqGL9n
        hTk30mpqltt0epQqyQ7EBaS2/WGlG+S2emAPtJcSAHyporS8Az54gi1J9TilGVus
        0m30MNm1drnInOyBWb8LnA+0U0vwAF17wABwBzUjP3w6ac32+A5TBAotdC3NrLBJ
        yHynOAzBv0IVYFVsLo3z9TBSQbAOZYV6hrZ1vZygFm6n6IIa5psKUC5BNXLgFkJd
        bz3AiOtbH98hVX5pJ8e2gxfyi1BdZFcuQsQAUlwdQOWI5eMrAgMBAAGjgaQwgaEw
        DgYDVR0PAQH/BAQDAgO4MB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAd
        BgNVHQ4EFgQUDjTbDO9tCTlpyAhLrugwMehN1qwwHwYDVR0jBBgwFoAUzJJ8qHy8
        rqI09P4ODZXR7e03h78wMAYDVR0RBCkwJ4IMMTAuMjQ0LjQuMjQ0ggsxMC4wLjE2
        LjI0NIcECvQE9IcECgAQ9DANBgkqhkiG9w0BAQsFAAOCAgEAceCt2FBqQ+/R6gTK
        oMlQWk5kE6dKR3b7gXWjkJy63pCnhN3g1EsEztOihj89E34WrFXm1uLvjUEZr3Tn
        J+R4rD0M6JDYMm2WAgrxAFMHg2SvY8ruR0WxrSmtjxX4aOP7DAqXU3GTWP8dBCrN
        xDsEGwOrB4kLe+mpJ0b6DGZUn2PLAyDP7DcJ05WkIZHIVOit6ZqOwc7er8vsA3rQ
        6a18FkfGV86buolj7b3PirrtE+sklctmg9APsSOI2b+kF43Jmfst8dLbU4d2rYJV
        w0Db5vjfjS+CPkIrbuhy6UXAIrK9GiEi8qhqdDU9+QbPwloPveLsRa3gPqCTCJSR
        xaZte3CQ0RIjaudLevissLx0qOZR1toFyg6HMIh8al6NqKPSTHGp9DYpK3829wNE
        fw7Gb/d2gk8rD4xc8A0PePcQ8GlhcRO+YNPKhVi1enkYENEKXFRiPUhUuGQcN/02
        KxNOJcRQxfvGSPywpTMT6D94FWYhtqwA+Eco8uAerJL7Viy7vgzP7RHpYTlScVoD
        FQ73kH0bsXHBTIe1Ui6UJzz1iCnULeCzSng4vbIQgUF/jgWR+18iHCSINiIQs6bq
        OMUl9uqZ5rvqLRZJra8S5mOWEFA+Jkpsxtf1yYchclYM4N/9TjortAOw+z5XnsSM
        Qf86uzapOfTUGuLNRGUls8+6BkA=
        -----END CERTIFICATE-----
      cpi_job_name: warden_cpi
      director:
        ca_cert: |+
          -----BEGIN CERTIFICATE-----
          MIIDtzCCAp+gAwIBAgIJAMZ/qRdRamluMA0GCSqGSIb3DQEBBQUAMEUxCzAJBgNV
          BAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBX
          aWRnaXRzIFB0eSBMdGQwIBcNMTYwODI2MjIzMzE5WhgPMjI5MDA2MTAyMjMzMTla
          MEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJ
          bnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
          ggEKAoIBAQDN/bv70wDn6APMqiJZV7ESZhUyGu8OzuaeEfb+64SNvQIIME0s9+i7
          D9gKAZjtoC2Tr9bJBqsKdVhREd/X6ePTaopxL8shC9GxXmTqJ1+vKT6UxN4kHr3U
          +Y+LK2SGYUAvE44nv7sBbiLxDl580P00ouYTf6RJgW6gOuKpIGcvsTGA4+u0UTc+
          y4pj6sT0+e3xj//Y4wbLdeJ6cfcNTU63jiHpKc9Rgo4Tcy97WeEryXWz93rtRh8d
          pvQKHVDU/26EkNsPSsn9AHNgaa+iOA2glZ2EzZ8xoaMPrHgQhcxoi8maFzfM2dX2
          XB1BOswa/46yqfzc4xAwaW0MLZLg3NffAgMBAAGjgacwgaQwHQYDVR0OBBYEFNRJ
          PYFebixALIR2Ee+yFoSqurxqMHUGA1UdIwRuMGyAFNRJPYFebixALIR2Ee+yFoSq
          urxqoUmkRzBFMQswCQYDVQQGEwJBVTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8G
          A1UEChMYSW50ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkggkAxn+pF1FqaW4wDAYDVR0T
          BAUwAwEB/zANBgkqhkiG9w0BAQUFAAOCAQEAoPTwU2rm0ca5b8xMni3vpjYmB9NW
          oSpGcWENbvu/p7NpiPAe143c5EPCuEHue/AbHWWxBzNAZvhVZBeFirYNB3HYnCla
          jP4WI3o2Q0MpGy3kMYigEYG76WeZAM5ovl0qDP6fKuikZofeiygb8lPs7Hv4/88x
          pSsZYBm7UPTS3Pl044oZfRJdqTpyHVPDqwiYD5KQcI0yHUE9v5KC0CnqOrU/83PE
          b0lpHA8bE9gQTQjmIa8MIpaP3UNTxvmKfEQnk5UAZ5xY2at5mmyj3t8woGdzoL98
          yDd2GtrGsguQXM2op+4LqEdHef57g7vwolZejJqN776Xu/lZtCTp01+HTA==
          -----END CERTIFICATE-----
        host: REPLACE_ME
        password: REPLACE_ME
        username: REPLACE_ME
      password: turbulence-password
      private_key: |+
        -----BEGIN RSA PRIVATE KEY-----
        MIIEowIBAAKCAQEAnkkHpnp7whIG6vlD9x3YoM4dfin0RvsrVdGYoiEOcYIcV8Xd
        h6rpl0HssiwurrDHZxjLvA45ooTRQc65Zwum+OEtoz7tJ7bTahi/Z4U5N9Jqapbb
        dHqUKskOxAWktv1hpRvktnpgD7SXEgB8qaK0vAM+eIItSfU4pRlbrNJt9DDZtXa5
        yJzsgVm/C5wPtFNL8ABde8AAcAc1Iz98OmnN9vgOUwQKLXQtzaywSch8pzgMwb9C
        FWBVbC6N8/UwUkGwDmWFeoa2db2coBZup+iCGuabClAuQTVy4BZCXW89wIjrWx/f
        IVV+aSfHtoMX8otQXWRXLkLEAFJcHUDliOXjKwIDAQABAoIBAE6x0lrjlHoFSmky
        aqcGsLgqIaGjLC7KG158aV1Z//tRC9RbrGsR/zbTcOGYc9EoVMseGdSyYdc6H2uP
        YsAnm+kq0uzRkUjoba3XBfpq0uq882lw/USo2Nd4xJ2SjLTpvs0+0/QhXXcRevZZ
        RaF1IlRDbKCvX+LgRzxWIi0HJeF8S79F/8s+M20JxqMWOOS8X6KqYJjzp/0NhYXC
        g39VCtfi6peWm+Yvo5rF3UJRU0APVKAdIccZesf7faoOMJmQH/pOQn2JJQ+xMRDv
        YrxntShd/aNFaGWSuxVhlTbHfIgOr9ogk+nfvkIVxrHpj/vQ57QIMrHnC+sz377M
        zYzaFfECgYEAwXZ00T2de/84idWJZnWzABIfh549WzQGV+0olvb/7h5zF79w1jHo
        6JlWpNYxF/N6x0hZUh2x+aPBV4zU7FkUwW0aVEI2K0zhgy4pqMFMvG2GcGjonIO2
        nzQVUNyD1j2gMLxKquGe6zpXZrqiLynbAYC7d2fk2P2Jq5pk5mkLdZMCgYEA0XOM
        WoE5UtTCOpYdwVGBC/Vu2pmyH/n8baZI+yNWb2YkIarYbxFZtf6HVqHjIwoQ4PE1
        lk2gD1sKkos0tDPgUtM+6eGd995N1lClxfEWBz//yVI9DvaE/UPlJCFd5l8spPXm
        LrTMfJaAK38yMHeJPVt8YAF/5nu97Q8FemCu2wkCgYB4mvhIWTkMTBdbFhwKG+Xz
        bVjqmuN1MAGkXtynAGScda8aZuZZIdQo7S7uo/kHDWrFQX0tjAWfs06c3db/YKln
        zDRVwtEyPUN5HBYsdhT4gu8EtOIOcK4woa+IMXCe5twuhbOmw/DmhABosoDZFibJ
        0Q8NaV9pRuXEbQPqACJ8sQKBgDiMd9JfnTht1Nq4eOQeuzadVwaSBHN5rNt1z7Ju
        QgHlk8+7LqAeERh/1c5f+tEVAKWauhsQbix3Kg2So/IbJ291NUEz9tBbJqy8LWWZ
        x5bBgq+6El4d1J4EXLM6hv2RqJ4I/dKSYbspbwVPXB+VxmnYb2YEQaHautZr/dCi
        ldLxAoGBAIQi+Q6Q5eV2RyPLQpqBP82cJeVS/q//DMNATGlEtbkSmE41x9ZRaZjL
        GJdB+xzKRKvI06YgEkqHNkBBxfu2/+VNXC3x0J1ZlHich8wewQo1N6AMASmhLb6/
        Df5iij5skYM4xaXrfgiPEFQgblA2rQEyKiXa0VXEmhItuf6t7XTL
        -----END RSA PRIVATE KEY-----
    warden_cpi:
      agent:
        blobstore:
          options:
            endpoint: http://10.254.50.4:25251
            password: agent-password
            user: agent
          provider: dav
        mbus: nats://nats:nats-password@10.254.50.4:4222
      warden:
        connect_address: 10.254.50.4:7777
        connect_network: tcp

update:
  canaries: 1
  canary_watch_time: 1000-180000
  max_in_flight: 1
  serial: true
  update_watch_time: 1000-180000
`
