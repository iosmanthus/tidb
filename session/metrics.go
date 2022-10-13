// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package session

import (
	"github.com/pingcap/log"
	"github.com/pingcap/tidb/br/pkg/task"
	"github.com/pingcap/tidb/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	statementPerTransactionPessimisticOK    prometheus.Observer
	statementPerTransactionPessimisticError prometheus.Observer
	statementPerTransactionOptimisticOK     prometheus.Observer
	statementPerTransactionOptimisticError  prometheus.Observer
	transactionDurationPessimisticCommit    prometheus.Observer
	transactionDurationPessimisticAbort     prometheus.Observer
	transactionDurationOptimisticCommit     prometheus.Observer
	transactionDurationOptimisticAbort      prometheus.Observer

	sessionExecuteCompileDurationInternal prometheus.Observer
	sessionExecuteCompileDurationGeneral  prometheus.Observer
	sessionExecuteParseDurationInternal   prometheus.Observer
	sessionExecuteParseDurationGeneral    prometheus.Observer

	telemetryCTEUsage               = metrics.TelemetrySQLCTECnt
	telemetryMultiSchemaChangeUsage = metrics.TelemetryMultiSchemaChangeCnt
)

func init() {
	InitMetricsVars()
	task.RegisterSessionMetrics(InitMetricsVars)
}

// InitMetricsVars init session metrics counter
func InitMetricsVars() {
	log.Info("init session metrics")
	statementPerTransactionPessimisticOK = metrics.StatementPerTransaction.WithLabelValues(metrics.LblPessimistic, metrics.LblOK)
	statementPerTransactionPessimisticError = metrics.StatementPerTransaction.WithLabelValues(metrics.LblPessimistic, metrics.LblError)
	statementPerTransactionOptimisticOK = metrics.StatementPerTransaction.WithLabelValues(metrics.LblOptimistic, metrics.LblOK)
	statementPerTransactionOptimisticError = metrics.StatementPerTransaction.WithLabelValues(metrics.LblOptimistic, metrics.LblError)
	transactionDurationPessimisticCommit = metrics.TransactionDuration.WithLabelValues(metrics.LblPessimistic, metrics.LblCommit)
	transactionDurationPessimisticAbort = metrics.TransactionDuration.WithLabelValues(metrics.LblPessimistic, metrics.LblAbort)
	transactionDurationOptimisticCommit = metrics.TransactionDuration.WithLabelValues(metrics.LblOptimistic, metrics.LblCommit)
	transactionDurationOptimisticAbort = metrics.TransactionDuration.WithLabelValues(metrics.LblOptimistic, metrics.LblAbort)

	sessionExecuteCompileDurationInternal = metrics.SessionExecuteCompileDuration.WithLabelValues(metrics.LblInternal)
	sessionExecuteCompileDurationGeneral = metrics.SessionExecuteCompileDuration.WithLabelValues(metrics.LblGeneral)
	sessionExecuteParseDurationInternal = metrics.SessionExecuteParseDuration.WithLabelValues(metrics.LblInternal)
	sessionExecuteParseDurationGeneral = metrics.SessionExecuteParseDuration.WithLabelValues(metrics.LblGeneral)
}
