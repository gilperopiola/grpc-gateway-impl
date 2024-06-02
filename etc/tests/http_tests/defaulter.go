package http_tests

var _ Defaulter[string] = (*assertDefaulter[string])(nil)

// Chain this with .OrDefaultTo to get a type-assertion-or-default in 1 line
func assertTo[T any](value any) Defaulter[T] {
	return &assertDefaulter[T]{value}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - Value Defaulter -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

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
