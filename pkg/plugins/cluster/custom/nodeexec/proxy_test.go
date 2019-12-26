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
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	restfake "k8s.io/client-go/rest/fake"
	"k8s.io/client-go/tools/remotecommand"
	cmdtesting "k8s.io/kubectl/pkg/cmd/testing"
	"net/http"
	"net/url"
	"testing"
	"tkestack.io/kube-jarvis/pkg/logger"
)

type fakeStream struct {
	err    error
	out    string
	outErr string
}

func (f *fakeStream) Stream(options remotecommand.StreamOptions) error {
	if f.err != nil {
		return f.err
	}

	_, _ = options.Stdout.Write([]byte(f.out))
	_, _ = options.Stderr.Write([]byte(f.outErr))
	return nil
}

type fakeExecutor struct {
	err    error
	out    string
	outErr string
}

func (f *fakeExecutor) doCmdOnPod(cli kubernetes.Interface, config *rest.Config, namespace string, podName string, cmd []string) (string, string, error) {
	if f.err != nil {
		return "", "", f.err
	}
	return f.out, f.outErr, nil
}

func TestDoCmd(t *testing.T) {
	pod := &v1.Pod{}
	pod.Name = "kube-jarvis-agent-1"
	pod.Namespace = "kube-system"
	pod.Spec.NodeName = "node1"
	pod.Labels = map[string]string{
		"k8s-app": "kube-jarvis-agent",
	}
	pod.Status.Conditions = []v1.PodCondition{
		{
			Type:   v1.PodReady,
			Status: v1.ConditionTrue,
		},
	}

	cases := []struct {
		err    error
		out    string
		outErr string
		pod    *v1.Pod
	}{
		// success
		{
			err:    nil,
			out:    "out",
			outErr: "outErr",
			pod:    pod,
		},

		// remote err
		{
			err: fmt.Errorf("err"),
			pod: pod,
		},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%+v", cs), func(t *testing.T) {
			cli := fake.NewSimpleClientset()
			if cs.pod != nil {
				if _, err := cli.CoreV1().Pods(pod.Namespace).Create(cs.pod); err != nil {
					t.Fatalf(err.Error())
				}
			}

			p, err := NewDaemonSetProxy(logger.NewLogger(), cli, nil, "kube-system", "kube-jarvis-agent")
			if err != nil {
				t.Fatalf(err.Error())
			}

			p.remoteExecutor = &fakeExecutor{
				err:    cs.err,
				out:    cs.out,
				outErr: cs.outErr,
			}

			out, outErr, err := p.DoCmd(pod.Spec.NodeName, []string{"test"})
			if cs.err != nil {
				if err == nil {
					t.Fatalf("should return an err if stream must return an err")
				}
				return
			}

			if cs.pod == nil {
				if err == nil {
					t.Fatalf("should return an err if no pod found")
				}
				return
			}

			if out != cs.out {
				t.Fatalf("want out %s,get %s", cs.out, out)
			}

			if outErr != cs.outErr {
				t.Fatalf("want out %s,get %s", cs.out, out)
			}

			if err := p.Finish(); err != nil {
				t.Fatalf(err.Error())
			}
		})
	}
}

func TestDefaultExecutor(t *testing.T) {
	pod := &v1.Pod{}
	pod.Name = "kube-jarvis-agent-1"
	pod.Namespace = "kube-system"
	pod.Spec.NodeName = "node1"
	pod.Labels = map[string]string{
		"k8s-app": "kube-jarvis-agent",
	}

	podPath := "/api/v1/namespaces/kube-system/pods/kube-jarvis-agent-1"
	fetchPodPath := "/namespaces/kube-system/pods/kube-jarvis-agent-1"
	execPath := "/api/v1/namespaces/kube-system/pods/kube-jarvis-agent-1/exec"

	cases := []struct {
		err    error
		out    string
		outErr string
	}{
		{
			err: fmt.Errorf("err"),
		},
		{
			out:    "out",
			outErr: "outErr",
		},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%v", cs), func(t *testing.T) {
			tf := cmdtesting.NewTestFactory().WithNamespace("test")
			defer tf.Cleanup()
			codec := scheme.Codecs.LegacyCodec(scheme.Scheme.PrioritizedVersionsAllGroups()...)
			ns := scheme.Codecs.WithoutConversion()

			tf.Client = &restfake.RESTClient{
				GroupVersion:         schema.GroupVersion{Group: "", Version: "v1"},
				NegotiatedSerializer: ns,
				Client: restfake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
					switch p, m := req.URL.Path, req.Method; {
					case p == podPath && m == "GET":
						body := cmdtesting.ObjBody(codec, pod)
						return &http.Response{StatusCode: http.StatusOK, Header: cmdtesting.DefaultHeader(), Body: body}, nil
					case p == fetchPodPath && m == "GET":
						body := cmdtesting.ObjBody(codec, pod)
						return &http.Response{StatusCode: http.StatusOK, Header: cmdtesting.DefaultHeader(), Body: body}, nil
					case p == execPath && m == "GET":
						body := cmdtesting.ObjBody(codec, pod)
						return &http.Response{StatusCode: http.StatusOK, Header: cmdtesting.DefaultHeader(), Body: body}, nil
					default:
						t.Errorf("%s: unexpected request: %s %#v\n%#v", pod.Name, req.Method, req.URL, req)
						return nil, fmt.Errorf("unexpected request")
					}
				}),
			}

			f := defaultExecutor{newSPDYExecutor: func(config *rest.Config, method string, url *url.URL) (executor remotecommand.Executor, e error) {
				return &fakeStream{
					err:    cs.err,
					out:    cs.out,
					outErr: cs.outErr,
				}, nil
			}}

			cli, err := tf.KubernetesClientSet()
			if err != nil {
				t.Fatalf(err.Error())
			}

			out, outErr, err := f.doCmdOnPod(cli, nil, pod.Namespace, pod.Name, []string{"test"})
			if cs.err != nil {
				if err == nil {
					t.Fatalf("should return an err")
				}
				return
			}

			if out != cs.out {
				t.Fatalf("want %s get %s", cs.out, out)
			}

			if outErr != cs.outErr {
				t.Fatalf("want %s get %s", cs.outErr, outErr)
			}
		})
	}
}
