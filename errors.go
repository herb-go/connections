package connections

import "errors"

//ErrConnIDDuplicated errors reaised when generated connection id duplicated
var ErrConnIDDuplicated = errors.New("generated connection id duplicated")
