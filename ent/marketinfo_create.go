// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/markethours"
	"github.com/nathanhack/sibyl/ent/marketinfo"
)

// MarketInfoCreate is the builder for creating a MarketInfo entity.
type MarketInfoCreate struct {
	config
	mutation *MarketInfoMutation
	hooks    []Hook
}

// SetHoursStart sets the "hours_start" field.
func (mic *MarketInfoCreate) SetHoursStart(t time.Time) *MarketInfoCreate {
	mic.mutation.SetHoursStart(t)
	return mic
}

// SetHoursEnd sets the "hours_end" field.
func (mic *MarketInfoCreate) SetHoursEnd(t time.Time) *MarketInfoCreate {
	mic.mutation.SetHoursEnd(t)
	return mic
}

// AddHourIDs adds the "hours" edge to the MarketHours entity by IDs.
func (mic *MarketInfoCreate) AddHourIDs(ids ...int) *MarketInfoCreate {
	mic.mutation.AddHourIDs(ids...)
	return mic
}

// AddHours adds the "hours" edges to the MarketHours entity.
func (mic *MarketInfoCreate) AddHours(m ...*MarketHours) *MarketInfoCreate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return mic.AddHourIDs(ids...)
}

// Mutation returns the MarketInfoMutation object of the builder.
func (mic *MarketInfoCreate) Mutation() *MarketInfoMutation {
	return mic.mutation
}

// Save creates the MarketInfo in the database.
func (mic *MarketInfoCreate) Save(ctx context.Context) (*MarketInfo, error) {
	return withHooks(ctx, mic.sqlSave, mic.mutation, mic.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (mic *MarketInfoCreate) SaveX(ctx context.Context) *MarketInfo {
	v, err := mic.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (mic *MarketInfoCreate) Exec(ctx context.Context) error {
	_, err := mic.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mic *MarketInfoCreate) ExecX(ctx context.Context) {
	if err := mic.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mic *MarketInfoCreate) check() error {
	if _, ok := mic.mutation.HoursStart(); !ok {
		return &ValidationError{Name: "hours_start", err: errors.New(`ent: missing required field "MarketInfo.hours_start"`)}
	}
	if _, ok := mic.mutation.HoursEnd(); !ok {
		return &ValidationError{Name: "hours_end", err: errors.New(`ent: missing required field "MarketInfo.hours_end"`)}
	}
	return nil
}

func (mic *MarketInfoCreate) sqlSave(ctx context.Context) (*MarketInfo, error) {
	if err := mic.check(); err != nil {
		return nil, err
	}
	_node, _spec := mic.createSpec()
	if err := sqlgraph.CreateNode(ctx, mic.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	mic.mutation.id = &_node.ID
	mic.mutation.done = true
	return _node, nil
}

func (mic *MarketInfoCreate) createSpec() (*MarketInfo, *sqlgraph.CreateSpec) {
	var (
		_node = &MarketInfo{config: mic.config}
		_spec = sqlgraph.NewCreateSpec(marketinfo.Table, sqlgraph.NewFieldSpec(marketinfo.FieldID, field.TypeInt))
	)
	if value, ok := mic.mutation.HoursStart(); ok {
		_spec.SetField(marketinfo.FieldHoursStart, field.TypeTime, value)
		_node.HoursStart = value
	}
	if value, ok := mic.mutation.HoursEnd(); ok {
		_spec.SetField(marketinfo.FieldHoursEnd, field.TypeTime, value)
		_node.HoursEnd = value
	}
	if nodes := mic.mutation.HoursIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   marketinfo.HoursTable,
			Columns: []string{marketinfo.HoursColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(markethours.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// MarketInfoCreateBulk is the builder for creating many MarketInfo entities in bulk.
type MarketInfoCreateBulk struct {
	config
	err      error
	builders []*MarketInfoCreate
}

// Save creates the MarketInfo entities in the database.
func (micb *MarketInfoCreateBulk) Save(ctx context.Context) ([]*MarketInfo, error) {
	if micb.err != nil {
		return nil, micb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(micb.builders))
	nodes := make([]*MarketInfo, len(micb.builders))
	mutators := make([]Mutator, len(micb.builders))
	for i := range micb.builders {
		func(i int, root context.Context) {
			builder := micb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*MarketInfoMutation)
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
					_, err = mutators[i+1].Mutate(root, micb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, micb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, micb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (micb *MarketInfoCreateBulk) SaveX(ctx context.Context) []*MarketInfo {
	v, err := micb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (micb *MarketInfoCreateBulk) Exec(ctx context.Context) error {
	_, err := micb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (micb *MarketInfoCreateBulk) ExecX(ctx context.Context) {
	if err := micb.Exec(ctx); err != nil {
		panic(err)
	}
}
