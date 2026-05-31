package service

import "errors"

// ErrValidation is the category for all input-validation failures across the
// service package. Specific sentinels wrap it (via %w) so the handler maps any
// of them to connect.CodeInvalidArgument with a single table entry.
var ErrValidation = errors.New("validation error")
