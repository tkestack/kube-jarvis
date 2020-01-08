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
      memory: "3.75Gi"
      cpu: "1000m"
    - maxnodetotal: 10
      memory: "7.5Gi"
      cpu: "2000m"
    - maxnodetotal: 100
      memory: "15Gi"
      cpu: "4000m"
    - maxnodetotal: 250
      memory: "30Gi"
      cpu: "8000m"
    - maxnodetotal: 500
      memory: "60Gi"
      cpu: "16000m"
    - maxnodetotal: 100000
      memory: "120Gi"
      cpu: "32000m"    
```
# supported cluster type 
* all