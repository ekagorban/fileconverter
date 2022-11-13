package converterror

import "errors"

var (
	ErrInvalidArgsCount    = errors.New("invalid args count")
	ErrEqualExtentions     = errors.New("src and dst extentions must be different")
	ErrNotImplementedRule  = errors.New("not implemented rule")
	ErrNotAllowedExtention = errors.New("not allowed extention")
)
