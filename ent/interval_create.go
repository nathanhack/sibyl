// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/bartimerange"
	"github.com/nathanhack/sibyl/ent/datasource"
	"github.com/nathanhack/sibyl/ent/entity"
	"github.com/nathanhack/sibyl/ent/interval"
	"github.com/nathanhack/sibyl/ent/tradetimerange"
)

// IntervalCreate is the builder for creating a Interval entity.
type IntervalCreate struct {
	config
	mutation *IntervalMutation
	hooks    []Hook
}

// SetActive sets the "active" field.
func (ic *IntervalCreate) SetActive(b bool) *IntervalCreate {
	ic.mutation.SetActive(b)
	return ic
}

// SetNillableActive sets the "active" field if the given value is not nil.
func (ic *IntervalCreate) SetNillableActive(b *bool) *IntervalCreate {
	if b != nil {
		ic.SetActive(*b)
	}
	return ic
}

// SetInterval sets the "interval" field.
func (ic *IntervalCreate) SetInterval(i interval.Interval) *IntervalCreate {
	ic.mutation.SetInterval(i)
	return ic
}

// SetStockID sets the "stock_id" field.
func (ic *IntervalCreate) SetStockID(i int) *IntervalCreate {
	ic.mutation.SetStockID(i)
	return ic
}

// SetDataSourceID sets the "data_source_id" field.
func (ic *IntervalCreate) SetDataSourceID(i int) *IntervalCreate {
	ic.mutation.SetDataSourceID(i)
	return ic
}

// SetDataSource sets the "data_source" edge to the DataSource entity.
func (ic *IntervalCreate) SetDataSource(d *DataSource) *IntervalCreate {
	return ic.SetDataSourceID(d.ID)
}

// SetStock sets the "stock" edge to the Entity entity.
func (ic *IntervalCreate) SetStock(e *Entity) *IntervalCreate {
	return ic.SetStockID(e.ID)
}

// AddBarIDs adds the "bars" edge to the BarTimeRange entity by IDs.
func (ic *IntervalCreate) AddBarIDs(ids ...int) *IntervalCreate {
	ic.mutation.AddBarIDs(ids...)
	return ic
}

// AddBars adds the "bars" edges to the BarTimeRange entity.
func (ic *IntervalCreate) AddBars(b ...*BarTimeRange) *IntervalCreate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return ic.AddBarIDs(ids...)
}

// AddTradeIDs adds the "trades" edge to the TradeTimeRange entity by IDs.
func (ic *IntervalCreate) AddTradeIDs(ids ...int) *IntervalCreate {
	ic.mutation.AddTradeIDs(ids...)
	return ic
}

// AddTrades adds the "trades" edges to the TradeTimeRange entity.
func (ic *IntervalCreate) AddTrades(t ...*TradeTimeRange) *IntervalCreate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return ic.AddTradeIDs(ids...)
}

// Mutation returns the IntervalMutation object of the builder.
func (ic *IntervalCreate) Mutation() *IntervalMutation {
	return ic.mutation
}

// Save creates the Interval in the database.
func (ic *IntervalCreate) Save(ctx context.Context) (*Interval, error) {
	ic.defaults()
	return withHooks[*Interval, IntervalMutation](ctx, ic.sqlSave, ic.mutation, ic.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (ic *IntervalCreate) SaveX(ctx context.Context) *Interval {
	v, err := ic.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ic *IntervalCreate) Exec(ctx context.Context) error {
	_, err := ic.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ic *IntervalCreate) ExecX(ctx context.Context) {
	if err := ic.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (ic *IntervalCreate) defaults() {
	if _, ok := ic.mutation.Active(); !ok {
		v := interval.DefaultActive
		ic.mutation.SetActive(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (ic *IntervalCreate) check() error {
	if _, ok := ic.mutation.Active(); !ok {
		return &ValidationError{Name: "active", err: errors.New(`ent: missing required field "Interval.active"`)}
	}
	if _, ok := ic.mutation.Interval(); !ok {
		return &ValidationError{Name: "interval", err: errors.New(`ent: missing required field "Interval.interval"`)}
	}
	if v, ok := ic.mutation.Interval(); ok {
		if err := interval.IntervalValidator(v); err != nil {
			return &ValidationError{Name: "interval", err: fmt.Errorf(`ent: validator failed for field "Interval.interval": %w`, err)}
		}
	}
	if _, ok := ic.mutation.StockID(); !ok {
		return &ValidationError{Name: "stock_id", err: errors.New(`ent: missing required field "Interval.stock_id"`)}
	}
	if _, ok := ic.mutation.DataSourceID(); !ok {
		return &ValidationError{Name: "data_source_id", err: errors.New(`ent: missing required field "Interval.data_source_id"`)}
	}
	if _, ok := ic.mutation.DataSourceID(); !ok {
		return &ValidationError{Name: "data_source", err: errors.New(`ent: missing required edge "Interval.data_source"`)}
	}
	if _, ok := ic.mutation.StockID(); !ok {
		return &ValidationError{Name: "stock", err: errors.New(`ent: missing required edge "Interval.stock"`)}
	}
	return nil
}

func (ic *IntervalCreate) sqlSave(ctx context.Context) (*Interval, error) {
	if err := ic.check(); err != nil {
		return nil, err
	}
	_node, _spec := ic.createSpec()
	if err := sqlgraph.CreateNode(ctx, ic.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	ic.mutation.id = &_node.ID
	ic.mutation.done = true
	return _node, nil
}

func (ic *IntervalCreate) createSpec() (*Interval, *sqlgraph.CreateSpec) {
	var (
		_node = &Interval{config: ic.config}
		_spec = &sqlgraph.CreateSpec{
			Table: interval.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: interval.FieldID,
			},
		}
	)
	if value, ok := ic.mutation.Active(); ok {
		_spec.SetField(interval.FieldActive, field.TypeBool, value)
		_node.Active = value
	}
	if value, ok := ic.mutation.Interval(); ok {
		_spec.SetField(interval.FieldInterval, field.TypeEnum, value)
		_node.Interval = value
	}
	if nodes := ic.mutation.DataSourceIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   interval.DataSourceTable,
			Columns: []string{interval.DataSourceColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: datasource.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.DataSourceID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ic.mutation.StockIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   interval.StockTable,
			Columns: []string{interval.StockColumn},
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
		_node.StockID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ic.mutation.BarsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   interval.BarsTable,
			Columns: []string{interval.BarsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: bartimerange.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := ic.mutation.TradesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   interval.TradesTable,
			Columns: []string{interval.TradesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: tradetimerange.FieldID,
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

// IntervalCreateBulk is the builder for creating many Interval entities in bulk.
type IntervalCreateBulk struct {
	config
	builders []*IntervalCreate
}

// Save creates the Interval entities in the database.
func (icb *IntervalCreateBulk) Save(ctx context.Context) ([]*Interval, error) {
	specs := make([]*sqlgraph.CreateSpec, len(icb.builders))
	nodes := make([]*Interval, len(icb.builders))
	mutators := make([]Mutator, len(icb.builders))
	for i := range icb.builders {
		func(i int, root context.Context) {
			builder := icb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*IntervalMutation)
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
					_, err = mutators[i+1].Mutate(root, icb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, icb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, icb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (icb *IntervalCreateBulk) SaveX(ctx context.Context) []*Interval {
	v, err := icb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (icb *IntervalCreateBulk) Exec(ctx context.Context) error {
	_, err := icb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (icb *IntervalCreateBulk) ExecX(ctx context.Context) {
	if err := icb.Exec(ctx); err != nil {
		panic(err)
	}
}
