package components

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/translate"
)

func TestDiagnostic_StartDiagnose(t *testing.T) {
	cases := []struct {
		running  bool
		restart  bool
		err      error
		notExist bool
	}{
		{
			notExist: true,
		},
		{
			err: fmt.Errorf("err"),
		},
		{
			running: false,
		},
		{
			running: true,
			restart: false,
		},
		{
			running: true,
			restart: true,
		},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("%+v", cs), func(t *testing.T) {
			pod := &v1.Pod{}
			pod.Name = "pod1"
			pod.Status.ContainerStatuses = []v1.ContainerStatus{
				{
					State: v1.ContainerState{
						Running: &v1.ContainerStateRunning{},
					},
				},
			}

			if cs.restart {
				pod.Status.ContainerStatuses[0].RestartCount = 1
				pod.Status.ContainerStatuses[0].State.Running.StartedAt = metav1.NewTime(time.Now())
			} else {
				pod.Status.ContainerStatuses[0].RestartCount = 1
				pod.Status.ContainerStatuses[0].State.Running.StartedAt = metav1.NewTime(time.Now().Add(-1 * time.Hour * 25))
			}

			cmpName := "test"
			res := cluster.NewResources()
			if !cs.notExist {
				res.CoreComponents[cmpName] = []cluster.Component{{
					Name:      cmpName,
					Node:      "node1",
					Error:     cs.err,
					IsRunning: cs.running,
					Pod:       pod,
				}}

			}

			param := &diagnose.MetaData{}
			param.Translator = translate.NewFake()

			d := NewDiagnostic(param).(*Diagnostic)
			d.Components = []string{cmpName}

			result, err := d.StartDiagnose(context.Background(), diagnose.StartDiagnoseParam{
				CloudType: "test",
				Resources: res,
			})

			if err != nil {
				t.Fatalf(err.Error())
			}

			total := 0
			for {
				r, ok := <-result
				if !ok {
					break
				}

				total++
				if total > 1 {
					t.Fatalf("should return only one result")
				}

				if cs.notExist {
					if r.Level != diagnose.HealthyLevelFailed {
						t.Fatalf("should return an failed result")
					}
					continue
				}

				if cs.err != nil {
					if r.Level != diagnose.HealthyLevelFailed {
						t.Fatalf("should return an failed result")
					}
					continue
				}

				if !cs.running {
					if r.Level != diagnose.HealthyLevelRisk {
						t.Fatalf("should return an risk result")
					}

					if r.Desc != "not-run-desc" {
						t.Fatalf("should get not-run-desc")
					}

					continue
				}

				if cs.running {
					if cs.restart {
						if r.Level != diagnose.HealthyLevelRisk {
							t.Fatalf("should return an result result")
						}

						if r.Desc != "restart-desc" {
							t.Fatalf("should get restart-desc")
						}
					} else {
						if r.Level != diagnose.HealthyLevelGood {
							t.Fatalf("should return an good result")
						}

						if r.Desc != "good-desc" {
							t.Fatalf("should get good-desc")
						}
					}
				}
			}
		})
	}
}
