#!/bin/bash
echo "Download kube-jarvis...."
wget https://kube-jarvis-1251707795.cos.ap-guangzhou.myqcloud.com/kube-jarvis.tar.gz
tar xf kube-jarvis.tar.gz
cd kube-jarvis
echo "Creating Namespace kube-jarvis..."
kubectl apply -f manifests/common/namespace.yaml
echo "Creating Daemonset kube-jarvis-agent..."
kubectl apply -f manifests/common/agent-ds.yaml
while true
do
     echo "Waiting all kube-jarvis-agent to running..."
     check=`kubectl get po -n kube-jarvis | grep -v NAME | grep -v Running | grep -v grep`
     if [ "$check" == "" ]
     then
             echo "Done"
             break
     fi
     sleep 3
done
./kube-jarvis
