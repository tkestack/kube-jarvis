# master-components diagnostic 

Check if the core components are working properly (include k8s node components)
Also check if they have been restarted within a special time 

# config
```yaml
diagnostics:
 - type: "master-components"
   # default values
   name: "master-components"
   catalogue: ["master"]
   config: 
     restarttime: "24h"  
     components:
       - "kube-apiserver"
       - "kube-scheduler"
	   - "kube-controller-manager"
	   - "etcd"
	   - "kube-proxy"
	   - "coredns"
	   - "kube-dns"
	   - "kubelet"
	   - "dockerd"
	   - "containerd"
```
# supported cluster type 
* all