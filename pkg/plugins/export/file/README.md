# file exporter
file exporter save result into files and provides a simple query interface if config "server" is true 

* POST "/exporter/file/query : query results
param:
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

# config

```yaml
exporters:
  - type: "file"
    config:
      server: "true"  # set true to open a query interface
      maxremain: 7  # how many result file be saved
      path: "results/" # the path of data files
```

# supported cluster type 
* all