// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/markethours"
	"github.com/nathanhack/sibyl/ent/marketinfo"
	"github.com/nathanhack/sibyl/ent/predicate"
)

// MarketInfoUpdate is the builder for updating MarketInfo entities.
type MarketInfoUpdate struct {
	config
	hooks    []Hook
	mutation *MarketInfoMutation
}

// Where appends a list predicates to the MarketInfoUpdate builder.
func (miu *MarketInfoUpdate) Where(ps ...predicate.MarketInfo) *MarketInfoUpdate {
	miu.mutation.Where(ps...)
	return miu
}

// SetHoursStart sets the "hours_start" field.
func (miu *MarketInfoUpdate) SetHoursStart(t time.Time) *MarketInfoUpdate {
	miu.mutation.SetHoursStart(t)
	return miu
}

// SetNillableHoursStart sets the "hours_start" field if the given value is not nil.
func (miu *MarketInfoUpdate) SetNillableHoursStart(t *time.Time) *MarketInfoUpdate {
	if t != nil {
		miu.SetHoursStart(*t)
	}
	return miu
}

// SetHoursEnd sets the "hours_end" field.
func (miu *MarketInfoUpdate) SetHoursEnd(t time.Time) *MarketInfoUpdate {
	miu.mutation.SetHoursEnd(t)
	return miu
}

// SetNillableHoursEnd sets the "hours_end" field if the given value is not nil.
func (miu *MarketInfoUpdate) SetNillableHoursEnd(t *time.Time) *MarketInfoUpdate {
	if t != nil {
		miu.SetHoursEnd(*t)
	}
	return miu
}

// AddHourIDs adds the "hours" edge to the MarketHours entity by IDs.
func (miu *MarketInfoUpdate) AddHourIDs(ids ...int) *MarketInfoUpdate {
	miu.mutation.AddHourIDs(ids...)
	return miu
}

// AddHours adds the "hours" edges to the MarketHours entity.
func (miu *MarketInfoUpdate) AddHours(m ...*MarketHours) *MarketInfoUpdate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return miu.AddHourIDs(ids...)
}

// Mutation returns the MarketInfoMutation object of the builder.
func (miu *MarketInfoUpdate) Mutation() *MarketInfoMutation {
	return miu.mutation
}

// ClearHours clears all "hours" edges to the MarketHours entity.
func (miu *MarketInfoUpdate) ClearHours() *MarketInfoUpdate {
	miu.mutation.ClearHours()
	return miu
}

// RemoveHourIDs removes the "hours" edge to MarketHours entities by IDs.
func (miu *MarketInfoUpdate) RemoveHourIDs(ids ...int) *MarketInfoUpdate {
	miu.mutation.RemoveHourIDs(ids...)
	return miu
}

