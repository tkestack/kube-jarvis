# kube-jarvis
[![Build Status](https://travis-ci.org/RayHuangCN/kube-jarvis.svg?branch=master)](https://travis-ci.org/RayHuangCN/kube-jarvis)
[![Go Report Card](https://goreportcard.com/badge/github.com/RayHuangCN/kube-jarvis)](https://goreportcard.com/report/github.com/RayHuangCN/kube-jarvis)
[![codecov](https://codecov.io/gh/RayHuangCN/kube-jarvis/branch/master/graph/badge.svg)](https://codecov.io/gh/RayHuangCN/kube-jarvis)

kube-jarvis is a tool used to check the health of the kubernetes cluster

# Features

* Comprehensively check the cluster health status
* Support a variety of cloud manufacturers
* Highly configurable
* Highly extensible
* Description statements can be customized

# Quick start
On any node has "/$HOME/.kube/config"
```bash
wget -O -  https://kube-jarvis-1251707795.cos.ap-guangzhou.myqcloud.com/run.sh | bash
```

# Config struct
```yaml
global:
  trans: "translation"
  lang: "en"

cluster:
  # see detail of custom cluster here:
  type: "custom"
  kubeconfig: ""

diagnostics:
  - type: "master-capacity"
  - type: "master-apiserver"
  - type: "node-sys"
  - type: "requests-limits"

exporters:
  - type: "stdout"

  - type: "file"
    name: "for json"
    config:
      format: "json"
      path: "result.json"


```

# Run in docker
login any node of your cluster and exec cmd:
```bash
docker run  -i -t docker.io/raylhuang110/kube-jarvis:latest
```
> [you can found all docker images here](https://hub.docker.com/r/raylhuang110/kube-jarvis/tags)

# Run as job or cronjob
create common resource (Namespaces, ServiceAccount ...)
```bash
kubectl apply -f manifests/ 
```
run as job
```bash
kubectl apply -f manifests/workload/job.yaml
```
run as cronjob (default run at 00:00 every day)
```bash
kubectl apply -f manifests/workload/cronjob.yaml
```
# Plugins
we call coordinator, diagnostics, evaluators and exporters as "plugins"
> [you can found all plugins lists here](https://github.com/RayHuangCN/kube-jarvis/tree/master/pkg/plugins)

# License
Apache License 2.0 - see LICENSE.md for more details
