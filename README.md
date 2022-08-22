# speedtest-clash

test your proxy servers speed over www.speedtest.net

## How

```bash
❯ ./dist/speedtest-clash -h
use speedtest-clash to test your proxy server speed over www.speedtest.net

Usage:
  speedtest-clash [flags]

Flags:
      --concurrency int     tcp concurrency, note that more concurrency will use more traffic (default 1)
      --config string       clash config file path (also support download from http such as your clash subscribe link)
      --debug               enable debug log file
      --duration duration   time duration for speed test (default 3s)
  -h, --help                help for speedtest-clash
```

* speed test by clash config yaml

```bash
 ./dist/speedtest-clash --config test/proxies.yaml
```

```yaml
# proxies.yaml
proxies:
  # Shadowsocks
  # The supported ciphers (encryption methods):
  #   aes-128-gcm aes-192-gcm aes-256-gcm
  #   aes-128-cfb aes-192-cfb aes-256-cfb
  #   aes-128-ctr aes-192-ctr aes-256-ctr
  #   rc4-md5 chacha20-ietf xchacha20
  #   chacha20-ietf-poly1305 xchacha20-ietf-poly1305
  - name: "ss1"
    type: ss
    server: server
    port: 443
    cipher: chacha20-ietf-poly1305
    password: "password"
    # udp: true

  - name: "ss2"
    type: ss
    server: server
    port: 443
    cipher: chacha20-ietf-poly1305
    password: "password"
    plugin: obfs
    plugin-opts:
      mode: tls # or http
      # host: bing.com
```

see more yaml config in clash [wiki](https://github.com/Dreamacro/clash/wiki/configuration)

* speed test by clash config from http

```bash
 ./dist/speedtest-clash --config https://www.youconfig.com
```

![Example](./docs/example.svg)

## Install

you can download pre-build binary from github release [page](https://github.com/Ehco1996/speedtest-clash/releases)
