package capacity

import (
	"context"
	"fmt"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"tkestack.io/kube-jarvis/pkg/plugins/diagnose"
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
	// MaxNodeTotal indicate the max node number of this master scale
	MaxNodeTotal int
}

// Diagnostic check whether the resources are sufficient for a specific size cluster
type Diagnostic struct {
	*diagnose.MetaData
	result     chan *diagnose.Result
	Capacities []Capacity
	param      *diagnose.StartDiagnoseParam
}

// NewDiagnostic return a master-node diagnostic
func NewDiagnostic(meta *diagnose.MetaData) diagnose.Diagnostic {
	return &Diagnostic{
		result:     make(chan *diagnose.Result, 100),
		MetaData:   meta,
		Capacities: DefCapacities,
	}
}

// Init do initialization
func (d *Diagnostic) Init() error {
	return nil
}

// StartDiagnose return a result chan that will output results
func (d *Diagnostic) StartDiagnose(ctx context.Context, param diagnose.StartDiagnoseParam) chan *diagnose.Result {
	d.param = &param
	go func() {
		defer diagnose.CommonDeafer(d.result)
		d.diagnoseCapacity(ctx)
	}()
	return d.result
}

func (d *Diagnostic) diagnoseCapacity(ctx context.Context) {
	nodeTotal := 0
	masters := make([]v12.Node, 0)
	for _, n := range d.param.Resources.Nodes.Items {
		for k := range n.Labels {
			if k == "node-role.kubernetes.io/master" {
				masters = append(masters, n)
				continue
			}
		}
		nodeTotal++
	}

	scale, err := d.targetCapacity(nodeTotal)
	if err != nil {
		d.result <- &diagnose.Result{Error: err}
		return
	}

	for _, m := range masters {
		if m.Status.Capacity.Cpu().Cmp(scale.CpuCore) < 0 {
			d.sendCapacityWarnResult(m.Name, "Cpu", nodeTotal, m.Status.Capacity.Cpu().String(), scale.CpuCore.String())
		} else {
			d.sendCapacityGoodResult(m.Name, "Cpu", nodeTotal, m.Status.Capacity.Cpu().String(), scale.CpuCore.String())
		}

		if m.Status.Capacity.Memory().Cmp(scale.Memory) < 0 {
			d.sendCapacityWarnResult(m.Name, "Memory", nodeTotal, m.Status.Capacity.Memory().String(), scale.Memory.String())
		} else {
			d.sendCapacityGoodResult(m.Name, "Memory", nodeTotal, m.Status.Capacity.Memory().String(), scale.Memory.String())
		}
	}
}

func (d *Diagnostic) sendCapacityWarnResult(name string, resource string, nTotal int, curVal, targetVal string) {
	d.result <- &diagnose.Result{
		ObjName: name,
		Level:   diagnose.HealthyLevelWarn,

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

func (d *Diagnostic) sendCapacityGoodResult(name string, resource string, nTotal int, curVal, targetVal string) {
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

func (d *Diagnostic) targetCapacity(nTotal int) (*Capacity, error) {
	for _, scale := range d.Capacities {
		if scale.MaxNodeTotal > nTotal {
			return &scale, nil
		}
	}
	return nil, fmt.Errorf("no target capacity found")
}
