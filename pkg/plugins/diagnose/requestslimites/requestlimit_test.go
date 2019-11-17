package requestslimites

import (
	"context"
	"testing"

	"github.com/RayHuangCN/Jarvis/pkg/plugins/diagnose"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes/fake"
)

func TestRequestLimitDiagnostic_StartDiagnose(t *testing.T) {
	cli := fake.NewSimpleClientset()
	kubeDNSlimit := make(map[v1.ResourceName]resource.Quantity)
	kubeDNSRequest := make(map[v1.ResourceName]resource.Quantity)

	kubeDNSlimit[v1.ResourceCPU] = resource.MustParse("100m")
	kubeDNSlimit[v1.ResourceMemory] = resource.MustParse("170M")
	kubeDNSRequest[v1.ResourceCPU] = resource.MustParse("100m")
	kubeDNSRequest[v1.ResourceMemory] = resource.MustParse("30M")

	pod := &v1.Pod{}
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

	if _, err := cli.CoreV1().Pods("default").Create(pod); err != nil {
		t.Fatalf(err.Error())
	}

	pod = &v1.Pod{}
	pod.Name = "pod2"
	pod.Namespace = "default"
	pod.Spec.Containers = []v1.Container{
		{
			Name:  "kubedns",
			Image: "1",
		},
	}

	if _, err := cli.CoreV1().Pods("default").Create(pod); err != nil {
		t.Fatalf(err.Error())
	}

	d := NewDiagnostic(&diagnose.CreateParam{
		Score:  10,
		Weight: 10,
		Cli:    cli,
	})
	result := d.StartDiagnose(context.Background())

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
	if total != 4 {
		t.Fatalf("should return 4 result")
	}
}
