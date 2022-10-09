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

package executor

import (
	"github.com/pingcap/tidb/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// metrics option
var (
	executorCounterMergeJoinExec            prometheus.Counter
	executorCountHashJoinExec               prometheus.Counter
	executorCounterHashAggExec              prometheus.Counter
	executorStreamAggExec                   prometheus.Counter
	executorCounterSortExec                 prometheus.Counter
	executorCounterTopNExec                 prometheus.Counter
	executorCounterNestedLoopApplyExec      prometheus.Counter
	executorCounterIndexLookUpJoin          prometheus.Counter
	executorCounterIndexLookUpExecutor      prometheus.Counter
	executorCounterIndexMergeReaderExecutor prometheus.Counter

	fastAnalyzeHistogramSample        prometheus.Observer
	fastAnalyzeHistogramAccessRegions prometheus.Observer
	fastAnalyzeHistogramScanKeys      prometheus.Observer

	sessionExecuteRunDurationInternal prometheus.Observer
	sessionExecuteRunDurationGeneral  prometheus.Observer
	totalTiFlashQuerySuccCounter      prometheus.Counter

	stmtNodeCounterUse       prometheus.Counter
	stmtNodeCounterShow      prometheus.Counter
	stmtNodeCounterBegin     prometheus.Counter
	stmtNodeCounterCommit    prometheus.Counter
	stmtNodeCounterRollback  prometheus.Counter
	stmtNodeCounterInsert    prometheus.Counter
	stmtNodeCounterReplace   prometheus.Counter
	stmtNodeCounterDelete    prometheus.Counter
	stmtNodeCounterUpdate    prometheus.Counter
	stmtNodeCounterSelect    prometheus.Counter
	stmtNodeCounterSavepoint prometheus.Counter

	totalQueryProcHistogramGeneral  prometheus.Observer
	totalCopProcHistogramGeneral    prometheus.Observer
	totalCopWaitHistogramGeneral    prometheus.Observer
	totalQueryProcHistogramInternal prometheus.Observer
	totalCopProcHistogramInternal   prometheus.Observer
	totalCopWaitHistogramInternal   prometheus.Observer

	transactionDurationPessimisticRollback prometheus.Observer
	transactionDurationOptimisticRollback  prometheus.Observer
)

func init() {
	InitMetricsVars()
}

// InitMetricsVars init executor metrics counter
func InitMetricsVars() {
	executorCounterMergeJoinExec = metrics.ExecutorCounter.WithLabelValues("MergeJoinExec")
	executorCountHashJoinExec = metrics.ExecutorCounter.WithLabelValues("HashJoinExec")
	executorCounterHashAggExec = metrics.ExecutorCounter.WithLabelValues("HashAggExec")
	executorStreamAggExec = metrics.ExecutorCounter.WithLabelValues("StreamAggExec")
	executorCounterSortExec = metrics.ExecutorCounter.WithLabelValues("SortExec")
	executorCounterTopNExec = metrics.ExecutorCounter.WithLabelValues("TopNExec")
	executorCounterNestedLoopApplyExec = metrics.ExecutorCounter.WithLabelValues("NestedLoopApplyExec")
	executorCounterIndexLookUpJoin = metrics.ExecutorCounter.WithLabelValues("IndexLookUpJoin")
	executorCounterIndexLookUpExecutor = metrics.ExecutorCounter.WithLabelValues("IndexLookUpExecutor")
	executorCounterIndexMergeReaderExecutor = metrics.ExecutorCounter.WithLabelValues("IndexMergeReaderExecutor")

	fastAnalyzeHistogramSample = metrics.FastAnalyzeHistogram.WithLabelValues(metrics.LblGeneral, "sample")
	fastAnalyzeHistogramAccessRegions = metrics.FastAnalyzeHistogram.WithLabelValues(metrics.LblGeneral, "access_regions")
	fastAnalyzeHistogramScanKeys = metrics.FastAnalyzeHistogram.WithLabelValues(metrics.LblGeneral, "scan_keys")

	sessionExecuteRunDurationInternal = metrics.SessionExecuteRunDuration.WithLabelValues(metrics.LblInternal)
	sessionExecuteRunDurationGeneral = metrics.SessionExecuteRunDuration.WithLabelValues(metrics.LblGeneral)
	totalTiFlashQuerySuccCounter = metrics.TiFlashQueryTotalCounter.WithLabelValues("", metrics.LblOK)

	stmtNodeCounterUse = metrics.StmtNodeCounter.WithLabelValues("Use")
	stmtNodeCounterShow = metrics.StmtNodeCounter.WithLabelValues("Show")
	stmtNodeCounterBegin = metrics.StmtNodeCounter.WithLabelValues("Begin")
	stmtNodeCounterCommit = metrics.StmtNodeCounter.WithLabelValues("Commit")
	stmtNodeCounterRollback = metrics.StmtNodeCounter.WithLabelValues("Rollback")
	stmtNodeCounterInsert = metrics.StmtNodeCounter.WithLabelValues("Insert")
	stmtNodeCounterReplace = metrics.StmtNodeCounter.WithLabelValues("Replace")
	stmtNodeCounterDelete = metrics.StmtNodeCounter.WithLabelValues("Delete")
	stmtNodeCounterUpdate = metrics.StmtNodeCounter.WithLabelValues("Update")
	stmtNodeCounterSelect = metrics.StmtNodeCounter.WithLabelValues("Select")
	stmtNodeCounterSavepoint = metrics.StmtNodeCounter.WithLabelValues("Savepoint")

	totalQueryProcHistogramGeneral = metrics.TotalQueryProcHistogram.WithLabelValues(metrics.LblGeneral)
	totalCopProcHistogramGeneral = metrics.TotalCopProcHistogram.WithLabelValues(metrics.LblGeneral)
	totalCopWaitHistogramGeneral = metrics.TotalCopWaitHistogram.WithLabelValues(metrics.LblGeneral)
	totalQueryProcHistogramInternal = metrics.TotalQueryProcHistogram.WithLabelValues(metrics.LblInternal)
	totalCopProcHistogramInternal = metrics.TotalCopProcHistogram.WithLabelValues(metrics.LblInternal)
	totalCopWaitHistogramInternal = metrics.TotalCopWaitHistogram.WithLabelValues(metrics.LblInternal)

	transactionDurationPessimisticRollback = metrics.TransactionDuration.WithLabelValues(metrics.LblPessimistic, metrics.LblRollback)
	transactionDurationOptimisticRollback = metrics.TransactionDuration.WithLabelValues(metrics.LblOptimistic, metrics.LblRollback)
}
