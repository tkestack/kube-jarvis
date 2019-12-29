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
      # join host network namespace
      hostNetwork: true
      # join host process namespace
      hostPID: true
      # join host IPC namespace
      hostIPC: true 
      # tolerations
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
        # Allows systemctl to communicate with the systemd running on the host
        - name: dbus
          mountPath: /var/run/dbus
        - name: run-systemd
          mountPath: /run/systemd
        # Allows to peek into systemd units that are baked into the official EKS AMI
        - name: etc-systemd
          mountPath: /etc/systemd
        # This is needed in order to fetch logs NOT managed by journald
        # journallog is stored only in memory by default, so we need
        #
        # If all you need is access to persistent journals, /var/log/journal/* would be enough
        # FYI, the volatile log store /var/run/journal was empty on my nodes. Perhaps it isn't used in Amazon Linux 2 / EKS AMI?
        # See https://askubuntu.com/a/1082910 for more background
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
      # for systemctl to systemd access
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
      # mainly for dockerd access via /var/run/docker.sock
      - name: var-run
        hostPath:
          path: /var/run
          type: Directory
      # var-run implies you also need this, because
      # /var/run is a synmlink to /run
      # sh-4.2$ ls -lah /var/run
      # lrwxrwxrwx 1 root root 6 Nov 14 07:22 /var/run -> ../run
      - name: run
        hostPath:
          path: /run
          type: Directory
      - name: usr-lib-systemd
        hostPath:
          path: /usr/lib/systemd
          type: Directory
      # Required by journalctl to locate the current boot.
      # If omitted, journalctl is unable to locate host's current boot journal
      - name: etc-machine-id
        hostPath:
          path: /etc/machine-id
          type: File
      # Avoid this error > ERROR [MessageGatewayService] Failed to add ssm-user to sudoers file: open /etc/sudoers.d/ssm-agent-users: no such file or directory
      - name: etc-sudoers
        hostPath:
          path: /etc/sudoers.d
          type: Directory`