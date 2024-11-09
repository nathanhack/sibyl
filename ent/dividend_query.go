// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/dividend"
	"github.com/nathanhack/sibyl/ent/entity"
	"github.com/nathanhack/sibyl/ent/predicate"
)

// DividendQuery is the builder for querying Dividend entities.
type DividendQuery struct {
	config
	ctx        *QueryContext
	order      []dividend.OrderOption
	inters     []Interceptor
	predicates []predicate.Dividend
	withStock  *EntityQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the DividendQuery builder.
func (dq *DividendQuery) Where(ps ...predicate.Dividend) *DividendQuery {
	dq.predicates = append(dq.predicates, ps...)
	return dq
}

// Limit the number of records to be returned by this query.
func (dq *DividendQuery) Limit(limit int) *DividendQuery {
	dq.ctx.Limit = &limit
	return dq
}

// Offset to start from.
func (dq *DividendQuery) Offset(offset int) *DividendQuery {
	dq.ctx.Offset = &offset
	return dq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (dq *DividendQuery) Unique(unique bool) *DividendQuery {
	dq.ctx.Unique = &unique
	return dq
}

// Order specifies how the records should be ordered.
func (dq *DividendQuery) Order(o ...dividend.OrderOption) *DividendQuery {
	dq.order = append(dq.order, o...)
	return dq
}

// QueryStock chains the current query on the "stock" edge.
func (dq *DividendQuery) QueryStock() *EntityQuery {
	query := (&EntityClient{config: dq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := dq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := dq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(dividend.Table, dividend.FieldID, selector),
			sqlgraph.To(entity.Table, entity.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, dividend.StockTable, dividend.StockPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(dq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Dividend entity from the query.
// Returns a *NotFoundError when no Dividend was found.
func (dq *DividendQuery) First(ctx context.Context) (*Dividend, error) {
	nodes, err := dq.Limit(1).All(setContextOp(ctx, dq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{dividend.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (dq *DividendQuery) FirstX(ctx context.Context) *Dividend {
	node, err := dq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Dividend ID from the query.
// Returns a *NotFoundError when no Dividend ID was found.
func (dq *DividendQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = dq.Limit(1).IDs(setContextOp(ctx, dq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{dividend.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (dq *DividendQuery) FirstIDX(ctx context.Context) int {
	id, err := dq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Dividend entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Dividend entity is found.
// Returns a *NotFoundError when no Dividend entities are found.
func (dq *DividendQuery) Only(ctx context.Context) (*Dividend, error) {
	nodes, err := dq.Limit(2).All(setContextOp(ctx, dq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{dividend.Label}
	default:
		return nil, &NotSingularError{dividend.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (dq *DividendQuery) OnlyX(ctx context.Context) *Dividend {
	node, err := dq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Dividend ID in the query.
// Returns a *NotSingularError when more than one Dividend ID is found.
// Returns a *NotFoundError when no entities are found.
func (dq *DividendQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = dq.Limit(2).IDs(setContextOp(ctx, dq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{dividend.Label}
	default:
		err = &NotSingularError{dividend.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (dq *DividendQuery) OnlyIDX(ctx context.Context) int {
	id, err := dq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Dividends.
func (dq *DividendQuery) All(ctx context.Context) ([]*Dividend, error) {
	ctx = setContextOp(ctx, dq.ctx, ent.OpQueryAll)
	if err := dq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Dividend, *DividendQuery]()
	return withInterceptors[[]*Dividend](ctx, dq, qr, dq.inters)
}

// AllX is like All, but panics if an error occurs.
func (dq *DividendQuery) AllX(ctx context.Context) []*Dividend {
	nodes, err := dq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Dividend IDs.
func (dq *DividendQuery) IDs(ctx context.Context) (ids []int, err error) {
	if dq.ctx.Unique == nil && dq.path != nil {
		dq.Unique(true)
	}
	ctx = setContextOp(ctx, dq.ctx, ent.OpQueryIDs)
	if err = dq.Select(dividend.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (dq *DividendQuery) IDsX(ctx context.Context) []int {
	ids, err := dq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (dq *DividendQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, dq.ctx, ent.OpQueryCount)
	if err := dq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, dq, querierCount[*DividendQuery](), dq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (dq *DividendQuery) CountX(ctx context.Context) int {
	count, err := dq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (dq *DividendQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, dq.ctx, ent.OpQueryExist)
	switch _, err := dq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (dq *DividendQuery) ExistX(ctx context.Context) bool {
	exist, err := dq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the DividendQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (dq *DividendQuery) Clone() *DividendQuery {
	if dq == nil {
		return nil
	}
	return &DividendQuery{
		config:     dq.config,
		ctx:        dq.ctx.Clone(),
		order:      append([]dividend.OrderOption{}, dq.order...),
		inters:     append([]Interceptor{}, dq.inters...),
		predicates: append([]predicate.Dividend{}, dq.predicates...),
		withStock:  dq.withStock.Clone(),
		// clone intermediate query.
		sql:  dq.sql.Clone(),
		path: dq.path,
	}
}

// WithStock tells the query-builder to eager-load the nodes that are connected to
// the "stock" edge. The optional arguments are used to configure the query builder of the edge.
func (dq *DividendQuery) WithStock(opts ...func(*EntityQuery)) *DividendQuery {
	query := (&EntityClient{config: dq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	dq.withStock = query
	return dq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Rate float64 `json:"rate,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Dividend.Query().
//		GroupBy(dividend.FieldRate).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (dq *DividendQuery) GroupBy(field string, fields ...string) *DividendGroupBy {
	dq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &DividendGroupBy{build: dq}
	grbuild.flds = &dq.ctx.Fields
	grbuild.label = dividend.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Rate float64 `json:"rate,omitempty"`
//	}
//
//	client.Dividend.Query().
//		Select(dividend.FieldRate).
//		Scan(ctx, &v)
func (dq *DividendQuery) Select(fields ...string) *DividendSelect {
	dq.ctx.Fields = append(dq.ctx.Fields, fields...)
	sbuild := &DividendSelect{DividendQuery: dq}
	sbuild.label = dividend.Label
	sbuild.flds, sbuild.scan = &dq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a DividendSelect configured with the given aggregations.
func (dq *DividendQuery) Aggregate(fns ...AggregateFunc) *DividendSelect {
	return dq.Select().Aggregate(fns...)
}

func (dq *DividendQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range dq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, dq); err != nil {
				return err
			}
		}
	}
	for _, f := range dq.ctx.Fields {
		if !dividend.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if dq.path != nil {
		prev, err := dq.path(ctx)
		if err != nil {
			return err
		}
		dq.sql = prev
	}
	return nil
}

func (dq *DividendQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Dividend, error) {
	var (
		nodes       = []*Dividend{}
		_spec       = dq.querySpec()
		loadedTypes = [1]bool{
			dq.withStock != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Dividend).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Dividend{config: dq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, dq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := dq.withStock; query != nil {
		if err := dq.loadStock(ctx, query, nodes,
			func(n *Dividend) { n.Edges.Stock = []*Entity{} },
			func(n *Dividend, e *Entity) { n.Edges.Stock = append(n.Edges.Stock, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (dq *DividendQuery) loadStock(ctx context.Context, query *EntityQuery, nodes []*Dividend, init func(*Dividend), assign func(*Dividend, *Entity)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[int]*Dividend)
	nids := make(map[int]map[*Dividend]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(dividend.StockTable)
		s.Join(joinT).On(s.C(entity.FieldID), joinT.C(dividend.StockPrimaryKey[0]))
		s.Where(sql.InValues(joinT.C(dividend.StockPrimaryKey[1]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(dividend.StockPrimaryKey[1]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	qr := QuerierFunc(func(ctx context.Context, q Query) (Value, error) {
		return query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
			assign := spec.Assign
			values := spec.ScanValues
			spec.ScanValues = func(columns []string) ([]any, error) {
				values, err := values(columns[1:])
				if err != nil {
					return nil, err
				}
				return append([]any{new(sql.NullInt64)}, values...), nil
			}
			spec.Assign = func(columns []string, values []any) error {
				outValue := int(values[0].(*sql.NullInt64).Int64)
				inValue := int(values[1].(*sql.NullInt64).Int64)
				if nids[inValue] == nil {
					nids[inValue] = map[*Dividend]struct{}{byID[outValue]: {}}
					return assign(columns[1:], values[1:])
				}
				nids[inValue][byID[outValue]] = struct{}{}
				return nil
			}
		})
	})
	neighbors, err := withInterceptors[[]*Entity](ctx, query, qr, query.inters)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "stock" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}

func (dq *DividendQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := dq.querySpec()
	_spec.Node.Columns = dq.ctx.Fields
	if len(dq.ctx.Fields) > 0 {
		_spec.Unique = dq.ctx.Unique != nil && *dq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, dq.driver, _spec)
}

func (dq *DividendQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(dividend.Table, dividend.Columns, sqlgraph.NewFieldSpec(dividend.FieldID, field.TypeInt))
	_spec.From = dq.sql
	if unique := dq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if dq.path != nil {
		_spec.Unique = true
	}
	if fields := dq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, dividend.FieldID)
		for i := range fields {
			if fields[i] != dividend.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := dq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := dq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := dq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := dq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (dq *DividendQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(dq.driver.Dialect())
	t1 := builder.Table(dividend.Table)
	columns := dq.ctx.Fields
	if len(columns) == 0 {
		columns = dividend.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if dq.sql != nil {
		selector = dq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if dq.ctx.Unique != nil && *dq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range dq.predicates {
		p(selector)
	}
	for _, p := range dq.order {
		p(selector)
	}
	if offset := dq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := dq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// DividendGroupBy is the group-by builder for Dividend entities.
type DividendGroupBy struct {
	selector
	build *DividendQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (dgb *DividendGroupBy) Aggregate(fns ...AggregateFunc) *DividendGroupBy {
	dgb.fns = append(dgb.fns, fns...)
	return dgb
}

// Scan applies the selector query and scans the result into the given value.
func (dgb *DividendGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, dgb.build.ctx, ent.OpQueryGroupBy)
	if err := dgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*DividendQuery, *DividendGroupBy](ctx, dgb.build, dgb, dgb.build.inters, v)
}

func (dgb *DividendGroupBy) sqlScan(ctx context.Context, root *DividendQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(dgb.fns))
	for _, fn := range dgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*dgb.flds)+len(dgb.fns))
		for _, f := range *dgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*dgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := dgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// DividendSelect is the builder for selecting fields of Dividend entities.
type DividendSelect struct {
	*DividendQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (ds *DividendSelect) Aggregate(fns ...AggregateFunc) *DividendSelect {
	ds.fns = append(ds.fns, fns...)
	return ds
}

// Scan applies the selector query and scans the result into the given value.
func (ds *DividendSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, ds.ctx, ent.OpQuerySelect)
	if err := ds.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*DividendQuery, *DividendSelect](ctx, ds.DividendQuery, ds, ds.inters, v)
}

func (ds *DividendSelect) sqlScan(ctx context.Context, root *DividendQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(ds.fns))
	for _, fn := range ds.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*ds.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ds.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
