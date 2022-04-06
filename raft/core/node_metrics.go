package core

import "github.com/AllenShaw19/raft/raft/metrics"

type NodeMetrics struct {
	metrics *metrics.MetricRegistry
}

func NewNodeMetrics(enableMetrics bool) *NodeMetrics {
	metrics := &NodeMetrics{}

	return metrics
}

func (n *NodeMetrics) Metrics() *metrics.MetricRegistry {
	return n.metrics
}

func (n *NodeMetrics) SetMetrics(metrics *metrics.MetricRegistry) {
	n.metrics = metrics
}
