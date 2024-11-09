// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/exchange"
	"github.com/nathanhack/sibyl/ent/tradecondition"
	"github.com/nathanhack/sibyl/ent/tradecorrection"
	"github.com/nathanhack/sibyl/ent/traderecord"
	"github.com/nathanhack/sibyl/ent/tradetimerange"
)

// TradeRecordCreate is the builder for creating a TradeRecord entity.
type TradeRecordCreate struct {
	config
	mutation *TradeRecordMutation
	hooks    []Hook
}

// SetPrice sets the "price" field.
func (trc *TradeRecordCreate) SetPrice(f float64) *TradeRecordCreate {
	trc.mutation.SetPrice(f)
	return trc
}

// SetTimestamp sets the "timestamp" field.
func (trc *TradeRecordCreate) SetTimestamp(t time.Time) *TradeRecordCreate {
	trc.mutation.SetTimestamp(t)
	return trc
}

// SetVolume sets the "volume" field.
func (trc *TradeRecordCreate) SetVolume(i int32) *TradeRecordCreate {
	trc.mutation.SetVolume(i)
	return trc
}

// SetTimeRangeID sets the "time_range_id" field.
func (trc *TradeRecordCreate) SetTimeRangeID(i int) *TradeRecordCreate {
	trc.mutation.SetTimeRangeID(i)
	return trc
}

// SetTimeRange sets the "time_range" edge to the TradeTimeRange entity.
func (trc *TradeRecordCreate) SetTimeRange(t *TradeTimeRange) *TradeRecordCreate {
	return trc.SetTimeRangeID(t.ID)
}

// AddConditionIDs adds the "conditions" edge to the TradeCondition entity by IDs.
func (trc *TradeRecordCreate) AddConditionIDs(ids ...int) *TradeRecordCreate {
	trc.mutation.AddConditionIDs(ids...)
	return trc
}

// AddConditions adds the "conditions" edges to the TradeCondition entity.
func (trc *TradeRecordCreate) AddConditions(t ...*TradeCondition) *TradeRecordCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return trc.AddConditionIDs(ids...)
}

// AddCorrectionIDs adds the "correction" edge to the TradeCorrection entity by IDs.
func (trc *TradeRecordCreate) AddCorrectionIDs(ids ...int) *TradeRecordCreate {
	trc.mutation.AddCorrectionIDs(ids...)
	return trc
}

// AddCorrection adds the "correction" edges to the TradeCorrection entity.
func (trc *TradeRecordCreate) AddCorrection(t ...*TradeCorrection) *TradeRecordCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return trc.AddCorrectionIDs(ids...)
}

// AddExchangeIDs adds the "exchange" edge to the Exchange entity by IDs.
func (trc *TradeRecordCreate) AddExchangeIDs(ids ...int) *TradeRecordCreate {
	trc.mutation.AddExchangeIDs(ids...)
	return trc
}

// AddExchange adds the "exchange" edges to the Exchange entity.
func (trc *TradeRecordCreate) AddExchange(e ...*Exchange) *TradeRecordCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return trc.AddExchangeIDs(ids...)
}

// Mutation returns the TradeRecordMutation object of the builder.
func (trc *TradeRecordCreate) Mutation() *TradeRecordMutation {
	return trc.mutation
}

// Save creates the TradeRecord in the database.
func (trc *TradeRecordCreate) Save(ctx context.Context) (*TradeRecord, error) {
	return withHooks(ctx, trc.sqlSave, trc.mutation, trc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (trc *TradeRecordCreate) SaveX(ctx context.Context) *TradeRecord {
	v, err := trc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (trc *TradeRecordCreate) Exec(ctx context.Context) error {
	_, err := trc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (trc *TradeRecordCreate) ExecX(ctx context.Context) {
	if err := trc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (trc *TradeRecordCreate) check() error {
	if _, ok := trc.mutation.Price(); !ok {
		return &ValidationError{Name: "price", err: errors.New(`ent: missing required field "TradeRecord.price"`)}
	}
	if _, ok := trc.mutation.Timestamp(); !ok {
		return &ValidationError{Name: "timestamp", err: errors.New(`ent: missing required field "TradeRecord.timestamp"`)}
	}
	if _, ok := trc.mutation.Volume(); !ok {
		return &ValidationError{Name: "volume", err: errors.New(`ent: missing required field "TradeRecord.volume"`)}
	}
	if _, ok := trc.mutation.TimeRangeID(); !ok {
		return &ValidationError{Name: "time_range_id", err: errors.New(`ent: missing required field "TradeRecord.time_range_id"`)}
	}
	if len(trc.mutation.TimeRangeIDs()) == 0 {
		return &ValidationError{Name: "time_range", err: errors.New(`ent: missing required edge "TradeRecord.time_range"`)}
	}
	return nil
}

func (trc *TradeRecordCreate) sqlSave(ctx context.Context) (*TradeRecord, error) {
	if err := trc.check(); err != nil {
		return nil, err
	}
	_node, _spec := trc.createSpec()
	if err := sqlgraph.CreateNode(ctx, trc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	trc.mutation.id = &_node.ID
	trc.mutation.done = true
	return _node, nil
}

func (trc *TradeRecordCreate) createSpec() (*TradeRecord, *sqlgraph.CreateSpec) {
	var (
		_node = &TradeRecord{config: trc.config}
		_spec = sqlgraph.NewCreateSpec(traderecord.Table, sqlgraph.NewFieldSpec(traderecord.FieldID, field.TypeInt))
	)
	if value, ok := trc.mutation.Price(); ok {
		_spec.SetField(traderecord.FieldPrice, field.TypeFloat64, value)
		_node.Price = value
	}
	if value, ok := trc.mutation.Timestamp(); ok {
		_spec.SetField(traderecord.FieldTimestamp, field.TypeTime, value)
		_node.Timestamp = value
	}
	if value, ok := trc.mutation.Volume(); ok {
		_spec.SetField(traderecord.FieldVolume, field.TypeInt32, value)
		_node.Volume = value
	}
	if nodes := trc.mutation.TimeRangeIDs(); len(nodes) > 0 {
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
		_node.TimeRangeID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := trc.mutation.ConditionsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := trc.mutation.CorrectionIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := trc.mutation.ExchangeIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// TradeRecordCreateBulk is the builder for creating many TradeRecord entities in bulk.
type TradeRecordCreateBulk struct {
	config
	err      error
	builders []*TradeRecordCreate
}

// Save creates the TradeRecord entities in the database.
func (trcb *TradeRecordCreateBulk) Save(ctx context.Context) ([]*TradeRecord, error) {
	if trcb.err != nil {
		return nil, trcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(trcb.builders))
	nodes := make([]*TradeRecord, len(trcb.builders))
	mutators := make([]Mutator, len(trcb.builders))
	for i := range trcb.builders {
		func(i int, root context.Context) {
			builder := trcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TradeRecordMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, trcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, trcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, trcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (trcb *TradeRecordCreateBulk) SaveX(ctx context.Context) []*TradeRecord {
	v, err := trcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (trcb *TradeRecordCreateBulk) Exec(ctx context.Context) error {
	_, err := trcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (trcb *TradeRecordCreateBulk) ExecX(ctx context.Context) {
	if err := trcb.Exec(ctx); err != nil {
		panic(err)
	}
}
