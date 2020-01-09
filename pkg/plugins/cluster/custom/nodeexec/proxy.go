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

import (
	"bytes"
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/util"
)

type remoteExecutor interface {
	doCmdOnPod(cli kubernetes.Interface, config *restclient.Config, namespace string, podName string, cmd []string) (string, string, error)
}

type defaultExecutor struct {
	newSPDYExecutor func(config *restclient.Config, method string, url *url.URL) (remotecommand.Executor, error)
}

// doCmdOnPod do remote cmd exec on target pod
func (d *defaultExecutor) doCmdOnPod(cli kubernetes.Interface, config *restclient.Config, namespace string, podName string, cmd []string) (string, string, error) {
	req := cli.CoreV1().RESTClient().Post().Resource("pods").Name(podName).
		Namespace(namespace).SubResource("exec")
	option := &v1.PodExecOptions{
		Command: cmd,
		//Stdin:   true,
		Stdout: true,
		Stderr: true,
	}

	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)

	exec, err := d.newSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", "", err
	}

	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	var sizeQueue remotecommand.TerminalSizeQueue
	err = exec.Stream(remotecommand.StreamOptions{
		//Stdin:  os.Stdin,
		Stdout:            stdout,
		Stderr:            stderr,
		TerminalSizeQueue: sizeQueue,
	})
	if err != nil {
		return "", "", err
	}

	return stdout.String(), stderr.String(), nil
}

// DaemonSetProxy will create a DaemonSet to do cmd on node
type DaemonSetProxy struct {
	logger         logger.Logger
	cli            kubernetes.Interface
	namespace      string
	dsName         string
	config         *restclient.Config
	remoteExecutor remoteExecutor
	image          string
	autoCreate     bool
	deleteNs       bool
}

// NewDaemonSetProxy create and init a new DaemonSetProxy
func NewDaemonSetProxy(logger logger.Logger, cli kubernetes.Interface, config *restclient.Config, namespace string, ds string, image string, autoCreate, deleteNs bool) (*DaemonSetProxy, error) {
	d := &DaemonSetProxy{
		cli:        cli,
		namespace:  namespace,
		dsName:     ds,
		logger:     logger,
		config:     config,
		image:      image,
		autoCreate: autoCreate,
		deleteNs:   deleteNs,
		remoteExecutor: &defaultExecutor{
			newSPDYExecutor: remotecommand.NewSPDYExecutor,
		},
	}

	if d.autoCreate {
		return d, d.tryCreateProxy()
	}
	return d, nil
}

func (d *DaemonSetProxy) tryCreateProxy() error {
	// create namespace
	ns := &v1.Namespace{}
	ns.Name = d.namespace
	if _, err := d.cli.CoreV1().Namespaces().Create(ns); err != nil {
		if !k8serr.IsAlreadyExists(err) {
			return errors.Wrapf(err, "create namespace %s failed", d.namespace)
		}
	}

	// try create DaemonSet
	dsYaml := fmt.Sprintf(proxyYaml, d.dsName, d.namespace, d.image)
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode([]byte(dsYaml), nil, nil)
	if err != nil {
		return errors.Wrapf(err, "decode proxy yaml failed")
	}

	ds, ok := obj.(*v12.DaemonSet)
	if !ok {
		return fmt.Errorf("covert to app/v1 DaemonSet failed")
	}

	if _, err := d.cli.AppsV1().DaemonSets(d.namespace).Create(ds); err != nil {
		if !k8serr.IsAlreadyExists(err) {
			return errors.Wrapf(err, "create namespace %s failed", d.namespace)
		}
	}

	// wait DaemonSet scheduled
	for {
		ds, err := d.cli.AppsV1().DaemonSets(d.namespace).Get(d.dsName, metav1.GetOptions{})
		if err != nil {
			return errors.Wrapf(err, "get DaemonSet %s/%s failed", d.namespace, d.dsName)
		}

		if ds.Status.DesiredNumberScheduled == ds.Status.CurrentNumberScheduled {
			return nil
		}

		d.logger.Infof("wait for agent DesiredNumberScheduled = CurrentNumberScheduled")
		time.Sleep(time.Second)
	}
}

// Machine get machine information
func (d *DaemonSetProxy) DoCmd(nodeName string, cmd []string) (string, string, error) {
	retStdout, retStderr := "", ""
	err := util.RetryUntilTimeout(time.Second*10, time.Minute, func() error {
		pods, err := d.cli.CoreV1().Pods(d.namespace).List(metav1.ListOptions{
			FieldSelector: "spec.nodeName=" + nodeName,
			LabelSelector: labels.SelectorFromSet(map[string]string{
				"k8s-app": "kube-jarvis-agent",
			}).String(),
		})
		if err != nil {
			return errors.Wrapf(err, "get agent pod failed")
		}

		if len(pods.Items) != 1 {
			d.logger.Infof("target agent pod on node %s not found, it may be scheduled later", nodeName)
			return util.RetryAbleErr
		}
		pod := pods.Items[0]

		// check pod status ,it must be ready
		isReady := false
		for _, c := range pod.Status.Conditions {
			if c.Type == v1.PodReady && c.Status == v1.ConditionTrue {
				isReady = true
				break
			}
		}

		if !isReady {
			d.logger.Infof("pod %s on node %s is not ready", pod.Name, nodeName)
			return util.RetryAbleErr
		}

		retStdout, retStderr, err = d.remoteExecutor.doCmdOnPod(d.cli, d.config, d.namespace, pod.Name, cmd)
		return err
	})

	return retStdout, retStderr, err
}

// Finish do clean for ComponentExecutor
func (d *DaemonSetProxy) Finish() error {
	if d.autoCreate {
		if err := d.cli.AppsV1().DaemonSets(d.namespace).Delete(d.dsName, &metav1.DeleteOptions{}); err != nil {
			return err
		}
	}

	if d.deleteNs {
		if err := d.cli.CoreV1().Namespaces().Delete(d.namespace, &metav1.DeleteOptions{}); err != nil {
			return err
		}
	}

	return nil
}
