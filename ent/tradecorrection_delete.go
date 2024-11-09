// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/nathanhack/sibyl/ent/predicate"
	"github.com/nathanhack/sibyl/ent/tradecorrection"
)

// TradeCorrectionDelete is the builder for deleting a TradeCorrection entity.
type TradeCorrectionDelete struct {
	config
	hooks    []Hook
	mutation *TradeCorrectionMutation
}

// Where appends a list predicates to the TradeCorrectionDelete builder.
func (tcd *TradeCorrectionDelete) Where(ps ...predicate.TradeCorrection) *TradeCorrectionDelete {
	tcd.mutation.Where(ps...)
	return tcd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (tcd *TradeCorrectionDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, tcd.sqlExec, tcd.mutation, tcd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (tcd *TradeCorrectionDelete) ExecX(ctx context.Context) int {
	n, err := tcd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (tcd *TradeCorrectionDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(tradecorrection.Table, sqlgraph.NewFieldSpec(tradecorrection.FieldID, field.TypeInt))
	if ps := tcd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, tcd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	tcd.mutation.done = true
	return affected, err
}

// TradeCorrectionDeleteOne is the builder for deleting a single TradeCorrection entity.
type TradeCorrectionDeleteOne struct {
	tcd *TradeCorrectionDelete
}

// Where appends a list predicates to the TradeCorrectionDelete builder.
func (tcdo *TradeCorrectionDeleteOne) Where(ps ...predicate.TradeCorrection) *TradeCorrectionDeleteOne {
	tcdo.tcd.mutation.Where(ps...)
	return tcdo
}

// Exec executes the deletion query.
func (tcdo *TradeCorrectionDeleteOne) Exec(ctx context.Context) error {
	n, err := tcdo.tcd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{tradecorrection.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (tcdo *TradeCorrectionDeleteOne) ExecX(ctx context.Context) {
	if err := tcdo.Exec(ctx); err != nil {
		panic(err)
	}
}
