package errors

import stderrors "errors"

var (
	ErrTaskNotFound  = stderrors.New("task not found")
	ErrInvalidStatus = stderrors.New("status must be one of: pending in-progress done")
)
