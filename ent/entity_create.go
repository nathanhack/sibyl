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
	"github.com/nathanhack/sibyl/ent/exchange"
	"github.com/nathanhack/sibyl/ent/financial"
	"github.com/nathanhack/sibyl/ent/interval"
	"github.com/nathanhack/sibyl/ent/split"
)

// EntityCreate is the builder for creating a Entity entity.
type EntityCreate struct {
	config
	mutation *EntityMutation
	hooks    []Hook
}

// SetActive sets the "active" field.
func (ec *EntityCreate) SetActive(b bool) *EntityCreate {
	ec.mutation.SetActive(b)
	return ec
}

// SetTicker sets the "ticker" field.
func (ec *EntityCreate) SetTicker(s string) *EntityCreate {
	ec.mutation.SetTicker(s)
	return ec
}

// SetName sets the "name" field.
func (ec *EntityCreate) SetName(s string) *EntityCreate {
	ec.mutation.SetName(s)
	return ec
}

// SetDescription sets the "description" field.
func (ec *EntityCreate) SetDescription(s string) *EntityCreate {
	ec.mutation.SetDescription(s)
	return ec
}

// SetListDate sets the "list_date" field.
func (ec *EntityCreate) SetListDate(t time.Time) *EntityCreate {
	ec.mutation.SetListDate(t)
	return ec
}

// SetDelisted sets the "delisted" field.
func (ec *EntityCreate) SetDelisted(t time.Time) *EntityCreate {
	ec.mutation.SetDelisted(t)
	return ec
}

// SetNillableDelisted sets the "delisted" field if the given value is not nil.
func (ec *EntityCreate) SetNillableDelisted(t *time.Time) *EntityCreate {
	if t != nil {
		ec.SetDelisted(*t)
	}
	return ec
}

// AddExchangeIDs adds the "exchanges" edge to the Exchange entity by IDs.
func (ec *EntityCreate) AddExchangeIDs(ids ...int) *EntityCreate {
	ec.mutation.AddExchangeIDs(ids...)
	return ec
}

// AddExchanges adds the "exchanges" edges to the Exchange entity.
func (ec *EntityCreate) AddExchanges(e ...*Exchange) *EntityCreate {
	ids := make([]int, len(e))
	for i := range e {
		ids[i] = e[i].ID
	}
	return ec.AddExchangeIDs(ids...)
}

// AddIntervalIDs adds the "intervals" edge to the Interval entity by IDs.
func (ec *EntityCreate) AddIntervalIDs(ids ...int) *EntityCreate {
	ec.mutation.AddIntervalIDs(ids...)
	return ec
}

// AddIntervals adds the "intervals" edges to the Interval entity.
func (ec *EntityCreate) AddIntervals(i ...*Interval) *EntityCreate {
	ids := make([]int, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return ec.AddIntervalIDs(ids...)
}

// AddDividendIDs adds the "dividends" edge to the Dividend entity by IDs.
func (ec *EntityCreate) AddDividendIDs(ids ...int) *EntityCreate {
	ec.mutation.AddDividendIDs(ids...)
	return ec
}

// AddDividends adds the "dividends" edges to the Dividend entity.
func (ec *EntityCreate) AddDividends(d ...*Dividend) *EntityCreate {
	ids := make([]int, len(d))
	for i := range d {
		ids[i] = d[i].ID
	}
	return ec.AddDividendIDs(ids...)
}

// AddSplitIDs adds the "splits" edge to the Split entity by IDs.
func (ec *EntityCreate) AddSplitIDs(ids ...int) *EntityCreate {
	ec.mutation.AddSplitIDs(ids...)
	return ec
}

// AddSplits adds the "splits" edges to the Split entity.
func (ec *EntityCreate) AddSplits(s ...*Split) *EntityCreate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return ec.AddSplitIDs(ids...)
}

// AddFinancialIDs adds the "financials" edge to the Financial entity by IDs.
func (ec *EntityCreate) AddFinancialIDs(ids ...int) *EntityCreate {
	ec.mutation.AddFinancialIDs(ids...)
	return ec
}

// AddFinancials adds the "financials" edges to the Financial entity.
func (ec *EntityCreate) AddFinancials(f ...*Financial) *EntityCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return ec.AddFinancialIDs(ids...)
}

// Mutation returns the EntityMutation object of the builder.
func (ec *EntityCreate) Mutation() *EntityMutation {
	return ec.mutation
}

