// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/interval"
	"github.com/nathanhack/sibyl/ent/predicate"
	"github.com/nathanhack/sibyl/ent/traderecord"
	"github.com/nathanhack/sibyl/ent/tradetimerange"
)

// TradeTimeRangeUpdate is the builder for updating TradeTimeRange entities.
type TradeTimeRangeUpdate struct {
	config
	hooks    []Hook
	mutation *TradeTimeRangeMutation
}

// Where appends a list predicates to the TradeTimeRangeUpdate builder.
func (ttru *TradeTimeRangeUpdate) Where(ps ...predicate.TradeTimeRange) *TradeTimeRangeUpdate {
	ttru.mutation.Where(ps...)
	return ttru
}

// SetStart sets the "start" field.
func (ttru *TradeTimeRangeUpdate) SetStart(t time.Time) *TradeTimeRangeUpdate {
	ttru.mutation.SetStart(t)
	return ttru
}

// SetNillableStart sets the "start" field if the given value is not nil.
func (ttru *TradeTimeRangeUpdate) SetNillableStart(t *time.Time) *TradeTimeRangeUpdate {
	if t != nil {
		ttru.SetStart(*t)
	}
	return ttru
}

// SetEnd sets the "end" field.
func (ttru *TradeTimeRangeUpdate) SetEnd(t time.Time) *TradeTimeRangeUpdate {
	ttru.mutation.SetEnd(t)
	return ttru
}

// SetNillableEnd sets the "end" field if the given value is not nil.
func (ttru *TradeTimeRangeUpdate) SetNillableEnd(t *time.Time) *TradeTimeRangeUpdate {
	if t != nil {
		ttru.SetEnd(*t)
	}
	return ttru
}

// SetIntervalID sets the "interval_id" field.
func (ttru *TradeTimeRangeUpdate) SetIntervalID(i int) *TradeTimeRangeUpdate {
	ttru.mutation.SetIntervalID(i)
	return ttru
}

// SetNillableIntervalID sets the "interval_id" field if the given value is not nil.
func (ttru *TradeTimeRangeUpdate) SetNillableIntervalID(i *int) *TradeTimeRangeUpdate {
	if i != nil {
		ttru.SetIntervalID(*i)
	}
	return ttru
}

// SetInterval sets the "interval" edge to the Interval entity.
func (ttru *TradeTimeRangeUpdate) SetInterval(i *Interval) *TradeTimeRangeUpdate {
	return ttru.SetIntervalID(i.ID)
}

// AddRecordIDs adds the "records" edge to the TradeRecord entity by IDs.
func (ttru *TradeTimeRangeUpdate) AddRecordIDs(ids ...int) *TradeTimeRangeUpdate {
	ttru.mutation.AddRecordIDs(ids...)
	return ttru
}

// AddRecords adds the "records" edges to the TradeRecord entity.
func (ttru *TradeTimeRangeUpdate) AddRecords(t ...*TradeRecord) *TradeTimeRangeUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return ttru.AddRecordIDs(ids...)
}

// Mutation returns the TradeTimeRangeMutation object of the builder.
func (ttru *TradeTimeRangeUpdate) Mutation() *TradeTimeRangeMutation {
	return ttru.mutation
}

// ClearInterval clears the "interval" edge to the Interval entity.
func (ttru *TradeTimeRangeUpdate) ClearInterval() *TradeTimeRangeUpdate {
	ttru.mutation.ClearInterval()
	return ttru
}

// ClearRecords clears all "records" edges to the TradeRecord entity.
func (ttru *TradeTimeRangeUpdate) ClearRecords() *TradeTimeRangeUpdate {
	ttru.mutation.ClearRecords()
	return ttru
}

// RemoveRecordIDs removes the "records" edge to TradeRecord entities by IDs.
func (ttru *TradeTimeRangeUpdate) RemoveRecordIDs(ids ...int) *TradeTimeRangeUpdate {
	ttru.mutation.RemoveRecordIDs(ids...)
	return ttru
}

