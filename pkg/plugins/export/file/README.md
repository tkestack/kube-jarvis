# file exporter
file exporter just print result to file with a target format

# config
```yaml
exporters:
  - type: "file"
    name: "json file"
    config:
      format: "json" #  json or yaml
      path: "result.json"
```

# supported cluster type 
* all