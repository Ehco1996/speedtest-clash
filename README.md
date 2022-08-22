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
      --duration duration   time duration for speed test (default 3s)
  -h, --help                help for speedtest-clash
```

* speed test by clash config yaml

```bash
 ./dist/speedtest-clash --config test/proxies.yaml
```

* speed test by clash config from http

```bash
 ./dist/speedtest-clash --config https://www.youconfig.com
```

![Example](./docs/example.svg)
