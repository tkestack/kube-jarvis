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
package ha

import (
	"context"
	"fmt"
	"k8s.io/api/core/v1"
	"testing"
	"tkestack.io/kube-jarvis/pkg/logger"
	"tkestack.io/kube-jarvis/pkg/plugins"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
	"tkestack.io/kube-jarvis/pkg/translate"
)

func TestDiagnostic_diagnoseNodeNum(t *testing.T) {
	var cases = []struct {
		Err       bool
		Pass      bool
		TotalNode int
	}{
		{
			Err:       false,
			Pass:      false,
			TotalNode: 0,
		},
		{
			Err:       false,
			Pass:      false,
			TotalNode: 1,
		},
		{
			Err:       false,
			Pass:      true,
			TotalNode: 2,
		},
		{
			Err:       false,
			Pass:      true,
			TotalNode: 100,
		},
	}
	total := 0
	for _, cs := range cases {
		t.Run(fmt.Sprintf("Pass=%v", cs.Pass), func(t *testing.T) {
			res := cluster.NewResources()
			var nodes []v1.Node
			for i := 0; i < cs.TotalNode; i++ {
				node := v1.Node{}
				node.Name = fmt.Sprintf("node-%d", i)
				nodes = append(nodes, node)
			}
			res.Nodes = &v1.NodeList{Items: nodes}
			// start diagnostic
			d := NewDiagnostic(&diagnose.MetaData{
				MetaData: plugins.MetaData{
					Translator: translate.NewFake(),
					Logger:     logger.NewLogger(),
					Type:       DiagnosticType,
					Name:       DiagnosticType,
				},
				Catalogue: diagnose.CatalogueMaster,
			}).(*Diagnostic)

			if err := d.Complete(); err != nil {
				t.Fatalf(err.Error())
			}

			_, _ = d.StartDiagnose(context.Background(), diagnose.StartDiagnoseParam{
				CloudType: "fake",
				Resources: res,
			})
			for {
				r, ok := <-d.result
				if !ok {
					break
				}
				if r.ObjName != "node-num" {
					continue
				}
				total++
				if cs.Err {
					if r.Level != diagnose.HealthyLevelFailed {
						t.Fatalf("want failed result")
					}
					return
				}
				if cs.Pass && r.Level != diagnose.HealthyLevelGood {
					t.Fatalf("should return HealthyLevelGood")
				}

				if !cs.Pass && r.Level == diagnose.HealthyLevelGood {
					t.Fatalf("should return not HealthyLevelGood")
				}
			}
		})
	}
	if total != len(cases) {
		t.Fatalf("should return %d results but get %d ", len(cases), total)
	}
}

func TestDiagnostic_diagnoseNodeZone(t *testing.T) {
	var cases = []struct {
		Pass          bool
		Err           bool
		TotalNode     int
		TotalZoneNum  int
		ZoneNodeRatio float64
	}{
		{
			Err:           true,
			Pass:          false,
			TotalNode:     1,
			TotalZoneNum:  1,
			ZoneNodeRatio: 0.0,
		},
		{
			Err:           false,
			Pass:          false,
			TotalNode:     1,
			TotalZoneNum:  1,
			ZoneNodeRatio: 0.0,
		},
		{
			Err:           false,
			Pass:          false,
			TotalNode:     2,
			TotalZoneNum:  1,
			ZoneNodeRatio: 0.0,
		},
		{
			Err:           false,
			Pass:          false,
			TotalNode:     100,
			TotalZoneNum:  1,
			ZoneNodeRatio: 0.0,
		},
		{
			Err:           false,
			Pass:          true,
			TotalNode:     100,
			TotalZoneNum:  2,
			ZoneNodeRatio: 0.0,
		},
		{
			Err:           false,
			Pass:          true,
			TotalNode:     100,
			TotalZoneNum:  4,
			ZoneNodeRatio: 0.0,
		},
	}

	for _, cs := range cases {
		t.Run(fmt.Sprintf("Pass=%v", cs.Pass), func(t *testing.T) {
			res := cluster.NewResources()
			var nodes []v1.Node
			for i := 0; i < cs.TotalNode; i++ {
				node := v1.Node{}
				node.Name = fmt.Sprintf("node-%d", i)
				if !cs.Err {
					id := i % cs.TotalZoneNum
					node.Labels = map[string]string{
						FailureDomainKey: fmt.Sprintf("zone-%d", id),
					}
				}
				nodes = append(nodes, node)
			}
			res.Nodes = &v1.NodeList{Items: nodes}
			// start diagnostic
			d := NewDiagnostic(&diagnose.MetaData{
				MetaData: plugins.MetaData{
					Translator: translate.NewFake(),
					Logger:     logger.NewLogger(),
					Type:       DiagnosticType,
					Name:       DiagnosticType,
				},
				Catalogue: diagnose.CatalogueMaster,
			}).(*Diagnostic)

			if err := d.Complete(); err != nil {
				t.Fatalf(err.Error())
			}

			_, _ = d.StartDiagnose(context.Background(), diagnose.StartDiagnoseParam{
				CloudType: "fake",
				Resources: res,
			})
			for {
				r, ok := <-d.result
				if !ok {
					break
				}
				if r.ObjName != "node-zone" {
					continue
				}
				if cs.Err {
					if r.Level != diagnose.HealthyLevelFailed {
						t.Fatalf("want failed result")
					}
					return
				}
				if cs.Pass && r.Level != diagnose.HealthyLevelGood {
					t.Fatalf("should return HealthyLevelGood")
				}

				if !cs.Pass && r.Level == diagnose.HealthyLevelGood {
					t.Fatalf("should return not HealthyLevelGood")
				}
			}
		})
	}
}
