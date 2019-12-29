# stdout exporter
stdout exporter just print result to stdout with a simple format

# config
```yaml
exporters:
  - type: "stdout"
    name: "stdout"
    config:
      format: "fmt"  # use "json" to print a json
```

# supported cluster type 
* all