# Plugs
## Cluster
Cluster is the abstraction of a particular type of cluster, and it is responsible for probing and discovering the core components of the cluster and collecting cluster-related information
* [custom](./cluster/custom/README.md)
## Coordinator
Coordinator is responsible for coordinating the work of the other plug-ins, executing the various diagnostics, and distributing the output to the exporters
* [default](./coordinate/basic/README.md)
## Diagnostic
Diagnostic is responsible for diagnosing an aspect of the cluster, outputting diagnostic results and repair recommendations
* [master-capacity](./diagnose/master/capacity/README.md)  
* [requests-limits](./diagnose/resource/workload/requestslimits/README.md)
* [example](./diagnose/other/example/README.md) 
## Exporter
Exporter is responsible for formatting the output or
storage
* [stdout](./export/stdout/README.md) 
* [config-map](./export/configmap/README.md) 
