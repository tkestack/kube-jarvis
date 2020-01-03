package status

import (
	"context"
	"fmt"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
)

const (
	// DiagnosticType is type Name of this Diagnostic
	DiagnosticType = "workload-status"
)

var WorkloadType = [...]string{"Deployment", "DaemonSet", "StatefulSet"}

// Diagnostic report the healthy of pods's resources requests limits configuration
type Diagnostic struct {
	*diagnose.MetaData
	result chan *diagnose.Result
	param  *diagnose.StartDiagnoseParam
}

// NewDiagnostic return a requests-limits Diagnostic
func NewDiagnostic(meta *diagnose.MetaData) diagnose.Diagnostic {
	return &Diagnostic{
		MetaData: meta,
		result:   make(chan *diagnose.Result, 1000),
	}
}

// Init do initialization
func (d *Diagnostic) Init() error {
	return nil
}

// Complete check and complete config items
func (d *Diagnostic) Complete() error {
	return nil
}

// ResourceItem is a inner struct
type ResourceItem struct {
	// Name is the name of workload.
	Name string
	// Replicas is the replicas of workload.
	Replicas int32
	// Available is the available replicas of workload.
	Available int32
}

type ResourceItemList []ResourceItem

// ResourceMap is a map of namespaces to three-ResourceItemList, they are Deployment, DaemonSet and  StatefulSet list.
type ResourceMap map[string]*[3]ResourceItemList

func appendWhatever(rsMap ResourceMap, namespace string, idx int, rsItem ResourceItem) {
	if idx >= 3 {
		return
	}
	if rsMap[namespace] == nil {
		rsMap[namespace] = &[3]ResourceItemList{}
	}
	if (*rsMap[namespace])[idx] == nil {
		(*rsMap[namespace])[idx] = make(ResourceItemList, 0)
	}
	(*rsMap[namespace])[idx] = append((*rsMap[namespace])[idx], rsItem)
}

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context, param diagnose.StartDiagnoseParam) (chan *diagnose.Result, error) {
	d.param = &param
	d.result = make(chan *diagnose.Result, 1000)
	if d.param.Resources == nil {
		return nil, fmt.Errorf("diagnose param is nil")
	}
	go func() {
		defer diagnose.CommonDeafer(d.result)
		rsMap := make(ResourceMap)
		if d.param.Resources.Deployments != nil {
			for _, deploy := range d.param.Resources.Deployments.Items {
				replicas := int32(1)
				if deploy.Spec.Replicas != nil {
					replicas = *deploy.Spec.Replicas
				}
				appendWhatever(rsMap, deploy.Namespace, 0, ResourceItem{
					Name:      deploy.Name,
					Replicas:  replicas,
					Available: deploy.Status.AvailableReplicas,
				})
			}
		}

		if d.param.Resources.DaemonSets != nil {
			for _, ds := range d.param.Resources.DaemonSets.Items {
				appendWhatever(rsMap, ds.Namespace, 1, ResourceItem{
					Name:      ds.Name,
					Replicas:  ds.Status.DesiredNumberScheduled,
					Available: ds.Status.NumberReady,
				})
			}
		}

		if d.param.Resources.StatefulSets != nil {
			for _, sts := range d.param.Resources.StatefulSets.Items {
				replicas := int32(1)
				if sts.Spec.Replicas != nil {
					replicas = *sts.Spec.Replicas
				}
				appendWhatever(rsMap, sts.Namespace, 2, ResourceItem{
					Name:      sts.Name,
					Replicas:  replicas,
					Available: sts.Status.ReadyReplicas,
				})
			}
		}

		d.diagnose(rsMap)

	}()
	return d.result, nil
}

func (d *Diagnostic) diagnose(rsMap ResourceMap) {
	for namespace, rsLists := range rsMap {
		for typ, rsList := range rsLists {
			if rsList == nil {
				continue
			}
			for idx, rs := range rsList {
				level := diagnose.HealthyLevelGood
				descId := "workload-status-good-desc"
				proposalId := "workload-status-good-proposal"
				if rs.Replicas != rs.Available {
					level = diagnose.HealthyLevelWarn
					descId = "workload-status-desc"
					proposalId = "workload-status-proposal"
				}
				d.result <- &diagnose.Result{
					Level:   level,
					Title:   d.Translator.Message("workload-status-title", nil),
					ObjName: fmt.Sprintf("%s:%s", namespace, rsLists[typ][idx].Name),
					Desc: d.Translator.Message(descId, map[string]interface{}{
						"Name":      rsLists[typ][idx].Name,
						"Namespace": namespace,
						"Workload":  WorkloadType[typ],
						"Replicas":  rs.Replicas,
						"Available": rs.Available,
					}),
					Proposal: d.Translator.Message(proposalId, nil),
				}
			}
		}
	}
	return
}
