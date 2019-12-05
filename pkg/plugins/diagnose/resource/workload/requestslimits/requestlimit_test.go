package requestslimits

import (
	"context"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/cluster"
	"testing"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins"

	"github.com/RayHuangCN/kube-jarvis/pkg/translate"

	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestRequestLimitDiagnostic_StartDiagnose(t *testing.T) {
	res := cluster.NewResources()
	res.Pods = &v1.PodList{}

	kubeDNSlimit := make(map[v1.ResourceName]resource.Quantity)
	kubeDNSRequest := make(map[v1.ResourceName]resource.Quantity)

	kubeDNSlimit[v1.ResourceCPU] = resource.MustParse("100m")
	kubeDNSlimit[v1.ResourceMemory] = resource.MustParse("170M")
	kubeDNSRequest[v1.ResourceCPU] = resource.MustParse("100m")
	kubeDNSRequest[v1.ResourceMemory] = resource.MustParse("30M")

	pod := v1.Pod{}
	pod.Name = "pod1"
	pod.Namespace = "default"
	pod.Spec.Containers = []v1.Container{
		{
			Name:  "kubedns",
			Image: "1",
			Resources: v1.ResourceRequirements{
				Limits:   kubeDNSlimit,
				Requests: kubeDNSRequest,
			},
		},
	}
	res.Pods.Items = append(res.Pods.Items, pod)

	pod = v1.Pod{}
	pod.Name = "pod2"
	pod.Namespace = "default"
	pod.Spec.Containers = []v1.Container{
		{
			Name:  "kubedns",
			Image: "1",
		},
	}
	res.Pods.Items = append(res.Pods.Items, pod)

	d := NewDiagnostic(&diagnose.MetaData{
		CommonMetaData: plugins.CommonMetaData{
			Translator: translate.NewFake(),
		},
	})

	if err := d.Init(); err != nil {
		t.Fatalf(err.Error())
	}

	result := d.StartDiagnose(context.Background(), diagnose.StartDiagnoseParam{
		Resources: res,
	})

	total := 0
	for {
		s, ok := <-result
		if !ok {
			break
		}
		total++

		if s.Error != nil {
			t.Fatalf(s.Error.Error())
		}

		t.Logf("%+v", *s)
	}
	if total != 1 {
		t.Fatalf("should return 1 result")
	}
}
