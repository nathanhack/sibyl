// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/barrecord"
	"github.com/nathanhack/sibyl/ent/predicate"
)

// BarRecordDelete is the builder for deleting a BarRecord entity.
type BarRecordDelete struct {
	config
	hooks    []Hook
	mutation *BarRecordMutation
}

// Where appends a list predicates to the BarRecordDelete builder.
func (brd *BarRecordDelete) Where(ps ...predicate.BarRecord) *BarRecordDelete {
	brd.mutation.Where(ps...)
	return brd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (brd *BarRecordDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, BarRecordMutation](ctx, brd.sqlExec, brd.mutation, brd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (brd *BarRecordDelete) ExecX(ctx context.Context) int {
	n, err := brd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (brd *BarRecordDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: barrecord.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: barrecord.FieldID,
			},
		},
	}
	if ps := brd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, brd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	brd.mutation.done = true
	return affected, err
}

// BarRecordDeleteOne is the builder for deleting a single BarRecord entity.
type BarRecordDeleteOne struct {
	brd *BarRecordDelete
}

// Exec executes the deletion query.
func (brdo *BarRecordDeleteOne) Exec(ctx context.Context) error {
	n, err := brdo.brd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{barrecord.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (brdo *BarRecordDeleteOne) ExecX(ctx context.Context) {
	brdo.brd.ExecX(ctx)
}