package nodeexec

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"net/url"
	"tkestack.io/kube-jarvis/pkg/logger"
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
}

// NewDaemonSetProxy create and init a new DaemonSetProxy
func NewDaemonSetProxy(logger logger.Logger, cli kubernetes.Interface, config *restclient.Config, namespace string, ds string) (*DaemonSetProxy, error) {
	d := &DaemonSetProxy{
		cli:       cli,
		namespace: namespace,
		dsName:    ds,
		logger:    logger,
		config:    config,
		remoteExecutor: &defaultExecutor{
			newSPDYExecutor: remotecommand.NewSPDYExecutor,
		},
	}

	return d, nil
}

// Machine get machine information
func (d *DaemonSetProxy) DoCmd(nodeName string, cmd []string) (string, string, error) {
	pod, err := d.cli.CoreV1().Pods(d.namespace).List(metav1.ListOptions{
		FieldSelector: "spec.nodeName=" + nodeName,
		LabelSelector: labels.SelectorFromSet(map[string]string{
			"k8s-app": "kube-jarvis-agent",
		}).String(),
	})
	if err != nil {
		return "", "", errors.Wrap(err, "found agent pod failed")
	}

	if len(pod.Items) != 1 {
		return "", "", fmt.Errorf("agent pod not found")
	}

	return d.remoteExecutor.doCmdOnPod(d.cli, d.config, d.namespace, pod.Items[0].Name, cmd)
}

// Finish do clean for ComponentExecutor
func (d *DaemonSetProxy) Finish() error {
	return nil //d.cli.ExtensionsV1beta1().DaemonSets(d.namespace).Delete(d.dsName, &metav1.DeleteOptions{})
}
