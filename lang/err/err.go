package err

import "github.com/muzin/go_rt/try"

// errors
var OutOfMemoryError = try.DeclareException("OutOfMemoryError")

// exception
var ArrayIndexOutOfBoundsException = try.DeclareException("ArrayIndexOutOfBoundsException")

var NoSuchElementException = try.DeclareException("NoSuchElementException")
