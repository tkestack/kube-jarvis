package status

import (
	"context"
	v12 "k8s.io/api/apps/v1"
	"k8s.io/utils/pointer"
	"testing"
	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/translate"
)

func TestStatusDiagnostic_StartDiagnose(t *testing.T) {
	res := cluster.NewResources()
	res.Deployments = &v12.DeploymentList{}

	d := NewDiagnostic(&diagnose.MetaData{
		MetaData: plugins.MetaData{
			Translator: translate.NewFake(),
		},
	})

	res.StatefulSets = &v12.StatefulSetList{}
	res.StatefulSets.Items = make([]v12.StatefulSet, 1)
	sts := v12.StatefulSet{}
	sts.Name = "sts1"
	sts.Spec.Replicas = pointer.Int32Ptr(1)
	sts.Status.Replicas = 1
	sts.Status.ReadyReplicas = 0
	res.StatefulSets.Items[0] = sts

	res.Deployments = &v12.DeploymentList{}
	res.Deployments.Items = make([]v12.Deployment, 1)
	deploy := v12.Deployment{}
	deploy.Name = "deploy1"
	deploy.Spec.Replicas = pointer.Int32Ptr(1)
	deploy.Status.Replicas = 1
	deploy.Status.AvailableReplicas = 0
	res.Deployments.Items[0] = deploy

	res.DaemonSets = &v12.DaemonSetList{}
	res.DaemonSets.Items = make([]v12.DaemonSet, 1)
	ds := v12.DaemonSet{}
	ds.Name = "ds1"
	ds.Status.DesiredNumberScheduled = 1
	ds.Status.NumberReady = 0
	res.DaemonSets.Items[0] = ds



	if err := d.Complete(); err != nil {
		t.Fatalf(err.Error())
	}

	result, _ := d.StartDiagnose(context.Background(), diagnose.StartDiagnoseParam{
		Resources: res,
	})

	total := 0
	for {
		s, ok := <-result
		if !ok {
			break
		}
		total++

		t.Logf("%+v", *s)
	}
	if total != 3 {
		t.Fatalf("should return 1 result")
	}
}