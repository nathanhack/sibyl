// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/tradecorrection"
	"github.com/nathanhack/sibyl/ent/traderecord"
)

// TradeCorrectionCreate is the builder for creating a TradeCorrection entity.
type TradeCorrectionCreate struct {
	config
	mutation *TradeCorrectionMutation
	hooks    []Hook
}

// SetCorrection sets the "correction" field.
func (tcc *TradeCorrectionCreate) SetCorrection(s string) *TradeCorrectionCreate {
	tcc.mutation.SetCorrection(s)
	return tcc
}

// AddRecordIDs adds the "record" edge to the TradeRecord entity by IDs.
func (tcc *TradeCorrectionCreate) AddRecordIDs(ids ...int) *TradeCorrectionCreate {
	tcc.mutation.AddRecordIDs(ids...)
	return tcc
}

// AddRecord adds the "record" edges to the TradeRecord entity.
func (tcc *TradeCorrectionCreate) AddRecord(t ...*TradeRecord) *TradeCorrectionCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tcc.AddRecordIDs(ids...)
}

// Mutation returns the TradeCorrectionMutation object of the builder.
func (tcc *TradeCorrectionCreate) Mutation() *TradeCorrectionMutation {
	return tcc.mutation
}

// Save creates the TradeCorrection in the database.
func (tcc *TradeCorrectionCreate) Save(ctx context.Context) (*TradeCorrection, error) {
	return withHooks[*TradeCorrection, TradeCorrectionMutation](ctx, tcc.sqlSave, tcc.mutation, tcc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tcc *TradeCorrectionCreate) SaveX(ctx context.Context) *TradeCorrection {
	v, err := tcc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcc *TradeCorrectionCreate) Exec(ctx context.Context) error {
	_, err := tcc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcc *TradeCorrectionCreate) ExecX(ctx context.Context) {
	if err := tcc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tcc *TradeCorrectionCreate) check() error {
	if _, ok := tcc.mutation.Correction(); !ok {
		return &ValidationError{Name: "correction", err: errors.New(`ent: missing required field "TradeCorrection.correction"`)}
	}
	return nil
}

func (tcc *TradeCorrectionCreate) sqlSave(ctx context.Context) (*TradeCorrection, error) {
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

func (tcc *TradeCorrectionCreate) createSpec() (*TradeCorrection, *sqlgraph.CreateSpec) {
	var (
		_node = &TradeCorrection{config: tcc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: tradecorrection.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: tradecorrection.FieldID,
			},
		}
	)
	if value, ok := tcc.mutation.Correction(); ok {
		_spec.SetField(tradecorrection.FieldCorrection, field.TypeString, value)
		_node.Correction = value
	}
	if nodes := tcc.mutation.RecordIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   tradecorrection.RecordTable,
			Columns: tradecorrection.RecordPrimaryKey,
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

// TradeCorrectionCreateBulk is the builder for creating many TradeCorrection entities in bulk.
type TradeCorrectionCreateBulk struct {
	config
	builders []*TradeCorrectionCreate
}

// Save creates the TradeCorrection entities in the database.
func (tccb *TradeCorrectionCreateBulk) Save(ctx context.Context) ([]*TradeCorrection, error) {
	specs := make([]*sqlgraph.CreateSpec, len(tccb.builders))
	nodes := make([]*TradeCorrection, len(tccb.builders))
	mutators := make([]Mutator, len(tccb.builders))
	for i := range tccb.builders {
		func(i int, root context.Context) {
			builder := tccb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TradeCorrectionMutation)
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
func (tccb *TradeCorrectionCreateBulk) SaveX(ctx context.Context) []*TradeCorrection {
	v, err := tccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tccb *TradeCorrectionCreateBulk) Exec(ctx context.Context) error {
	_, err := tccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tccb *TradeCorrectionCreateBulk) ExecX(ctx context.Context) {
	if err := tccb.Exec(ctx); err != nil {
		panic(err)
	}
}