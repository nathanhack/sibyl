// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/markethours"
	"github.com/nathanhack/sibyl/ent/predicate"
)

// MarketHoursDelete is the builder for deleting a MarketHours entity.
type MarketHoursDelete struct {
	config
	hooks    []Hook
	mutation *MarketHoursMutation
}

// Where appends a list predicates to the MarketHoursDelete builder.
func (mhd *MarketHoursDelete) Where(ps ...predicate.MarketHours) *MarketHoursDelete {
	mhd.mutation.Where(ps...)
	return mhd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (mhd *MarketHoursDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, MarketHoursMutation](ctx, mhd.sqlExec, mhd.mutation, mhd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (mhd *MarketHoursDelete) ExecX(ctx context.Context) int {
	n, err := mhd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (mhd *MarketHoursDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: markethours.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: markethours.FieldID,
			},
		},
	}
	if ps := mhd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, mhd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	mhd.mutation.done = true
	return affected, err
}

// MarketHoursDeleteOne is the builder for deleting a single MarketHours entity.
type MarketHoursDeleteOne struct {
	mhd *MarketHoursDelete
}

// Exec executes the deletion query.
func (mhdo *MarketHoursDeleteOne) Exec(ctx context.Context) error {
	n, err := mhdo.mhd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{markethours.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (mhdo *MarketHoursDeleteOne) ExecX(ctx context.Context) {
	mhdo.mhd.ExecX(ctx)
}
