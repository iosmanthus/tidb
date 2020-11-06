// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package executor

import (
	"context"

	"github.com/pingcap/errors"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/tidb/bindinfo"
	"github.com/pingcap/tidb/domain"
	plannercore "github.com/pingcap/tidb/planner/core"
	"github.com/pingcap/tidb/util/chunk"
	"github.com/pingcap/tidb/util/hint"
)

// SQLBindExec represents a bind executor.
type SQLBindExec struct {
	baseExecutor

	sqlBindOp    plannercore.SQLBindOpType
	bindingTp    ast.BindingType
	normdOrigSQL string
	stmtDigest   string
	bindSQL      string
	hints        *hint.HintsSet
	charset      string
	collation    string
	db           string
	isGlobal     bool
	bindAst      ast.StmtNode
}

// Next implements the Executor Next interface.
func (e *SQLBindExec) Next(ctx context.Context, req *chunk.Chunk) error {
	req.Reset()
	switch e.sqlBindOp {
	case plannercore.OpSQLBindCreate:
		return e.createSQLBind()
	case plannercore.OpSQLBindDrop:
		return e.dropSQLBind()
	case plannercore.OpFlushBindings:
		return e.flushBindings()
	case plannercore.OpCaptureBindings:
		e.captureBindings()
	case plannercore.OpEvolveBindings:
		return e.evolveBindings()
	case plannercore.OpReloadBindings:
		return e.reloadBindings()
	default:
		return errors.Errorf("unsupported SQL bind operation: %v", e.sqlBindOp)
	}
	return nil
}

func (e *SQLBindExec) dropSQLBind() error {
	var bindInfo *bindinfo.Binding
	if e.bindSQL != "" || e.hints != nil {
		bindInfo = &bindinfo.Binding{
			BindSQL:   e.bindSQL,
			BindingTp: e.bindingTp,
			Hint:      e.hints,
			Charset:   e.charset,
			Collation: e.collation,
		}
	}
	record := &bindinfo.BindRecord{
		StmtDigest:  e.stmtDigest,
		OriginalSQL: e.normdOrigSQL,
		Db:          e.db,
	}
	if !e.isGlobal {
		handle := e.ctx.Value(bindinfo.SessionBindInfoKeyType).(*bindinfo.SessionHandle)
		return handle.DropBindRecord(e.ctx, record, bindInfo)
	}
	return domain.GetDomain(e.ctx).BindHandle().DropBindRecord(record, bindInfo)
}

func (e *SQLBindExec) createSQLBind() error {
	bindInfo := bindinfo.Binding{
		BindSQL:   e.bindSQL,
		BindingTp: e.bindingTp,
		Hint:      e.hints,
		Charset:   e.charset,
		Collation: e.collation,
		Status:    bindinfo.Using,
		Source:    bindinfo.Manual,
	}
	record := &bindinfo.BindRecord{
		StmtDigest:  e.stmtDigest,
		OriginalSQL: e.normdOrigSQL,
		Db:          e.db,
		Bindings:    []bindinfo.Binding{bindInfo},
	}
	if !e.isGlobal {
		handle := e.ctx.Value(bindinfo.SessionBindInfoKeyType).(*bindinfo.SessionHandle)
		return handle.CreateBindRecord(e.ctx, record)
	}
	return domain.GetDomain(e.ctx).BindHandle().CreateBindRecord(e.ctx, record)
}

func (e *SQLBindExec) flushBindings() error {
	return domain.GetDomain(e.ctx).BindHandle().FlushBindings()
}

func (e *SQLBindExec) captureBindings() {
	domain.GetDomain(e.ctx).BindHandle().CaptureBaselines()
}

func (e *SQLBindExec) evolveBindings() error {
	return domain.GetDomain(e.ctx).BindHandle().HandleEvolvePlanTask(e.ctx, true)
}

func (e *SQLBindExec) reloadBindings() error {
	return domain.GetDomain(e.ctx).BindHandle().ReloadBindings()
}
