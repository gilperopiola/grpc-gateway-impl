package http_tests

import "fmt"

var _ Defaulter[string] = (*assertDefaulter[string])(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - Value Defaulter -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// -> Haha I was so high when I wrote this I don't even know if it makes sense.
// Let's see:
func doesThisMakeSense() {
	var pageSizeAny = any(100)
	var pageSizeAny2 = any("100")

	pageSize := AssertToType[int](pageSizeAny).OrDefaultTo(50)
	pageSize2 := AssertToType[int](pageSizeAny2).OrDefaultTo(50)
	fmt.Println(pageSize, pageSize2)
}

// Chain this with .OrDefaultTo to get a type-assertion-or-default-value logic in 1 line
func AssertToType[T any](value any) Defaulter[T] {
	return &assertDefaulter[T]{value}
}

// Example -> pageSize := assertTo[int](pageSizeAny).OrDefaultTo(100)
type Defaulter[T any] interface {
	OrDefaultTo(defaultValue T) T
}

// Holds a value that needs assertion to type T
type assertDefaulter[T any] struct {
	value any
}

// Assert the value to type T.
// If it works, return it. If not, return the default
func (d *assertDefaulter[T]) OrDefaultTo(defaultValue T) T {
	if value, ok := d.value.(T); ok {
		return value
	}
	return defaultValue
}
