// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/NpoolPlatform/staker-manager/pkg/db/ent/predicate"
	"github.com/NpoolPlatform/staker-manager/pkg/db/ent/pubsubmessage"
	"github.com/google/uuid"
)

// PubsubMessageUpdate is the builder for updating PubsubMessage entities.
type PubsubMessageUpdate struct {
	config
	hooks     []Hook
	mutation  *PubsubMessageMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the PubsubMessageUpdate builder.
func (pmu *PubsubMessageUpdate) Where(ps ...predicate.PubsubMessage) *PubsubMessageUpdate {
	pmu.mutation.Where(ps...)
	return pmu
}

// SetCreatedAt sets the "created_at" field.
func (pmu *PubsubMessageUpdate) SetCreatedAt(u uint32) *PubsubMessageUpdate {
	pmu.mutation.ResetCreatedAt()
	pmu.mutation.SetCreatedAt(u)
	return pmu
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (pmu *PubsubMessageUpdate) SetNillableCreatedAt(u *uint32) *PubsubMessageUpdate {
	if u != nil {
		pmu.SetCreatedAt(*u)
	}
	return pmu
}

// AddCreatedAt adds u to the "created_at" field.
func (pmu *PubsubMessageUpdate) AddCreatedAt(u int32) *PubsubMessageUpdate {
	pmu.mutation.AddCreatedAt(u)
	return pmu
}

// SetUpdatedAt sets the "updated_at" field.
func (pmu *PubsubMessageUpdate) SetUpdatedAt(u uint32) *PubsubMessageUpdate {
	pmu.mutation.ResetUpdatedAt()
	pmu.mutation.SetUpdatedAt(u)
	return pmu
}

// AddUpdatedAt adds u to the "updated_at" field.
func (pmu *PubsubMessageUpdate) AddUpdatedAt(u int32) *PubsubMessageUpdate {
	pmu.mutation.AddUpdatedAt(u)
	return pmu
}

// SetDeletedAt sets the "deleted_at" field.
func (pmu *PubsubMessageUpdate) SetDeletedAt(u uint32) *PubsubMessageUpdate {
	pmu.mutation.ResetDeletedAt()
	pmu.mutation.SetDeletedAt(u)
	return pmu
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (pmu *PubsubMessageUpdate) SetNillableDeletedAt(u *uint32) *PubsubMessageUpdate {
	if u != nil {
		pmu.SetDeletedAt(*u)
	}
	return pmu
}

// AddDeletedAt adds u to the "deleted_at" field.
func (pmu *PubsubMessageUpdate) AddDeletedAt(u int32) *PubsubMessageUpdate {
	pmu.mutation.AddDeletedAt(u)
	return pmu
}

// SetMessageID sets the "message_id" field.
func (pmu *PubsubMessageUpdate) SetMessageID(s string) *PubsubMessageUpdate {
	pmu.mutation.SetMessageID(s)
	return pmu
}

// SetNillableMessageID sets the "message_id" field if the given value is not nil.
func (pmu *PubsubMessageUpdate) SetNillableMessageID(s *string) *PubsubMessageUpdate {
	if s != nil {
		pmu.SetMessageID(*s)
	}
	return pmu
}

// ClearMessageID clears the value of the "message_id" field.
func (pmu *PubsubMessageUpdate) ClearMessageID() *PubsubMessageUpdate {
	pmu.mutation.ClearMessageID()
	return pmu
}

// SetState sets the "state" field.
func (pmu *PubsubMessageUpdate) SetState(s string) *PubsubMessageUpdate {
	pmu.mutation.SetState(s)
	return pmu
}

// SetNillableState sets the "state" field if the given value is not nil.
func (pmu *PubsubMessageUpdate) SetNillableState(s *string) *PubsubMessageUpdate {
	if s != nil {
		pmu.SetState(*s)
	}
	return pmu
}

// ClearState clears the value of the "state" field.
func (pmu *PubsubMessageUpdate) ClearState() *PubsubMessageUpdate {
	pmu.mutation.ClearState()
	return pmu
}

// SetRespToID sets the "resp_to_id" field.
func (pmu *PubsubMessageUpdate) SetRespToID(u uuid.UUID) *PubsubMessageUpdate {
	pmu.mutation.SetRespToID(u)
	return pmu
}

// SetNillableRespToID sets the "resp_to_id" field if the given value is not nil.
func (pmu *PubsubMessageUpdate) SetNillableRespToID(u *uuid.UUID) *PubsubMessageUpdate {
	if u != nil {
		pmu.SetRespToID(*u)
	}
	return pmu
}

// ClearRespToID clears the value of the "resp_to_id" field.
func (pmu *PubsubMessageUpdate) ClearRespToID() *PubsubMessageUpdate {
	pmu.mutation.ClearRespToID()
	return pmu
}

// SetUndoID sets the "undo_id" field.
func (pmu *PubsubMessageUpdate) SetUndoID(u uuid.UUID) *PubsubMessageUpdate {
	pmu.mutation.SetUndoID(u)
	return pmu
}

// SetNillableUndoID sets the "undo_id" field if the given value is not nil.
func (pmu *PubsubMessageUpdate) SetNillableUndoID(u *uuid.UUID) *PubsubMessageUpdate {
	if u != nil {
		pmu.SetUndoID(*u)
	}
	return pmu
}

// ClearUndoID clears the value of the "undo_id" field.
func (pmu *PubsubMessageUpdate) ClearUndoID() *PubsubMessageUpdate {
	pmu.mutation.ClearUndoID()
	return pmu
}

// SetArguments sets the "arguments" field.
func (pmu *PubsubMessageUpdate) SetArguments(s string) *PubsubMessageUpdate {
	pmu.mutation.SetArguments(s)
	return pmu
}

// SetNillableArguments sets the "arguments" field if the given value is not nil.
func (pmu *PubsubMessageUpdate) SetNillableArguments(s *string) *PubsubMessageUpdate {
	if s != nil {
		pmu.SetArguments(*s)
	}
	return pmu
}

// ClearArguments clears the value of the "arguments" field.
func (pmu *PubsubMessageUpdate) ClearArguments() *PubsubMessageUpdate {
	pmu.mutation.ClearArguments()
	return pmu
}

// Mutation returns the PubsubMessageMutation object of the builder.
func (pmu *PubsubMessageUpdate) Mutation() *PubsubMessageMutation {
	return pmu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pmu *PubsubMessageUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if err := pmu.defaults(); err != nil {
		return 0, err
	}
	if len(pmu.hooks) == 0 {
		affected, err = pmu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PubsubMessageMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			pmu.mutation = mutation
			affected, err = pmu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(pmu.hooks) - 1; i >= 0; i-- {
			if pmu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = pmu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, pmu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (pmu *PubsubMessageUpdate) SaveX(ctx context.Context) int {
	affected, err := pmu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pmu *PubsubMessageUpdate) Exec(ctx context.Context) error {
	_, err := pmu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pmu *PubsubMessageUpdate) ExecX(ctx context.Context) {
	if err := pmu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pmu *PubsubMessageUpdate) defaults() error {
	if _, ok := pmu.mutation.UpdatedAt(); !ok {
		if pubsubmessage.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized pubsubmessage.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := pubsubmessage.UpdateDefaultUpdatedAt()
		pmu.mutation.SetUpdatedAt(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (pmu *PubsubMessageUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PubsubMessageUpdate {
	pmu.modifiers = append(pmu.modifiers, modifiers...)
	return pmu
}

func (pmu *PubsubMessageUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   pubsubmessage.Table,
			Columns: pubsubmessage.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: pubsubmessage.FieldID,
			},
		},
	}
	if ps := pmu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pmu.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: pubsubmessage.FieldCreatedAt,
		})
	}
	if value, ok := pmu.mutation.AddedCreatedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: pubsubmessage.FieldCreatedAt,
		})
	}
	if value, ok := pmu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: pubsubmessage.FieldUpdatedAt,
		})
	}
	if value, ok := pmu.mutation.AddedUpdatedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: pubsubmessage.FieldUpdatedAt,
		})
	}
	if value, ok := pmu.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: pubsubmessage.FieldDeletedAt,
		})
	}
	if value, ok := pmu.mutation.AddedDeletedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: pubsubmessage.FieldDeletedAt,
		})
	}
	if value, ok := pmu.mutation.MessageID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: pubsubmessage.FieldMessageID,
		})
	}
	if pmu.mutation.MessageIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: pubsubmessage.FieldMessageID,
		})
	}
	if value, ok := pmu.mutation.State(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: pubsubmessage.FieldState,
		})
	}
	if pmu.mutation.StateCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: pubsubmessage.FieldState,
		})
	}
	if value, ok := pmu.mutation.RespToID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Value:  value,
			Column: pubsubmessage.FieldRespToID,
		})
	}
	if pmu.mutation.RespToIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Column: pubsubmessage.FieldRespToID,
		})
	}
	if value, ok := pmu.mutation.UndoID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Value:  value,
			Column: pubsubmessage.FieldUndoID,
		})
	}
	if pmu.mutation.UndoIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Column: pubsubmessage.FieldUndoID,
		})
	}
	if value, ok := pmu.mutation.Arguments(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: pubsubmessage.FieldArguments,
		})
	}
	if pmu.mutation.ArgumentsCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: pubsubmessage.FieldArguments,
		})
	}
	_spec.Modifiers = pmu.modifiers
	if n, err = sqlgraph.UpdateNodes(ctx, pmu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{pubsubmessage.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	return n, nil
}

// PubsubMessageUpdateOne is the builder for updating a single PubsubMessage entity.
type PubsubMessageUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *PubsubMessageMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetCreatedAt sets the "created_at" field.
func (pmuo *PubsubMessageUpdateOne) SetCreatedAt(u uint32) *PubsubMessageUpdateOne {
	pmuo.mutation.ResetCreatedAt()
	pmuo.mutation.SetCreatedAt(u)
	return pmuo
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (pmuo *PubsubMessageUpdateOne) SetNillableCreatedAt(u *uint32) *PubsubMessageUpdateOne {
	if u != nil {
		pmuo.SetCreatedAt(*u)
	}
	return pmuo
}

// AddCreatedAt adds u to the "created_at" field.
func (pmuo *PubsubMessageUpdateOne) AddCreatedAt(u int32) *PubsubMessageUpdateOne {
	pmuo.mutation.AddCreatedAt(u)
	return pmuo
}

// SetUpdatedAt sets the "updated_at" field.
func (pmuo *PubsubMessageUpdateOne) SetUpdatedAt(u uint32) *PubsubMessageUpdateOne {
	pmuo.mutation.ResetUpdatedAt()
	pmuo.mutation.SetUpdatedAt(u)
	return pmuo
}

// AddUpdatedAt adds u to the "updated_at" field.
func (pmuo *PubsubMessageUpdateOne) AddUpdatedAt(u int32) *PubsubMessageUpdateOne {
	pmuo.mutation.AddUpdatedAt(u)
	return pmuo
}

// SetDeletedAt sets the "deleted_at" field.
func (pmuo *PubsubMessageUpdateOne) SetDeletedAt(u uint32) *PubsubMessageUpdateOne {
	pmuo.mutation.ResetDeletedAt()
	pmuo.mutation.SetDeletedAt(u)
	return pmuo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (pmuo *PubsubMessageUpdateOne) SetNillableDeletedAt(u *uint32) *PubsubMessageUpdateOne {
	if u != nil {
		pmuo.SetDeletedAt(*u)
	}
	return pmuo
}

// AddDeletedAt adds u to the "deleted_at" field.
func (pmuo *PubsubMessageUpdateOne) AddDeletedAt(u int32) *PubsubMessageUpdateOne {
	pmuo.mutation.AddDeletedAt(u)
	return pmuo
}

// SetMessageID sets the "message_id" field.
func (pmuo *PubsubMessageUpdateOne) SetMessageID(s string) *PubsubMessageUpdateOne {
	pmuo.mutation.SetMessageID(s)
	return pmuo
}

// SetNillableMessageID sets the "message_id" field if the given value is not nil.
func (pmuo *PubsubMessageUpdateOne) SetNillableMessageID(s *string) *PubsubMessageUpdateOne {
	if s != nil {
		pmuo.SetMessageID(*s)
	}
	return pmuo
}

// ClearMessageID clears the value of the "message_id" field.
func (pmuo *PubsubMessageUpdateOne) ClearMessageID() *PubsubMessageUpdateOne {
	pmuo.mutation.ClearMessageID()
	return pmuo
}

// SetState sets the "state" field.
func (pmuo *PubsubMessageUpdateOne) SetState(s string) *PubsubMessageUpdateOne {
	pmuo.mutation.SetState(s)
	return pmuo
}

// SetNillableState sets the "state" field if the given value is not nil.
func (pmuo *PubsubMessageUpdateOne) SetNillableState(s *string) *PubsubMessageUpdateOne {
	if s != nil {
		pmuo.SetState(*s)
	}
	return pmuo
}

// ClearState clears the value of the "state" field.
func (pmuo *PubsubMessageUpdateOne) ClearState() *PubsubMessageUpdateOne {
	pmuo.mutation.ClearState()
	return pmuo
}

// SetRespToID sets the "resp_to_id" field.
func (pmuo *PubsubMessageUpdateOne) SetRespToID(u uuid.UUID) *PubsubMessageUpdateOne {
	pmuo.mutation.SetRespToID(u)
	return pmuo
}

// SetNillableRespToID sets the "resp_to_id" field if the given value is not nil.
func (pmuo *PubsubMessageUpdateOne) SetNillableRespToID(u *uuid.UUID) *PubsubMessageUpdateOne {
	if u != nil {
		pmuo.SetRespToID(*u)
	}
	return pmuo
}

// ClearRespToID clears the value of the "resp_to_id" field.
func (pmuo *PubsubMessageUpdateOne) ClearRespToID() *PubsubMessageUpdateOne {
	pmuo.mutation.ClearRespToID()
	return pmuo
}

// SetUndoID sets the "undo_id" field.
func (pmuo *PubsubMessageUpdateOne) SetUndoID(u uuid.UUID) *PubsubMessageUpdateOne {
	pmuo.mutation.SetUndoID(u)
	return pmuo
}

// SetNillableUndoID sets the "undo_id" field if the given value is not nil.
func (pmuo *PubsubMessageUpdateOne) SetNillableUndoID(u *uuid.UUID) *PubsubMessageUpdateOne {
	if u != nil {
		pmuo.SetUndoID(*u)
	}
	return pmuo
}

// ClearUndoID clears the value of the "undo_id" field.
func (pmuo *PubsubMessageUpdateOne) ClearUndoID() *PubsubMessageUpdateOne {
	pmuo.mutation.ClearUndoID()
	return pmuo
}

// SetArguments sets the "arguments" field.
func (pmuo *PubsubMessageUpdateOne) SetArguments(s string) *PubsubMessageUpdateOne {
	pmuo.mutation.SetArguments(s)
	return pmuo
}

// SetNillableArguments sets the "arguments" field if the given value is not nil.
func (pmuo *PubsubMessageUpdateOne) SetNillableArguments(s *string) *PubsubMessageUpdateOne {
	if s != nil {
		pmuo.SetArguments(*s)
	}
	return pmuo
}

// ClearArguments clears the value of the "arguments" field.
func (pmuo *PubsubMessageUpdateOne) ClearArguments() *PubsubMessageUpdateOne {
	pmuo.mutation.ClearArguments()
	return pmuo
}

// Mutation returns the PubsubMessageMutation object of the builder.
func (pmuo *PubsubMessageUpdateOne) Mutation() *PubsubMessageMutation {
	return pmuo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (pmuo *PubsubMessageUpdateOne) Select(field string, fields ...string) *PubsubMessageUpdateOne {
	pmuo.fields = append([]string{field}, fields...)
	return pmuo
}

// Save executes the query and returns the updated PubsubMessage entity.
func (pmuo *PubsubMessageUpdateOne) Save(ctx context.Context) (*PubsubMessage, error) {
	var (
		err  error
		node *PubsubMessage
	)
	if err := pmuo.defaults(); err != nil {
		return nil, err
	}
	if len(pmuo.hooks) == 0 {
		node, err = pmuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*PubsubMessageMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			pmuo.mutation = mutation
			node, err = pmuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(pmuo.hooks) - 1; i >= 0; i-- {
			if pmuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = pmuo.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, pmuo.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*PubsubMessage)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from PubsubMessageMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (pmuo *PubsubMessageUpdateOne) SaveX(ctx context.Context) *PubsubMessage {
	node, err := pmuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (pmuo *PubsubMessageUpdateOne) Exec(ctx context.Context) error {
	_, err := pmuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pmuo *PubsubMessageUpdateOne) ExecX(ctx context.Context) {
	if err := pmuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pmuo *PubsubMessageUpdateOne) defaults() error {
	if _, ok := pmuo.mutation.UpdatedAt(); !ok {
		if pubsubmessage.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized pubsubmessage.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := pubsubmessage.UpdateDefaultUpdatedAt()
		pmuo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (pmuo *PubsubMessageUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *PubsubMessageUpdateOne {
	pmuo.modifiers = append(pmuo.modifiers, modifiers...)
	return pmuo
}

func (pmuo *PubsubMessageUpdateOne) sqlSave(ctx context.Context) (_node *PubsubMessage, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   pubsubmessage.Table,
			Columns: pubsubmessage.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: pubsubmessage.FieldID,
			},
		},
	}
	id, ok := pmuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "PubsubMessage.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := pmuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, pubsubmessage.FieldID)
		for _, f := range fields {
			if !pubsubmessage.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != pubsubmessage.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := pmuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pmuo.mutation.CreatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: pubsubmessage.FieldCreatedAt,
		})
	}
	if value, ok := pmuo.mutation.AddedCreatedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: pubsubmessage.FieldCreatedAt,
		})
	}
	if value, ok := pmuo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: pubsubmessage.FieldUpdatedAt,
		})
	}
	if value, ok := pmuo.mutation.AddedUpdatedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: pubsubmessage.FieldUpdatedAt,
		})
	}
	if value, ok := pmuo.mutation.DeletedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: pubsubmessage.FieldDeletedAt,
		})
	}
	if value, ok := pmuo.mutation.AddedDeletedAt(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeUint32,
			Value:  value,
			Column: pubsubmessage.FieldDeletedAt,
		})
	}
	if value, ok := pmuo.mutation.MessageID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: pubsubmessage.FieldMessageID,
		})
	}
	if pmuo.mutation.MessageIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: pubsubmessage.FieldMessageID,
		})
	}
	if value, ok := pmuo.mutation.State(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: pubsubmessage.FieldState,
		})
	}
	if pmuo.mutation.StateCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: pubsubmessage.FieldState,
		})
	}
	if value, ok := pmuo.mutation.RespToID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Value:  value,
			Column: pubsubmessage.FieldRespToID,
		})
	}
	if pmuo.mutation.RespToIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Column: pubsubmessage.FieldRespToID,
		})
	}
	if value, ok := pmuo.mutation.UndoID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Value:  value,
			Column: pubsubmessage.FieldUndoID,
		})
	}
	if pmuo.mutation.UndoIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeUUID,
			Column: pubsubmessage.FieldUndoID,
		})
	}
	if value, ok := pmuo.mutation.Arguments(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: pubsubmessage.FieldArguments,
		})
	}
	if pmuo.mutation.ArgumentsCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: pubsubmessage.FieldArguments,
		})
	}
	_spec.Modifiers = pmuo.modifiers
	_node = &PubsubMessage{config: pmuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, pmuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{pubsubmessage.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	return _node, nil
}
