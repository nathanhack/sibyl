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
	"github.com/nathanhack/sibyl/ent/exchange"
	"github.com/nathanhack/sibyl/ent/predicate"
	"github.com/nathanhack/sibyl/ent/tradecondition"
	"github.com/nathanhack/sibyl/ent/tradecorrection"
	"github.com/nathanhack/sibyl/ent/traderecord"
	"github.com/nathanhack/sibyl/ent/tradetimerange"
)

// TradeRecordUpdate is the builder for updating TradeRecord entities.
type TradeRecordUpdate struct {
	config
	hooks    []Hook
	mutation *TradeRecordMutation
}

// Where appends a list predicates to the TradeRecordUpdate builder.
func (tru *TradeRecordUpdate) Where(ps ...predicate.TradeRecord) *TradeRecordUpdate {
	tru.mutation.Where(ps...)
	return tru
}

// SetPrice sets the "price" field.
func (tru *TradeRecordUpdate) SetPrice(f float64) *TradeRecordUpdate {
	tru.mutation.ResetPrice()
	tru.mutation.SetPrice(f)
	return tru
}

// SetNillablePrice sets the "price" field if the given value is not nil.
func (tru *TradeRecordUpdate) SetNillablePrice(f *float64) *TradeRecordUpdate {
	if f != nil {
		tru.SetPrice(*f)
	}
	return tru
}

// AddPrice adds f to the "price" field.
func (tru *TradeRecordUpdate) AddPrice(f float64) *TradeRecordUpdate {
	tru.mutation.AddPrice(f)
	return tru
}

// SetTimestamp sets the "timestamp" field.
func (tru *TradeRecordUpdate) SetTimestamp(t time.Time) *TradeRecordUpdate {
	tru.mutation.SetTimestamp(t)
	return tru
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (tru *TradeRecordUpdate) SetNillableTimestamp(t *time.Time) *TradeRecordUpdate {
	if t != nil {
		tru.SetTimestamp(*t)
	}
	return tru
}

// SetVolume sets the "volume" field.
func (tru *TradeRecordUpdate) SetVolume(i int32) *TradeRecordUpdate {
	tru.mutation.ResetVolume()
	tru.mutation.SetVolume(i)
	return tru
}

// SetNillableVolume sets the "volume" field if the given value is not nil.
func (tru *TradeRecordUpdate) SetNillableVolume(i *int32) *TradeRecordUpdate {
	if i != nil {
		tru.SetVolume(*i)
	}
	return tru
}

// AddVolume adds i to the "volume" field.
func (tru *TradeRecordUpdate) AddVolume(i int32) *TradeRecordUpdate {
	tru.mutation.AddVolume(i)
	return tru
}

// SetTimeRangeID sets the "time_range_id" field.
func (tru *TradeRecordUpdate) SetTimeRangeID(i int) *TradeRecordUpdate {
	tru.mutation.SetTimeRangeID(i)
	return tru
}

// SetNillableTimeRangeID sets the "time_range_id" field if the given value is not nil.
func (tru *TradeRecordUpdate) SetNillableTimeRangeID(i *int) *TradeRecordUpdate {
	if i != nil {
		tru.SetTimeRangeID(*i)
	}
	return tru
}

// SetTimeRange sets the "time_range" edge to the TradeTimeRange entity.
func (tru *TradeRecordUpdate) SetTimeRange(t *TradeTimeRange) *TradeRecordUpdate {
	return tru.SetTimeRangeID(t.ID)
}

// AddConditionIDs adds the "conditions" edge to the TradeCondition entity by IDs.
func (tru *TradeRecordUpdate) AddConditionIDs(ids ...int) *TradeRecordUpdate {
	tru.mutation.AddConditionIDs(ids...)
	return tru
}

// AddConditions adds the "conditions" edges to the TradeCondition entity.
func (tru *TradeRecordUpdate) AddConditions(t ...*TradeCondition) *TradeRecordUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tru.AddConditionIDs(ids...)
}

// AddCorrectionIDs adds the "correction" edge to the TradeCorrection entity by IDs.
func (tru *TradeRecordUpdate) AddCorrectionIDs(ids ...int) *TradeRecordUpdate {
	tru.mutation.AddCorrectionIDs(ids...)
	return tru
}

// AddCorrection adds the "correction" edges to the TradeCorrection entity.
func (tru *TradeRecordUpdate) AddCorrection(t ...*TradeCorrection) *TradeRecordUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tru.AddCorrectionIDs(ids...)
}

