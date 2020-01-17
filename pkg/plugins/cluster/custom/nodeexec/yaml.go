/*
* Tencent is pleased to support the open source community by making TKEStack
* available.
*
* Copyright (C) 2012-2019 Tencent. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the “License”); you may not use
* this file except in compliance with the License. You may obtain a copy of the
* License at
*
* https://opensource.org/licenses/Apache-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an “AS IS” BASIS, WITHOUT
* WARRANTIES OF ANY KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations under the License.
 */
package nodeexec

var proxyYaml = `apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: %s
  labels:
    k8s-app: kube-jarvis-agent
  namespace: %s
spec:
  selector:
    matchLabels:
      k8s-app: kube-jarvis-agent
  template:
    metadata:
      labels:
        k8s-app: kube-jarvis-agent
    spec:
      hostNetwork: true
      hostPID: true
      hostIPC: true
      tolerations:
      - effect: NoExecute
        operator: Exists
      - effect: NoSchedule
        operator: Exists
      containers:
      - image: %s
        command: ["sleep","1000000000d"]
        imagePullPolicy: Always
        name: proxy
        securityContext:
          runAsUser: 0
          privileged: true
        volumeMounts:
        - name: dbus
          mountPath: /var/run/dbus
        - name: run-systemd
          mountPath: /run/systemd
        - name: etc-systemd
          mountPath: /etc/systemd
        - name: var-log
          mountPath: /var/log
        - name: var-run
          mountPath: /var/run
        - name: run
          mountPath: /run
        - name: usr-lib-systemd
          mountPath: /usr/lib/systemd
        - name: etc-machine-id
          mountPath: /etc/machine-id
        - name: etc-sudoers
          mountPath: /etc/sudoers.d
      volumes:
      - name: dbus
        hostPath:
          path: /var/run/dbus
          type: Directory
      - name: run-systemd
        hostPath:
          path: /run/systemd
          type: Directory
      - name: etc-systemd
        hostPath:
          path: /etc/systemd
          type: Directory
      - name: var-log
        hostPath:
          path: /var/log
          type: Directory
      - name: var-run
        hostPath:
          path: /var/run
          type: Directory
      - name: run
        hostPath:
          path: /run
          type: Directory
      - name: usr-lib-systemd
        hostPath:
          path: /usr/lib/systemd
          type: Directory
      - name: etc-machine-id
        hostPath:
          path: /etc/machine-id
          type: File
      - name: etc-sudoers
        hostPath:
          path: /etc/sudoers.d
          type: Directory`
