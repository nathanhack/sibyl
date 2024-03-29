// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/dividend"
	"github.com/nathanhack/sibyl/ent/entity"
)

// DividendCreate is the builder for creating a Dividend entity.
type DividendCreate struct {
	config
	mutation *DividendMutation
	hooks    []Hook
}

// SetCashAmount sets the "cash_amount" field.
func (dc *DividendCreate) SetCashAmount(f float64) *DividendCreate {
	dc.mutation.SetCashAmount(f)
	return dc
}

// SetDeclarationDate sets the "declaration_date" field.
func (dc *DividendCreate) SetDeclarationDate(t time.Time) *DividendCreate {
	dc.mutation.SetDeclarationDate(t)
	return dc
}

// SetDividendType sets the "dividend_type" field.
func (dc *DividendCreate) SetDividendType(dt dividend.DividendType) *DividendCreate {
	dc.mutation.SetDividendType(dt)
	return dc
}

// SetExDividendDate sets the "ex_dividend_date" field.
func (dc *DividendCreate) SetExDividendDate(t time.Time) *DividendCreate {
	dc.mutation.SetExDividendDate(t)
	return dc
}

// SetFrequency sets the "frequency" field.
func (dc *DividendCreate) SetFrequency(i int) *DividendCreate {
	dc.mutation.SetFrequency(i)
	return dc
}

// SetRecordDate sets the "record_date" field.
func (dc *DividendCreate) SetRecordDate(t time.Time) *DividendCreate {
	dc.mutation.SetRecordDate(t)
	return dc
}

// SetPayDate sets the "pay_date" field.
func (dc *DividendCreate) SetPayDate(t time.Time) *DividendCreate {
	dc.mutation.SetPayDate(t)
	return dc
}

// AddStockIDs adds the "stock" edge to the Entity entity by IDs.
func (dc *DividendCreate) AddStockIDs(ids ...int) *DividendCreate {
	dc.mutation.AddStockIDs(ids...)
	return dc
}

// AddStock adds the "stock" edges to the Entity entity.
func (dc *DividendCreate) AddStock(e ...*Entity) *DividendCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return dc.AddStockIDs(ids...)
}

// Mutation returns the DividendMutation object of the builder.
func (dc *DividendCreate) Mutation() *DividendMutation {
	return dc.mutation
}

// Save creates the Dividend in the database.
func (dc *DividendCreate) Save(ctx context.Context) (*Dividend, error) {
	return withHooks[*Dividend, DividendMutation](ctx, dc.sqlSave, dc.mutation, dc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (dc *DividendCreate) SaveX(ctx context.Context) *Dividend {
	v, err := dc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (dc *DividendCreate) Exec(ctx context.Context) error {
	_, err := dc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (dc *DividendCreate) ExecX(ctx context.Context) {
	if err := dc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (dc *DividendCreate) check() error {
	if _, ok := dc.mutation.CashAmount(); !ok {
		return &ValidationError{Name: "cash_amount", err: errors.New(`ent: missing required field "Dividend.cash_amount"`)}
	}
	if _, ok := dc.mutation.DeclarationDate(); !ok {
		return &ValidationError{Name: "declaration_date", err: errors.New(`ent: missing required field "Dividend.declaration_date"`)}
	}
	if _, ok := dc.mutation.DividendType(); !ok {
		return &ValidationError{Name: "dividend_type", err: errors.New(`ent: missing required field "Dividend.dividend_type"`)}
	}
	if v, ok := dc.mutation.DividendType(); ok {
		if err := dividend.DividendTypeValidator(v); err != nil {
			return &ValidationError{Name: "dividend_type", err: fmt.Errorf(`ent: validator failed for field "Dividend.dividend_type": %w`, err)}
		}
	}
	if _, ok := dc.mutation.ExDividendDate(); !ok {
		return &ValidationError{Name: "ex_dividend_date", err: errors.New(`ent: missing required field "Dividend.ex_dividend_date"`)}
	}
	if _, ok := dc.mutation.Frequency(); !ok {
		return &ValidationError{Name: "frequency", err: errors.New(`ent: missing required field "Dividend.frequency"`)}
	}
	if _, ok := dc.mutation.RecordDate(); !ok {
		return &ValidationError{Name: "record_date", err: errors.New(`ent: missing required field "Dividend.record_date"`)}
	}
	if _, ok := dc.mutation.PayDate(); !ok {
		return &ValidationError{Name: "pay_date", err: errors.New(`ent: missing required field "Dividend.pay_date"`)}
	}
	return nil
}

func (dc *DividendCreate) sqlSave(ctx context.Context) (*Dividend, error) {
	if err := dc.check(); err != nil {
		return nil, err
	}
	_node, _spec := dc.createSpec()
	if err := sqlgraph.CreateNode(ctx, dc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	dc.mutation.id = &_node.ID
	dc.mutation.done = true
	return _node, nil
}

func (dc *DividendCreate) createSpec() (*Dividend, *sqlgraph.CreateSpec) {
	var (
		_node = &Dividend{config: dc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: dividend.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: dividend.FieldID,
			},
		}
	)
	if value, ok := dc.mutation.CashAmount(); ok {
		_spec.SetField(dividend.FieldCashAmount, field.TypeFloat64, value)
		_node.CashAmount = value
	}
	if value, ok := dc.mutation.DeclarationDate(); ok {
		_spec.SetField(dividend.FieldDeclarationDate, field.TypeTime, value)
		_node.DeclarationDate = value
	}
	if value, ok := dc.mutation.DividendType(); ok {
		_spec.SetField(dividend.FieldDividendType, field.TypeEnum, value)
		_node.DividendType = value
	}
	if value, ok := dc.mutation.ExDividendDate(); ok {
		_spec.SetField(dividend.FieldExDividendDate, field.TypeTime, value)
		_node.ExDividendDate = value
	}
	if value, ok := dc.mutation.Frequency(); ok {
		_spec.SetField(dividend.FieldFrequency, field.TypeInt, value)
		_node.Frequency = value
	}
	if value, ok := dc.mutation.RecordDate(); ok {
		_spec.SetField(dividend.FieldRecordDate, field.TypeTime, value)
		_node.RecordDate = value
	}
	if value, ok := dc.mutation.PayDate(); ok {
		_spec.SetField(dividend.FieldPayDate, field.TypeTime, value)
		_node.PayDate = value
	}
	if nodes := dc.mutation.StockIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   dividend.StockTable,
			Columns: dividend.StockPrimaryKey,
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

// DividendCreateBulk is the builder for creating many Dividend entities in bulk.
type DividendCreateBulk struct {
	config
	builders []*DividendCreate
}

// Save creates the Dividend entities in the database.
func (dcb *DividendCreateBulk) Save(ctx context.Context) ([]*Dividend, error) {
	specs := make([]*sqlgraph.CreateSpec, len(dcb.builders))
	nodes := make([]*Dividend, len(dcb.builders))
	mutators := make([]Mutator, len(dcb.builders))
	for i := range dcb.builders {
		func(i int, root context.Context) {
			builder := dcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*DividendMutation)
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
					_, err = mutators[i+1].Mutate(root, dcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, dcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, dcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (dcb *DividendCreateBulk) SaveX(ctx context.Context) []*Dividend {
	v, err := dcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (dcb *DividendCreateBulk) Exec(ctx context.Context) error {
	_, err := dcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (dcb *DividendCreateBulk) ExecX(ctx context.Context) {
	if err := dcb.Exec(ctx); err != nil {
		panic(err)
	}
}
