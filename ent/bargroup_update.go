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
	"github.com/nathanhack/sibyl/ent/barrecord"
	"github.com/nathanhack/sibyl/ent/bartimerange"
	"github.com/nathanhack/sibyl/ent/predicate"
)

// BarGroupUpdate is the builder for updating BarGroup entities.
type BarGroupUpdate struct {
	config
	hooks    []Hook
	mutation *BarGroupMutation
}

// Where appends a list predicates to the BarGroupUpdate builder.
func (bgu *BarGroupUpdate) Where(ps ...predicate.BarGroup) *BarGroupUpdate {
	bgu.mutation.Where(ps...)
	return bgu
}

// SetFirst sets the "first" field.
func (bgu *BarGroupUpdate) SetFirst(t time.Time) *BarGroupUpdate {
	bgu.mutation.SetFirst(t)
	return bgu
}

// SetNillableFirst sets the "first" field if the given value is not nil.
func (bgu *BarGroupUpdate) SetNillableFirst(t *time.Time) *BarGroupUpdate {
	if t != nil {
		bgu.SetFirst(*t)
	}
	return bgu
}

// SetLast sets the "last" field.
func (bgu *BarGroupUpdate) SetLast(t time.Time) *BarGroupUpdate {
	bgu.mutation.SetLast(t)
	return bgu
}

// SetNillableLast sets the "last" field if the given value is not nil.
func (bgu *BarGroupUpdate) SetNillableLast(t *time.Time) *BarGroupUpdate {
	if t != nil {
		bgu.SetLast(*t)
	}
	return bgu
}

// SetCount sets the "count" field.
func (bgu *BarGroupUpdate) SetCount(i int) *BarGroupUpdate {
	bgu.mutation.ResetCount()
	bgu.mutation.SetCount(i)
	return bgu
}

// SetNillableCount sets the "count" field if the given value is not nil.
func (bgu *BarGroupUpdate) SetNillableCount(i *int) *BarGroupUpdate {
	if i != nil {
		bgu.SetCount(*i)
	}
	return bgu
}

// AddCount adds i to the "count" field.
func (bgu *BarGroupUpdate) AddCount(i int) *BarGroupUpdate {
	bgu.mutation.AddCount(i)
	return bgu
}

// SetTimeRangeID sets the "time_range_id" field.
func (bgu *BarGroupUpdate) SetTimeRangeID(i int) *BarGroupUpdate {
	bgu.mutation.SetTimeRangeID(i)
	return bgu
}

// SetNillableTimeRangeID sets the "time_range_id" field if the given value is not nil.
func (bgu *BarGroupUpdate) SetNillableTimeRangeID(i *int) *BarGroupUpdate {
	if i != nil {
		bgu.SetTimeRangeID(*i)
	}
	return bgu
}

// SetTimeRange sets the "time_range" edge to the BarTimeRange entity.
func (bgu *BarGroupUpdate) SetTimeRange(b *BarTimeRange) *BarGroupUpdate {
	return bgu.SetTimeRangeID(b.ID)
}

// AddRecordIDs adds the "records" edge to the BarRecord entity by IDs.
func (bgu *BarGroupUpdate) AddRecordIDs(ids ...int) *BarGroupUpdate {
	bgu.mutation.AddRecordIDs(ids...)
	return bgu
}

// AddRecords adds the "records" edges to the BarRecord entity.
func (bgu *BarGroupUpdate) AddRecords(b ...*BarRecord) *BarGroupUpdate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return bgu.AddRecordIDs(ids...)
}

// Mutation returns the BarGroupMutation object of the builder.
func (bgu *BarGroupUpdate) Mutation() *BarGroupMutation {
	return bgu.mutation
}

// ClearTimeRange clears the "time_range" edge to the BarTimeRange entity.
func (bgu *BarGroupUpdate) ClearTimeRange() *BarGroupUpdate {
	bgu.mutation.ClearTimeRange()
	return bgu
}

// ClearRecords clears all "records" edges to the BarRecord entity.
func (bgu *BarGroupUpdate) ClearRecords() *BarGroupUpdate {
	bgu.mutation.ClearRecords()
	return bgu
}

// RemoveRecordIDs removes the "records" edge to BarRecord entities by IDs.
func (bgu *BarGroupUpdate) RemoveRecordIDs(ids ...int) *BarGroupUpdate {
	bgu.mutation.RemoveRecordIDs(ids...)
	return bgu
}

