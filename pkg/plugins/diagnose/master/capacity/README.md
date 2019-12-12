# master-capacity diagnostic 

check whether the capacity are sufficient for a specific size cluster  

# config
```yaml
diagnostics:
- type: "master-capacity" 
  # default values
  name: "master-capacity"
  catalogue: ["master"]
  config:
    Capacities: 
    - maxnodetotal: 5
      memory: "8000000Ki"
      cpu: "4000m"
    - maxnodetotal: 20
      memory: "16000000Ki"
      cpu: "4000m"
    - maxnodetotal: 100
      memory: "32000000Ki"
      cpu: "8000m"
    - maxnodetotal: 200
      memory: "64000000Ki"
      cpu: "16000m"
    - maxnodetotal: 100000
      memory: "128000000Ki"
      cpu: "16000m"    
```
# supported cluster type 
* all