// AddExchangeIDs adds the "exchange" edge to the Exchange entity by IDs.
func (tru *TradeRecordUpdate) AddExchangeIDs(ids ...int) *TradeRecordUpdate {
	tru.mutation.AddExchangeIDs(ids...)
	return tru
}

// AddExchange adds the "exchange" edges to the Exchange entity.
func (tru *TradeRecordUpdate) AddExchange(e ...*Exchange) *TradeRecordUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return tru.AddExchangeIDs(ids...)
}

// Mutation returns the TradeRecordMutation object of the builder.
func (tru *TradeRecordUpdate) Mutation() *TradeRecordMutation {
	return tru.mutation
}

// ClearTimeRange clears the "time_range" edge to the TradeTimeRange entity.
func (tru *TradeRecordUpdate) ClearTimeRange() *TradeRecordUpdate {
	tru.mutation.ClearTimeRange()
	return tru
}

// ClearConditions clears all "conditions" edges to the TradeCondition entity.
func (tru *TradeRecordUpdate) ClearConditions() *TradeRecordUpdate {
	tru.mutation.ClearConditions()
	return tru
}

// RemoveConditionIDs removes the "conditions" edge to TradeCondition entities by IDs.
func (tru *TradeRecordUpdate) RemoveConditionIDs(ids ...int) *TradeRecordUpdate {
	tru.mutation.RemoveConditionIDs(ids...)
	return tru
}

// RemoveConditions removes "conditions" edges to TradeCondition entities.
func (tru *TradeRecordUpdate) RemoveConditions(t ...*TradeCondition) *TradeRecordUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tru.RemoveConditionIDs(ids...)
}

// ClearCorrection clears all "correction" edges to the TradeCorrection entity.
func (tru *TradeRecordUpdate) ClearCorrection() *TradeRecordUpdate {
	tru.mutation.ClearCorrection()
	return tru
}

// RemoveCorrectionIDs removes the "correction" edge to TradeCorrection entities by IDs.
func (tru *TradeRecordUpdate) RemoveCorrectionIDs(ids ...int) *TradeRecordUpdate {
	tru.mutation.RemoveCorrectionIDs(ids...)
	return tru
}

// RemoveCorrection removes "correction" edges to TradeCorrection entities.
func (tru *TradeRecordUpdate) RemoveCorrection(t ...*TradeCorrection) *TradeRecordUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tru.RemoveCorrectionIDs(ids...)
}

// ClearExchange clears all "exchange" edges to the Exchange entity.
func (tru *TradeRecordUpdate) ClearExchange() *TradeRecordUpdate {
	tru.mutation.ClearExchange()
	return tru
}

// RemoveExchangeIDs removes the "exchange" edge to Exchange entities by IDs.
func (tru *TradeRecordUpdate) RemoveExchangeIDs(ids ...int) *TradeRecordUpdate {
	tru.mutation.RemoveExchangeIDs(ids...)
	return tru
}

