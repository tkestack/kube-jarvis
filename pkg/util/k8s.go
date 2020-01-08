package util

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/resource"
)

func MemQuantityStr(q *resource.Quantity) string {
	return fmt.Sprintf("%.2f", float64(q.Value())/1024/1024/1024) + "GB"
}

func CpuQuantityStr(q *resource.Quantity) string {
	return fmt.Sprintf("%.2f", float64(q.MilliValue())/1000) + "Core"
}
