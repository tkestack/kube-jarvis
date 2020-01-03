package status

import (
	"context"
	"k8s.io/kubectl/pkg/describe"
	"testing"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"

	"tkestack.io/kube-jarvis/pkg/plugins"

	"tkestack.io/kube-jarvis/pkg/translate"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRequestLimitDiagnostic_StartDiagnose(t *testing.T) {
	res := cluster.NewResources()
	res.Nodes = &v1.NodeList{}

	condition := v1.NodeCondition{
		Type:    v1.NodeReady,
		Status:  v1.ConditionTrue,
		Reason:  "123",
		Message: "123",
	}
	conditions := make([]v1.NodeCondition, 0)
	conditions = append(conditions, condition)
	condition.Type = v1.NodePIDPressure
	condition.Status = v1.ConditionUnknown
	conditions = append(conditions, condition)
	condition.Type = v1.NodeDiskPressure
	condition.Status = v1.ConditionTrue
	conditions = append(conditions, condition)
	condition.Type = v1.NodeMemoryPressure
	condition.Status = v1.ConditionTrue
	conditions = append(conditions, condition)
	condition.Type = v1.NodeNetworkUnavailable
	condition.Status = v1.ConditionTrue
	conditions = append(conditions, condition)

	node1 := v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node1",
			Labels: map[string]string{},
		},
		Spec: v1.NodeSpec{},
		Status: v1.NodeStatus{
			Conditions: conditions,
		},
	}

	node2 := v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:   "node1",
			Labels: map[string]string{
				describe.LabelNodeRolePrefix+"master":"true",
			},
		},
		Spec: v1.NodeSpec{},
		Status: v1.NodeStatus{
			Conditions: conditions,
		},
	}


	res.Nodes.Items = append(res.Nodes.Items, node1)
	res.Nodes.Items = append(res.Nodes.Items, node2)

	d := NewDiagnostic(&diagnose.MetaData{
		MetaData: plugins.MetaData{
			Translator: translate.NewFake(),
			Logger:     logger.NewLogger(),
			Type:       DiagnosticType,
			Name:       DiagnosticType,
		},
	})

	if err := d.Complete(); err != nil {
		t.Fatalf(err.Error())
	}

	result, _ := d.StartDiagnose(context.Background(), diagnose.StartDiagnoseParam{
		Resources: res,
	})

	for {
		s, ok := <-result
		if !ok {
			break
		}
		t.Logf("%+v", s)
	}
}
