// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/entity"
	"github.com/nathanhack/sibyl/ent/split"
)

// SplitCreate is the builder for creating a Split entity.
type SplitCreate struct {
	config
	mutation *SplitMutation
	hooks    []Hook
}

// SetExecutionDate sets the "execution_date" field.
func (sc *SplitCreate) SetExecutionDate(t time.Time) *SplitCreate {
	sc.mutation.SetExecutionDate(t)
	return sc
}

// SetFrom sets the "from" field.
func (sc *SplitCreate) SetFrom(f float64) *SplitCreate {
	sc.mutation.SetFrom(f)
	return sc
}

// SetTo sets the "to" field.
func (sc *SplitCreate) SetTo(f float64) *SplitCreate {
	sc.mutation.SetTo(f)
	return sc
}

// SetStockID sets the "stock" edge to the Entity entity by ID.
func (sc *SplitCreate) SetStockID(id int) *SplitCreate {
	sc.mutation.SetStockID(id)
	return sc
}

// SetNillableStockID sets the "stock" edge to the Entity entity by ID if the given value is not nil.
func (sc *SplitCreate) SetNillableStockID(id *int) *SplitCreate {
	if id != nil {
		sc = sc.SetStockID(*id)
	}
	return sc
}

// SetStock sets the "stock" edge to the Entity entity.
func (sc *SplitCreate) SetStock(e *Entity) *SplitCreate {
	return sc.SetStockID(e.ID)
}

// Mutation returns the SplitMutation object of the builder.
func (sc *SplitCreate) Mutation() *SplitMutation {
	return sc.mutation
}

// Save creates the Split in the database.
func (sc *SplitCreate) Save(ctx context.Context) (*Split, error) {
	return withHooks[*Split, SplitMutation](ctx, sc.sqlSave, sc.mutation, sc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SplitCreate) SaveX(ctx context.Context) *Split {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *SplitCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *SplitCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *SplitCreate) check() error {
	if _, ok := sc.mutation.ExecutionDate(); !ok {
		return &ValidationError{Name: "execution_date", err: errors.New(`ent: missing required field "Split.execution_date"`)}
	}
	if _, ok := sc.mutation.From(); !ok {
		return &ValidationError{Name: "from", err: errors.New(`ent: missing required field "Split.from"`)}
	}
	if _, ok := sc.mutation.To(); !ok {
		return &ValidationError{Name: "to", err: errors.New(`ent: missing required field "Split.to"`)}
	}
	return nil
}

func (sc *SplitCreate) sqlSave(ctx context.Context) (*Split, error) {
	if err := sc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	sc.mutation.id = &_node.ID
	sc.mutation.done = true
	return _node, nil
}

func (sc *SplitCreate) createSpec() (*Split, *sqlgraph.CreateSpec) {
	var (
		_node = &Split{config: sc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: split.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: split.FieldID,
			},
		}
	)
	if value, ok := sc.mutation.ExecutionDate(); ok {
		_spec.SetField(split.FieldExecutionDate, field.TypeTime, value)
		_node.ExecutionDate = value
	}
	if value, ok := sc.mutation.From(); ok {
		_spec.SetField(split.FieldFrom, field.TypeFloat64, value)
		_node.From = value
	}
	if value, ok := sc.mutation.To(); ok {
		_spec.SetField(split.FieldTo, field.TypeFloat64, value)
		_node.To = value
	}
	if nodes := sc.mutation.StockIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   split.StockTable,
			Columns: []string{split.StockColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: entity.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.entity_splits = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// SplitCreateBulk is the builder for creating many Split entities in bulk.
type SplitCreateBulk struct {
	config
	builders []*SplitCreate
}

// Save creates the Split entities in the database.
func (scb *SplitCreateBulk) Save(ctx context.Context) ([]*Split, error) {
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Split, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SplitMutation)
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
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *SplitCreateBulk) SaveX(ctx context.Context) []*Split {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *SplitCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *SplitCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}
