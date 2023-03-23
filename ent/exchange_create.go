// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/entity"
	"github.com/nathanhack/sibyl/ent/exchange"
)

// ExchangeCreate is the builder for creating a Exchange entity.
type ExchangeCreate struct {
	config
	mutation *ExchangeMutation
	hooks    []Hook
}

// SetCode sets the "code" field.
func (ec *ExchangeCreate) SetCode(s string) *ExchangeCreate {
	ec.mutation.SetCode(s)
	return ec
}

// SetName sets the "name" field.
func (ec *ExchangeCreate) SetName(s string) *ExchangeCreate {
	ec.mutation.SetName(s)
	return ec
}

// AddStockIDs adds the "stocks" edge to the Entity entity by IDs.
func (ec *ExchangeCreate) AddStockIDs(ids ...int) *ExchangeCreate {
	ec.mutation.AddStockIDs(ids...)
	return ec
}

// AddStocks adds the "stocks" edges to the Entity entity.
func (ec *ExchangeCreate) AddStocks(e ...*Entity) *ExchangeCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ec.AddStockIDs(ids...)
}

// Mutation returns the ExchangeMutation object of the builder.
func (ec *ExchangeCreate) Mutation() *ExchangeMutation {
	return ec.mutation
}

// Save creates the Exchange in the database.
func (ec *ExchangeCreate) Save(ctx context.Context) (*Exchange, error) {
	return withHooks[*Exchange, ExchangeMutation](ctx, ec.sqlSave, ec.mutation, ec.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ec *ExchangeCreate) SaveX(ctx context.Context) *Exchange {
	v, err := ec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ec *ExchangeCreate) Exec(ctx context.Context) error {
	_, err := ec.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ec *ExchangeCreate) ExecX(ctx context.Context) {
	if err := ec.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ec *ExchangeCreate) check() error {
	if _, ok := ec.mutation.Code(); !ok {
		return &ValidationError{Name: "code", err: errors.New(`ent: missing required field "Exchange.code"`)}
	}
	if _, ok := ec.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Exchange.name"`)}
	}
	return nil
}

func (ec *ExchangeCreate) sqlSave(ctx context.Context) (*Exchange, error) {
	if err := ec.check(); err != nil {
		return nil, err
	}
	_node, _spec := ec.createSpec()
	if err := sqlgraph.CreateNode(ctx, ec.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	ec.mutation.id = &_node.ID
	ec.mutation.done = true
	return _node, nil
}

func (ec *ExchangeCreate) createSpec() (*Exchange, *sqlgraph.CreateSpec) {
	var (
		_node = &Exchange{config: ec.config}
		_spec = &sqlgraph.CreateSpec{
			Table: exchange.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: exchange.FieldID,
			},
		}
	)
	if value, ok := ec.mutation.Code(); ok {
		_spec.SetField(exchange.FieldCode, field.TypeString, value)
		_node.Code = value
	}
	if value, ok := ec.mutation.Name(); ok {
		_spec.SetField(exchange.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if nodes := ec.mutation.StocksIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   exchange.StocksTable,
			Columns: exchange.StocksPrimaryKey,
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

// ExchangeCreateBulk is the builder for creating many Exchange entities in bulk.
type ExchangeCreateBulk struct {
	config
	builders []*ExchangeCreate
}

// Save creates the Exchange entities in the database.
func (ecb *ExchangeCreateBulk) Save(ctx context.Context) ([]*Exchange, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ecb.builders))
	nodes := make([]*Exchange, len(ecb.builders))
	mutators := make([]Mutator, len(ecb.builders))
	for i := range ecb.builders {
		func(i int, root context.Context) {
			builder := ecb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*ExchangeMutation)
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
					_, err = mutators[i+1].Mutate(root, ecb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ecb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ecb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ecb *ExchangeCreateBulk) SaveX(ctx context.Context) []*Exchange {
	v, err := ecb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ecb *ExchangeCreateBulk) Exec(ctx context.Context) error {
	_, err := ecb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecb *ExchangeCreateBulk) ExecX(ctx context.Context) {
	if err := ecb.Exec(ctx); err != nil {
		panic(err)
	}
}
