# cf-ddns
![GitHub](https://img.shields.io/github/license/exepirit/cf-ddns?style=flat-square)

cf-ddns synchronizes domains provided in configuration (e. g. config, Kubernetes Ingresses) with external DNS provider.

## Installation

Use go compiler to install cf-ddns.

```shell
go get github.com/exepirit/cf-ddns
```

## Usage

1. Set `DDNS_PROVIDER` and `DDNS_SOURCE` environment variables with chosen DNS and configuration provider.
2. Set DNS provider credentials. To do this, fill appropriate environment variables, e.g `DDNS_CFEMAIL`, `DDNS_CFAPIKEY`.
3. Run `cf-ddns`.

## Roadmap

- Write comrehensive documentation. ;)
- Add `/etc/hosts` file as DNS provider support (for mDNS support).
- Add PowerDNS provider support.
- Add Bind as provider support.

## License

![MIT](https://choosealicense.com/licenses/mit/)
