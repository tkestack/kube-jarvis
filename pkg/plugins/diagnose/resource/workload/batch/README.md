# job/cronjob diagnostic

This diagnostic report the healthy of job/cronjob's restart strategy configuration  

# config
```yaml
diagnostics:
  - type: "batch-check"
    # default values
    name: "job/cronjob check"
    catalogue: ["resource"]
```
# supported cluster type 
* all