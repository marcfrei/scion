// Copyright 2019 Anapaya Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pathdb

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/ctrl/seg"
	dblib "github.com/scionproto/scion/go/lib/infra/modules/db"
	"github.com/scionproto/scion/go/lib/pathdb"
	"github.com/scionproto/scion/go/lib/pathdb/query"
	"github.com/scionproto/scion/go/lib/prom"
	"github.com/scionproto/scion/go/lib/tracing"
	"github.com/scionproto/scion/go/pkg/storage"
)

const (
	promNamespace = "pathdb"

	promDBName = "db"
)

type promOp string

const (
	promOpInsert          promOp = "insert"
	promOpInsertHpCfg     promOp = "insert_with_hpcfg"
	promOpDeleteExpired   promOp = "delete_expired"
	promOpGet             promOp = "get"
	promOpGetAll          promOp = "get_all"
	promOpInsertNextQuery promOp = "insert_next_query"
	promOpGetNextQuery    promOp = "get_next_query"

	promOpBeginTx    promOp = "tx_begin"
	promOpCommitTx   promOp = "tx_commit"
	promOpRollbackTx promOp = "tx_rollback"
)

var (
	queriesTotal *prometheus.CounterVec
	resultsTotal *prometheus.CounterVec

	initMetricsOnce sync.Once
)

type Config struct {
	Driver string
}

func initMetrics() {
	initMetricsOnce.Do(func() {
		// Cardinality: X (dbName) * 13 (len(all ops))
		queriesTotal = prom.NewCounterVec(promNamespace, "", "queries_total",
			"Total queries to the database.", []string{promDBName, prom.LabelOperation})
		// Cardinality: X (dbNmae) * 13 (len(all ops)) * Y (len(all results))
		resultsTotal = prom.NewCounterVec(promNamespace, "", "results_total",
			"The results of the pathdb ops.",
			[]string{promDBName, prom.LabelResult, prom.LabelOperation})
	})
}

// WrapDB wraps the given PathDB into one that also exports metrics. dbName will
// be added as a label to all metrics, so that multiple path DBs can be
// differentiated.
func WrapDB(pathDB storage.PathDB, cfg Config) storage.PathDB {
	initMetrics()
	labels := prometheus.Labels{promDBName: cfg.Driver}
	return &metricsPathDB{
		metricsExecutor: &metricsExecutor{
			pathDB: pathDB,
			metrics: &counters{
				queriesTotal: queriesTotal.MustCurryWith(labels),
				resultsTotal: resultsTotal.MustCurryWith(labels),
			},
		},
		db: pathDB,
	}
}

type counters struct {
	queriesTotal *prometheus.CounterVec
	resultsTotal *prometheus.CounterVec
}

func (c *counters) Observe(ctx context.Context, op promOp, action func(ctx context.Context) error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("pathdb.%s", string(op)))
	defer span.Finish()
	c.queriesTotal.WithLabelValues(string(op)).Inc()
	err := action(ctx)

	label := dblib.ErrToMetricLabel(err)
	tracing.Error(span, err)
	tracing.ResultLabel(span, label)

	c.resultsTotal.WithLabelValues(label, string(op)).Inc()
}

var _ (storage.PathDB) = (*metricsPathDB)(nil)

// metricsPathDB is a PathDB wrapper that exports the counts of operations as prometheus metrics.
type metricsPathDB struct {
	*metricsExecutor
	// db is only needed to have BeginTransaction method.
	db storage.PathDB
}

