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
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/utils/pointer"
	"testing"
	"tkestack.io/kube-jarvis/pkg/plugins/cluster"

	"tkestack.io/kube-jarvis/pkg/plugins"

	"tkestack.io/kube-jarvis/pkg/translate"

	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"

	v1 "k8s.io/api/core/v1"
)

func TestBatchCheckDiagnostic_StartDiagnose(t *testing.T) {
	res := cluster.NewResources()
	res.Jobs = &batchv1.JobList{}
	res.CronJobs = &batchv1beta1.CronJobList{}

	job := batchv1.Job{}
	job.Name = "job1"
	job.Namespace = "default"

	job.Spec.BackoffLimit = pointer.Int32Ptr(20)
	job.Spec.Template = v1.PodTemplateSpec{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "jobContainer",
					Image: "1",
				},
			},
		},
	}

	res.Jobs.Items = append(res.Jobs.Items, job)

	job = batchv1.Job{}
	job.Name = "job2"
	job.Namespace = "default"

	job.Spec.BackoffLimit = pointer.Int32Ptr(5)
	job.Spec.Template = v1.PodTemplateSpec{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "jobContainer",
					Image: "1",
				},
			},
		},
	}

	cronJob := batchv1beta1.CronJob{}
	cronJob.Name = "cronJob1"
	cronJob.Namespace = "default"

	cronJob.Spec.ConcurrencyPolicy = batchv1beta1.ForbidConcurrent
	cronJob.Spec.SuccessfulJobsHistoryLimit = pointer.Int32Ptr(3)
	cronJob.Spec.FailedJobsHistoryLimit = pointer.Int32Ptr(3)
	cronJob.Spec.Schedule = "*/5 * * * *"
	cronJob.Spec.JobTemplate = batchv1beta1.JobTemplateSpec{
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "jobContainer",
							Image: "1",
						},
					},
				},
			},
		},
	}

	res.CronJobs.Items = append(res.CronJobs.Items, cronJob)

	cronJob = batchv1beta1.CronJob{}
	cronJob.Name = "cronJob2"
	cronJob.Namespace = "default"

	cronJob.Spec.ConcurrencyPolicy = batchv1beta1.AllowConcurrent
	cronJob.Spec.SuccessfulJobsHistoryLimit = pointer.Int32Ptr(20)
	cronJob.Spec.FailedJobsHistoryLimit = pointer.Int32Ptr(30)
	cronJob.Spec.Schedule = "*/5 * * * *"
	cronJob.Spec.JobTemplate = batchv1beta1.JobTemplateSpec{
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "jobContainer",
							Image: "1",
						},
					},
				},
			},
		},
	}

	res.CronJobs.Items = append(res.CronJobs.Items, cronJob)

	d := NewDiagnostic(&diagnose.MetaData{
		MetaData: plugins.MetaData{
			Translator: translate.NewFake(),
		},
	})

	if err := d.Complete(); err != nil {
		t.Fatalf(err.Error())
	}

	result, _ := d.StartDiagnose(context.Background(), diagnose.StartDiagnoseParam{
		Resources: res,
	})

	total := 0
	outputs := make(map[string][]diagnose.Result, 0)
	for {
		s, ok := <-result
		if !ok {
			break
		}
		total++

		if _, ok := outputs[(*s).ObjName]; !ok {
			outputs[(*s).ObjName] = []diagnose.Result{*s}
		} else {
			outputs[(*s).ObjName] = append(outputs[(*s).ObjName], *s)
		}

		t.Logf("%+v", *s)
	}
	t.Logf("total entry: %d", total)

	if _, ok := outputs["default:job1"]; !ok {
		t.Fatalf("job1 should not pass check")
	}
	if _, ok := outputs["default:job2"]; ok {
		t.Fatalf("job2 should pass check")
	}
	if _, ok := outputs["default:cronJob1"]; ok {
		t.Fatalf("cronJob1 should pass check")
	}
	if _, ok := outputs["default:cronJob2"]; !ok {
		t.Fatalf("cronJob2 should not pass check")
	}
}
