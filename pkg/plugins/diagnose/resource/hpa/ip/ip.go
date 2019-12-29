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
package ip

import (
	"context"
	"math"
	"net"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "hpa-ip"
)

// Diagnostic is a example diagnostic shows how to write a diagnostic
type Diagnostic struct {
	*diagnose.MetaData
	result chan *diagnose.Result
	param  *diagnose.StartDiagnoseParam
}

// NewDiagnostic return a example diagnostic
func NewDiagnostic(meta *diagnose.MetaData) diagnose.Diagnostic {
	return &Diagnostic{
		result:   make(chan *diagnose.Result, 1000),
		MetaData: meta,
	}
}

// Complete check and complete config items
func (d *Diagnostic) Complete() error {
	return nil
}

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context, param diagnose.StartDiagnoseParam) (chan *diagnose.Result, error) {
	d.param = &param
	go func() {
		defer close(d.result)
		var totalIPCount int
		var curIPCount int
		var hpaMaxIPCount int
		for _, node := range d.param.Resources.Nodes.Items {
			podCIDR := node.Spec.PodCIDR
			if podCIDR == "" {
				continue
			}
			_, netCIDR, _ := net.ParseCIDR(podCIDR)
			cur, total := netCIDR.Mask.Size()
			totalIPCount += int(math.Pow(2, float64(total-cur))) - 2
		}
		for _, pod := range d.param.Resources.Pods.Items {
			if pod.Spec.HostNetwork {
				continue
			}
			curIPCount += 1
		}
		deploySet := make(map[string]int)
		for _, deploy := range d.param.Resources.Deployments.Items {
			deploySet[deploy.Namespace+"/"+deploy.Name] = int(*deploy.Spec.Replicas)
		}
		hpaMaxIPCount = curIPCount
		for _, hpa := range d.param.Resources.HPAs.Items {
			if hpa.Spec.ScaleTargetRef.Kind != "Deployment" {
				d.Logger.Errorf("hpa %s/%s related %v, not deployment", hpa.Namespace, hpa.Name, hpa.Spec.ScaleTargetRef)
				continue
			}
			key := hpa.Namespace + "/" + hpa.Spec.ScaleTargetRef.Name
			replicas, ok := deploySet[key]
			if ok {
				hpaMaxIPCount += int(hpa.Spec.MaxReplicas) - replicas
			}
		}

		d.result <- &diagnose.Result{
			Level:   diagnose.HealthyLevelGood,
			Title:   d.Translator.Message("hpa-ip-title", nil),
			ObjName: "*",
			Desc: d.Translator.Message("hpa-ip-desc", map[string]interface{}{
				"CurrentIPCount": curIPCount,
				"HPAMaxIPCount":  hpaMaxIPCount,
				"ClusterIPCount": totalIPCount,
			}),
		}
	}()
	return d.result, nil
}
