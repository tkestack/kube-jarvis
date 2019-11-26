package capacity

import "k8s.io/apimachinery/pkg/api/resource"

var (
	DefCapacities = []Capacity{
		{
			MaxNodeTotal: 5,

			Memory:  resource.MustParse("8000000Ki"),
			CpuCore: resource.MustParse("4000m"),
		},
		{
			MaxNodeTotal: 20,

			Memory:  resource.MustParse("16000000Ki"),
			CpuCore: resource.MustParse("4000m"),
		},
		{
			MaxNodeTotal: 100,

			Memory:  resource.MustParse("32000000Ki"),
			CpuCore: resource.MustParse("8000m"),
		},
		{
			MaxNodeTotal: 200,

			Memory:  resource.MustParse("64000000Ki"),
			CpuCore: resource.MustParse("16000m"),
		},
		{
			MaxNodeTotal: 100000,
			Memory:       resource.MustParse("128000000Ki"),
			CpuCore:      resource.MustParse("16000m"),
		},
	}
)
