package consul

const manifestV2 = `
name: REPLACE_ME

stemcells:
- alias: default
  os: ubuntu-xenial
  version: latest

releases:
- name: consul
  version: latest

instance_groups:
- name: consul
  instances: 1
  azs: []
  jobs:
  - name: consul_agent
    release: consul
    consumes:
      consul_common: { from: common_link }
    provides:
      consul_common: { as: common_link }
  vm_type: default
  stemcell: default
  persistent_disk_type: 1GB
  networks:
  - name: default
  properties:
    consul:
      server:
        tls: ((/consul_server_tls))
      agent:
        tls: ((/consul_agent_tls))
        mode: server
        log_level: info
        domain: cf.internal
        datacenter: dc1
        services:
          router:
            name: gorouter
            check:
              name: router-check
              args: ["/var/vcap/jobs/router/bin/script"]
              interval: 1m
            tags: [routing]
          cloud_controller: {}
      
      encrypt_keys:
      - ((/consul_encrypt_key))
     
- name: testconsumer
  instances: 1
  azs: []
  jobs:
  - name: consul_agent
    release: consul
    consumes:
      consul_common: { from: common_link }
  - name: consul-test-consumer
    release: consul
  vm_type: default
  stemcell: default
  persistent_disk_type: 1GB
  networks:
  - name: default
  properties:
    consul:
      encrypt_keys:
        - ((/consul_encrypt_key))
      server:
        tls: ((/consul_server_tls))
      agent:
        tls: ((/consul_agent_tls))

update:
  canaries: 1
  canary_watch_time: 1000-180000
  max_in_flight: 1
  serial: true
  update_watch_time: 1000-180000

variables:
  - name: /consul_encrypt_key
    type: password
  - name: /consul_agent_ca
    type: certificate
    options:
      is_ca: true
      common_name: consulCA
  - name: /consul_agent_tls
    type: certificate
    options:
      ca: /consul_agent_ca
      common_name: consul_agent
      extended_key_usage:
      - client_auth
      - server_auth
      alternative_names:
      - 127.0.0.1
  - name: /consul_server_tls
    type: certificate
    options:
      ca: /consul_agent_ca
      common_name: server.dc1.cf.internal
      extended_key_usage:
      - client_auth
      - server_auth
`
