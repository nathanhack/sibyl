// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/entity"
	"github.com/nathanhack/sibyl/ent/exchange"
	"github.com/nathanhack/sibyl/ent/predicate"
)

// ExchangeQuery is the builder for querying Exchange entities.
type ExchangeQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	inters     []Interceptor
	predicates []predicate.Exchange
	withStocks *EntityQuery
	withFKs    bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the ExchangeQuery builder.
func (eq *ExchangeQuery) Where(ps ...predicate.Exchange) *ExchangeQuery {
	eq.predicates = append(eq.predicates, ps...)
	return eq
}

// Limit the number of records to be returned by this query.
func (eq *ExchangeQuery) Limit(limit int) *ExchangeQuery {
	eq.limit = &limit
	return eq
}

// Offset to start from.
func (eq *ExchangeQuery) Offset(offset int) *ExchangeQuery {
	eq.offset = &offset
	return eq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (eq *ExchangeQuery) Unique(unique bool) *ExchangeQuery {
	eq.unique = &unique
	return eq
}

// Order specifies how the records should be ordered.
func (eq *ExchangeQuery) Order(o ...OrderFunc) *ExchangeQuery {
	eq.order = append(eq.order, o...)
	return eq
}

// QueryStocks chains the current query on the "stocks" edge.
func (eq *ExchangeQuery) QueryStocks() *EntityQuery {
	query := (&EntityClient{config: eq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := eq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := eq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(exchange.Table, exchange.FieldID, selector),
			sqlgraph.To(entity.Table, entity.FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, exchange.StocksTable, exchange.StocksPrimaryKey...),
		)
		fromU = sqlgraph.SetNeighbors(eq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first Exchange entity from the query.
// Returns a *NotFoundError when no Exchange was found.
func (eq *ExchangeQuery) First(ctx context.Context) (*Exchange, error) {
	nodes, err := eq.Limit(1).All(newQueryContext(ctx, TypeExchange, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{exchange.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (eq *ExchangeQuery) FirstX(ctx context.Context) *Exchange {
	node, err := eq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Exchange ID from the query.
// Returns a *NotFoundError when no Exchange ID was found.
func (eq *ExchangeQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = eq.Limit(1).IDs(newQueryContext(ctx, TypeExchange, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{exchange.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (eq *ExchangeQuery) FirstIDX(ctx context.Context) int {
	id, err := eq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Exchange entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Exchange entity is found.
// Returns a *NotFoundError when no Exchange entities are found.
func (eq *ExchangeQuery) Only(ctx context.Context) (*Exchange, error) {
	nodes, err := eq.Limit(2).All(newQueryContext(ctx, TypeExchange, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{exchange.Label}
	default:
		return nil, &NotSingularError{exchange.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (eq *ExchangeQuery) OnlyX(ctx context.Context) *Exchange {
	node, err := eq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Exchange ID in the query.
// Returns a *NotSingularError when more than one Exchange ID is found.
// Returns a *NotFoundError when no entities are found.
func (eq *ExchangeQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = eq.Limit(2).IDs(newQueryContext(ctx, TypeExchange, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{exchange.Label}
	default:
		err = &NotSingularError{exchange.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (eq *ExchangeQuery) OnlyIDX(ctx context.Context) int {
	id, err := eq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Exchanges.
func (eq *ExchangeQuery) All(ctx context.Context) ([]*Exchange, error) {
	ctx = newQueryContext(ctx, TypeExchange, "All")
	if err := eq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Exchange, *ExchangeQuery]()
	return withInterceptors[[]*Exchange](ctx, eq, qr, eq.inters)
}

// AllX is like All, but panics if an error occurs.
func (eq *ExchangeQuery) AllX(ctx context.Context) []*Exchange {
	nodes, err := eq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Exchange IDs.
func (eq *ExchangeQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	ctx = newQueryContext(ctx, TypeExchange, "IDs")
	if err := eq.Select(exchange.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (eq *ExchangeQuery) IDsX(ctx context.Context) []int {
	ids, err := eq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (eq *ExchangeQuery) Count(ctx context.Context) (int, error) {
	ctx = newQueryContext(ctx, TypeExchange, "Count")
	if err := eq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, eq, querierCount[*ExchangeQuery](), eq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (eq *ExchangeQuery) CountX(ctx context.Context) int {
	count, err := eq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (eq *ExchangeQuery) Exist(ctx context.Context) (bool, error) {
	ctx = newQueryContext(ctx, TypeExchange, "Exist")
	switch _, err := eq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (eq *ExchangeQuery) ExistX(ctx context.Context) bool {
	exist, err := eq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the ExchangeQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (eq *ExchangeQuery) Clone() *ExchangeQuery {
	if eq == nil {
		return nil
	}
	return &ExchangeQuery{
		config:     eq.config,
		limit:      eq.limit,
		offset:     eq.offset,
		order:      append([]OrderFunc{}, eq.order...),
		inters:     append([]Interceptor{}, eq.inters...),
		predicates: append([]predicate.Exchange{}, eq.predicates...),
		withStocks: eq.withStocks.Clone(),
		// clone intermediate query.
		sql:    eq.sql.Clone(),
		path:   eq.path,
		unique: eq.unique,
	}
}

// WithStocks tells the query-builder to eager-load the nodes that are connected to
// the "stocks" edge. The optional arguments are used to configure the query builder of the edge.
func (eq *ExchangeQuery) WithStocks(opts ...func(*EntityQuery)) *ExchangeQuery {
	query := (&EntityClient{config: eq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	eq.withStocks = query
	return eq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Code string `json:"code,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Exchange.Query().
//		GroupBy(exchange.FieldCode).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (eq *ExchangeQuery) GroupBy(field string, fields ...string) *ExchangeGroupBy {
	eq.fields = append([]string{field}, fields...)
	grbuild := &ExchangeGroupBy{build: eq}
	grbuild.flds = &eq.fields
	grbuild.label = exchange.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Code string `json:"code,omitempty"`
//	}
//
//	client.Exchange.Query().
//		Select(exchange.FieldCode).
//		Scan(ctx, &v)
func (eq *ExchangeQuery) Select(fields ...string) *ExchangeSelect {
	eq.fields = append(eq.fields, fields...)
	sbuild := &ExchangeSelect{ExchangeQuery: eq}
	sbuild.label = exchange.Label
	sbuild.flds, sbuild.scan = &eq.fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a ExchangeSelect configured with the given aggregations.
func (eq *ExchangeQuery) Aggregate(fns ...AggregateFunc) *ExchangeSelect {
	return eq.Select().Aggregate(fns...)
}

func (eq *ExchangeQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range eq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, eq); err != nil {
				return err
			}
		}
	}
	for _, f := range eq.fields {
		if !exchange.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if eq.path != nil {
		prev, err := eq.path(ctx)
		if err != nil {
			return err
		}
		eq.sql = prev
	}
	return nil
}

func (eq *ExchangeQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Exchange, error) {
	var (
		nodes       = []*Exchange{}
		withFKs     = eq.withFKs
		_spec       = eq.querySpec()
		loadedTypes = [1]bool{
			eq.withStocks != nil,
		}
	)
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, exchange.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Exchange).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Exchange{config: eq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, eq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := eq.withStocks; query != nil {
		if err := eq.loadStocks(ctx, query, nodes,
			func(n *Exchange) { n.Edges.Stocks = []*Entity{} },
			func(n *Exchange, e *Entity) { n.Edges.Stocks = append(n.Edges.Stocks, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (eq *ExchangeQuery) loadStocks(ctx context.Context, query *EntityQuery, nodes []*Exchange, init func(*Exchange), assign func(*Exchange, *Entity)) error {
	edgeIDs := make([]driver.Value, len(nodes))
	byID := make(map[int]*Exchange)
	nids := make(map[int]map[*Exchange]struct{})
	for i, node := range nodes {
		edgeIDs[i] = node.ID
		byID[node.ID] = node
		if init != nil {
			init(node)
		}
	}
	query.Where(func(s *sql.Selector) {
		joinT := sql.Table(exchange.StocksTable)
		s.Join(joinT).On(s.C(entity.FieldID), joinT.C(exchange.StocksPrimaryKey[0]))
		s.Where(sql.InValues(joinT.C(exchange.StocksPrimaryKey[1]), edgeIDs...))
		columns := s.SelectedColumns()
		s.Select(joinT.C(exchange.StocksPrimaryKey[1]))
		s.AppendSelect(columns...)
		s.SetDistinct(false)
	})
	if err := query.prepareQuery(ctx); err != nil {
		return err
	}
	neighbors, err := query.sqlAll(ctx, func(_ context.Context, spec *sqlgraph.QuerySpec) {
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
				nids[inValue] = map[*Exchange]struct{}{byID[outValue]: {}}
				return assign(columns[1:], values[1:])
			}
			nids[inValue][byID[outValue]] = struct{}{}
			return nil
		}
	})
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected "stocks" node returned %v`, n.ID)
		}
		for kn := range nodes {
			assign(kn, n)
		}
	}
	return nil
}

func (eq *ExchangeQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := eq.querySpec()
	_spec.Node.Columns = eq.fields
	if len(eq.fields) > 0 {
		_spec.Unique = eq.unique != nil && *eq.unique
	}
	return sqlgraph.CountNodes(ctx, eq.driver, _spec)
}

func (eq *ExchangeQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   exchange.Table,
			Columns: exchange.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: exchange.FieldID,
			},
		},
		From:   eq.sql,
		Unique: true,
	}
	if unique := eq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := eq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, exchange.FieldID)
		for i := range fields {
			if fields[i] != exchange.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := eq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := eq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := eq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := eq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (eq *ExchangeQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(eq.driver.Dialect())
	t1 := builder.Table(exchange.Table)
	columns := eq.fields
	if len(columns) == 0 {
		columns = exchange.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if eq.sql != nil {
		selector = eq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if eq.unique != nil && *eq.unique {
		selector.Distinct()
	}
	for _, p := range eq.predicates {
		p(selector)
	}
	for _, p := range eq.order {
		p(selector)
	}
	if offset := eq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := eq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// ExchangeGroupBy is the group-by builder for Exchange entities.
type ExchangeGroupBy struct {
	selector
	build *ExchangeQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (egb *ExchangeGroupBy) Aggregate(fns ...AggregateFunc) *ExchangeGroupBy {
	egb.fns = append(egb.fns, fns...)
	return egb
}

// Scan applies the selector query and scans the result into the given value.
func (egb *ExchangeGroupBy) Scan(ctx context.Context, v any) error {
	ctx = newQueryContext(ctx, TypeExchange, "GroupBy")
	if err := egb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ExchangeQuery, *ExchangeGroupBy](ctx, egb.build, egb, egb.build.inters, v)
}

func (egb *ExchangeGroupBy) sqlScan(ctx context.Context, root *ExchangeQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(egb.fns))
	for _, fn := range egb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*egb.flds)+len(egb.fns))
		for _, f := range *egb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*egb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := egb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// ExchangeSelect is the builder for selecting fields of Exchange entities.
type ExchangeSelect struct {
	*ExchangeQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (es *ExchangeSelect) Aggregate(fns ...AggregateFunc) *ExchangeSelect {
	es.fns = append(es.fns, fns...)
	return es
}

// Scan applies the selector query and scans the result into the given value.
func (es *ExchangeSelect) Scan(ctx context.Context, v any) error {
	ctx = newQueryContext(ctx, TypeExchange, "Select")
	if err := es.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*ExchangeQuery, *ExchangeSelect](ctx, es.ExchangeQuery, es, es.inters, v)
}

func (es *ExchangeSelect) sqlScan(ctx context.Context, root *ExchangeQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(es.fns))
	for _, fn := range es.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*es.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := es.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}