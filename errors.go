package rhombifer

import "github.com/chapgx/rhombifer/pkg/errs"

var (
	ErroNoRootRunFunc = errs.NewError("root has no Run function. Therefore it can't interpret any flags.")

	ErroFlagUndefined   = errs.NewError("flag is undefined.")
	ErroFlagHasNoValues = errs.NewError("flag has no values.")
	ErroFoundFlagsIsNil = errs.NewError("found flags is <nil>.")
	ErroFlagNotFound    = errs.NewError("flag not found.")

	ErrNoSubCommandPassed = errs.NewError("no subcommand was passed")
)
