// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/entity"
	"github.com/nathanhack/sibyl/ent/financial"
)

// FinancialCreate is the builder for creating a Financial entity.
type FinancialCreate struct {
	config
	mutation *FinancialMutation
	hooks    []Hook
}

// AddStockIDs adds the "stock" edge to the Entity entity by IDs.
func (fc *FinancialCreate) AddStockIDs(ids ...int) *FinancialCreate {
	fc.mutation.AddStockIDs(ids...)
	return fc
}

// AddStock adds the "stock" edges to the Entity entity.
func (fc *FinancialCreate) AddStock(e ...*Entity) *FinancialCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return fc.AddStockIDs(ids...)
}

// Mutation returns the FinancialMutation object of the builder.
func (fc *FinancialCreate) Mutation() *FinancialMutation {
	return fc.mutation
}

// Save creates the Financial in the database.
func (fc *FinancialCreate) Save(ctx context.Context) (*Financial, error) {
	return withHooks[*Financial, FinancialMutation](ctx, fc.sqlSave, fc.mutation, fc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (fc *FinancialCreate) SaveX(ctx context.Context) *Financial {
	v, err := fc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fc *FinancialCreate) Exec(ctx context.Context) error {
	_, err := fc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fc *FinancialCreate) ExecX(ctx context.Context) {
	if err := fc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fc *FinancialCreate) check() error {
	return nil
}

func (fc *FinancialCreate) sqlSave(ctx context.Context) (*Financial, error) {
	if err := fc.check(); err != nil {
		return nil, err
	}
	_node, _spec := fc.createSpec()
	if err := sqlgraph.CreateNode(ctx, fc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	fc.mutation.id = &_node.ID
	fc.mutation.done = true
	return _node, nil
}

func (fc *FinancialCreate) createSpec() (*Financial, *sqlgraph.CreateSpec) {
	var (
		_node = &Financial{config: fc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: financial.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: financial.FieldID,
			},
		}
	)
	if nodes := fc.mutation.StockIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   financial.StockTable,
			Columns: financial.StockPrimaryKey,
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// FinancialCreateBulk is the builder for creating many Financial entities in bulk.
type FinancialCreateBulk struct {
	config
	builders []*FinancialCreate
}

// Save creates the Financial entities in the database.
func (fcb *FinancialCreateBulk) Save(ctx context.Context) ([]*Financial, error) {
	specs := make([]*sqlgraph.CreateSpec, len(fcb.builders))
	nodes := make([]*Financial, len(fcb.builders))
	mutators := make([]Mutator, len(fcb.builders))
	for i := range fcb.builders {
		func(i int, root context.Context) {
			builder := fcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*FinancialMutation)
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
					_, err = mutators[i+1].Mutate(root, fcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, fcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, fcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (fcb *FinancialCreateBulk) SaveX(ctx context.Context) []*Financial {
	v, err := fcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fcb *FinancialCreateBulk) Exec(ctx context.Context) error {
	_, err := fcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fcb *FinancialCreateBulk) ExecX(ctx context.Context) {
	if err := fcb.Exec(ctx); err != nil {
		panic(err)
	}
}
