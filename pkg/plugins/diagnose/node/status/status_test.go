package status

import (
	"context"
	"reflect"
	"strings"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/describe"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/translate"
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
			Name: "node1",
			Labels: map[string]string{
				describe.LabelNodeRolePrefix + "master": "true",
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

func TestDiagnostic_StartDiagnose_Issue1(t *testing.T) {
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
			Name: "node2",
			Labels: map[string]string{
				describe.LabelNodeRolePrefix + "master": "true",
			},
		},
		Spec: v1.NodeSpec{},
		Status: v1.NodeStatus{
			Conditions: conditions,
		},
	}

	node3 := v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node3",
			Labels: map[string]string{
				describe.LabelNodeRolePrefix + "master": "",
			},
		},
		Spec: v1.NodeSpec{},
		Status: v1.NodeStatus{
			Conditions: conditions,
		},
	}

	res.Nodes.Items = append(res.Nodes.Items, node1)
	res.Nodes.Items = append(res.Nodes.Items, node2)
	res.Nodes.Items = append(res.Nodes.Items, node3)

	type args struct {
		ctx   context.Context
		param diagnose.StartDiagnoseParam
	}
	type want struct {
		numMaster int
		numNode   int
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "master mistake for node",
			args: args{
				ctx: context.TODO(),
				param: diagnose.StartDiagnoseParam{
					Resources: res,
				},
			},
			want: want{
				numMaster: 2,
				numNode:   1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDiagnostic(&diagnose.MetaData{
				MetaData: plugins.MetaData{
					Translator: translate.NewFake(),
					Logger:     logger.NewLogger(),
					Type:       DiagnosticType,
					Name:       DiagnosticType,
				},
			})
			ch, err := d.StartDiagnose(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("Diagnostic.StartDiagnose() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := want{}
			for res := range ch {
				if strings.Contains(string(res.Title), "master") {
					got.numMaster++
					continue
				}

				got.numNode++
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Diagnostic.StartDiagnose() = %v, want %v", got, tt.want)
			}
		})
	}
}