// RemoveExchange removes "exchange" edges to Exchange entities.
func (tru *TradeRecordUpdate) RemoveExchange(e ...*Exchange) *TradeRecordUpdate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return tru.RemoveExchangeIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tru *TradeRecordUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, tru.sqlSave, tru.mutation, tru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tru *TradeRecordUpdate) SaveX(ctx context.Context) int {
	affected, err := tru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tru *TradeRecordUpdate) Exec(ctx context.Context) error {
	_, err := tru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tru *TradeRecordUpdate) ExecX(ctx context.Context) {
	if err := tru.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tru *TradeRecordUpdate) check() error {
	if tru.mutation.TimeRangeCleared() && len(tru.mutation.TimeRangeIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "TradeRecord.time_range"`)
	}
	return nil
}

func (tru *TradeRecordUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := tru.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(traderecord.Table, traderecord.Columns, sqlgraph.NewFieldSpec(traderecord.FieldID, field.TypeInt))
	if ps := tru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tru.mutation.Price(); ok {
		_spec.SetField(traderecord.FieldPrice, field.TypeFloat64, value)
	}
	if value, ok := tru.mutation.AddedPrice(); ok {
		_spec.AddField(traderecord.FieldPrice, field.TypeFloat64, value)
	}
	if value, ok := tru.mutation.Timestamp(); ok {
		_spec.SetField(traderecord.FieldTimestamp, field.TypeTime, value)
	}
	if value, ok := tru.mutation.Volume(); ok {
		_spec.SetField(traderecord.FieldVolume, field.TypeInt32, value)
	}
	if value, ok := tru.mutation.AddedVolume(); ok {
		_spec.AddField(traderecord.FieldVolume, field.TypeInt32, value)
	}
	if tru.mutation.TimeRangeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   traderecord.TimeRangeTable,
			Columns: []string{traderecord.TimeRangeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradetimerange.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tru.mutation.TimeRangeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   traderecord.TimeRangeTable,
			Columns: []string{traderecord.TimeRangeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradetimerange.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if tru.mutation.ConditionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   traderecord.ConditionsTable,
			Columns: traderecord.ConditionsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradecondition.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tru.mutation.RemovedConditionsIDs(); len(nodes) > 0 && !tru.mutation.ConditionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   traderecord.ConditionsTable,
			Columns: traderecord.ConditionsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradecondition.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tru.mutation.ConditionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   traderecord.ConditionsTable,
			Columns: traderecord.ConditionsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradecondition.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if tru.mutation.CorrectionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   traderecord.CorrectionTable,
			Columns: traderecord.CorrectionPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradecorrection.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tru.mutation.RemovedCorrectionIDs(); len(nodes) > 0 && !tru.mutation.CorrectionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   traderecord.CorrectionTable,
			Columns: traderecord.CorrectionPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradecorrection.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tru.mutation.CorrectionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   traderecord.CorrectionTable,
			Columns: traderecord.CorrectionPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradecorrection.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if tru.mutation.ExchangeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   traderecord.ExchangeTable,
			Columns: []string{traderecord.ExchangeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exchange.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tru.mutation.RemovedExchangeIDs(); len(nodes) > 0 && !tru.mutation.ExchangeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   traderecord.ExchangeTable,
			Columns: []string{traderecord.ExchangeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exchange.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tru.mutation.ExchangeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   traderecord.ExchangeTable,
			Columns: []string{traderecord.ExchangeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exchange.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{traderecord.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	tru.mutation.done = true
	return n, nil
}

// TradeRecordUpdateOne is the builder for updating a single TradeRecord entity.
type TradeRecordUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *TradeRecordMutation
}

// SetPrice sets the "price" field.
func (truo *TradeRecordUpdateOne) SetPrice(f float64) *TradeRecordUpdateOne {
	truo.mutation.ResetPrice()
	truo.mutation.SetPrice(f)
	return truo
}

// SetNillablePrice sets the "price" field if the given value is not nil.
func (truo *TradeRecordUpdateOne) SetNillablePrice(f *float64) *TradeRecordUpdateOne {
	if f != nil {
		truo.SetPrice(*f)
	}
	return truo
}

// AddPrice adds f to the "price" field.
func (truo *TradeRecordUpdateOne) AddPrice(f float64) *TradeRecordUpdateOne {
	truo.mutation.AddPrice(f)
	return truo
}

// SetTimestamp sets the "timestamp" field.
func (truo *TradeRecordUpdateOne) SetTimestamp(t time.Time) *TradeRecordUpdateOne {
	truo.mutation.SetTimestamp(t)
	return truo
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (truo *TradeRecordUpdateOne) SetNillableTimestamp(t *time.Time) *TradeRecordUpdateOne {
	if t != nil {
		truo.SetTimestamp(*t)
	}
	return truo
}

// SetVolume sets the "volume" field.
func (truo *TradeRecordUpdateOne) SetVolume(i int32) *TradeRecordUpdateOne {
	truo.mutation.ResetVolume()
	truo.mutation.SetVolume(i)
	return truo
}

// SetNillableVolume sets the "volume" field if the given value is not nil.
func (truo *TradeRecordUpdateOne) SetNillableVolume(i *int32) *TradeRecordUpdateOne {
	if i != nil {
		truo.SetVolume(*i)
	}
	return truo
}

// AddVolume adds i to the "volume" field.
func (truo *TradeRecordUpdateOne) AddVolume(i int32) *TradeRecordUpdateOne {
	truo.mutation.AddVolume(i)
	return truo
}

// SetTimeRangeID sets the "time_range_id" field.
func (truo *TradeRecordUpdateOne) SetTimeRangeID(i int) *TradeRecordUpdateOne {
	truo.mutation.SetTimeRangeID(i)
	return truo
}

// SetNillableTimeRangeID sets the "time_range_id" field if the given value is not nil.
func (truo *TradeRecordUpdateOne) SetNillableTimeRangeID(i *int) *TradeRecordUpdateOne {
	if i != nil {
		truo.SetTimeRangeID(*i)
	}
	return truo
}

// SetTimeRange sets the "time_range" edge to the TradeTimeRange entity.
func (truo *TradeRecordUpdateOne) SetTimeRange(t *TradeTimeRange) *TradeRecordUpdateOne {
	return truo.SetTimeRangeID(t.ID)
}

// AddConditionIDs adds the "conditions" edge to the TradeCondition entity by IDs.
func (truo *TradeRecordUpdateOne) AddConditionIDs(ids ...int) *TradeRecordUpdateOne {
	truo.mutation.AddConditionIDs(ids...)
	return truo
}

// AddConditions adds the "conditions" edges to the TradeCondition entity.
func (truo *TradeRecordUpdateOne) AddConditions(t ...*TradeCondition) *TradeRecordUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return truo.AddConditionIDs(ids...)
}

// AddCorrectionIDs adds the "correction" edge to the TradeCorrection entity by IDs.
func (truo *TradeRecordUpdateOne) AddCorrectionIDs(ids ...int) *TradeRecordUpdateOne {
	truo.mutation.AddCorrectionIDs(ids...)
	return truo
}

// AddCorrection adds the "correction" edges to the TradeCorrection entity.
func (truo *TradeRecordUpdateOne) AddCorrection(t ...*TradeCorrection) *TradeRecordUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return truo.AddCorrectionIDs(ids...)
}

// AddExchangeIDs adds the "exchange" edge to the Exchange entity by IDs.
func (truo *TradeRecordUpdateOne) AddExchangeIDs(ids ...int) *TradeRecordUpdateOne {
	truo.mutation.AddExchangeIDs(ids...)
	return truo
}

// AddExchange adds the "exchange" edges to the Exchange entity.
func (truo *TradeRecordUpdateOne) AddExchange(e ...*Exchange) *TradeRecordUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return truo.AddExchangeIDs(ids...)
}

// Mutation returns the TradeRecordMutation object of the builder.
func (truo *TradeRecordUpdateOne) Mutation() *TradeRecordMutation {
	return truo.mutation
}

// ClearTimeRange clears the "time_range" edge to the TradeTimeRange entity.
func (truo *TradeRecordUpdateOne) ClearTimeRange() *TradeRecordUpdateOne {
	truo.mutation.ClearTimeRange()
	return truo
}

// ClearConditions clears all "conditions" edges to the TradeCondition entity.
func (truo *TradeRecordUpdateOne) ClearConditions() *TradeRecordUpdateOne {
	truo.mutation.ClearConditions()
	return truo
}

// RemoveConditionIDs removes the "conditions" edge to TradeCondition entities by IDs.
func (truo *TradeRecordUpdateOne) RemoveConditionIDs(ids ...int) *TradeRecordUpdateOne {
	truo.mutation.RemoveConditionIDs(ids...)
	return truo
}

// RemoveConditions removes "conditions" edges to TradeCondition entities.
func (truo *TradeRecordUpdateOne) RemoveConditions(t ...*TradeCondition) *TradeRecordUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return truo.RemoveConditionIDs(ids...)
}

// ClearCorrection clears all "correction" edges to the TradeCorrection entity.
func (truo *TradeRecordUpdateOne) ClearCorrection() *TradeRecordUpdateOne {
	truo.mutation.ClearCorrection()
	return truo
}

// RemoveCorrectionIDs removes the "correction" edge to TradeCorrection entities by IDs.
func (truo *TradeRecordUpdateOne) RemoveCorrectionIDs(ids ...int) *TradeRecordUpdateOne {
	truo.mutation.RemoveCorrectionIDs(ids...)
	return truo
}

// RemoveCorrection removes "correction" edges to TradeCorrection entities.
func (truo *TradeRecordUpdateOne) RemoveCorrection(t ...*TradeCorrection) *TradeRecordUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return truo.RemoveCorrectionIDs(ids...)
}

// ClearExchange clears all "exchange" edges to the Exchange entity.
func (truo *TradeRecordUpdateOne) ClearExchange() *TradeRecordUpdateOne {
	truo.mutation.ClearExchange()
	return truo
}

// RemoveExchangeIDs removes the "exchange" edge to Exchange entities by IDs.
func (truo *TradeRecordUpdateOne) RemoveExchangeIDs(ids ...int) *TradeRecordUpdateOne {
	truo.mutation.RemoveExchangeIDs(ids...)
	return truo
}

// RemoveExchange removes "exchange" edges to Exchange entities.
func (truo *TradeRecordUpdateOne) RemoveExchange(e ...*Exchange) *TradeRecordUpdateOne {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return truo.RemoveExchangeIDs(ids...)
}

// Where appends a list predicates to the TradeRecordUpdate builder.
func (truo *TradeRecordUpdateOne) Where(ps ...predicate.TradeRecord) *TradeRecordUpdateOne {
	truo.mutation.Where(ps...)
	return truo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (truo *TradeRecordUpdateOne) Select(field string, fields ...string) *TradeRecordUpdateOne {
	truo.fields = append([]string{field}, fields...)
	return truo
}

// Save executes the query and returns the updated TradeRecord entity.
func (truo *TradeRecordUpdateOne) Save(ctx context.Context) (*TradeRecord, error) {
	return withHooks(ctx, truo.sqlSave, truo.mutation, truo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (truo *TradeRecordUpdateOne) SaveX(ctx context.Context) *TradeRecord {
	node, err := truo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (truo *TradeRecordUpdateOne) Exec(ctx context.Context) error {
	_, err := truo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (truo *TradeRecordUpdateOne) ExecX(ctx context.Context) {
	if err := truo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (truo *TradeRecordUpdateOne) check() error {
	if truo.mutation.TimeRangeCleared() && len(truo.mutation.TimeRangeIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "TradeRecord.time_range"`)
	}
	return nil
}