// RemoveRecords removes "records" edges to TradeRecord entities.
func (ttru *TradeTimeRangeUpdate) RemoveRecords(t ...*TradeRecord) *TradeTimeRangeUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return ttru.RemoveRecordIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (ttru *TradeTimeRangeUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, ttru.sqlSave, ttru.mutation, ttru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ttru *TradeTimeRangeUpdate) SaveX(ctx context.Context) int {
	affected, err := ttru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (ttru *TradeTimeRangeUpdate) Exec(ctx context.Context) error {
	_, err := ttru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ttru *TradeTimeRangeUpdate) ExecX(ctx context.Context) {
	if err := ttru.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ttru *TradeTimeRangeUpdate) check() error {
	if ttru.mutation.IntervalCleared() && len(ttru.mutation.IntervalIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "TradeTimeRange.interval"`)
	}
	return nil
}

func (ttru *TradeTimeRangeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := ttru.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(tradetimerange.Table, tradetimerange.Columns, sqlgraph.NewFieldSpec(tradetimerange.FieldID, field.TypeInt))
	if ps := ttru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ttru.mutation.Start(); ok {
		_spec.SetField(tradetimerange.FieldStart, field.TypeTime, value)
	}
	if value, ok := ttru.mutation.End(); ok {
		_spec.SetField(tradetimerange.FieldEnd, field.TypeTime, value)
	}
	if ttru.mutation.IntervalCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   tradetimerange.IntervalTable,
			Columns: []string{tradetimerange.IntervalColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(interval.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ttru.mutation.IntervalIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   tradetimerange.IntervalTable,
			Columns: []string{tradetimerange.IntervalColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(interval.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ttru.mutation.RecordsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tradetimerange.RecordsTable,
			Columns: []string{tradetimerange.RecordsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(traderecord.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ttru.mutation.RemovedRecordsIDs(); len(nodes) > 0 && !ttru.mutation.RecordsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tradetimerange.RecordsTable,
			Columns: []string{tradetimerange.RecordsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(traderecord.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ttru.mutation.RecordsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tradetimerange.RecordsTable,
			Columns: []string{tradetimerange.RecordsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(traderecord.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, ttru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{tradetimerange.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	ttru.mutation.done = true
	return n, nil
}

// TradeTimeRangeUpdateOne is the builder for updating a single TradeTimeRange entity.
type TradeTimeRangeUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *TradeTimeRangeMutation
}

// SetStart sets the "start" field.
func (ttruo *TradeTimeRangeUpdateOne) SetStart(t time.Time) *TradeTimeRangeUpdateOne {
	ttruo.mutation.SetStart(t)
	return ttruo
}

// SetNillableStart sets the "start" field if the given value is not nil.
func (ttruo *TradeTimeRangeUpdateOne) SetNillableStart(t *time.Time) *TradeTimeRangeUpdateOne {
	if t != nil {
		ttruo.SetStart(*t)
	}
	return ttruo
}

// SetEnd sets the "end" field.
func (ttruo *TradeTimeRangeUpdateOne) SetEnd(t time.Time) *TradeTimeRangeUpdateOne {
	ttruo.mutation.SetEnd(t)
	return ttruo
}

// SetNillableEnd sets the "end" field if the given value is not nil.
func (ttruo *TradeTimeRangeUpdateOne) SetNillableEnd(t *time.Time) *TradeTimeRangeUpdateOne {
	if t != nil {
		ttruo.SetEnd(*t)
	}
	return ttruo
}

// SetIntervalID sets the "interval_id" field.
func (ttruo *TradeTimeRangeUpdateOne) SetIntervalID(i int) *TradeTimeRangeUpdateOne {
	ttruo.mutation.SetIntervalID(i)
	return ttruo
}

// SetNillableIntervalID sets the "interval_id" field if the given value is not nil.
func (ttruo *TradeTimeRangeUpdateOne) SetNillableIntervalID(i *int) *TradeTimeRangeUpdateOne {
	if i != nil {
		ttruo.SetIntervalID(*i)
	}
	return ttruo
}

// SetInterval sets the "interval" edge to the Interval entity.
func (ttruo *TradeTimeRangeUpdateOne) SetInterval(i *Interval) *TradeTimeRangeUpdateOne {
	return ttruo.SetIntervalID(i.ID)
}

// AddRecordIDs adds the "records" edge to the TradeRecord entity by IDs.
func (ttruo *TradeTimeRangeUpdateOne) AddRecordIDs(ids ...int) *TradeTimeRangeUpdateOne {
	ttruo.mutation.AddRecordIDs(ids...)
	return ttruo
}

// AddRecords adds the "records" edges to the TradeRecord entity.
func (ttruo *TradeTimeRangeUpdateOne) AddRecords(t ...*TradeRecord) *TradeTimeRangeUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return ttruo.AddRecordIDs(ids...)
}

// Mutation returns the TradeTimeRangeMutation object of the builder.
func (ttruo *TradeTimeRangeUpdateOne) Mutation() *TradeTimeRangeMutation {
	return ttruo.mutation
}

// ClearInterval clears the "interval" edge to the Interval entity.
func (ttruo *TradeTimeRangeUpdateOne) ClearInterval() *TradeTimeRangeUpdateOne {
	ttruo.mutation.ClearInterval()
	return ttruo
}

// ClearRecords clears all "records" edges to the TradeRecord entity.
func (ttruo *TradeTimeRangeUpdateOne) ClearRecords() *TradeTimeRangeUpdateOne {
	ttruo.mutation.ClearRecords()
	return ttruo
}

// RemoveRecordIDs removes the "records" edge to TradeRecord entities by IDs.
func (ttruo *TradeTimeRangeUpdateOne) RemoveRecordIDs(ids ...int) *TradeTimeRangeUpdateOne {
	ttruo.mutation.RemoveRecordIDs(ids...)
	return ttruo
}

// RemoveRecords removes "records" edges to TradeRecord entities.
func (ttruo *TradeTimeRangeUpdateOne) RemoveRecords(t ...*TradeRecord) *TradeTimeRangeUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return ttruo.RemoveRecordIDs(ids...)
}

// Where appends a list predicates to the TradeTimeRangeUpdate builder.
func (ttruo *TradeTimeRangeUpdateOne) Where(ps ...predicate.TradeTimeRange) *TradeTimeRangeUpdateOne {
	ttruo.mutation.Where(ps...)
	return ttruo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (ttruo *TradeTimeRangeUpdateOne) Select(field string, fields ...string) *TradeTimeRangeUpdateOne {
	ttruo.fields = append([]string{field}, fields...)
	return ttruo
}

// Save executes the query and returns the updated TradeTimeRange entity.
func (ttruo *TradeTimeRangeUpdateOne) Save(ctx context.Context) (*TradeTimeRange, error) {
	return withHooks(ctx, ttruo.sqlSave, ttruo.mutation, ttruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (ttruo *TradeTimeRangeUpdateOne) SaveX(ctx context.Context) *TradeTimeRange {
	node, err := ttruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (ttruo *TradeTimeRangeUpdateOne) Exec(ctx context.Context) error {
	_, err := ttruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ttruo *TradeTimeRangeUpdateOne) ExecX(ctx context.Context) {
	if err := ttruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ttruo *TradeTimeRangeUpdateOne) check() error {
	if ttruo.mutation.IntervalCleared() && len(ttruo.mutation.IntervalIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "TradeTimeRange.interval"`)
	}
	return nil
}

