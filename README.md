## go-proxy
Simple proxy with `x-http-method-override` header support

### Configuration

Arguments precedence:

1. CLI flags
2. ENV vars
3. Config file `config.yml`

Example `config.yml`:

```
from: "localhost:8888"
to: "localhost:9999"
```

Example ENV vars:

```
export GO_PROXY_FROM=localhost:8888
export GO_PROXY_TO=localhost:9999
```

Example CLI args:

```
./go-proxy --from localhost:7000 --to localhost:7777
```