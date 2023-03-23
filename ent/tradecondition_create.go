// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/tradecondition"
	"github.com/nathanhack/sibyl/ent/traderecord"
)

// TradeConditionCreate is the builder for creating a TradeCondition entity.
type TradeConditionCreate struct {
	config
	mutation *TradeConditionMutation
	hooks    []Hook
}

// SetCondition sets the "condition" field.
func (tcc *TradeConditionCreate) SetCondition(s string) *TradeConditionCreate {
	tcc.mutation.SetCondition(s)
	return tcc
}

// AddRecordIDs adds the "record" edge to the TradeRecord entity by IDs.
func (tcc *TradeConditionCreate) AddRecordIDs(ids ...int) *TradeConditionCreate {
	tcc.mutation.AddRecordIDs(ids...)
	return tcc
}

// AddRecord adds the "record" edges to the TradeRecord entity.
func (tcc *TradeConditionCreate) AddRecord(t ...*TradeRecord) *TradeConditionCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tcc.AddRecordIDs(ids...)
}

// Mutation returns the TradeConditionMutation object of the builder.
func (tcc *TradeConditionCreate) Mutation() *TradeConditionMutation {
	return tcc.mutation
}

// Save creates the TradeCondition in the database.
func (tcc *TradeConditionCreate) Save(ctx context.Context) (*TradeCondition, error) {
	return withHooks[*TradeCondition, TradeConditionMutation](ctx, tcc.sqlSave, tcc.mutation, tcc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tcc *TradeConditionCreate) SaveX(ctx context.Context) *TradeCondition {
	v, err := tcc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcc *TradeConditionCreate) Exec(ctx context.Context) error {
	_, err := tcc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcc *TradeConditionCreate) ExecX(ctx context.Context) {
	if err := tcc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tcc *TradeConditionCreate) check() error {
	if _, ok := tcc.mutation.Condition(); !ok {
		return &ValidationError{Name: "condition", err: errors.New(`ent: missing required field "TradeCondition.condition"`)}
	}
	return nil
}

func (tcc *TradeConditionCreate) sqlSave(ctx context.Context) (*TradeCondition, error) {
	if err := tcc.check(); err != nil {
		return nil, err
	}
	_node, _spec := tcc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tcc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	tcc.mutation.id = &_node.ID
	tcc.mutation.done = true
	return _node, nil
}

func (tcc *TradeConditionCreate) createSpec() (*TradeCondition, *sqlgraph.CreateSpec) {
	var (
		_node = &TradeCondition{config: tcc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: tradecondition.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: tradecondition.FieldID,
			},
		}
	)
	if value, ok := tcc.mutation.Condition(); ok {
		_spec.SetField(tradecondition.FieldCondition, field.TypeString, value)
		_node.Condition = value
	}
	if nodes := tcc.mutation.RecordIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   tradecondition.RecordTable,
			Columns: tradecondition.RecordPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: traderecord.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// TradeConditionCreateBulk is the builder for creating many TradeCondition entities in bulk.
type TradeConditionCreateBulk struct {
	config
	builders []*TradeConditionCreate
}

// Save creates the TradeCondition entities in the database.
func (tccb *TradeConditionCreateBulk) Save(ctx context.Context) ([]*TradeCondition, error) {
	specs := make([]*sqlgraph.CreateSpec, len(tccb.builders))
	nodes := make([]*TradeCondition, len(tccb.builders))
	mutators := make([]Mutator, len(tccb.builders))
	for i := range tccb.builders {
		func(i int, root context.Context) {
			builder := tccb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TradeConditionMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, tccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tccb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, tccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tccb *TradeConditionCreateBulk) SaveX(ctx context.Context) []*TradeCondition {
	v, err := tccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tccb *TradeConditionCreateBulk) Exec(ctx context.Context) error {
	_, err := tccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tccb *TradeConditionCreateBulk) ExecX(ctx context.Context) {
	if err := tccb.Exec(ctx); err != nil {
		panic(err)
	}
}
