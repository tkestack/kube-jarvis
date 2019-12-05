# stdout exporter
stdout exporter just print result to stdout with a simple format

# config
```yaml
exporters:
  - type: "stdout"
    name: "stdout 1"
    config:
      format: "fmt" # fmt or json or yaml
```

# supported cloud providers
* all