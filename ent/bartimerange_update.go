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
	"github.com/nathanhack/sibyl/ent/bargroup"
	"github.com/nathanhack/sibyl/ent/bartimerange"
	"github.com/nathanhack/sibyl/ent/interval"
	"github.com/nathanhack/sibyl/ent/predicate"
)

// BarTimeRangeUpdate is the builder for updating BarTimeRange entities.
type BarTimeRangeUpdate struct {
	config
	hooks    []Hook
	mutation *BarTimeRangeMutation
}

// Where appends a list predicates to the BarTimeRangeUpdate builder.
func (btru *BarTimeRangeUpdate) Where(ps ...predicate.BarTimeRange) *BarTimeRangeUpdate {
	btru.mutation.Where(ps...)
	return btru
}

// SetStart sets the "start" field.
func (btru *BarTimeRangeUpdate) SetStart(t time.Time) *BarTimeRangeUpdate {
	btru.mutation.SetStart(t)
	return btru
}

// SetNillableStart sets the "start" field if the given value is not nil.
func (btru *BarTimeRangeUpdate) SetNillableStart(t *time.Time) *BarTimeRangeUpdate {
	if t != nil {
		btru.SetStart(*t)
	}
	return btru
}

// SetEnd sets the "end" field.
func (btru *BarTimeRangeUpdate) SetEnd(t time.Time) *BarTimeRangeUpdate {
	btru.mutation.SetEnd(t)
	return btru
}

// SetNillableEnd sets the "end" field if the given value is not nil.
func (btru *BarTimeRangeUpdate) SetNillableEnd(t *time.Time) *BarTimeRangeUpdate {
	if t != nil {
		btru.SetEnd(*t)
	}
	return btru
}

// SetCount sets the "count" field.
func (btru *BarTimeRangeUpdate) SetCount(i int) *BarTimeRangeUpdate {
	btru.mutation.ResetCount()
	btru.mutation.SetCount(i)
	return btru
}

// SetNillableCount sets the "count" field if the given value is not nil.
func (btru *BarTimeRangeUpdate) SetNillableCount(i *int) *BarTimeRangeUpdate {
	if i != nil {
		btru.SetCount(*i)
	}
	return btru
}

// AddCount adds i to the "count" field.
func (btru *BarTimeRangeUpdate) AddCount(i int) *BarTimeRangeUpdate {
	btru.mutation.AddCount(i)
	return btru
}

// SetIntervalID sets the "interval_id" field.
func (btru *BarTimeRangeUpdate) SetIntervalID(i int) *BarTimeRangeUpdate {
	btru.mutation.SetIntervalID(i)
	return btru
}

// SetNillableIntervalID sets the "interval_id" field if the given value is not nil.
func (btru *BarTimeRangeUpdate) SetNillableIntervalID(i *int) *BarTimeRangeUpdate {
	if i != nil {
		btru.SetIntervalID(*i)
	}
	return btru
}

// SetStatus sets the "status" field.
func (btru *BarTimeRangeUpdate) SetStatus(b bartimerange.Status) *BarTimeRangeUpdate {
	btru.mutation.SetStatus(b)
	return btru
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (btru *BarTimeRangeUpdate) SetNillableStatus(b *bartimerange.Status) *BarTimeRangeUpdate {
	if b != nil {
		btru.SetStatus(*b)
	}
	return btru
}

// SetUpdateTime sets the "update_time" field.
func (btru *BarTimeRangeUpdate) SetUpdateTime(t time.Time) *BarTimeRangeUpdate {
	btru.mutation.SetUpdateTime(t)
	return btru
}

// SetInterval sets the "interval" edge to the Interval entity.
func (btru *BarTimeRangeUpdate) SetInterval(i *Interval) *BarTimeRangeUpdate {
	return btru.SetIntervalID(i.ID)
}

// AddGroupIDs adds the "groups" edge to the BarGroup entity by IDs.
func (btru *BarTimeRangeUpdate) AddGroupIDs(ids ...int) *BarTimeRangeUpdate {
	btru.mutation.AddGroupIDs(ids...)
	return btru
}

// AddGroups adds the "groups" edges to the BarGroup entity.
func (btru *BarTimeRangeUpdate) AddGroups(b ...*BarGroup) *BarTimeRangeUpdate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return btru.AddGroupIDs(ids...)
}