func (ttruo *TradeTimeRangeUpdateOne) sqlSave(ctx context.Context) (_node *TradeTimeRange, err error) {
	if err := ttruo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(tradetimerange.Table, tradetimerange.Columns, sqlgraph.NewFieldSpec(tradetimerange.FieldID, field.TypeInt))
	id, ok := ttruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "TradeTimeRange.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := ttruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, tradetimerange.FieldID)
		for _, f := range fields {
			if !tradetimerange.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != tradetimerange.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := ttruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := ttruo.mutation.Start(); ok {
		_spec.SetField(tradetimerange.FieldStart, field.TypeTime, value)
	}
	if value, ok := ttruo.mutation.End(); ok {
		_spec.SetField(tradetimerange.FieldEnd, field.TypeTime, value)
	}
	if ttruo.mutation.IntervalCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   tradetimerange.IntervalTable,
			Columns: []string{tradetimerange.IntervalColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(interval.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ttruo.mutation.IntervalIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   tradetimerange.IntervalTable,
			Columns: []string{tradetimerange.IntervalColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(interval.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if ttruo.mutation.RecordsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tradetimerange.RecordsTable,
			Columns: []string{tradetimerange.RecordsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(traderecord.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ttruo.mutation.RemovedRecordsIDs(); len(nodes) > 0 && !ttruo.mutation.RecordsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tradetimerange.RecordsTable,
			Columns: []string{tradetimerange.RecordsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(traderecord.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := ttruo.mutation.RecordsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tradetimerange.RecordsTable,
			Columns: []string{tradetimerange.RecordsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(traderecord.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &TradeTimeRange{config: ttruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, ttruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{tradetimerange.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	ttruo.mutation.done = true
	return _node, nil
}
