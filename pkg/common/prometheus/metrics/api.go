package metrics

import "github.com/yanglunara/im/pkg/common/prometheus/ginprometheus"

var (
	ImApiCustomCnt = &ginprometheus.Metric{
		Name:        "custom_cnt",
		Description: "custom counter events for im-Api",
		Type:        "counter_vec",
		Args:        []string{"label_one", "label_two"},
	}
)
