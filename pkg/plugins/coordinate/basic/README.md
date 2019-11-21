# basic coordinator

basic coordinator is the default coordinator, it just run diagnostics one by one and the, run evaluators one by one.
any result will be send to all exporters

# config
```yaml
coordinate:
  type: "basic"
```