func (db *metricsPathDB) BeginTransaction(ctx context.Context,
	opts *sql.TxOptions) (pathdb.Transaction, error) {

	var tx pathdb.Transaction
	var err error
	db.metricsExecutor.metrics.Observe(ctx, promOpBeginTx, func(ctx context.Context) error {
		tx, err = db.db.BeginTransaction(ctx, opts)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &metricsTransaction{
		tx:  tx,
		ctx: ctx,
		metricsExecutor: &metricsExecutor{
			pathDB:  tx,
			metrics: db.metricsExecutor.metrics,
		},
	}, err
}

func (db *metricsPathDB) Close() error {
	return db.db.Close()
}

var _ (pathdb.Transaction) = (*metricsTransaction)(nil)

type metricsTransaction struct {
	*metricsExecutor
	tx  pathdb.Transaction
	ctx context.Context
}

func (tx *metricsTransaction) Commit() error {
	var err error
	tx.metrics.Observe(tx.ctx, promOpCommitTx, func(_ context.Context) error {
		err = tx.tx.Commit()
		return err
	})
	return err
}

func (tx *metricsTransaction) Rollback() error {
	var err error
	tx.metrics.Observe(tx.ctx, promOpRollbackTx, func(_ context.Context) error {
		err = tx.tx.Rollback()
		if err == sql.ErrTxDone {
			return nil
		}
		return err
	})
	return err
}

var _ (pathdb.ReadWrite) = (*metricsExecutor)(nil)

type metricsExecutor struct {
	pathDB  pathdb.ReadWrite
	metrics *counters
}

func (db *metricsExecutor) Insert(ctx context.Context, meta *seg.Meta) (pathdb.InsertStats, error) {
	var cnt pathdb.InsertStats
	var err error
	db.metrics.Observe(ctx, promOpInsert, func(ctx context.Context) error {
		cnt, err = db.pathDB.Insert(ctx, meta)
		return err
	})
	return cnt, err
}

func (db *metricsExecutor) InsertWithHPGroupIDs(ctx context.Context,
	meta *seg.Meta, hpGroupIDs []uint64) (pathdb.InsertStats, error) {

	var cnt pathdb.InsertStats
	var err error
	db.metrics.Observe(ctx, promOpInsertHpCfg, func(ctx context.Context) error {
		cnt, err = db.pathDB.InsertWithHPGroupIDs(ctx, meta, hpGroupIDs)
		return err
	})
	return cnt, err
}

func (db *metricsExecutor) DeleteExpired(ctx context.Context, now time.Time) (int, error) {
	var cnt int
	var err error
	db.metrics.Observe(ctx, promOpDeleteExpired, func(ctx context.Context) error {
		cnt, err = db.pathDB.DeleteExpired(ctx, now)
		return err
	})
	return cnt, err
}

func (db *metricsExecutor) Get(ctx context.Context, params *query.Params) (query.Results, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("pathdb.%s", string(promOpGet)))
	defer span.Finish()
	if params != nil {
		span.SetTag("query.starts_at", params.StartsAt)
		span.SetTag("query.ends_at", params.EndsAt)
	}

	db.metrics.queriesTotal.WithLabelValues(string(promOpGet)).Inc()
	res, err := db.pathDB.Get(ctx, params)

	label := dblib.ErrToMetricLabel(err)
	tracing.Error(span, err)
	tracing.ResultLabel(span, label)
	span.SetTag("result.size", len(res))
	db.metrics.resultsTotal.WithLabelValues(label, string(promOpGet)).Inc()
	return res, err
}

func (db *metricsExecutor) GetAll(ctx context.Context) (query.Results, error) {
	var res query.Results
	var err error
	db.metrics.Observe(ctx, promOpGetAll, func(ctx context.Context) error {
		res, err = db.pathDB.GetAll(ctx)
		return err
	})
	return res, err
}

func (db *metricsExecutor) InsertNextQuery(ctx context.Context,
	src, dst addr.IA, nextQuery time.Time) (bool, error) {

	var ok bool
	var err error
	db.metrics.Observe(ctx, promOpInsertNextQuery, func(ctx context.Context) error {
		ok, err = db.pathDB.InsertNextQuery(ctx, src, dst, nextQuery)
		return err
	})
	return ok, err
}

func (db *metricsExecutor) GetNextQuery(ctx context.Context, src, dst addr.IA) (time.Time, error) {

	var t time.Time
	var err error
	db.metrics.Observe(ctx, promOpGetNextQuery, func(ctx context.Context) error {
		t, err = db.pathDB.GetNextQuery(ctx, src, dst)
		return err
	})
	return t, err
}