func (truo *TradeRecordUpdateOne) sqlSave(ctx context.Context) (_node *TradeRecord, err error) {
	if err := truo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(traderecord.Table, traderecord.Columns, sqlgraph.NewFieldSpec(traderecord.FieldID, field.TypeInt))
	id, ok := truo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "TradeRecord.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := truo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, traderecord.FieldID)
		for _, f := range fields {
			if !traderecord.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != traderecord.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := truo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := truo.mutation.Price(); ok {
		_spec.SetField(traderecord.FieldPrice, field.TypeFloat64, value)
	}
	if value, ok := truo.mutation.AddedPrice(); ok {
		_spec.AddField(traderecord.FieldPrice, field.TypeFloat64, value)
	}
	if value, ok := truo.mutation.Timestamp(); ok {
		_spec.SetField(traderecord.FieldTimestamp, field.TypeTime, value)
	}
	if value, ok := truo.mutation.Volume(); ok {
		_spec.SetField(traderecord.FieldVolume, field.TypeInt32, value)
	}
	if value, ok := truo.mutation.AddedVolume(); ok {
		_spec.AddField(traderecord.FieldVolume, field.TypeInt32, value)
	}
	if truo.mutation.TimeRangeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   traderecord.TimeRangeTable,
			Columns: []string{traderecord.TimeRangeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradetimerange.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := truo.mutation.TimeRangeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   traderecord.TimeRangeTable,
			Columns: []string{traderecord.TimeRangeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradetimerange.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if truo.mutation.ConditionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   traderecord.ConditionsTable,
			Columns: traderecord.ConditionsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradecondition.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := truo.mutation.RemovedConditionsIDs(); len(nodes) > 0 && !truo.mutation.ConditionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   traderecord.ConditionsTable,
			Columns: traderecord.ConditionsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradecondition.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := truo.mutation.ConditionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   traderecord.ConditionsTable,
			Columns: traderecord.ConditionsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradecondition.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if truo.mutation.CorrectionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   traderecord.CorrectionTable,
			Columns: traderecord.CorrectionPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradecorrection.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := truo.mutation.RemovedCorrectionIDs(); len(nodes) > 0 && !truo.mutation.CorrectionCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   traderecord.CorrectionTable,
			Columns: traderecord.CorrectionPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradecorrection.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := truo.mutation.CorrectionIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   traderecord.CorrectionTable,
			Columns: traderecord.CorrectionPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tradecorrection.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if truo.mutation.ExchangeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   traderecord.ExchangeTable,
			Columns: []string{traderecord.ExchangeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exchange.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := truo.mutation.RemovedExchangeIDs(); len(nodes) > 0 && !truo.mutation.ExchangeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   traderecord.ExchangeTable,
			Columns: []string{traderecord.ExchangeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exchange.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := truo.mutation.ExchangeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   traderecord.ExchangeTable,
			Columns: []string{traderecord.ExchangeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(exchange.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &TradeRecord{config: truo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, truo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{traderecord.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	truo.mutation.done = true
	return _node, nil
}
