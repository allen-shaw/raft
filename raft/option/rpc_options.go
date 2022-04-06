package option

import "github.com/AllenShaw19/raft/raft/metrics"

type RpcOptions struct {
	metricRegistry *metrics.MetricRegistry
}

func (r *RpcOptions) MetricRegistry() *metrics.MetricRegistry {
	return r.metricRegistry
}

func (r *RpcOptions) SetMetricRegistry(metricRegistry *metrics.MetricRegistry) {
	r.metricRegistry = metricRegistry
}
