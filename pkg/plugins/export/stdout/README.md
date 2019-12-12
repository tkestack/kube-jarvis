# stdout exporter
stdout exporter just print result to stdout with a simple format

# config
```yaml
exporters:
  - type: "stdout"
    name: "stdout"
    config:
      format: "fmt" # fmt or json or yaml
```

# supported cluster type 
* all