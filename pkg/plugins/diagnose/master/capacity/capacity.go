package capacity

import (
	"context"
	"fmt"
	"github.com/RayHuangCN/kube-jarvis/pkg/plugins/diagnose"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

const (
	// DiagnosticType is type name of this Diagnostic
	DiagnosticType = "master-capacity"
)

// Capacity define a health master node resource status
type Capacity struct {
	// Memory is total memory of master node
	Memory resource.Quantity
	// CpuCore is total core number of master node
	CpuCore resource.Quantity
	// 	MaxNodeTotal indicate the max node number of this master scale
	MaxNodeTotal int
}

// Diagnostic check whether the resources are sufficient for a specific size cluster
type Diagnostic struct {
	*diagnose.MetaData
	result     chan *diagnose.Result
	Capacities []Capacity
}

// NewDiagnostic return a master-node diagnostic
func NewDiagnostic(meta *diagnose.MetaData) diagnose.Diagnostic {
	return &Diagnostic{
		result:     make(chan *diagnose.Result, 100),
		MetaData:   meta,
		Capacities: DefCapacities,
	}

}

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context) chan *diagnose.Result {
	go func() {
		defer close(d.result)
		defer func() {
			if err := recover(); err != nil {
				d.result <- &diagnose.Result{
					Error: fmt.Errorf("%v", err),
				}
				d.Score = 0
			}
		}()

		d.diagnoseCapacity(ctx)
	}()
	return d.result
}

func (d *Diagnostic) diagnoseCapacity(ctx context.Context) {
	label := labels.NewSelector()
	req, err := labels.NewRequirement("node-role.kubernetes.io/master", selection.Exists, nil)
	if err != nil {
		d.result <- &diagnose.Result{Error: err}
		d.Score = 0
		return
	}
	label = label.Add(*req)

	masters, err := d.Cli.CoreV1().Nodes().List(v1.ListOptions{
		LabelSelector: label.String(),
	})

	if err != nil {
		d.result <- &diagnose.Result{Error: err}
		d.Score = 0
		return
	}

	scale, nTotal, err := d.targetCapacity()
	if err != nil {
		d.result <- &diagnose.Result{Error: err}
		d.Score = 0
		return
	}

	score := d.Score / float64(len(masters.Items))
	for _, m := range masters.Items {
		if m.Status.Capacity.Cpu().Cmp(scale.CpuCore) < 0 {
			d.sendCapacityWarnResult(m.Name, "Cpu", nTotal, m.Status.Capacity.Cpu().String(), scale.CpuCore.String(), score)
		} else {
			d.sendCapacityGoodResult(m.Name, "Cpu", nTotal, m.Status.Capacity.Cpu().String(), scale.CpuCore.String(), score)
		}

		if m.Status.Capacity.Memory().Cmp(scale.Memory) < 0 {
			d.sendCapacityWarnResult(m.Name, "Memory", nTotal, m.Status.Capacity.Memory().String(), scale.Memory.String(), score)
		} else {
			d.sendCapacityGoodResult(m.Name, "Memory", nTotal, m.Status.Capacity.Memory().String(), scale.Memory.String(), score)
		}
	}
}

func (d *Diagnostic) sendCapacityWarnResult(name string, resource string, nTotal int, curVal, targetVal string, score float64) {
	d.Score -= score
	d.result <- &diagnose.Result{
		ObjName: name,
		Level:   diagnose.HealthyLevelWarn,
		Score:   score,

		Title: d.Translator.Message("capacity-title", map[string]interface{}{
			"Resource": resource,
		}),

		Desc: d.Translator.Message("capacity-desc", map[string]interface{}{
			"NodeName":  name,
			"Resource":  resource,
			"CurValue":  curVal,
			"NodeTotal": nTotal,
		}),

		Proposal: d.Translator.Message("capacity-proposal", map[string]interface{}{
			"NodeName":    name,
			"Resource":    resource,
			"TargetValue": targetVal,
			"NodeTotal":   nTotal,
		}),
	}
}

func (d *Diagnostic) sendCapacityGoodResult(name string, resource string, nTotal int, curVal, targetVal string, score float64) {
	d.result <- &diagnose.Result{
		ObjName: name,
		Level:   diagnose.HealthyLevelGood,

		Title: d.Translator.Message("capacity-title", map[string]interface{}{
			"Resource": resource,
		}),

		Desc: d.Translator.Message("capacity-good-desc", map[string]interface{}{
			"NodeName":  name,
			"Resource":  resource,
			"CurValue":  curVal,
			"NodeTotal": nTotal,
		}),
	}
}

func (d *Diagnostic) targetCapacity() (Capacity, int, error) {
	label := labels.NewSelector()
	req, err := labels.NewRequirement("node-role.kubernetes.io/master", selection.DoesNotExist, nil)
	if err != nil {
		return Capacity{}, 0, err
	}

	label = label.Add(*req)
	nodes, err := d.Cli.CoreV1().Nodes().List(v1.ListOptions{
		LabelSelector: label.String(),
	})
	if err != nil {
		return Capacity{}, 0, err
	}

	nTotal := len(nodes.Items)
	for _, scale := range d.Capacities {
		if scale.MaxNodeTotal > nTotal {
			return scale, nTotal, nil
		}
	}
	return Capacity{}, 0, fmt.Errorf("no target capacity found")
}
