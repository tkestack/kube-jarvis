# kube-jarvis
[![Build Status](https://travis-ci.org/RayHuangCN/kube-jarvis.svg?branch=master)](https://travis-ci.org/RayHuangCN/kube-jarvis)
[![Go Report Card](https://goreportcard.com/badge/github.com/RayHuangCN/kube-jarvis)](https://goreportcard.com/report/github.com/RayHuangCN/kube-jarvis)
[![codecov](https://codecov.io/gh/RayHuangCN/kube-jarvis/branch/master/graph/badge.svg)](https://codecov.io/gh/RayHuangCN/kube-jarvis)

kube-jarvis is a tool used to check the health of the kubernetes cluster

# quick start
```bash
go build -o kube-jarvis cmd/kube-jarvis/*.go
./kube-jarvis --config conf/default.yaml
```

# config struct
```yaml
global:
  trans: "translation" #translation file root director
  lang: "en"  #target language
  cluster:
    kubeconfig: "fake" #cluster kubeconfig filepath,use empty string to enable in cluster model

# coordinator knows how to run all diagnostics, evaluators and exporters
coordinator:
  type: "default" 

# diagnostics diagnose special aspects of cluster
diagnostics: #
  - type: "example"
    name: "example 1"
    score: 10
    weight: 10
    config:
      message: "message"

# evaluators evaluate all diagnose results
evaluators:
  - type: "sum"
    name: "sum 1"

// exporters exporte all diagnostic result and evaluation results
exporters:
  - type: "stdout"
    name: "stdout 1"
```

# run in docker
on any node of your cluster and exec cmd:
```bash
docker run  -i -t docker.io/raylhuang110/kube-jarvis:latest
```
> [you can found all docker images here](https://hub.docker.com/r/raylhuang110/kube-jarvis/tags)


# all plugins
we call coordinator, diagnostics, evaluators and exporters as "plugins"
> [you can found all plugins lists here](https://github.com/RayHuangCN/kube-jarvis/tree/master/pkg/plugins)
