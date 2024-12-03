## go-proxy

Simple proxy with `x-http-method-override` header support

### Configuration

Arguments precedence:

1. CLI flags
2. ENV vars
3. Config file `config.yml`

Example `config.yml`:

```yaml
from: "localhost:8888"
to: "localhost:9999"
tls_skip_verify: false
```

Example ENV vars:

```sh
export GO_PROXY_FROM=localhost:8888
export GO_PROXY_TO=localhost:9999
export GO_PROXY_TLS_SKIP_VERIFY=false
```

Example CLI args:

```sh
./go-proxy --from localhost:7000 --to localhost:7777 --tls_skip_verify false
```