// RemoveRecords removes "records" edges to BarRecord entities.
func (bgu *BarGroupUpdate) RemoveRecords(b ...*BarRecord) *BarGroupUpdate {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return bgu.RemoveRecordIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (bgu *BarGroupUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, bgu.sqlSave, bgu.mutation, bgu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (bgu *BarGroupUpdate) SaveX(ctx context.Context) int {
	affected, err := bgu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (bgu *BarGroupUpdate) Exec(ctx context.Context) error {
	_, err := bgu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bgu *BarGroupUpdate) ExecX(ctx context.Context) {
	if err := bgu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (bgu *BarGroupUpdate) check() error {
	if bgu.mutation.TimeRangeCleared() && len(bgu.mutation.TimeRangeIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "BarGroup.time_range"`)
	}
	return nil
}

func (bgu *BarGroupUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := bgu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(bargroup.Table, bargroup.Columns, sqlgraph.NewFieldSpec(bargroup.FieldID, field.TypeInt))
	if ps := bgu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := bgu.mutation.First(); ok {
		_spec.SetField(bargroup.FieldFirst, field.TypeTime, value)
	}
	if value, ok := bgu.mutation.Last(); ok {
		_spec.SetField(bargroup.FieldLast, field.TypeTime, value)
	}
	if value, ok := bgu.mutation.Count(); ok {
		_spec.SetField(bargroup.FieldCount, field.TypeInt, value)
	}
	if value, ok := bgu.mutation.AddedCount(); ok {
		_spec.AddField(bargroup.FieldCount, field.TypeInt, value)
	}
	if bgu.mutation.TimeRangeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   bargroup.TimeRangeTable,
			Columns: []string{bargroup.TimeRangeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bartimerange.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bgu.mutation.TimeRangeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   bargroup.TimeRangeTable,
			Columns: []string{bargroup.TimeRangeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bartimerange.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if bgu.mutation.RecordsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   bargroup.RecordsTable,
			Columns: []string{bargroup.RecordsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(barrecord.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bgu.mutation.RemovedRecordsIDs(); len(nodes) > 0 && !bgu.mutation.RecordsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   bargroup.RecordsTable,
			Columns: []string{bargroup.RecordsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(barrecord.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bgu.mutation.RecordsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   bargroup.RecordsTable,
			Columns: []string{bargroup.RecordsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(barrecord.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, bgu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{bargroup.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	bgu.mutation.done = true
	return n, nil
}

// BarGroupUpdateOne is the builder for updating a single BarGroup entity.
type BarGroupUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *BarGroupMutation
}

// SetFirst sets the "first" field.
func (bguo *BarGroupUpdateOne) SetFirst(t time.Time) *BarGroupUpdateOne {
	bguo.mutation.SetFirst(t)
	return bguo
}

// SetNillableFirst sets the "first" field if the given value is not nil.
func (bguo *BarGroupUpdateOne) SetNillableFirst(t *time.Time) *BarGroupUpdateOne {
	if t != nil {
		bguo.SetFirst(*t)
	}
	return bguo
}

// SetLast sets the "last" field.
func (bguo *BarGroupUpdateOne) SetLast(t time.Time) *BarGroupUpdateOne {
	bguo.mutation.SetLast(t)
	return bguo
}

// SetNillableLast sets the "last" field if the given value is not nil.
func (bguo *BarGroupUpdateOne) SetNillableLast(t *time.Time) *BarGroupUpdateOne {
	if t != nil {
		bguo.SetLast(*t)
	}
	return bguo
}

// SetCount sets the "count" field.
func (bguo *BarGroupUpdateOne) SetCount(i int) *BarGroupUpdateOne {
	bguo.mutation.ResetCount()
	bguo.mutation.SetCount(i)
	return bguo
}

// SetNillableCount sets the "count" field if the given value is not nil.
func (bguo *BarGroupUpdateOne) SetNillableCount(i *int) *BarGroupUpdateOne {
	if i != nil {
		bguo.SetCount(*i)
	}
	return bguo
}

// AddCount adds i to the "count" field.
func (bguo *BarGroupUpdateOne) AddCount(i int) *BarGroupUpdateOne {
	bguo.mutation.AddCount(i)
	return bguo
}

// SetTimeRangeID sets the "time_range_id" field.
func (bguo *BarGroupUpdateOne) SetTimeRangeID(i int) *BarGroupUpdateOne {
	bguo.mutation.SetTimeRangeID(i)
	return bguo
}

// SetNillableTimeRangeID sets the "time_range_id" field if the given value is not nil.
func (bguo *BarGroupUpdateOne) SetNillableTimeRangeID(i *int) *BarGroupUpdateOne {
	if i != nil {
		bguo.SetTimeRangeID(*i)
	}
	return bguo
}

// SetTimeRange sets the "time_range" edge to the BarTimeRange entity.
func (bguo *BarGroupUpdateOne) SetTimeRange(b *BarTimeRange) *BarGroupUpdateOne {
	return bguo.SetTimeRangeID(b.ID)
}

// AddRecordIDs adds the "records" edge to the BarRecord entity by IDs.
func (bguo *BarGroupUpdateOne) AddRecordIDs(ids ...int) *BarGroupUpdateOne {
	bguo.mutation.AddRecordIDs(ids...)
	return bguo
}

// AddRecords adds the "records" edges to the BarRecord entity.
func (bguo *BarGroupUpdateOne) AddRecords(b ...*BarRecord) *BarGroupUpdateOne {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return bguo.AddRecordIDs(ids...)
}

// Mutation returns the BarGroupMutation object of the builder.
func (bguo *BarGroupUpdateOne) Mutation() *BarGroupMutation {
	return bguo.mutation
}

// ClearTimeRange clears the "time_range" edge to the BarTimeRange entity.
func (bguo *BarGroupUpdateOne) ClearTimeRange() *BarGroupUpdateOne {
	bguo.mutation.ClearTimeRange()
	return bguo
}

// ClearRecords clears all "records" edges to the BarRecord entity.
func (bguo *BarGroupUpdateOne) ClearRecords() *BarGroupUpdateOne {
	bguo.mutation.ClearRecords()
	return bguo
}

// RemoveRecordIDs removes the "records" edge to BarRecord entities by IDs.
func (bguo *BarGroupUpdateOne) RemoveRecordIDs(ids ...int) *BarGroupUpdateOne {
	bguo.mutation.RemoveRecordIDs(ids...)
	return bguo
}

// RemoveRecords removes "records" edges to BarRecord entities.
func (bguo *BarGroupUpdateOne) RemoveRecords(b ...*BarRecord) *BarGroupUpdateOne {
	ids := make([]int, len(b))
	for i := range b {
		ids[i] = b[i].ID
	}
	return bguo.RemoveRecordIDs(ids...)
}

// Where appends a list predicates to the BarGroupUpdate builder.
func (bguo *BarGroupUpdateOne) Where(ps ...predicate.BarGroup) *BarGroupUpdateOne {
	bguo.mutation.Where(ps...)
	return bguo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (bguo *BarGroupUpdateOne) Select(field string, fields ...string) *BarGroupUpdateOne {
	bguo.fields = append([]string{field}, fields...)
	return bguo
}

// Save executes the query and returns the updated BarGroup entity.
func (bguo *BarGroupUpdateOne) Save(ctx context.Context) (*BarGroup, error) {
	return withHooks(ctx, bguo.sqlSave, bguo.mutation, bguo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (bguo *BarGroupUpdateOne) SaveX(ctx context.Context) *BarGroup {
	node, err := bguo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (bguo *BarGroupUpdateOne) Exec(ctx context.Context) error {
	_, err := bguo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bguo *BarGroupUpdateOne) ExecX(ctx context.Context) {
	if err := bguo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (bguo *BarGroupUpdateOne) check() error {
	if bguo.mutation.TimeRangeCleared() && len(bguo.mutation.TimeRangeIDs()) > 0 {
		return errors.New(`ent: clearing a required unique edge "BarGroup.time_range"`)
	}
	return nil
}

func (bguo *BarGroupUpdateOne) sqlSave(ctx context.Context) (_node *BarGroup, err error) {
	if err := bguo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(bargroup.Table, bargroup.Columns, sqlgraph.NewFieldSpec(bargroup.FieldID, field.TypeInt))
	id, ok := bguo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "BarGroup.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := bguo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, bargroup.FieldID)
		for _, f := range fields {
			if !bargroup.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != bargroup.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := bguo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := bguo.mutation.First(); ok {
		_spec.SetField(bargroup.FieldFirst, field.TypeTime, value)
	}
	if value, ok := bguo.mutation.Last(); ok {
		_spec.SetField(bargroup.FieldLast, field.TypeTime, value)
	}
	if value, ok := bguo.mutation.Count(); ok {
		_spec.SetField(bargroup.FieldCount, field.TypeInt, value)
	}
	if value, ok := bguo.mutation.AddedCount(); ok {
		_spec.AddField(bargroup.FieldCount, field.TypeInt, value)
	}
	if bguo.mutation.TimeRangeCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   bargroup.TimeRangeTable,
			Columns: []string{bargroup.TimeRangeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bartimerange.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bguo.mutation.TimeRangeIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   bargroup.TimeRangeTable,
			Columns: []string{bargroup.TimeRangeColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(bartimerange.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if bguo.mutation.RecordsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   bargroup.RecordsTable,
			Columns: []string{bargroup.RecordsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(barrecord.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bguo.mutation.RemovedRecordsIDs(); len(nodes) > 0 && !bguo.mutation.RecordsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   bargroup.RecordsTable,
			Columns: []string{bargroup.RecordsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(barrecord.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bguo.mutation.RecordsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   bargroup.RecordsTable,
			Columns: []string{bargroup.RecordsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(barrecord.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &BarGroup{config: bguo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, bguo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{bargroup.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	bguo.mutation.done = true
	return _node, nil
}
