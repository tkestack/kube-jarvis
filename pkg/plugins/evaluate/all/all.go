package all

import (
	"github.com/RayHuangCN/Jarvis/pkg/plugins/evaluate"
	"github.com/RayHuangCN/Jarvis/pkg/plugins/evaluate/sum"
)

func init() {
	evaluate.Add("sum", sum.NewEvaluator)
}
