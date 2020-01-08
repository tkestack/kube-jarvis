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
package batch

import (
	"context"
	"fmt"

	"tkestack.io/kube-jarvis/pkg/plugins/cluster"

	v1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "batch-check"
)

// Diagnostic report the healthy of pods's resources health check configuration
type Diagnostic struct {
	Filter cluster.ResourcesFilter
	*diagnose.MetaData
	result chan *diagnose.Result
	param  *diagnose.StartDiagnoseParam
}

// NewDiagnostic return a health check Diagnostic
func NewDiagnostic(meta *diagnose.MetaData) diagnose.Diagnostic {
	return &Diagnostic{
		MetaData: meta,
		result:   make(chan *diagnose.Result, 1000),
	}
}

// Complete check and complete config items
func (d *Diagnostic) Complete() error {
	return d.Filter.Compile()
}

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context, param diagnose.StartDiagnoseParam) (chan *diagnose.Result, error) {
	d.param = &param
	d.result = make(chan *diagnose.Result, 1000)
	go func() {
		defer diagnose.CommonDeafer(d.result)
		for _, job := range d.param.Resources.Jobs.Items {
			if d.Filter.Filtered(job.Namespace, "Job", job.Name) {
				continue
			}
			d.diagnoseJob(job)
		}
		for _, cronJob := range d.param.Resources.CronJobs.Items {
			if d.Filter.Filtered(cronJob.Namespace, "CronJob", cronJob.Name) {
				continue
			}
			d.diagnoseCronJob(cronJob)
		}
	}()
	return d.result, nil
}

func (d *Diagnostic) diagnoseJob(job v1.Job) {
	obj := map[string]interface{}{
		"Namespace":        job.Namespace,
		"Name":             job.Name,
		"RecommendedValue": 10,
	}

	if job.Spec.BackoffLimit != nil && *job.Spec.BackoffLimit > 10 {
		d.result <- &diagnose.Result{
			Level:    diagnose.HealthyLevelWarn,
			ObjName:  fmt.Sprintf("%s:%s", job.Namespace, job.Name),
			ObjInfo:  obj,
			Title:    d.Translator.Message("job-backofflimit-title", nil),
			Desc:     d.Translator.Message("job-backofflimit-desc", obj),
			Proposal: d.Translator.Message("job-backofflimit-proposal", obj),
		}
	}
}

func (d *Diagnostic) diagnoseCronJob(cronJob v1beta1.CronJob) {
	obj := map[string]interface{}{
		"Namespace":        cronJob.Namespace,
		"Name":             cronJob.Name,
		"RecommendedValue": 10,
	}

	if cronJob.Spec.FailedJobsHistoryLimit != nil && *cronJob.Spec.FailedJobsHistoryLimit > 10 {
		d.result <- &diagnose.Result{
			Level:    diagnose.HealthyLevelWarn,
			ObjName:  fmt.Sprintf("CronJob:%s:%s", cronJob.Namespace, cronJob.Name),
			ObjInfo:  obj,
			Title:    d.Translator.Message("cronjob-failedjobhistorylimit-title", nil),
			Desc:     d.Translator.Message("cronjob-failedjobhistorylimit-desc", obj),
			Proposal: d.Translator.Message("cronjob-failedjobhistorylimit-proposal", obj),
		}
	}

	if cronJob.Spec.SuccessfulJobsHistoryLimit != nil && *cronJob.Spec.SuccessfulJobsHistoryLimit > 10 {
		d.result <- &diagnose.Result{
			Level:    diagnose.HealthyLevelWarn,
			ObjName:  fmt.Sprintf("%s:%s", cronJob.Namespace, cronJob.Name),
			ObjInfo:  obj,
			Title:    d.Translator.Message("cronjob-successfuljobshistorylimit-title", nil),
			Desc:     d.Translator.Message("cronjob-failedjobhistorylimit-desc", obj),
			Proposal: d.Translator.Message("cronjob-failedjobhistorylimit-proposal", obj),
		}
	}

	obj2 := map[string]interface{}{
		"Namespace":                    cronJob.Namespace,
		"Name":                         cronJob.Name,
		"CurrentConcurrencyPolicy":     cronJob.Spec.ConcurrencyPolicy,
		"RecommendedConcurrencyPolicy": v1beta1.ForbidConcurrent,
	}

	if cronJob.Spec.ConcurrencyPolicy != v1beta1.ForbidConcurrent {
		d.result <- &diagnose.Result{
			Level:    diagnose.HealthyLevelWarn,
			ObjName:  fmt.Sprintf("%s:%s", cronJob.Namespace, cronJob.Name),
			ObjInfo:  obj2,
			Title:    d.Translator.Message("cronjob-concurrencypolicy-title", nil),
			Desc:     d.Translator.Message("cronjob-concurrencypolicy-desc", obj2),
			Proposal: d.Translator.Message("cronjob-concurrencypolicy-proposal", obj2),
		}
	}
}
