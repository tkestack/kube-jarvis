# apiserver-args diagnostic 

This is diagnostic detection of whether kube-apiserver'arguments are a best practice

# config
```yaml
diagnostics:
  - type: "kube-apiserver-args"
    catalogue: ["master"]    
```
# supported cluster type
* all