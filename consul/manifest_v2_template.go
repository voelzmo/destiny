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
      agent_cert: |-
        -----BEGIN CERTIFICATE-----
        MIIEOTCCAiGgAwIBAgIRAL6Jrp8hF4NVai1RMWu7vUUwDQYJKoZIhvcNAQELBQAw
        EzERMA8GA1UEAxMIY29uc3VsQ0EwHhcNMTkwMzA3MTQxMTE1WhcNMjAwOTA3MTQx
        MTEzWjAXMRUwEwYDVQQDEwxjb25zdWwgYWdlbnQwggEiMA0GCSqGSIb3DQEBAQUA
        A4IBDwAwggEKAoIBAQDwwHqJ+goEaz/aIuAY9aDPr1oiipHdhrxSMQQMXyEoO1/k
        /0MUKeZ1K76x4e7T1qY4wDnQT6KOxJOKzcuAoMIeUYze8G9WUQ9lz7omT4vFE0iG
        afmHWv4IxMODhyJEqMRGuuZojDJ/J5ZrZ7MV157UOAig6rz2g7wH2nrYzPlZAFua
        otLUynLljUaRcrRFjD4iqtYBpaZ5K2PiOl9/Kf3xSVuaGZ+ElsAB45UyvHODz7A1
        DkC0a7eY4SdVALw/6AJ4+vhhPK2TwZqlcAqQKm5sEJxGkM7QFa/deQAEfrTOlrfz
        Renyxe2iUD6wyUqNd+JlXdlOmPVwDre/aqo6J4yJAgMBAAGjgYMwgYAwDgYDVR0P
        AQH/BAQDAgO4MB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAdBgNVHQ4E
        FgQUXQUyNbrMYkFSeJAZyWaf1QyIxeYwHwYDVR0jBBgwFoAUnhwsCopQWAnhZb4N
        ZAcFzDDiy7wwDwYDVR0RBAgwBocEfwAAATANBgkqhkiG9w0BAQsFAAOCAgEAniQS
        L9ZHAuVohg6Zxgef+1zcSQAm8vUVx7bzu7tUU4E53HbC9j+4JcUIJN2NzpL9Lj/H
        DjXAqY3IIQJcawr2+ap5qORlG1o2WNqEz4OQW5ZpMgQDnV0QdDVsUdUj9ZAGoGvO
        az/AaoveOsZyeGpoXdiR9UJl23kSzUKQsJ2AXJV8ofqKP0aMkd7y4NfzNuldFHL5
        +T5pi5hOnvdHRO1Iri4zO5b2/Q44sfI3ACpQk4WEwHC2206fYuTnr/P4Bf6SbYlN
        EGJYWpA1C2ZJwGuhlVPksARbtHwkxWI1lr2T3Ayz1uZh0Fa6qkd3NUP5V9IC5xrO
        G6PrKFhOKxe8gbRHWbtmihoDNZIDawDv1icHMm8+ZTsIRMoS9usI7sS2HX9ruDsA
        WjhacjJeuOaRMdLHv+D+lm9kWJUPYdOFb9LPtP/oZFWdwswX1i3AdZA6+bqANSf1
        nHLN6DhDIVU9HZ3jX1PxX5zmBZdk4Sf1T5cQHC6sRJFu0f8JKs5DQBz2fdgV+EIT
        4Twhg4jClO6KdNQUzVj+ghhdS33ola1WTqlcc/zUipSSlwKvvYi6/Qj6x2tI770o
        V/OqI1/z51z0/1i2lztctWxJtWpGJAleWQVT1sFoelkdl7hiW5GxcfM8yAgt5XFM
        4JD6zEJkoRqiWZYG+o9BXl2o6lrcaZpWy/m0X/g=
        -----END CERTIFICATE-----
      
      agent_key: |-
        -----BEGIN RSA PRIVATE KEY-----
        MIIEpQIBAAKCAQEA8MB6ifoKBGs/2iLgGPWgz69aIoqR3Ya8UjEEDF8hKDtf5P9D
        FCnmdSu+seHu09amOMA50E+ijsSTis3LgKDCHlGM3vBvVlEPZc+6Jk+LxRNIhmn5
        h1r+CMTDg4ciRKjERrrmaIwyfyeWa2ezFdee1DgIoOq89oO8B9p62Mz5WQBbmqLS
        1Mpy5Y1GkXK0RYw+IqrWAaWmeStj4jpffyn98UlbmhmfhJbAAeOVMrxzg8+wNQ5A
        tGu3mOEnVQC8P+gCePr4YTytk8GapXAKkCpubBCcRpDO0BWv3XkABH60zpa380Xp
        8sXtolA+sMlKjXfiZV3ZTpj1cA63v2qqOieMiQIDAQABAoIBAQDlrRq6KtoHoTGH
        LyJPbXV+7LC2py/FANcEKlLMYqzFLu+rDYK258o+Gf+QwBQR8IMfPPNqsa5JOWvw
        TahJDBRkzDPyM7pjxG1GGchi6BxVZLIZ2Vv+L2aOgvhqsC6XBdJKD9/j+WvrNp8f
        1AxbWBrkJpjEu2yOWNq3O01bKyDuL7lmYBBzD9hYX0InyX4KTpPwlvVHq04lmb71
        ZuX+BQAr1Mpnzr94XbwgkomSrYxtingcMoGrHPcHJ09G88WErqQgcVBHrkplQ+YN
        QgBNU0O1uw+cziKONEIbj7zCVi2PRslCwZZuhg7PtoKFX66jjNZV3dp+SY5UzIZF
        iCcVR1DpAoGBAP30onEO8qxL1Bxc7xV9uMACwS7W1zlK1IPIFEPJQfBB/XcgFqYt
        +/GDMrPN++WYlF8bIxlMlUUluebQ/e9DlpZHBV7HKnAVaTNxmqlkNtHlJn030n0t
        I3GPfUU5mjnlDrXXd/bXAI6vYLT0niAj4/J16n457W2nKCWOLYbyay9zAoGBAPKw
        ohbL/d+9fJ7c7tdPlTJxufwT3cHcaF102iPpjRidIyJgLazMWeGbZjYjRHOzd+Sk
        T+e6hUalT28rAd+Y/dlxFVikRcYlrB+PMOtZc74IZc/76NTKwifVnLOq70pJmK3Y
        LvgtVap2NcRWOspbchNiFtu05eqRkkz4MRKeIR0TAoGBAP1mx1gs5DTMDUCn8uDs
        7BacKQuF0IgYhOliEeZ5wdPs5O+jEzaKl+UrVsJXfUxh7VrhByrNYfz3YgJQ87F1
        LaOSBmfGMDBbDPgKGZuApbrW/orf0qaZDc7YsNUMXzn5t6327HtfmezGTqcBl27W
        oTNkObHuN12897BRqFgJOK7FAoGBAMuT7bQP1Fud+O7OR6/nezEAg1H4XDolqIpU
        3jPn00sFbZaFdWsRVIhSsg/Rz7b6oiTyzHCHXwse2p5XRlAJZ0/Cc3STAFCyA0vJ
        8vBJbjTHmJg6KVpu5yVBJBz205nOWLvjr6rRZJ4EYR/ccZ6TzQKDcsdEXOVCzaWx
        QIDskrxvAoGAWLtqWDHf49A8QBF8bIKYcmsCLTQmnbmIJAO7o5nKDgy8Za2Jq/FU
        n0SKRCoQYCtVKbT7rIGoVcdAnR0dVRlLsovk/A3PvTl1/KMfZUzBHf9vZTDzaD4i
        A7T68ZkZEUEwe9wAmEeXmozffRNykO6wCmBnmluSduQxRuXyrRHf2bw=
        -----END RSA PRIVATE KEY-----
      
      ca_cert: |-
        -----BEGIN CERTIFICATE-----
        MIIE5jCCAs6gAwIBAgIBATANBgkqhkiG9w0BAQsFADATMREwDwYDVQQDEwhjb25z
        dWxDQTAeFw0xOTAzMDcxNDExMTNaFw0yMDA5MDcxNDExMTNaMBMxETAPBgNVBAMT
        CGNvbnN1bENBMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA2/krPZyx
        8NRnag40A19HftaRRX8QU6RZOtWay5CPbVcXWQ3Ib2uOQ6UV2aQCeJKgIwZQLlL2
        DuMv5C6vYwfbbJO6FaSINZ4Qc02NTUxRO7XfV9ZJpBxZNEQsjtQSUFNm6YHVixCX
        r0pG1yNAnTY/0AphBZb51OZd4POogSYmWa8YLPX+hw5oPzF3uQgNXUIXLct1NMww
        yq+/avZBcnRgUXaAgnYYspGEkY4XtVpsAPkaV0y/2nXolDmJqH0Aba9EbptdAFkz
        +5+OXKk3N1AUMWSzYJ2XaOcPJhfjnEfApYYTiwx4ikaSftYv88QUWizFYvTPk0Ie
        7QRqt7QYrGOKs+wiatk5ExnsvCqns8A+515XGplVHXacriYJDtCXS8hBrAZRMPU1
        15vxSBWabsKb0sGawxY4xTXX3E20YGfqDvd7v593N2atbQLHZBZ0A8k8ntYlZidD
        R1GHFgPm8X+Vjq/0SHpZ6VK8sPa+oJm202F65a6GFwgkaGe+nNfKLJmk0ltk+lTP
        Wg/8jhzbUOF/l+SWPkACng/OtW2CYviveLkyGKW33FAA0JNrixVELF5BRo3aIfUa
        4VzfwK81PXsLEjG1xKuN45AuicweuuzcCSotV15ekgboRHdsq90f/KtbH5sJwcIZ
        SiPuQgJxudwVfMzV/T0WKyO4wA5BBPIT/FECAwEAAaNFMEMwDgYDVR0PAQH/BAQD
        AgEGMBIGA1UdEwEB/wQIMAYBAf8CAQAwHQYDVR0OBBYEFJ4cLAqKUFgJ4WW+DWQH
        Bcww4su8MA0GCSqGSIb3DQEBCwUAA4ICAQAC1nDqOEioH1Nzc9EikAxtpzRsf6Xm
        ZDjKyxPCAYSM8zDaDlyEKVsA8/LR2jJdshsTUMyT+eXM+43fWf3nlBklChK/QTVc
        8jQn4xccV+jtI6/xDEnxq+PSYvQ2BWdsaLly1LWnh9MCnO9mF966x6HjOT1iV0Ex
        MT5UpQpBHxpenoN3BKfQZDEsS39maDV8orhu+wVqb2NRCsVAA6ZggqfFY5dgtgQk
        DoqLwejsPI+pmiokFYHylqKQ5ETktCIMB+7wp0reAwAWW0n4ZHUDzrk5KBzsBKvF
        HaBzbMmklMnV1DCEIr7cYJUnpozWQ61PLMo6yHf57RJfz9ZNHFJDW8WVfOEjTj7v
        LMJPfqBB2FEkMfnzBFEjTrXS+25CCXNqpZjmQZrRaN9UkH2echSdMwmk2mT4gPik
        fNL7hV1+oWv58nsWWOrRIWedzbIUtgL09We5BkaN4NAAt0ErLMPP9E7y8jlAmyZ6
        iFDmQXATKhi+zW6nNbjtsfeWhPffFfBRC1Ns3/UaSjWJ++eKi1EGjJvny37Rbbpn
        ++PhN5qwJ7ANiBFDbtNUmshdsO8ytBjbJ5ZbGU7dREPq3reBYtGKwTAtADNErHf8
        grSNxwCbgkAr/wr5YrOtyVu39oHPhrNODNwPUCbOKYxNM5EoHgg2Pb/DOjq/eW+R
        K7EsthmscNIdjA==
        -----END CERTIFICATE-----
      
      encrypt_keys:
      - Twas brillig, and the slithy toves Did gyre and gimble in the wabe; All mimsy were the borogoves, And the mome raths outgrabe.
      server_cert: |-
        -----BEGIN CERTIFICATE-----
        MIIEMDCCAhigAwIBAgIRAMAtqd3SiOhLe6V9DM1juLUwDQYJKoZIhvcNAQELBQAw
        EzERMA8GA1UEAxMIY29uc3VsQ0EwHhcNMTkwMzA3MTQxMTE1WhcNMjAwOTA3MTQx
        MTEzWjAhMR8wHQYDVQQDExZzZXJ2ZXIuZGMxLmNmLmludGVybmFsMIIBIjANBgkq
        hkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyvRHJiKAvgs0W3u2e+dLFtRT4hhuHkT6
        wofwbib4rbYksgDPqNHVMtisd6pbQPG4gLC6bm9WkPmOykEOXLggqA1KsHQKp3qU
        sYBW4tUkcCYetbGTISWOSl4jineUJKOt6+GvAkYptX/wT5ujCRsMDZ5FK8SF2f6W
        pWanuoFvKd+eXipYnd0yHWOozGVf5Qh2z1d1PZY/ypXo7LnfenB+DEthbdnElrB1
        qUMwiDnriI0TvkcQXjn45bsP4+QwQNbahYIJz8bDRs346XQn6DstvdmL6J7B+XKg
        cXQmLMr6COq3Ju2H9WLNEmwRbzOjxRoCIKu6ZPgLIJRmAGB/8H4PowIDAQABo3Ew
        bzAOBgNVHQ8BAf8EBAMCA7gwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMC
        MB0GA1UdDgQWBBQMRrEOdswphbISbH0bPwIcSNcuzjAfBgNVHSMEGDAWgBSeHCwK
        ilBYCeFlvg1kBwXMMOLLvDANBgkqhkiG9w0BAQsFAAOCAgEAI2jP1cWfUPXwVM+Q
        /zIJJXB0m0eMFVb2G0/hP+HLwWP/zpkQhhzEtjsnfPSSP6u3GhJTW8IpqtkMGsIX
        0KKxKcXQ0ydBE1TpVc1FUbbXOpMqWYg/wyQPV6m9zgM0jEcbMVJBpuEXciSD8tVL
        hiG9rJdJdZQ5LF3myR6tLCWzl5P9LfL2k4e2ZAh/cUAZv2snQqqfvkksW/KGRKVp
        5UnoIZ3+MsiYo7PO4qK3rQ1jusGCg+76lR5KwR0KC/LstQ7fpAy37rpxyfE93S/P
        KywewTWuSSdK3ExHPB8dL52W9A9CBVRQI3MDn/TJLLmEwzUMlIfjglYAGlUgT5Px
        DVZcIITxlz7MWpYwzxVC6oEbP24qbn76JS1SOhIrX93/P3sUFQyAclh13HqAQGxR
        Wmc0QkWJvMfQwseHMhomQ4dd8lrdKp/rrtNGVS9oJGw/LUehqLUaBtD7H39vz3qM
        n8A51RXcjYuAZ0YG+hivHOl7kKMkEDl4z7cUFVuAwGrFWaLO9ZJ7MBZFUW23EFy4
        E170fphcHnCDs97zM0qNdqJFuv2NfC1b9e8xEJJf6wQc3ovLZhBMxDj/zYiTL29D
        3QRuBZqNfpkQ0knGSdSki4QjUN9APeU1oKG6rKHxfYRzZ8ltJIFgztZP3jOd7BrK
        ZVqWQ5gSywymgWoQRMbEDdkmCzw=
        -----END CERTIFICATE-----
      server_key: |-
        -----BEGIN RSA PRIVATE KEY-----
        MIIEpQIBAAKCAQEAyvRHJiKAvgs0W3u2e+dLFtRT4hhuHkT6wofwbib4rbYksgDP
        qNHVMtisd6pbQPG4gLC6bm9WkPmOykEOXLggqA1KsHQKp3qUsYBW4tUkcCYetbGT
        ISWOSl4jineUJKOt6+GvAkYptX/wT5ujCRsMDZ5FK8SF2f6WpWanuoFvKd+eXipY
        nd0yHWOozGVf5Qh2z1d1PZY/ypXo7LnfenB+DEthbdnElrB1qUMwiDnriI0TvkcQ
        Xjn45bsP4+QwQNbahYIJz8bDRs346XQn6DstvdmL6J7B+XKgcXQmLMr6COq3Ju2H
        9WLNEmwRbzOjxRoCIKu6ZPgLIJRmAGB/8H4PowIDAQABAoIBAQCwRn9NlgYwCldN
        ZiRXdcr4Zy78X6+1RsBuGdgwNFUlu+AfYyj6PlZotV0HCsX7oHdZ/yXOKZIMlVT6
        QosQ6TQkAnduzzs5v/RKP2g32Fyvs8xUj4l07sOpwB4qdDYNpMS47eotlXdAl4DI
        BCwVCpLreR4nJ6gCcWey/XiNO0KsIGmpX8ohH4qzhcEbFLW9gAR9GTfrDOAaYkSv
        KziN+PZ2dBTl0jauwGCikMN5QWk7+mvXUq5Sx8IBtWIF51a1xqVH3+FhGUX3ENEk
        frfQyDaH6RU0M/bwO/cA8xebeY8lr+6C0TWv/RghgX4o9Ei3wKn3XG/UBH3Cznne
        2ttdwEPhAoGBAM0foGw5uBhrOqlDkva7LRxHTEJQQUYNcT+6mCA9ZUNUTPLTIHLE
        7wTJWPpcu0gbCJAxXf6SYHoEji9M5swSiDP4dAePRO0tS8VckXL4p6NLk97CELB6
        aoiDHtpv1mungvmMcPLesg3w2G2VmI3WTjKt4ngsDP/DaegBS0eA9fg7AoGBAP1K
        6LZkOjPwGJbCeqwTGec9kYNgAvldSMnAGwnFz7ofigZUGduv0Ej4Ov5P9Jo/NIPn
        AjYjNyL26M5en1t4DJjnUipdPhU9yjsClSump9hXKeeOjBNLVaryQSUjlt0o95xr
        HMt7zAuQ43zyFmoUItzp1syoobYIpCAVyEvVqze5AoGBAMYxp9TKVFmrygtgUnD5
        3BV1wnZUiy0/scwM5A5KpDxRCOSbIMAkDnqGfeWykfaSwExqltJx5qwfGK8VU++c
        fGQSzTG8ubGdUZgJ4DPBlGCQlvjmdC/AqIzsfHQ9GWX9fezXSQ8yI8KaktQXdkad
        6gLHxomsroa17u+PyIf3UDKfAoGBAN8wZqZw4qhpZAFUFOwTWLveEJ6Gt5grjrvX
        vvt5hnUm3WR+LtrZrNrfgHwe0BYqo4emwtgZZ7gzgSh3UEw1GESTcF9MEix9afld
        aTwxeay0AYS8oslNlIsxNB4ZohH2y1jVOWZEC6QVY57xYrbOT3oBwvhLj1LrglOT
        Xg8Uk+5hAoGAVsJqjSCEnmGnZdaCyjh/ox3hQo0+HkMTcii6yA5oWULjPx/DqKwu
        7uakOYtotJlkWMW2YWhd3ym6MubZRU4Z252RmU80D/rsJYiYOpDfqW50XHJ8QXyO
        WE3wgii7CIdTqwROuJLGCqdoZDedcYsRKCTxJanm6wSTZPoJ0EnGU7Q=
        -----END RSA PRIVATE KEY-----
      
- name: testconsumer
  instances: 1
  azs: []
  jobs:
  - name: consul_agent
    release: consul
    consumes:
      consul_common: { from: common_link }
      consul_server: nil
      consul_client: { from: client_link }
  - name: consul-test-consumer
    release: consul
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
