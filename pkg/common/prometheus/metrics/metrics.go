package metrics

import (
	"github.com/yanglunara/im/pkg/common/prometheus/ginprometheus"
)

func GetImApiGinMetrics() []*ginprometheus.Metric {
	return []*ginprometheus.Metric{ImApiCustomCnt}
}
