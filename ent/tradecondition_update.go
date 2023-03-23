// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/predicate"
	"github.com/nathanhack/sibyl/ent/tradecondition"
	"github.com/nathanhack/sibyl/ent/traderecord"
)

// TradeConditionUpdate is the builder for updating TradeCondition entities.
type TradeConditionUpdate struct {
	config
	hooks    []Hook
	mutation *TradeConditionMutation
}

// Where appends a list predicates to the TradeConditionUpdate builder.
func (tcu *TradeConditionUpdate) Where(ps ...predicate.TradeCondition) *TradeConditionUpdate {
	tcu.mutation.Where(ps...)
	return tcu
}

// SetCondition sets the "condition" field.
func (tcu *TradeConditionUpdate) SetCondition(s string) *TradeConditionUpdate {
	tcu.mutation.SetCondition(s)
	return tcu
}

// AddRecordIDs adds the "record" edge to the TradeRecord entity by IDs.
func (tcu *TradeConditionUpdate) AddRecordIDs(ids ...int) *TradeConditionUpdate {
	tcu.mutation.AddRecordIDs(ids...)
	return tcu
}

// AddRecord adds the "record" edges to the TradeRecord entity.
func (tcu *TradeConditionUpdate) AddRecord(t ...*TradeRecord) *TradeConditionUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tcu.AddRecordIDs(ids...)
}

// Mutation returns the TradeConditionMutation object of the builder.
func (tcu *TradeConditionUpdate) Mutation() *TradeConditionMutation {
	return tcu.mutation
}

// ClearRecord clears all "record" edges to the TradeRecord entity.
func (tcu *TradeConditionUpdate) ClearRecord() *TradeConditionUpdate {
	tcu.mutation.ClearRecord()
	return tcu
}

// RemoveRecordIDs removes the "record" edge to TradeRecord entities by IDs.
func (tcu *TradeConditionUpdate) RemoveRecordIDs(ids ...int) *TradeConditionUpdate {
	tcu.mutation.RemoveRecordIDs(ids...)
	return tcu
}

// RemoveRecord removes "record" edges to TradeRecord entities.
func (tcu *TradeConditionUpdate) RemoveRecord(t ...*TradeRecord) *TradeConditionUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tcu.RemoveRecordIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tcu *TradeConditionUpdate) Save(ctx context.Context) (int, error) {
	return withHooks[int, TradeConditionMutation](ctx, tcu.sqlSave, tcu.mutation, tcu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tcu *TradeConditionUpdate) SaveX(ctx context.Context) int {
	affected, err := tcu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tcu *TradeConditionUpdate) Exec(ctx context.Context) error {
	_, err := tcu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcu *TradeConditionUpdate) ExecX(ctx context.Context) {
	if err := tcu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tcu *TradeConditionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   tradecondition.Table,
			Columns: tradecondition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: tradecondition.FieldID,
			},
		},
	}
	if ps := tcu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tcu.mutation.Condition(); ok {
		_spec.SetField(tradecondition.FieldCondition, field.TypeString, value)
	}
	if tcu.mutation.RecordCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   tradecondition.RecordTable,
			Columns: tradecondition.RecordPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: traderecord.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tcu.mutation.RemovedRecordIDs(); len(nodes) > 0 && !tcu.mutation.RecordCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   tradecondition.RecordTable,
			Columns: tradecondition.RecordPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: traderecord.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tcu.mutation.RecordIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   tradecondition.RecordTable,
			Columns: tradecondition.RecordPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: traderecord.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tcu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{tradecondition.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	tcu.mutation.done = true
	return n, nil
}

// TradeConditionUpdateOne is the builder for updating a single TradeCondition entity.
type TradeConditionUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *TradeConditionMutation
}

// SetCondition sets the "condition" field.
func (tcuo *TradeConditionUpdateOne) SetCondition(s string) *TradeConditionUpdateOne {
	tcuo.mutation.SetCondition(s)
	return tcuo
}

// AddRecordIDs adds the "record" edge to the TradeRecord entity by IDs.
func (tcuo *TradeConditionUpdateOne) AddRecordIDs(ids ...int) *TradeConditionUpdateOne {
	tcuo.mutation.AddRecordIDs(ids...)
	return tcuo
}

// AddRecord adds the "record" edges to the TradeRecord entity.
func (tcuo *TradeConditionUpdateOne) AddRecord(t ...*TradeRecord) *TradeConditionUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tcuo.AddRecordIDs(ids...)
}

// Mutation returns the TradeConditionMutation object of the builder.
func (tcuo *TradeConditionUpdateOne) Mutation() *TradeConditionMutation {
	return tcuo.mutation
}

// ClearRecord clears all "record" edges to the TradeRecord entity.
func (tcuo *TradeConditionUpdateOne) ClearRecord() *TradeConditionUpdateOne {
	tcuo.mutation.ClearRecord()
	return tcuo
}

// RemoveRecordIDs removes the "record" edge to TradeRecord entities by IDs.
func (tcuo *TradeConditionUpdateOne) RemoveRecordIDs(ids ...int) *TradeConditionUpdateOne {
	tcuo.mutation.RemoveRecordIDs(ids...)
	return tcuo
}

// RemoveRecord removes "record" edges to TradeRecord entities.
func (tcuo *TradeConditionUpdateOne) RemoveRecord(t ...*TradeRecord) *TradeConditionUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tcuo.RemoveRecordIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tcuo *TradeConditionUpdateOne) Select(field string, fields ...string) *TradeConditionUpdateOne {
	tcuo.fields = append([]string{field}, fields...)
	return tcuo
}

// Save executes the query and returns the updated TradeCondition entity.
func (tcuo *TradeConditionUpdateOne) Save(ctx context.Context) (*TradeCondition, error) {
	return withHooks[*TradeCondition, TradeConditionMutation](ctx, tcuo.sqlSave, tcuo.mutation, tcuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tcuo *TradeConditionUpdateOne) SaveX(ctx context.Context) *TradeCondition {
	node, err := tcuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tcuo *TradeConditionUpdateOne) Exec(ctx context.Context) error {
	_, err := tcuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcuo *TradeConditionUpdateOne) ExecX(ctx context.Context) {
	if err := tcuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tcuo *TradeConditionUpdateOne) sqlSave(ctx context.Context) (_node *TradeCondition, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   tradecondition.Table,
			Columns: tradecondition.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: tradecondition.FieldID,
			},
		},
	}
	id, ok := tcuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "TradeCondition.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := tcuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, tradecondition.FieldID)
		for _, f := range fields {
			if !tradecondition.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != tradecondition.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tcuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tcuo.mutation.Condition(); ok {
		_spec.SetField(tradecondition.FieldCondition, field.TypeString, value)
	}
	if tcuo.mutation.RecordCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   tradecondition.RecordTable,
			Columns: tradecondition.RecordPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: traderecord.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tcuo.mutation.RemovedRecordIDs(); len(nodes) > 0 && !tcuo.mutation.RecordCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   tradecondition.RecordTable,
			Columns: tradecondition.RecordPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: traderecord.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tcuo.mutation.RecordIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   tradecondition.RecordTable,
			Columns: tradecondition.RecordPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: traderecord.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &TradeCondition{config: tcuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tcuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{tradecondition.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	tcuo.mutation.done = true
	return _node, nil
}