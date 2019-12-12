# Custom Cluster

A custom cluster is a user-built cluster that can be highly customized through a configuration file

# config
```yaml
cluster:
  type: "custom"
  # default config value
  config:
    kubeconfig: "" # the path of kubeconfig file, use "" to use $HOMEDIR/.kube/.config or use in-cluster way
    node: # the way to fetch node machine level data
      type: "proxy" # via the a agent DaemonSet 
      namespace: "kube-jarvis" # the namespace of agent 
      daemonset: "kube-jarvis-agent" # the name of agent DaemonSet 

    components:  # the components that should to explore their information 
      kube-apiserver: # this is the example of component "kube-apiserver"
                      # the default components also includes as follow
                      # "kube-apiserver", "kube-scheduler", "kube-controller-manager", "etcd", "kube-proxy"
                      # "coredns", "kube-dns", "kubelet", "kube-proxy", "dockerd", "containerd"

        type: "auto"  # the way used to explore this component, 
                      # Auto : try follow ways one by one
                      # Bare : use "ps" command to explore component on nodes 
                      # Label: select pods with label selector to explore component
                      # StaticPod: use static pod on node to explore component

        pretype: "Bare"   # the explore way that with highest priority
        name: "kube-apiserver" # the real name of target component
        namespace: "kube-system" # the namespace of the component if use "Lable" or "StaticPod" exploring
        masternodes: true # only explore component on master nodes 
                          # the master nodes are nodes with label item "node-role.kubernetes.io/master" exist

        nodes: [] # only explore component on target nodes
                  #  will try explore component on all nodes if "masternodes" is false and "nodes" is empty

        labels: # the labels that used to select pod when use "Label" exploring
          k8s-app : "kube-apiserver"

```
