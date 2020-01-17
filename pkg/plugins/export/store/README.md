# store exporter
store exporter save result into global store and provides a simple query interface if config "server" is true 

# config

```yaml
exporters:
  - type: "store"
    config:
      server: true
```

# supported cluster type 
* all

# http api
* POST "/exporter/store/history : "query history"

request:
```json
{
  "Offset":0,
  "Limit":0
}
```

response:
```json
{
  "Records": [
    {
      "ID": "1578920254692112293",
      "Overview": {
        "StartTime": "2020-01-13T20:57:34.692112293+08:00",
        "EndTime": "2020-01-13T20:57:34.69260182+08:00",
        "Statistics": {
          "warn": 1
        },
        "Diagnostics": null
      }
    }
  ]
}
```

* POST "/exporter/store/query : query results

request:
```json
{
  "ID":123, 
  "Type":"node-sys",
  "Name":"",
  "Level":"warn",
  "Offset":0,
  "Limit":0
}
```

response:
```json
{
  "StartTime": "2019-12-30T09:55:22.064492219+08:00",
  "EndTime": "2019-12-30T09:55:22.065226914+08:00",
  "Statistics": {
    "good": 4,
    "warn": 1
  },
  "Diagnostics": [
    {
      "StartTime": "2019-12-30T09:55:22.064973158+08:00",
      "EndTime": "0001-01-01T00:00:00Z",
      "Catalogue": [
        "node"
      ],
      "Type": "node-sys",
      "Name": "",
      "Desc": "",
      "Results": [
        {
          "Level": "warn",
          "ObjName": "10.0.2.4",
          "Title": "Kernel Parameters",
          "Desc": "Node 10.0.2.4 Parameters[ net.ipv4.tcp_tw_reuse=0 ] is not recommended",
          "Proposal": "Set net.ipv4.tcp_tw_reuse=1"
        }
      ],
      "Statistics": {
        "good": 4,
        "warn": 1
      }
    }
  ]
}
```

