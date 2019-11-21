# config-map exporter
config-map exporter save result to k8s ConfigMap 

# config
```yaml
evaluators:
  - type: "config-map"
    name: "" 
    config:
      namespace: "default" 
      name: "kube-jarvis" # ConfigMap name
      format: "json"
```