// RemoveHours removes "hours" edges to MarketHours entities.
func (miu *MarketInfoUpdate) RemoveHours(m ...*MarketHours) *MarketInfoUpdate {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return miu.RemoveHourIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (miu *MarketInfoUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, miu.sqlSave, miu.mutation, miu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (miu *MarketInfoUpdate) SaveX(ctx context.Context) int {
	affected, err := miu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (miu *MarketInfoUpdate) Exec(ctx context.Context) error {
	_, err := miu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (miu *MarketInfoUpdate) ExecX(ctx context.Context) {
	if err := miu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (miu *MarketInfoUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(marketinfo.Table, marketinfo.Columns, sqlgraph.NewFieldSpec(marketinfo.FieldID, field.TypeInt))
	if ps := miu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := miu.mutation.HoursStart(); ok {
		_spec.SetField(marketinfo.FieldHoursStart, field.TypeTime, value)
	}
	if value, ok := miu.mutation.HoursEnd(); ok {
		_spec.SetField(marketinfo.FieldHoursEnd, field.TypeTime, value)
	}
	if miu.mutation.HoursCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := miu.mutation.RemovedHoursIDs(); len(nodes) > 0 && !miu.mutation.HoursCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := miu.mutation.HoursIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, miu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{marketinfo.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	miu.mutation.done = true
	return n, nil
}

// MarketInfoUpdateOne is the builder for updating a single MarketInfo entity.
type MarketInfoUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *MarketInfoMutation
}

// SetHoursStart sets the "hours_start" field.
func (miuo *MarketInfoUpdateOne) SetHoursStart(t time.Time) *MarketInfoUpdateOne {
	miuo.mutation.SetHoursStart(t)
	return miuo
}

// SetNillableHoursStart sets the "hours_start" field if the given value is not nil.
func (miuo *MarketInfoUpdateOne) SetNillableHoursStart(t *time.Time) *MarketInfoUpdateOne {
	if t != nil {
		miuo.SetHoursStart(*t)
	}
	return miuo
}

// SetHoursEnd sets the "hours_end" field.
func (miuo *MarketInfoUpdateOne) SetHoursEnd(t time.Time) *MarketInfoUpdateOne {
	miuo.mutation.SetHoursEnd(t)
	return miuo
}

// SetNillableHoursEnd sets the "hours_end" field if the given value is not nil.
func (miuo *MarketInfoUpdateOne) SetNillableHoursEnd(t *time.Time) *MarketInfoUpdateOne {
	if t != nil {
		miuo.SetHoursEnd(*t)
	}
	return miuo
}

// AddHourIDs adds the "hours" edge to the MarketHours entity by IDs.
func (miuo *MarketInfoUpdateOne) AddHourIDs(ids ...int) *MarketInfoUpdateOne {
	miuo.mutation.AddHourIDs(ids...)
	return miuo
}

// AddHours adds the "hours" edges to the MarketHours entity.
func (miuo *MarketInfoUpdateOne) AddHours(m ...*MarketHours) *MarketInfoUpdateOne {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return miuo.AddHourIDs(ids...)
}

// Mutation returns the MarketInfoMutation object of the builder.
func (miuo *MarketInfoUpdateOne) Mutation() *MarketInfoMutation {
	return miuo.mutation
}

// ClearHours clears all "hours" edges to the MarketHours entity.
func (miuo *MarketInfoUpdateOne) ClearHours() *MarketInfoUpdateOne {
	miuo.mutation.ClearHours()
	return miuo
}

// RemoveHourIDs removes the "hours" edge to MarketHours entities by IDs.
func (miuo *MarketInfoUpdateOne) RemoveHourIDs(ids ...int) *MarketInfoUpdateOne {
	miuo.mutation.RemoveHourIDs(ids...)
	return miuo
}

// RemoveHours removes "hours" edges to MarketHours entities.
func (miuo *MarketInfoUpdateOne) RemoveHours(m ...*MarketHours) *MarketInfoUpdateOne {
	ids := make([]int, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return miuo.RemoveHourIDs(ids...)
}

// Where appends a list predicates to the MarketInfoUpdate builder.
func (miuo *MarketInfoUpdateOne) Where(ps ...predicate.MarketInfo) *MarketInfoUpdateOne {
	miuo.mutation.Where(ps...)
	return miuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (miuo *MarketInfoUpdateOne) Select(field string, fields ...string) *MarketInfoUpdateOne {
	miuo.fields = append([]string{field}, fields...)
	return miuo
}

// Save executes the query and returns the updated MarketInfo entity.
func (miuo *MarketInfoUpdateOne) Save(ctx context.Context) (*MarketInfo, error) {
	return withHooks(ctx, miuo.sqlSave, miuo.mutation, miuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (miuo *MarketInfoUpdateOne) SaveX(ctx context.Context) *MarketInfo {
	node, err := miuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (miuo *MarketInfoUpdateOne) Exec(ctx context.Context) error {
	_, err := miuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (miuo *MarketInfoUpdateOne) ExecX(ctx context.Context) {
	if err := miuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (miuo *MarketInfoUpdateOne) sqlSave(ctx context.Context) (_node *MarketInfo, err error) {
	_spec := sqlgraph.NewUpdateSpec(marketinfo.Table, marketinfo.Columns, sqlgraph.NewFieldSpec(marketinfo.FieldID, field.TypeInt))
	id, ok := miuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "MarketInfo.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := miuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, marketinfo.FieldID)
		for _, f := range fields {
			if !marketinfo.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != marketinfo.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := miuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := miuo.mutation.HoursStart(); ok {
		_spec.SetField(marketinfo.FieldHoursStart, field.TypeTime, value)
	}
	if value, ok := miuo.mutation.HoursEnd(); ok {
		_spec.SetField(marketinfo.FieldHoursEnd, field.TypeTime, value)
	}
	if miuo.mutation.HoursCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := miuo.mutation.RemovedHoursIDs(); len(nodes) > 0 && !miuo.mutation.HoursCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := miuo.mutation.HoursIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &MarketInfo{config: miuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, miuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{marketinfo.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	miuo.mutation.done = true
	return _node, nil
}
