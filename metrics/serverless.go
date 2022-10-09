// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	TagTenantId  = "tenant_id"
	TagProjectId = "project_id"
	TagClusterId = "cluster_id"
)

// Identify serverless cluster. Should be setup before register metrics.
var (
	ServerlessLabels prometheus.Labels
)

func SetServerlessLabels(ServerlessTenantID, ServerlessProjectID, ServerlessClusterID string) {
	ServerlessLabels = make(prometheus.Labels)
	ServerlessLabels[TagTenantId] = ServerlessTenantID
	ServerlessLabels[TagProjectId] = ServerlessProjectID
	ServerlessLabels[TagClusterId] = ServerlessClusterID
}

func NewCounter(opts prometheus.CounterOpts) prometheus.Counter {
	opts.ConstLabels = ServerlessLabels
	return prometheus.NewCounter(opts)
}

func NewCounterVec(opts prometheus.CounterOpts, labelNames []string) *prometheus.CounterVec {
	opts.ConstLabels = ServerlessLabels
	return prometheus.NewCounterVec(opts, labelNames)
}

func NewGauge(opts prometheus.GaugeOpts) prometheus.Gauge {
	opts.ConstLabels = ServerlessLabels
	return prometheus.NewGauge(opts)
}

func NewGaugeVec(opts prometheus.GaugeOpts, labelNames []string) *prometheus.GaugeVec {
	opts.ConstLabels = ServerlessLabels
	return prometheus.NewGaugeVec(opts, labelNames)
}

func NewHistogram(opts prometheus.HistogramOpts) prometheus.Histogram {
	opts.ConstLabels = ServerlessLabels
	return prometheus.NewHistogram(opts)
}

func NewHistogramVec(opts prometheus.HistogramOpts, labelNames []string) *prometheus.HistogramVec {
	opts.ConstLabels = ServerlessLabels
	return prometheus.NewHistogramVec(opts, labelNames)
}