// Mutation returns the BarTimeRangeMutation object of the builder.
func (btru *BarTimeRangeUpdate) Mutation() *BarTimeRangeMutation {
	return btru.mutation
}

// ClearInterval clears the "interval" edge to the Interval entity.
func (btru *BarTimeRangeUpdate) ClearInterval() *BarTimeRangeUpdate {
	btru.mutation.ClearInterval()
	return btru
}

// ClearGroups clears all "groups" edges to the BarGroup entity.
func (btru *BarTimeRangeUpdate) ClearGroups() *BarTimeRangeUpdate {
	btru.mutation.ClearGroups()
	return btru
}

// RemoveGroupIDs removes the "groups" edge to BarGroup entities by IDs.
func (btru *BarTimeRangeUpdate) RemoveGroupIDs(ids ...int) *BarTimeRangeUpdate {
	btru.mutation.RemoveGroupIDs(ids...)
	return btru
}

// RemoveGroups removes "groups" edges to BarGroup entities.
func (btru *BarTimeRangeUpdate) RemoveGroups(b ...*BarGroup) *BarTimeRangeUpdate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return btru.RemoveGroupIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (btru *BarTimeRangeUpdate) Save(ctx context.Context) (int, error) {
	btru.defaults()
	return withHooks(ctx, btru.sqlSave, btru.mutation, btru.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (btru *BarTimeRangeUpdate) SaveX(ctx context.Context) int {
	affected, err := btru.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (btru *BarTimeRangeUpdate) Exec(ctx context.Context) error {
	_, err := btru.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (btru *BarTimeRangeUpdate) ExecX(ctx context.Context) {
	if err := btru.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (btru *BarTimeRangeUpdate) defaults() {
	if _, ok := btru.mutation.UpdateTime(); !ok {
		v := bartimerange.UpdateDefaultUpdateTime()
		btru.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (btru *BarTimeRangeUpdate) check() error {
	if v, ok := btru.mutation.Status(); ok {
		if err := bartimerange.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "BarTimeRange.status": %w`, err)}
		}
	}
	if btru.mutation.IntervalCleared() && len(btru.mutation.IntervalIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "BarTimeRange.interval"`)
	}
	return nil
}

func (btru *BarTimeRangeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := btru.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(bartimerange.Table, bartimerange.Columns, sqlgraph.NewFieldSpec(bartimerange.FieldID, field.TypeInt))
	if ps := btru.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := btru.mutation.Start(); ok {
		_spec.SetField(bartimerange.FieldStart, field.TypeTime, value)
	}
	if value, ok := btru.mutation.End(); ok {
		_spec.SetField(bartimerange.FieldEnd, field.TypeTime, value)
	}
	if value, ok := btru.mutation.Count(); ok {
		_spec.SetField(bartimerange.FieldCount, field.TypeInt, value)
	}
	if value, ok := btru.mutation.AddedCount(); ok {
		_spec.AddField(bartimerange.FieldCount, field.TypeInt, value)
	}
	if value, ok := btru.mutation.Status(); ok {
		_spec.SetField(bartimerange.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := btru.mutation.UpdateTime(); ok {
		_spec.SetField(bartimerange.FieldUpdateTime, field.TypeTime, value)
	}
	if btru.mutation.IntervalCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   bartimerange.IntervalTable,
			Columns: []string{bartimerange.IntervalColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(interval.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := btru.mutation.IntervalIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   bartimerange.IntervalTable,
			Columns: []string{bartimerange.IntervalColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(interval.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if btru.mutation.GroupsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   bartimerange.GroupsTable,
			Columns: []string{bartimerange.GroupsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bargroup.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := btru.mutation.RemovedGroupsIDs(); len(nodes) > 0 && !btru.mutation.GroupsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   bartimerange.GroupsTable,
			Columns: []string{bartimerange.GroupsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bargroup.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := btru.mutation.GroupsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   bartimerange.GroupsTable,
			Columns: []string{bartimerange.GroupsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bargroup.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, btru.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{bartimerange.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	btru.mutation.done = true
	return n, nil
}

// BarTimeRangeUpdateOne is the builder for updating a single BarTimeRange entity.
type BarTimeRangeUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *BarTimeRangeMutation
}

// SetStart sets the "start" field.
func (btruo *BarTimeRangeUpdateOne) SetStart(t time.Time) *BarTimeRangeUpdateOne {
	btruo.mutation.SetStart(t)
	return btruo
}

// SetNillableStart sets the "start" field if the given value is not nil.
func (btruo *BarTimeRangeUpdateOne) SetNillableStart(t *time.Time) *BarTimeRangeUpdateOne {
	if t != nil {
		btruo.SetStart(*t)
	}
	return btruo
}

// SetEnd sets the "end" field.
func (btruo *BarTimeRangeUpdateOne) SetEnd(t time.Time) *BarTimeRangeUpdateOne {
	btruo.mutation.SetEnd(t)
	return btruo
}

// SetNillableEnd sets the "end" field if the given value is not nil.
func (btruo *BarTimeRangeUpdateOne) SetNillableEnd(t *time.Time) *BarTimeRangeUpdateOne {
	if t != nil {
		btruo.SetEnd(*t)
	}
	return btruo
}

// SetCount sets the "count" field.
func (btruo *BarTimeRangeUpdateOne) SetCount(i int) *BarTimeRangeUpdateOne {
	btruo.mutation.ResetCount()
	btruo.mutation.SetCount(i)
	return btruo
}

// SetNillableCount sets the "count" field if the given value is not nil.
func (btruo *BarTimeRangeUpdateOne) SetNillableCount(i *int) *BarTimeRangeUpdateOne {
	if i != nil {
		btruo.SetCount(*i)
	}
	return btruo
}

// AddCount adds i to the "count" field.
func (btruo *BarTimeRangeUpdateOne) AddCount(i int) *BarTimeRangeUpdateOne {
	btruo.mutation.AddCount(i)
	return btruo
}

// SetIntervalID sets the "interval_id" field.
func (btruo *BarTimeRangeUpdateOne) SetIntervalID(i int) *BarTimeRangeUpdateOne {
	btruo.mutation.SetIntervalID(i)
	return btruo
}

// SetNillableIntervalID sets the "interval_id" field if the given value is not nil.
func (btruo *BarTimeRangeUpdateOne) SetNillableIntervalID(i *int) *BarTimeRangeUpdateOne {
	if i != nil {
		btruo.SetIntervalID(*i)
	}
	return btruo
}

// SetStatus sets the "status" field.
func (btruo *BarTimeRangeUpdateOne) SetStatus(b bartimerange.Status) *BarTimeRangeUpdateOne {
	btruo.mutation.SetStatus(b)
	return btruo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (btruo *BarTimeRangeUpdateOne) SetNillableStatus(b *bartimerange.Status) *BarTimeRangeUpdateOne {
	if b != nil {
		btruo.SetStatus(*b)
	}
	return btruo
}

// SetUpdateTime sets the "update_time" field.
func (btruo *BarTimeRangeUpdateOne) SetUpdateTime(t time.Time) *BarTimeRangeUpdateOne {
	btruo.mutation.SetUpdateTime(t)
	return btruo
}

// SetInterval sets the "interval" edge to the Interval entity.
func (btruo *BarTimeRangeUpdateOne) SetInterval(i *Interval) *BarTimeRangeUpdateOne {
	return btruo.SetIntervalID(i.ID)
}

// AddGroupIDs adds the "groups" edge to the BarGroup entity by IDs.
func (btruo *BarTimeRangeUpdateOne) AddGroupIDs(ids ...int) *BarTimeRangeUpdateOne {
	btruo.mutation.AddGroupIDs(ids...)
	return btruo
}

// AddGroups adds the "groups" edges to the BarGroup entity.
func (btruo *BarTimeRangeUpdateOne) AddGroups(b ...*BarGroup) *BarTimeRangeUpdateOne {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return btruo.AddGroupIDs(ids...)
}

// Mutation returns the BarTimeRangeMutation object of the builder.
func (btruo *BarTimeRangeUpdateOne) Mutation() *BarTimeRangeMutation {
	return btruo.mutation
}

// ClearInterval clears the "interval" edge to the Interval entity.
func (btruo *BarTimeRangeUpdateOne) ClearInterval() *BarTimeRangeUpdateOne {
	btruo.mutation.ClearInterval()
	return btruo
}

// ClearGroups clears all "groups" edges to the BarGroup entity.
func (btruo *BarTimeRangeUpdateOne) ClearGroups() *BarTimeRangeUpdateOne {
	btruo.mutation.ClearGroups()
	return btruo
}

// RemoveGroupIDs removes the "groups" edge to BarGroup entities by IDs.
func (btruo *BarTimeRangeUpdateOne) RemoveGroupIDs(ids ...int) *BarTimeRangeUpdateOne {
	btruo.mutation.RemoveGroupIDs(ids...)
	return btruo
}

// RemoveGroups removes "groups" edges to BarGroup entities.
func (btruo *BarTimeRangeUpdateOne) RemoveGroups(b ...*BarGroup) *BarTimeRangeUpdateOne {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return btruo.RemoveGroupIDs(ids...)
}

// Where appends a list predicates to the BarTimeRangeUpdate builder.
func (btruo *BarTimeRangeUpdateOne) Where(ps ...predicate.BarTimeRange) *BarTimeRangeUpdateOne {
	btruo.mutation.Where(ps...)
	return btruo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (btruo *BarTimeRangeUpdateOne) Select(field string, fields ...string) *BarTimeRangeUpdateOne {
	btruo.fields = append([]string{field}, fields...)
	return btruo
}

// Save executes the query and returns the updated BarTimeRange entity.
func (btruo *BarTimeRangeUpdateOne) Save(ctx context.Context) (*BarTimeRange, error) {
	btruo.defaults()
	return withHooks(ctx, btruo.sqlSave, btruo.mutation, btruo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (btruo *BarTimeRangeUpdateOne) SaveX(ctx context.Context) *BarTimeRange {
	node, err := btruo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (btruo *BarTimeRangeUpdateOne) Exec(ctx context.Context) error {
	_, err := btruo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (btruo *BarTimeRangeUpdateOne) ExecX(ctx context.Context) {
	if err := btruo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (btruo *BarTimeRangeUpdateOne) defaults() {
	if _, ok := btruo.mutation.UpdateTime(); !ok {
		v := bartimerange.UpdateDefaultUpdateTime()
		btruo.mutation.SetUpdateTime(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (btruo *BarTimeRangeUpdateOne) check() error {
	if v, ok := btruo.mutation.Status(); ok {
		if err := bartimerange.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "BarTimeRange.status": %w`, err)}
		}
	}
	if btruo.mutation.IntervalCleared() && len(btruo.mutation.IntervalIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "BarTimeRange.interval"`)
	}
	return nil
}

func (btruo *BarTimeRangeUpdateOne) sqlSave(ctx context.Context) (_node *BarTimeRange, err error) {
	if err := btruo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(bartimerange.Table, bartimerange.Columns, sqlgraph.NewFieldSpec(bartimerange.FieldID, field.TypeInt))
	id, ok := btruo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "BarTimeRange.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := btruo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, bartimerange.FieldID)
		for _, f := range fields {
			if !bartimerange.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != bartimerange.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := btruo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := btruo.mutation.Start(); ok {
		_spec.SetField(bartimerange.FieldStart, field.TypeTime, value)
	}
	if value, ok := btruo.mutation.End(); ok {
		_spec.SetField(bartimerange.FieldEnd, field.TypeTime, value)
	}
	if value, ok := btruo.mutation.Count(); ok {
		_spec.SetField(bartimerange.FieldCount, field.TypeInt, value)
	}
	if value, ok := btruo.mutation.AddedCount(); ok {
		_spec.AddField(bartimerange.FieldCount, field.TypeInt, value)
	}
	if value, ok := btruo.mutation.Status(); ok {
		_spec.SetField(bartimerange.FieldStatus, field.TypeEnum, value)
	}
	if value, ok := btruo.mutation.UpdateTime(); ok {
		_spec.SetField(bartimerange.FieldUpdateTime, field.TypeTime, value)
	}
	if btruo.mutation.IntervalCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   bartimerange.IntervalTable,
			Columns: []string{bartimerange.IntervalColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(interval.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := btruo.mutation.IntervalIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   bartimerange.IntervalTable,
			Columns: []string{bartimerange.IntervalColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(interval.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if btruo.mutation.GroupsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   bartimerange.GroupsTable,
			Columns: []string{bartimerange.GroupsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bargroup.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := btruo.mutation.RemovedGroupsIDs(); len(nodes) > 0 && !btruo.mutation.GroupsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   bartimerange.GroupsTable,
			Columns: []string{bartimerange.GroupsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bargroup.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := btruo.mutation.GroupsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   bartimerange.GroupsTable,
			Columns: []string{bartimerange.GroupsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bargroup.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &BarTimeRange{config: btruo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, btruo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{bartimerange.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	btruo.mutation.done = true
	return _node, nil
}