// Save creates the Entity in the database.
func (ec *EntityCreate) Save(ctx context.Context) (*Entity, error) {
	return withHooks[*Entity, EntityMutation](ctx, ec.sqlSave, ec.mutation, ec.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ec *EntityCreate) SaveX(ctx context.Context) *Entity {
	v, err := ec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ec *EntityCreate) Exec(ctx context.Context) error {
	_, err := ec.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ec *EntityCreate) ExecX(ctx context.Context) {
	if err := ec.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ec *EntityCreate) check() error {
	if _, ok := ec.mutation.Active(); !ok {
		return &ValidationError{Name: "active", err: errors.New(`ent: missing required field "Entity.active"`)}
	}
	if _, ok := ec.mutation.Ticker(); !ok {
		return &ValidationError{Name: "ticker", err: errors.New(`ent: missing required field "Entity.ticker"`)}
	}
	if v, ok := ec.mutation.Ticker(); ok {
		if err := entity.TickerValidator(v); err != nil {
			return &ValidationError{Name: "ticker", err: fmt.Errorf(`ent: validator failed for field "Entity.ticker": %w`, err)}
		}
	}
	if _, ok := ec.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Entity.name"`)}
	}
	if v, ok := ec.mutation.Name(); ok {
		if err := entity.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Entity.name": %w`, err)}
		}
	}
	if _, ok := ec.mutation.Description(); !ok {
		return &ValidationError{Name: "description", err: errors.New(`ent: missing required field "Entity.description"`)}
	}
	if v, ok := ec.mutation.Description(); ok {
		if err := entity.DescriptionValidator(v); err != nil {
			return &ValidationError{Name: "description", err: fmt.Errorf(`ent: validator failed for field "Entity.description": %w`, err)}
		}
	}
	if _, ok := ec.mutation.ListDate(); !ok {
		return &ValidationError{Name: "list_date", err: errors.New(`ent: missing required field "Entity.list_date"`)}
	}
	return nil
}

func (ec *EntityCreate) sqlSave(ctx context.Context) (*Entity, error) {
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

func (ec *EntityCreate) createSpec() (*Entity, *sqlgraph.CreateSpec) {
	var (
		_node = &Entity{config: ec.config}
		_spec = &sqlgraph.CreateSpec{
			Table: entity.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: entity.FieldID,
			},
		}
	)
	if value, ok := ec.mutation.Active(); ok {
		_spec.SetField(entity.FieldActive, field.TypeBool, value)
		_node.Active = value
	}
	if value, ok := ec.mutation.Ticker(); ok {
		_spec.SetField(entity.FieldTicker, field.TypeString, value)
		_node.Ticker = value
	}
	if value, ok := ec.mutation.Name(); ok {
		_spec.SetField(entity.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := ec.mutation.Description(); ok {
		_spec.SetField(entity.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if value, ok := ec.mutation.ListDate(); ok {
		_spec.SetField(entity.FieldListDate, field.TypeTime, value)
		_node.ListDate = value
	}
	if value, ok := ec.mutation.Delisted(); ok {
		_spec.SetField(entity.FieldDelisted, field.TypeTime, value)
		_node.Delisted = &value
	}
	if nodes := ec.mutation.ExchangesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   entity.ExchangesTable,
			Columns: entity.ExchangesPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: exchange.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.IntervalsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   entity.IntervalsTable,
			Columns: []string{entity.IntervalsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: interval.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.DividendsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   entity.DividendsTable,
			Columns: entity.DividendsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: dividend.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.SplitsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   entity.SplitsTable,
			Columns: []string{entity.SplitsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: split.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ec.mutation.FinancialsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   entity.FinancialsTable,
			Columns: entity.FinancialsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: financial.FieldID,
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

// EntityCreateBulk is the builder for creating many Entity entities in bulk.
type EntityCreateBulk struct {
	config
	builders []*EntityCreate
}

// Save creates the Entity entities in the database.
func (ecb *EntityCreateBulk) Save(ctx context.Context) ([]*Entity, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ecb.builders))
	nodes := make([]*Entity, len(ecb.builders))
	mutators := make([]Mutator, len(ecb.builders))
	for i := range ecb.builders {
		func(i int, root context.Context) {
			builder := ecb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*EntityMutation)
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
func (ecb *EntityCreateBulk) SaveX(ctx context.Context) []*Entity {
	v, err := ecb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ecb *EntityCreateBulk) Exec(ctx context.Context) error {
	_, err := ecb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ecb *EntityCreateBulk) ExecX(ctx context.Context) {
	if err := ecb.Exec(ctx); err != nil {
		panic(err)
	}
}