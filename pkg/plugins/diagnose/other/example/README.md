# example diagnostic 

This is the example that shows how to write a diagnostic plugins 

# config
```yaml
diagnostics:
  - type: "example"
   # default values
    name: "example"
    catalogue: ["other"]    
    config:
      message: "message"     #extra message
```
# supported cluster type
* all