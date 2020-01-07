# node-ha diagnostic 

check whether the cluster nodes are evenly distributed in multiple zones and there are no single point of failure.  

# config
```yaml
diagnostics:
- type: "node-ha" 
  # default values
  name: "node-ha"
  catalogue: ["master"]
  config:
```
# supported cluster type 
* all
