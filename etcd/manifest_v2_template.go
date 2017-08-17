package etcd

const (
	manifestV2NonTLS = `
name: REPLACE_ME

stemcells:
- alias: default
  os: ubuntu-trusty
  version: latest

releases:
- name: etcd
  version: latest

instance_groups:
- name: etcd
  instances: 3
  azs:
  - z1
  - z2
  jobs:
  - name: etcd
    release: etcd
    consumes:
      etcd: { from: etcd_server }
    provides:
      etcd: { as: etcd_server }
  vm_type: large
  stemcell: default
  persistent_disk_type: 1GB
  networks:
  - name: private
  properties:
    etcd:
      enable_debug_logging: true
      heartbeat_interval_in_milliseconds: 50
      peer_require_ssl: false
      require_ssl: false

- name: testconsumer
  instances: 1
  azs:
  - z1
  - z2
  jobs:
  - name: etcd_testconsumer
    release: etcd
    consumes:
      etcd: { from: etcd_server }
  vm_type: default
  stemcell: default
  persistent_disk_type: 1GB
  networks:
  - name: private

update:
  canaries: 1
  canary_watch_time: 1000-180000
  max_in_flight: 1
  serial: true
  update_watch_time: 1000-180000
`
	manifestV2TLS = `
name: REPLACE_ME

stemcells:
- alias: default
  os: ubuntu-trusty
  version: latest

releases:
- name: etcd
  version: latest
- name: consul
  version: latest

instance_groups:
- name: consul
  instances: 1
  azs:
  - z1
  - z2
  jobs:
  - name: consul_agent
    release: consul
    consumes:
      consul_common: { from: common_link }
      consul_server: { from: server_link }
      consul_client: { from: client_link }
    provides:
      consul_common: { as: common_link }
      consul_server: { as: server_link }
      consul_client: { as: client_link }
  vm_type: default
  stemcell: default
  persistent_disk_type: 1GB
  networks:
  - name: private
  properties:
    consul:
      agent:
        mode: server
        domain: cf.internal
      encrypt_keys:
      - Atzo3VBv+YVDzQAzlQRPRA==
      require_ssl: true
      agent_cert: |+
        -----BEGIN CERTIFICATE-----
        MIIEODCCAiCgAwIBAgIQLzVpihAzFtaw+0Nr81S4czANBgkqhkiG9w0BAQsFADAT
        MREwDwYDVQQDEwhjb25zdWxDQTAeFw0xNzAxMTIxOTU1MDdaFw0xOTAxMTIxOTU1
        MDdaMBcxFTATBgNVBAMTDGNvbnN1bCBhZ2VudDCCASIwDQYJKoZIhvcNAQEBBQAD
        ggEPADCCAQoCggEBANTMBAUp8pBiqWQGOAMwULjwyiSwkCfK8PggdO0PotzjiMrB
        EpY3rKaMZm7QrB9JKSFjisNY4L3eXcy9JgdIR53/Bnzbm1dTz283820L9XmuHMjE
        OFzicgHK8oBw2qL2H4gzRNy6YYgpc5QiXsOlO15nqj3j5B+7GVRV4On6cfT9FwQQ
        +GxcB5UKDLjL/7+rzYGVJnO+DokGtg7E3V9Cgd7Xum54oDKDvBfCVsSv9Gn38rSj
        iyZ79VZ8FJxc8sFqjRnboz/cIaBm8K0BcU4Oq8GscHHxpAUvcfLpITWLZoFhckCj
        0ApvcsoN9ZkxgbJHeowbsKqDH3Fy29e0a4kDgN0CAwEAAaOBgzCBgDAOBgNVHQ8B
        Af8EBAMCA7gwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB0GA1UdDgQW
        BBTBXAohz/DCAzP6qd7rR25RIdJFCjAfBgNVHSMEGDAWgBTDBGov4zbZqKmeNsRA
        tMDZpQNUyjAPBgNVHREECDAGhwR/AAABMA0GCSqGSIb3DQEBCwUAA4ICAQAdAmKV
        H3J3jRewJizV8msbftbcqB/hdQRJhwcN+LNA14gGcTP2f34lO452MsdUUvRBz+Te
        D2HKUU0WobuIeCf0TvMreLP+kDBOrPGmS0BM0mA5OnExm5WNTLB2xujqhmbP4kOB
        Wm61XAXNANWrCkflm6axmSMUJd7RZIpZFyFTYeb2739IeEWY0s04GfLC1pmL++Tb
        OR7Auo+qvekgv+fGbaVo22qEpNdU9xxhbdMUMgT3WqQX49WyQBi5EDVkukXmoVRx
        bsjXtUpyY7i3vwZTvFeWLAyfwIpFv6KXtYZ9NLJ3/2pzahBOZX8S+PnfV3rnJTjk
        P9z96ZR+9hNckXz0Wyve47+ffliR2YU1ouKjNWumZh/d1x1oBlV+ZWKgBUVL2vmw
        jkWw2W6aVahjfukyjyiFgI6K7KDcFvl5ce/3ub8k1/jTIjs/+OiaM2uVesUEkd2Z
        u4dk8Cw6bxr3zB8L+wPX1IIDhi+aOQw8gWDGezmHnC+FruEQvrwbOVvn0pwe93SW
        CNLmvj+2Rbclk9nqB1KBwWQur+4OM04mfT15QptV8adTnaC0wJQKD7ZsluwHZYQF
        ALnJk09TXg+eZtMJ9zD4nPvssEOal3VzZ7tPPzbuartWSZMOlDS5HziB1OQ4Gh8X
        s9pKNZhW+cRPukmjNHHjU/LF6NmkSUmqmiI8Tg==
        -----END CERTIFICATE-----
      agent_key: |+
        -----BEGIN RSA PRIVATE KEY-----
        MIIEpAIBAAKCAQEA1MwEBSnykGKpZAY4AzBQuPDKJLCQJ8rw+CB07Q+i3OOIysES
        ljespoxmbtCsH0kpIWOKw1jgvd5dzL0mB0hHnf8GfNubV1PPbzfzbQv1ea4cyMQ4
        XOJyAcrygHDaovYfiDNE3LphiClzlCJew6U7XmeqPePkH7sZVFXg6fpx9P0XBBD4
        bFwHlQoMuMv/v6vNgZUmc74OiQa2DsTdX0KB3te6bnigMoO8F8JWxK/0affytKOL
        Jnv1VnwUnFzywWqNGdujP9whoGbwrQFxTg6rwaxwcfGkBS9x8ukhNYtmgWFyQKPQ
        Cm9yyg31mTGBskd6jBuwqoMfcXLb17RriQOA3QIDAQABAoIBAQCO9TCON4wZq+6Y
        oATpP4A7fqiO1X9C/He+ei+TQznqo4G2lNbjzCtVCGWYdN/tdL0JDVKfwgnaBJWH
        glsV8V0Lq9Sz9OT7Wfa1hSUoUSxsvqffyNMEs6xbv/gCic6YRDkSyz6r+xqi2xYm
        oqB/V3X3CjW4tmz/VDbEDZ24EuST7GnrumU4RHbGW1aqOpVhRX/3zIQ+f9YQGFMZ
        /BkusrSfVUqlV2VHGZkxsyA8+Xd8Jyi/TAksnxx/2HoXjcjGr5XUrvdM/XVdlAtr
        52TmUn0Fpyu7oYf4oKFknVRzXmxPy2WeCyaiGAq+m/F7JAF/CIqIusA9ASR+B0R4
        1Rw79H+BAoGBAO4Lxuo/5Q7j6cDRAMFQIVazrc/0ykwwKytAqjZoavzTL4zANt0G
        1ezu9tK7uzPXgtlFNiZ3gIpQDWAfPRE4gLXjmwB3TIh1EiBOwTgtfzGoUCp3wWmR
        ARHb933ik+xQnG4u6P6iHrzWAYZUprvumv3rDEK/C2bmAMsqj9WWJQHVAoGBAOTY
        ufeNL0XYZYGTXEJj8RIi7owuUQUq4tOcu3BoexLKbjs/faLlucMXstPs5ae5JKsP
        5weOEgqgiKZNHvLh619opdTJyggjcEyK48SqLmEAwJ+M8EV3LYyw6zZtKWQDR15N
        OBTNPTc6jQzqUOE8z/4ZP/CVkW95HqXOOtBptn7pAoGBAJsuHD8q5gTd+M1UsnxS
        41jlCyLs/k/KeunYXt3XFh+5AF9uEpXl1eF+KnNYJIJ4NHm1D8bl0mrYItANrT6j
        qexo8uvL2Z1/TBC5pmYb6rYRdikpJnHOMHdXATEUWsAMEN4XQJZ2Uzlg/V93obYT
        pwBukPCWIDW1LMFE/r0LAxb9AoGAO98fuE5twb49wErHZm8zUOVmt7IebFWuBmMI
        /v22xVHEySdxPT8Q/KOkm6Fs7BaaK077yJQ40CLz3V5r7GuC4vFEAYnRm5N5++yS
        bo9/ls1Vl+iNq/7kIdzfjNu+anYZI+jb9UVE8MAWyvw6sNLyL653dgALjriHdiWg
        aYpevpECgYA5Wh1o5TCjxBYaoPNMy432OwGiuUuP4pU46OBH5b9G+aHxPmFIQ1/b
        axQ9CBFypBMp1GVYym8k0xhV5NU6UO8sUwMz4+g3jpECLxO6z2JS5IdRsuLZWLxf
        5XnBawt8jmo/1wrttqr1EFuhNWOw1ypIKnTYk0Zfe42mwZuUvWxYpg==
        -----END RSA PRIVATE KEY-----
      ca_cert: |+
        -----BEGIN CERTIFICATE-----
        MIIE5jCCAs6gAwIBAgIBATANBgkqhkiG9w0BAQsFADATMREwDwYDVQQDEwhjb25z
        dWxDQTAeFw0xNzAxMTIxOTU1MDVaFw0yNzAxMTIxOTU1MDZaMBMxETAPBgNVBAMT
        CGNvbnN1bENBMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAvnVQDLCe
        xBR6qUaGk1vdRyVGQnHsJnOtrshthaoshqOfsjwbpQsXBCRv5R2vIudebBN4T28A
        3WwgMBhLTiAiFU+7fWMPtIOcPStkc13stv9n09QaWu+hVtYrq5s0qS/2e3ZpcFkv
        wkuHyLzb/vYX06r65gHFMmHbh3GuzzlhlKpqKO++3bfe4GGWWB4AeXRqhMPB8ONP
        ap035q3zarmqp19hx55vIvbXszUGNGC56tWT5KxKk8feF73Fg1Nzd7+HIJSgPVoc
        1SkToiKbfjTyf6iOPiJZ5yrD3tD47IQY0DgFNlGB2lLtE38k1ktZeRF+eh9oMC1A
        XKHz73oljji9YcrKp6V8DkXggFctSspIHcbL5vuEMEeiUWWINvZRZaakWkqArVN7
        UFI4A35Uaitxkjh2ngdVpOYKCxZh7VdyxAnOMOrk6bBsfq42gbPanFxw52SYc/n+
        o0/RXsOo0MNcGa+bDGSY3qMR8BRWKdlPUippsQxSjqRKF3XxsvqUNb0Ub0qP7EWq
        XUUPmWR2VZFmse4tQZE6OlZO5pG5DoW8rcgn/+p0q46PJM/d5F/rxP4OaSePmL9Z
        sNUy6Ahv1Ez0ponaky3diZK5sytB6NOJhIUyvOaKrNHTR0ZP4W/Brk6XylG+wB6f
        d1xm0pA9/1RxpkrDGtIIBA5bEWw0avtDUR0CAwEAAaNFMEMwDgYDVR0PAQH/BAQD
        AgEGMBIGA1UdEwEB/wQIMAYBAf8CAQAwHQYDVR0OBBYEFMMEai/jNtmoqZ42xEC0
        wNmlA1TKMA0GCSqGSIb3DQEBCwUAA4ICAQCt1vjN6ecz3EH3rQDoPxnS8wlbdwlt
        QqsFp664pHXmJLFyeNEqTVTUSpRLdWkuyl0p6u9cLriovltM7IxwVX45Xuto4TG/
        g4Lha1OBDjmtamM+kRU+bMADSDfpR+89A8pTIJISIA69/xNT7L3xFuNmSHWgwlJt
        VQBe3XagAPk4azRAyMx7J5nswpltTR7UompiNceMQGHrihY33nesHHEb/YPKYP2K
        tjDFAJXyg9N61dyej+TykAtwriaLl5fBdHETsdBZlhMaHZFjqRJhAwsodz8etvoI
        BjfMu6u9Xbovfs7KLUMtFF8dbk3I8Wg0tty/uXAdZrEWwl8Wzd1KVQIUaJRkncif
        qJ7a7NWSXA1EGbx+WK16bJe5lwimbEpO5IkcgLwx+xQPJjt9pzIA522llll7kMNK
        lz0DO+XI+OX+l8d0fKBrYaY/Ei3rkxcAKJkOjTbgnKY2BkYMA3Ly6ZngT95zY9Qc
        aLszk3Xqmy+uNvJQZvOQ54mcp0+i4eg3ZO4YzuvIy7GxktYidn2/Hm2A+x9mzh90
        4vxnDL2r4mo7r2IWMJMd0vw+XJO870lEkRMZ9rp2wfgU2ijZn2UFnVRsl/doJ/vF
        lkuVVwDnmyBJzPfLEZrkl/z3k0lfxieilrLeO68/TJ4+MbNOdyfUFqjF1+0y2ebk
        Yn0I5jH3LNeorA==
        -----END CERTIFICATE-----
      server_cert: |+
        -----BEGIN CERTIFICATE-----
        MIIEMDCCAhigAwIBAgIRANoqBUcDT6LcsutWb/sz3fYwDQYJKoZIhvcNAQELBQAw
        EzERMA8GA1UEAxMIY29uc3VsQ0EwHhcNMTcwMTEyMTk1NTA2WhcNMTkwMTEyMTk1
        NTA2WjAhMR8wHQYDVQQDExZzZXJ2ZXIuZGMxLmNmLmludGVybmFsMIIBIjANBgkq
        hkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0pS/oJQJm0WKU46sYGkJTIvptr3dSM53
        qCmkampDRc04PrfWHFOsl7mTCxUDoYtNalzk5P7iAS+NYZjeL3ceB+K7XoP35l2G
        Y858A5uNjQbjzsYpI1+vZvLi783L8NpnMecnBteZCn0X2pqEWP12detiyqiHdTSs
        e/AWnXzu1nt4VCEKU7UsCDF2zUVYtlZzYqY5MFoDFm35egNJQ0Gzw4hvRMZGObiK
        +3El+CATbgSKX6jIhiyiqeVEIc2zw/LNtHs1k3+vSMB0NONTG+ledat8zYVP4nQK
        3Mw5CwZ9ATMgSftKpC9qvJIcLOiGaoFqXFqmdprreYTNAEB47le5kQIDAQABo3Ew
        bzAOBgNVHQ8BAf8EBAMCA7gwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMC
        MB0GA1UdDgQWBBRgkvQ3qtsnyyjF+NP9+DiM2pL6EzAfBgNVHSMEGDAWgBTDBGov
        4zbZqKmeNsRAtMDZpQNUyjANBgkqhkiG9w0BAQsFAAOCAgEAuMuzWC6kQQNyWI3P
        EsCcON8GZ+dzmEfn864zwWTtxtiMi2Xey7RHd8lS7t8kIinWYLaVuaKK4a2uvB+T
        YPtLrpCu7S3E5MM/ADlW1qxZRM+nez1Ud5kHL6zg4HFOgHMvbF8GMq5+Gf1MDeS8
        ritAf0DnaKmJIKXbYpWt2BK0GaNQLxk7A6nyFCx8Vzq8Xq/RePCy2QugjauPlrMG
        tDrusUVL99aZEi0ZFlmNxplpiDlNI0exH2Q0vUvlwvo3J68GdR4GQIT7QQwJBqRR
        FT67NmkyMEn1vq6Y2aR+hPZlzlvJJzaP/PUA5x2qeK6vKG8CN3i91b18pAoD5sAi
        LuFPiajcu3ZipN+TfKW2zeGXZ1FeN6uxIvZdQU/++6VmBAii3TRtQoJ1NnHGbW1Z
        nNz69Z7oTpGa/0ErEeNIRo5drL29aMMBsYk/x94+z9WydgCbxkcCCR9fQKrrUUYU
        jPPALsXSnnaYKcqEnjEb1lSogCOLLihco2zOoeO7pCtIGwcLzlxU9Fs3t/P+TW90
        3az0e40QNFzz1Dxu7iGBdKT84icyR3vIVZV/a0DHyP2BMibGaA5EIlvlpAwCxz8Q
        rxBY+Mdvwdux9fncq6j8xr3nXO2s+SxW7cOl23Avq59Snpcr73ipt3+fJSMOPJdc
        Wm4+9+jefDizyfRIwml1Ok5UvDE=
        -----END CERTIFICATE-----
      server_key: |+
        -----BEGIN RSA PRIVATE KEY-----
        MIIEpQIBAAKCAQEA0pS/oJQJm0WKU46sYGkJTIvptr3dSM53qCmkampDRc04PrfW
        HFOsl7mTCxUDoYtNalzk5P7iAS+NYZjeL3ceB+K7XoP35l2GY858A5uNjQbjzsYp
        I1+vZvLi783L8NpnMecnBteZCn0X2pqEWP12detiyqiHdTSse/AWnXzu1nt4VCEK
        U7UsCDF2zUVYtlZzYqY5MFoDFm35egNJQ0Gzw4hvRMZGObiK+3El+CATbgSKX6jI
        hiyiqeVEIc2zw/LNtHs1k3+vSMB0NONTG+ledat8zYVP4nQK3Mw5CwZ9ATMgSftK
        pC9qvJIcLOiGaoFqXFqmdprreYTNAEB47le5kQIDAQABAoIBAQCdsuWa5KIZFMfV
        cVgnzyE2oOTChIdN+cjkN2M4iiGdCWWgml2O0x7CdSfoObGBbefoym5kC3jG+IyB
        VVC27RahQyucSWoBq3J0FfMLZJdp0IoTlJTEN+kMSMKoYU7kLTrwxTGVzyl+EFYn
        0GVim1X2UvOl3vWqUWsGWbMl96SJG4ttbywpTHhXyaZvijqiU24OgCKa3RjBMb/b
        mpSt5CZyU2UK7QnASxgA2DEO6scwvK4WZVINAy3XcpQy+KjHNPPl9HPTcQvjNo2k
        arCp+FJqi/HCVBk4Ww7rVUaLMi2j40IYrl2AW9GkAULEAFKD7KoIsjBNz1+UeT9/
        OG2UD6b9AoGBAPfnG8cZMY03HAvk8scQkCyopLlIM7o/ksthcnMT4nRChv9qAy8o
        Eg2VMiAF78tRB03QeBAJMh7dL+R8p4icQv3nSOk/x9Q/ve7oeb3pU5VT9+SEkhYA
        1bwxBfMMf+mjBCjYsCGpEa7AtZ5A1jL/7xK3DREHd5eMWVuOnP4as+x3AoGBANl1
        kQuUuml1yWGzPCo7DEd/ucB8tH9rOEAUQ/M9ykprMAUTEblll/IkDELGgTlpYR7F
        GKwVJrt5ClXAp3tSWGPBkr6fUMS8L2Piaiv7CK737eYLfAvWLEBFxzvbYD+uRuM5
        sBHFU4XE0Qcz6tu4u9IjZTpsiW+P4YliHa2ORnQ3AoGBAPL1d83rrRq/lic6HY5n
        d0WtirNkRf4VbGMTgD20kU5sHS6Z0cEXvom9XUDxUJCtO0FSPTlKKesB0HxYh0Fm
        FGoPkO+46LnmNtm80gQEdzx07RDztNEHxHIKgdAwwfRTJjJ6HDUBJClnCRiuZr/Z
        AZAQAyhbbyQCE1meLdMEjK4FAoGAD8OtGyjSBrkqOzHyJ6GWN0y0G5cuwpn0Pvj5
        IBYXpyN0HLoQK9+Ij147oU+gqJfSGZfyPO9fmnGg5SyNN6x1ie3LhJQqF8kIqnYM
        elm9fGmuzmGAwZ7qIFKuqdEyfgtVSj2xXOhwMJ9fA+WonfsbapV0TjL2F6dXk00Q
        l7dbtisCgYEAljP9kQY+g4B7DZ3H/ljy2shkcDssu2MqDFOQ3dM+MCv8vdjUmEUG
        3ICL7s3crk94rPTgWwfGf4TWHSpokPJDjANil6H4r0TfXiokL4E2PSy924yUwX6R
        yZyRpy84+Dn7BIcSHK/fJHl6TMiyObmBrYglUHDL7wmbq3y37LVTaNE=
        -----END RSA PRIVATE KEY-----
- name: etcd
  instances: 3
  azs:
  - z1
  - z2
  jobs:
  - name: consul_agent
    release: consul
    consumes:
      consul_common: { from: common_link }
      consul_server: nil
      consul_client: { from: client_link }
  - name: etcd
    release: etcd
    consumes:
      etcd: { from: etcd_server }
    provides:
      etcd: { as: etcd_server }
  vm_type: large
  stemcell: default
  persistent_disk_type: 1GB
  networks:
  - name: private
  properties:
    consul:
      agent:
        services:
          etcd: {}
    etcd:
      enable_debug_logging: true
      advertise_urls_dns_suffix: etcd.service.cf.internal
      cluster:
      - instances: 3
        name: etcd
      heartbeat_interval_in_milliseconds: 50
      peer_require_ssl: true
      require_ssl: true
      ca_cert: |+
        -----BEGIN CERTIFICATE-----
        MIIE5DCCAsygAwIBAgIBATANBgkqhkiG9w0BAQsFADASMRAwDgYDVQQDDAdldGNk
        X2NhMB4XDTE3MDgxNjIxMzEzMloXDTI3MDgxNjIxMzEzNVowEjEQMA4GA1UEAwwH
        ZXRjZF9jYTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAOYipx01n7nd
        7+FfO3Hz9OHRY51kvmVFeJjto9PMrpxz++P/l9Zr7pNu3y4nja94Tmhuz4SbBhT7
        nCGeK/Sxad0l2rFS3LEFeiRKGfMJnAMGsv0vVb0UIE/MqXner3Kt4nqCDAxejOAM
        W7D4DUQQZsrVl5uV9EYbrxyHKJV2wMQ1UwGJZjhmuv6Gtm7LvHHjyDKG+OLBXvC/
        /CLXV1mp9/Y15DL6AIUdd0HifX5eTWkEfXEJON8+HCu1VFNPbxZFBQVZApz5/Zyl
        jVUldfNPGCYRKwoCYE4anatdUgeV5qc9ZmpDwjDKR4bRV8UDvtGxVdc7poiU7Fy+
        imuIyw87YtPXsViRJ/iMQ39QGM319VB7cSG6WNxe3ejfxvc4u1hCi+9mWoLCkDou
        rYadQ5XcHDGPWnIpZoMaP//quRLucp2XD1Nmf7RtKoN/FnYRTqWNe5bB9F+BAwvq
        /IpGe0wMH6Uy3cTBKLS7ZejTj2V5P+yAd6z9yy4aIW57TzdIxAf6zCjxRzj6ht2e
        NN8AYbGxrxHXT/E8y33VwgONnG/l4Cx1EoCsKJIMUx34jP7ccNcHKyeLvD1Gw7rn
        nmJlmJJbh4NbKw7lbb+hZWfPKl9cPCiF5zfwsRkX4yG6vrLYYDDMXzIa9gx9awCH
        qWu0bwkJOq+4syk/sJr6YNeSZy3w8Qa/AgMBAAGjRTBDMA4GA1UdDwEB/wQEAwIB
        BjASBgNVHRMBAf8ECDAGAQH/AgEAMB0GA1UdDgQWBBRCgJlbHB4Z2BE4ML06uhr6
        /ArXajANBgkqhkiG9w0BAQsFAAOCAgEAW+7ievC5drPjY8bXYln5/gqYv+p8lj0p
        3Bjk5YnVuNpuKoKcM9FP5bUi0qCbJYN01A779M9WXW3pPGy8uk4NUGZaFnbG3ztj
        c5HlDrzA1YTeDAKe316Qj2wODVZ9cnEMoE+MpMicSGLHiumBUZ92ENxoR6sgNZR5
        TGIWfH7U5DIfEKolpwkUjbK0Ww0OoJXTk3LFstYnj0TuHbV7ULmUTI7i4FL4LRGp
        ZB3RDx/1gGNklf3ObgdgcZyBuoQ4WX2nYT4N+4LPAAta2uslWjOSoXcQ1pMJdqKM
        nUo8MWsLELYXpnq+fAyKyUfzsffgjbfZtQ/BLNwm7mQ3efahGLDRMMHaObqqz3wF
        FuCnOxXExhF7gkJMp0PPUygySxE5XysSXhPaFrQvFtmAPJuUVYUjCxP5Y1VNWYU5
        B2sxlGkkKEf3zm/CfAgVd3+xfDKc0AIDCsCAIkzKiQZi3XyBPKZVHokuO4ox4FnA
        RUnf/K8mOQq9MIsYvH23wuZwxpppU05Vdog+NvyXu0f6VBEMDx7rqYrkVPlukvQ/
        1YyDU8TdoCnRCOGgkxkWzMTKf8dEppGPa1BI2tilk9Sbi5gKRCQnnXY1A7t0W8vW
        pkRVZ5VuDrAaXBGbriTzKuATDKlPsChjWmiLSIzCUNuvQEURuDMbHgB4Ox32Yaes
        c2dyAUyfqrA=
        -----END CERTIFICATE-----
      client_cert: |+
        -----BEGIN CERTIFICATE-----
        MIIEKjCCAhKgAwIBAgIRAMfSzI1sYKr9nVfL76Da5EYwDQYJKoZIhvcNAQELBQAw
        EjEQMA4GA1UEAwwHZXRjZF9jYTAeFw0xNzA4MTYyMTQwMDdaFw0yNzA4MTYyMTMx
        MzRaMBwxGjAYBgNVBAMTEWRpZWdvIGV0Y2QgY2xpZW50MIIBIjANBgkqhkiG9w0B
        AQEFAAOCAQ8AMIIBCgKCAQEAvzImI9ywonQoh5TKYRzb5zRNKi8IwoaO4jBnPsv/
        mKwXpuyo7ifgwyjkHekpHMGnT21941Fk4dqo7TU4gp/mVGPB6+7nEXqxOBVUYSXF
        QX857408iGmp4mybiC8LmhfyHKIlMwufTSeGsD7al2Z1CAB1o0+aOdLy9OrOqhAp
        gQJ/QcODM2mHmoU0iH2GPlo//lcr7tKqPj3jnP3Z7MjLmepCws/4NiPtTRInadXY
        MBOVqOnWD3WqVXoEhFrRmdcEvJb6vl7VQEVn2N/xOy6YnkmbFmaGciM955V+yIFo
        2XeTiLlWHdeijlOL2A6MI8Oxc9oojlVOkmWE5iwyf6fyaQIDAQABo3EwbzAOBgNV
        HQ8BAf8EBAMCA7gwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB0GA1Ud
        DgQWBBTt+hUm3mhiSNe7LekxL2JCMvoxSTAfBgNVHSMEGDAWgBRCgJlbHB4Z2BE4
        ML06uhr6/ArXajANBgkqhkiG9w0BAQsFAAOCAgEA0wELf4HkUVVt0DFlWQp2jcPO
        svl2zqz0Z6kdYmm+z35U7asfEIqXRHjuGFSUY//MP4zxxztT3IErnvySDbL9oCjj
        oUL9Bfj6rBAipnBnHwYsHJrIxxWHWkBHzyLoP8y6BTOD944f4AAeq1DVeySbnfm/
        EQ9Y4MSoRHBlJUFmMjvo0Nl1DgyaLgGOBywHCJdbTq9vDUYe1rSzX6uIL4R8KgbQ
        Ma3x9N4dquXx+wPmZpTbQgGJCoS2rNiPhaX/d3MANDPH9tmwn6dSPmShMXqgcE5Z
        kvyyGBb8lC8lJZw9e+h7mfEfm6qv2wNhrsKFT7wGqaYHxB9mg2SXCMZC2m/oQAOs
        86uhFCs05bLZh1n0+q7b/4ac5nZNvle6ejKe+qAZx1/bVhLu46sB+Gq7u/UajO/g
        4jsYxZjQmBLrC2a7/HbZIsRHavHXINCAc2eQ5ugQRdhmOEC+2m89Bbdl/nIXasKp
        70Brys8slGd1x3pEBTJaWVivXxnCPBtZ+J3pPkbLDUAAIM5R5FIjPf1VdIWgQufU
        1v2oYfDGJ1j/eBI7l9RIDWeDSLMMozX4vXbATLDMXfeY5Guq0DUl6QH5s4j5Nhm2
        4Ezk1u8brza7CXqrBmNJ/negFY/0kSSqDyX4iUk80/KNSyx99G0LFoyrXEA4g4Ng
        aIxT2iFA+lWmv0kkpI4=
        -----END CERTIFICATE-----
      client_key: |+
        -----BEGIN RSA PRIVATE KEY-----
        MIIEpAIBAAKCAQEAvzImI9ywonQoh5TKYRzb5zRNKi8IwoaO4jBnPsv/mKwXpuyo
        7ifgwyjkHekpHMGnT21941Fk4dqo7TU4gp/mVGPB6+7nEXqxOBVUYSXFQX857408
        iGmp4mybiC8LmhfyHKIlMwufTSeGsD7al2Z1CAB1o0+aOdLy9OrOqhApgQJ/QcOD
        M2mHmoU0iH2GPlo//lcr7tKqPj3jnP3Z7MjLmepCws/4NiPtTRInadXYMBOVqOnW
        D3WqVXoEhFrRmdcEvJb6vl7VQEVn2N/xOy6YnkmbFmaGciM955V+yIFo2XeTiLlW
        HdeijlOL2A6MI8Oxc9oojlVOkmWE5iwyf6fyaQIDAQABAoIBAQCAztBTKML3Lzp+
        3QbbFg4wXVP/L2C/bNemGuXzsIup14a3toi4qbUKxempHQPNk8mcAS/mjVdhsWZN
        KKXBmugZwntK79BMPlRWbEhEiqWx0ny9nnFBla4WFQVTYh522dsK31IgaZwQ5qge
        5Llvdl8x3N1kAKTuf+eeiPJiMDFF77P1rUJpMCTOgLTIpsc6vJblyy5kHHjTeJ1r
        uABeVet0vJ0g4Pv+AzUxp/Tc2GIv4YnP0a+8WD7RsPrSV0c7WJyFWbom2tW6c41C
        xAIAzjut1MWiRuicIcoX9i7ubK2j06SkjygfkAjfRTxufYgAvyBj/ooQi11m76w4
        C10GH7gBAoGBAMGCcRJ/tmSYjbMtP9pjDbCxpyvSuHggBmOqEsq8ocuhCdPqBgHv
        b9le/UOwbnP8sw3kgKzLv7s/WlkkpjMwDERFenb6MgRMLZxV4CLorpH7CNoYgup3
        apL1t3WbruIPJBmhev6twgFZYOGrvL6WlwNyR9Px4EV1+eFJopLmhMdJAoGBAPzw
        b9kT40vVoxCtyTLgDfoEo6NALGWDX+VqJKVkn7rvgag9aM2aJ7KV63nkR5+Gne8+
        qXQkQvkYKnN13ZFdITbx5ziSj3N05bNa0oGenesbktXqwIZ8QufuYwdW33HGYDq5
        t7Ra1tJVMlSGP8SQDjKyNnPhTseaDHA34lQ2rjIhAoGAMA1GKsPP9Pb06PNpkb9b
        HO9ghb9T03CQZZtMA1AIFVqt6BOK3lwouB+gYHilVOQBSofddAs8VzEKLGyvYLKj
        uShPms/SL8MC6HliqQiCoPlnX1EK4VI6ArhFkEzShowf+MVil29qZ49cQW219tXK
        Ni7gqz665ETBgjIYzsWzXxECgYAralW4a/p6vMvFhB7h1aVwgbVYwx4buoYOSb7K
        iNAF9TBLIWdIyyn/NE572JwWnLOlKhtJ7SN1wBkhQlzqo5Kc7L6kbjujNLBsra0u
        RHyUq2Hzx9yN+Ow/BSMIUnf9/m/sBI6srV7sMWV3LqfpZFSbjQ1drJGqHx39cQov
        LEeQIQKBgQCiBLAFXleQDBCKCVaoL4qB3yfe+d6AFnpAMiDCwSIAb4Kpzgjtw7B8
        nQmewWKpanRER2R+AmLGO5duyvURajTVKhz+JyIhfn9f/hFxJbtCxLdzPC5mk5C1
        Q1Gr7kyCrleceBvDJkZuXhL9qV9eid5cEr1sb/7xFEEv6KB+4okdKw==
        -----END RSA PRIVATE KEY-----
      peer_ca_cert: |+
        -----BEGIN CERTIFICATE-----
        MIIE5DCCAsygAwIBAgIBATANBgkqhkiG9w0BAQsFADASMRAwDgYDVQQDDAdwZWVy
        X2NhMB4XDTE3MDgxNjIxNDE0NVoXDTI3MDgxNjIxNDE0NlowEjEQMA4GA1UEAwwH
        cGVlcl9jYTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAL3Ek6779qKz
        OgEwUWP6TFSrgP5KUPrzPOi4Ou1MeptiqdLk5skufsKAmLAJ9TjrbOflMOj2aszg
        PhfdW3nXFpJuQjnYg7iWwJ6XNR91tHhwNt4kxXAgktgsPV++0vRQ/pC3BVN7441C
        e8TRdHK86OJqaRptBeNkizDEt4srH/uUUzNJaW3lALtTGTfrbXxxfrwP2+pnCIZW
        Z0tc/IAIq0TWqjR07dAxwPF2E4Xm4NBDeoU4KIrPSWl+8mO2buZZp2nabIMfGHNV
        OF5qkIgfPYPv4O0CkuTumPRvJMZDG/b2tDWZnTOCWHLyRRstoWXWx/S50Ob7lZ3J
        FQJ6jdm0FyO4w9cR99MkruL+GJLqZ77AGSt2Hp8sQI2lZ7haFEYuA0/GZ8BaWn05
        0urfQKpfqt+kV5JtXypNBEpNH2vJZsuEPU53jhUzIedCF7k0ZlEqOM1/KE53uc3P
        6Tqx+Rlkz0lU7mGvk0dv8mMVcZqDt9AJ35BS+LdzBMaAQA6sFlL47LOOzRc7Zqfa
        qgrHxoAxwT/9FS2d3Noie3Nqn5x0YBiiSs11cOUI4qK+htc14GrmyQLpJkAJQ57u
        YxWd9stsRrC6xv+DfA/164HPpUk1mqYOF+8FwLrxuZpeYJ746IA6egjVVgRgGsVx
        FTVgTBQSUtnN6S/iNAN98ZsPE16s+x2bAgMBAAGjRTBDMA4GA1UdDwEB/wQEAwIB
        BjASBgNVHRMBAf8ECDAGAQH/AgEAMB0GA1UdDgQWBBQJNjhSM6PtC9iNtSmb/7a+
        3vdyhDANBgkqhkiG9w0BAQsFAAOCAgEAYufTRQ5/r8Zrm16VNe7eykYZrlm927wH
        dDxV2TxsPyqqV6+ccZKG7G1CgprZS7PYS104SKsi9TTcmfqpor3GxrLFufMVHvLx
        9Kbw+n5Z7EPbk5VD12WMoiPxmHGKy84bET0Ta5TXokIcCqxmryAeBlM6SIU93nPn
        kpvnM6VLAe6r+GsB5nyt2xcG4RPtellfS+vQ4z6drNK85jtiG38NIDXTy3PJHAsh
        n4l/uDdTp5ToBUPk/jX9jDuMfvifbQzGBd2vf2lv+SljJC0yxI74TLusHy2MHPQ5
        40TkXhPyFfbbj0kEYWBxxu/GBDnJoK0aKts2sRckA5M8u4JcBG295qKa7rj/C7Ng
        gWtPSC5iIXouOkUAQMn8VZNgYaMKzoyAOfXyHH3SHxB4Y9cQUEi2GdaFOyLW4O8H
        naVsjXbgudDE00fyySYt8Jy63L/1NVNQFirBohjKs5SYi2YWI/d9ByCw16j4zBDW
        1N3+iD1Y1B0ehFU//ni5UhRxePnPwPdq1dkCJlkep82HKJnKMAhHLbFskD8/POrL
        2PUMYm3qjW/Mn9cELB20QuUribuEV92Tzi7nABwCTZX+r51rkFSQ+mgWFILjlDR5
        eEbYez0V7k19b+xjdBacJGYwQMjjR2FN0/S63RINxzuIV4yVJkzFTe214cKnRvd6
        9fIgFfJ63rE=
        -----END CERTIFICATE-----
      peer_cert: |+
        -----BEGIN CERTIFICATE-----
        MIIEdDCCAlygAwIBAgIRAPlcOYV/t2EtlMpf+MIQ2e0wDQYJKoZIhvcNAQELBQAw
        EjEQMA4GA1UEAwwHcGVlcl9jYTAeFw0xNzA4MTcyMjQxNTBaFw0yNzA4MTYyMTQx
        NDVaMCMxITAfBgNVBAMTGGV0Y2Quc2VydmljZS5jZi5pbnRlcm5hbDCCASIwDQYJ
        KoZIhvcNAQEBBQADggEPADCCAQoCggEBAPCVjt6mU7gPLP39MruAvw3yqcGOhd4P
        AGih8YOT3ZZprX85ZCJQA/b2suMiOH3OcWIaA/zI7Abry5eK25AX0VaXYz1J9I4R
        JIxOMncHnsj2jL58WA0scXNTqt7tJZNykhXYADI1SEMsJP6JTZcQXDurbRY2dU/b
        1c2z1f2c+yuEZEGnKnPF+UDsykQrlRU4/FEdougiEpw48aa2yIB/RyA/nxZ0ai3l
        v9ss33eNJgtcWnxVNuEPra+ZKBn9jukKYMb5A14LTUaj0n5VK/RbCnxYm/iVRHFK
        Z5wbYx4F1EieCY0QroKfp7MBKwG9bpLWnhMxPDNOhFZHU6E5bvZoYTcCAwEAAaOB
        szCBsDAOBgNVHQ8BAf8EBAMCA7gwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUF
        BwMCMB0GA1UdDgQWBBRT9vdkvQnUslLo/VFCmN9wt1a4QjAfBgNVHSMEGDAWgBQJ
        NjhSM6PtC9iNtSmb/7a+3vdyhDA/BgNVHREEODA2ghoqLmV0Y2Quc2VydmljZS5j
        Zi5pbnRlcm5hbIIYZXRjZC5zZXJ2aWNlLmNmLmludGVybmFsMA0GCSqGSIb3DQEB
        CwUAA4ICAQAo9vsGy16/zZxAKrsbdKYFY+9OfCLK0FiCnPI8EgmdKzs/BrJPHbPV
        bTj3VMiSlfRveyhus/7WRWQabtv/YZm+wlWCT4YAmlv3tQvAxK7U5UtNtljtq7lw
        vImwtmDny74yB+pKA+l+o2NAJRhIWVNUPctE2xwtoaUutmwVLQClcIkF8Rt2Jhx6
        vxPuG+t5WGl1ryTWDQ1G8i+TSeLmr/BQm/B7Fuv/zwRwRSNx1f7ao6egxk8l28Rd
        qU/0+kOXJxKdjV85MRQS9Z/AdVFUOU85wGB7ZquDPAljPe6df1IhcNAkrGzwbfBW
        g+C5wKjHKXPGp6WYmAR/+5Oe6hDWKX9URj8yDl4k0Cp0z9BZ2l0mX1ehOjw2RtFE
        baPE9RbMSwE51UvHMjiqkndYUpjfpxTybyFP38U9PEvzUsXZ0IHQHkmgYzF2h0xG
        DGQBDAfjSIHYJCdMdDXLjQXg4/pOkX4olb+2eyupAoUI/8TJTIcoTPfp8KSOz/fR
        NpX3V/elQf6vNhp5Copt6H1ZgLF/tf8PXH7JcnWqm7LWVz4oNuA32P1dn6xH0Ka/
        Cl2WeSI2Ok0luZ9PUxV5a2F3fsSVahqVL8Jra+d724CU9Qz3hQUsbgdPRg2qBZv6
        HEiezWLuaYBYAdI8alkQubkUVONOuP9d5NBkNM06zbbH3257moER2A==
        -----END CERTIFICATE-----
      peer_key: |+
        -----BEGIN RSA PRIVATE KEY-----
        MIIEpAIBAAKCAQEA8JWO3qZTuA8s/f0yu4C/DfKpwY6F3g8AaKHxg5Pdlmmtfzlk
        IlAD9vay4yI4fc5xYhoD/MjsBuvLl4rbkBfRVpdjPUn0jhEkjE4ydweeyPaMvnxY
        DSxxc1Oq3u0lk3KSFdgAMjVIQywk/olNlxBcO6ttFjZ1T9vVzbPV/Zz7K4RkQacq
        c8X5QOzKRCuVFTj8UR2i6CISnDjxprbIgH9HID+fFnRqLeW/2yzfd40mC1xafFU2
        4Q+tr5koGf2O6QpgxvkDXgtNRqPSflUr9FsKfFib+JVEcUpnnBtjHgXUSJ4JjRCu
        gp+nswErAb1uktaeEzE8M06EVkdToTlu9mhhNwIDAQABAoIBAQC3q2YsQsztWuCd
        c1z02uCBFH5W36kBk3BbcS8BpbRorXsgAr+Yln/AXizJzIlWOnJDU9sxdG8FBaUj
        p4XiJtzRf7fqxXgnsZy2ZMiQKMgnYlqm3iUWwZRHWFu930xtme0/Me1MZ3MonR4N
        GOOcbYgMod4hNCgxdIJwjVfUS7FRUUGOTWVqEXQ8BcB+b9nv9iLN8SNSz9qV9QoM
        mUTsNkEMls24eik1a48Ba1/rCjDe2wOQH/J/D46VRAqwaHj6lah4Nc0CEtw3UwHg
        WMJzdwl8M6KsGUW5udyZ1z3DGmlIqniTaaO/uxCmgn921Dd/vbVtk8ipjh6/3V/0
        Y5ybulJ5AoGBAPd+LEnSNlRp5w3nNsQgaBROWrqdC7vE4ENAZOrTqDLJ3RYcwwFz
        AylHBmyklH2oNNQ9nawuWyWyxWS1CA6I/BOpQa/d3K1k3Cs2XpuypcPK7OGZ0JFg
        v2aeVSCdJVEq6b4EDdORxWlwlXWp0jWkbe1RSF7zSUmAUdsbH5rFC97DAoGBAPja
        l46+Bordkl5f+S3F1wbnkRnP7zHs/I0iHRM2ia3PtLYoqLpRFGQnjZPUTmMiMw9r
        wtda8yE56+3mpjJrBcwQtqETENGYGCG8RV2bJ4fU4N5d8zHzpspjYuMCl73YH+Qr
        dnwnS1ltACGqcP1dIGGoeprcM4s/quwqdEPkXzR9AoGASa9PEEt90XQWTpVgQNRF
        KIaLjLPlImpjOqKZaTDLCxP+tu6pQG01q7xxtTbq3t6NnAMcRn8ms/qdunYLiAhQ
        xKnH3Mx5P0agJl1xnXl60OhBzok1B5N+aNcLEUK4MYpNPT2HwE3OMK5MUVPWOhJS
        iC2DFoHod/G0bT+OEU5JUnECgYEAjM5Slwfad9RrL08qlMWup097gJlxBFTNiaXV
        wbtIJ7qwy6kx30plOU5QA2dLezgsn/sfYe8qRpCZeCDbxQddXlvOmlFJYO6oKN54
        eUCDG45ONkP+iTMOGtIlb8FVzqttUBNvlUw+jDjqrCHekN2Spu9HgDw8RfweYEad
        RpT/cZkCgYB2S3WdTv5fhgPqsBupTUSw3N1gfdc53xn9lkxQ71Iue8QJJpN/Qy+E
        +YEx8zqSF4eMZl+w5Y3UNNMpL14+T7AOMfnNU5Q2b+h6h0vsEV+WBRWFTHyhVNOU
        bb5lt5SZoR/swze5RKvTrAVRl2Ze6LgwcOrjfB1ZrOrLqPslj0qE1A==
        -----END RSA PRIVATE KEY-----
      server_cert: |+
        -----BEGIN CERTIFICATE-----
        MIIEczCCAlugAwIBAgIQFq2HkXX2eQiUIbxa4sZeWTANBgkqhkiG9w0BAQsFADAS
        MRAwDgYDVQQDDAdldGNkX2NhMB4XDTE3MDgxNzIyMzkyNFoXDTI3MDgxNjIxMzEz
        NFowIzEhMB8GA1UEAxMYZXRjZC5zZXJ2aWNlLmNmLmludGVybmFsMIIBIjANBgkq
        hkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsBk0X3I3fcME7LiTdmifScwrPS/Bqenr
        khvmIx6qLxxEcIkKiZU2EVnWK8/r2U+5iW0iwVGMa0alLvnWSlh5ibS/BEl8iyaj
        8oQXC+uXefBlniv7PHJ6wzCg2+t1NRmFwx8Bm3zqFpXiXD/rICwBQOdBVe82FlQs
        ylAiWm84fM05yh39aEV+WyG+j/qIEdpBdP3G10qtKtSXBZR5ErMZQVaxD9c82rAt
        Rge/vZLEdUpXmkymUljDvs/VNH+vyR0zAEW221GQX2mGoSpX68fmbxtOmJC+LXcr
        rj1PgJY19ZVqc0t+kC0M5tQlVmMrea+40DaMzOK+ZpazP4WE54tV5QIDAQABo4Gz
        MIGwMA4GA1UdDwEB/wQEAwIDuDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUH
        AwIwHQYDVR0OBBYEFLJDZkA6PmE+hVlKNG1XbmCxYTfSMB8GA1UdIwQYMBaAFEKA
        mVscHhnYETgwvTq6Gvr8CtdqMD8GA1UdEQQ4MDaCGiouZXRjZC5zZXJ2aWNlLmNm
        LmludGVybmFsghhldGNkLnNlcnZpY2UuY2YuaW50ZXJuYWwwDQYJKoZIhvcNAQEL
        BQADggIBAELJqJe0BLbu4HiXoEPSiCRxkDfUcwGqdv1c7DSgaMFFpVCGEHvRyJtA
        f6BDa1MW5SAef0trCq/vrsI2Xk9p2jUtecVM5mBNsejZIz0t8DdRSZcCnFboaGGo
        tcW/oUTNbNMly+9bya7VXA8hDzyB8UaRwZ/cvBvZBRVkCZcr1YwIJ9JHMHcqZwtc
        XP2PmMBiPPyS6b4UbW8mphutVo/kRPqsrDYkfG8/yqthatP42WqXWC2Lfg+dKql7
        c2BtiPanBfc6TAgy8tIklmAmXux3GcZxHc+fmULPUXEPDNDbVo9ssBAvz9jT4Aq+
        fkMhI8UjBu2BtYhiQva6hglhteSrTaem+IL3ld75Z8ZiSDHgwpsEEVUZVdL12iHv
        fmvlzHD9eDIlmRxCH1/BhC5PGwqiklPFA8oFhkfrvR+EwReqnW3IAY8XDruiYw48
        SdsyZK+ZrmCACfZ+W6JcBENgXEkjrx78UMkD8bJCY87Z4iRGuXxpf1CoWE8ZiHBa
        b0QSqap/Y51UR5O4sL0+Kb7YMD1E/SU07pP3mVbCj8pgGRzlEZH4OXy5nnotvA24
        MavIVPDBYJCY+4TPXrrFuGTwxOO0y+tnRWVhpj/aNMV4+HXfoiOq5RRPcMmj+01A
        Rmwcf2C6PXb0gNvcYH/a7cXZUxO4To/ZTk+sJRpClI4ZPO9eGIuR
        -----END CERTIFICATE-----
      server_key: |+
        -----BEGIN RSA PRIVATE KEY-----
        MIIEpAIBAAKCAQEAsBk0X3I3fcME7LiTdmifScwrPS/BqenrkhvmIx6qLxxEcIkK
        iZU2EVnWK8/r2U+5iW0iwVGMa0alLvnWSlh5ibS/BEl8iyaj8oQXC+uXefBlniv7
        PHJ6wzCg2+t1NRmFwx8Bm3zqFpXiXD/rICwBQOdBVe82FlQsylAiWm84fM05yh39
        aEV+WyG+j/qIEdpBdP3G10qtKtSXBZR5ErMZQVaxD9c82rAtRge/vZLEdUpXmkym
        UljDvs/VNH+vyR0zAEW221GQX2mGoSpX68fmbxtOmJC+LXcrrj1PgJY19ZVqc0t+
        kC0M5tQlVmMrea+40DaMzOK+ZpazP4WE54tV5QIDAQABAoIBAGbZwboBVqGwNn8f
        6pic0HPkuFhbPSxFQF4sx0Q508ICK4Lit3HV4SdPJgSewqxAoSe/wy3PuEirkSyX
        pO31ML8Z/vq9BO2s2tJ5DZDbv7PrkR5Jp8oNPuAj1b+8jM8/od1tjZ3H3lzkm7mv
        Z1959B15M7LCLP/rl+Ft3jKdhQt2RS0yok4SQVrOjduyPKntxFqYLhGuXAmQ6TRa
        q3lI0fp8WaTYVbRHztOQEQ73xtiYyMVyc8mgbC6/5H1mrRQD6SGwXZNAdfjECPRR
        tO3bMPT602WF/Kavwk92tvB/rlp8LzeEJtNOQNNs1Dj3K8qEwoRYbkCTXOgn9OJ8
        ikE6cSECgYEA45hXUwnoMwL7BDBvKoZQHxRxfQfvtr+WrsPxJ4k07QOQbLxIe/k0
        L0gGu2ewaXmh5a6flP/2SkTh1/3SpQ963zddCK9hYi+KK1cO3AsYaClWJwlgEb8x
        7dbhfgHiSktRzupIqDr3fCrUNCBrAbvqgsqzbu7/VGGEllGkM5ERSzkCgYEAxhOM
        EVqRi6ZgM3WblptTyPR4ulg+iIwnfYT3wbydWdayCNkQvwRz883DdKKNBffeSdE5
        IMi+DJKuh+p4gqvgc7P6glBOy5phWSIwPEG+BWJZjB5+RjgG/tum+fTZt0hp6z7w
        IpLyS7IESrZTFtl5tX03SHn9rM9u5YdgKm9KpA0CgYEAlmGneXe3VFVo1JjIKzn6
        IL8aSbn/uymWf716T1xKezz6pc42uqurvn4B7LwThW3X+nJKgWIrM2GWNGhDUcsL
        rgff0ghH+V9eFUr9x4kRRGnjwgFg1/kUHYn5DpBiHCLuWCDXh0kHE6Uc96Bf9BJd
        XrReoTMLxI520/f33hbBbYkCgYAGlrfeC+kzgAFLNOpMBDaxRJCPgkfyOtdFcZrc
        Mu35Aw9BBBdugzNoNLv/sTiHrksSoYcI9CR+PpLXqpD/p7/7mU0H8KvuUeBTGrQI
        DRfJDhB0fL8ujsaMy7muLtrfIeWEEb/jJogwxGcoJRB2fh1yUAv6uTQa/3ts3yfv
        wWv2MQKBgQCdg4vFRlaylKrmhtgq8prQkfXZYVrwgdush8A2U73AotJcqV2IJ0U4
        3H/06l5iPfwGc5dUxeLtIO3livRthF3ndXbfhfxwaoiFiMQW29fHyH6ZmISBlhws
        ABY9tvEeUiuVzDAB9h+cXZwMIko6rmT5/sKuIwhEGRsfWOtenh0v6g==
        -----END RSA PRIVATE KEY-----
- name: testconsumer
  instances: 1
  azs:
  - z1
  - z2
  jobs:
  - name: consul_agent
    release: consul
    consumes:
      consul_common: { from: common_link }
      consul_server: nil
      consul_client: { from: client_link }
  - name: etcd_testconsumer
    release: etcd
    consumes:
      etcd: { from: etcd_server }
  vm_type: default
  stemcell: default
  persistent_disk_type: 1GB
  networks:
  - name: private

update:
  canaries: 1
  canary_watch_time: 1000-180000
  max_in_flight: 1
  serial: true
  update_watch_time: 1000-180000
`
)
