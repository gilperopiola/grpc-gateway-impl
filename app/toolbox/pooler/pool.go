package pooler

type PoolOf[T any] interface {
	TakeOne() (T, func())
	Len() int
}

type poolOf[T any] struct {
	units     []unit[T]
	newUnitFn func() T
}

type unit[T any] struct {
	value     T
	available bool
}

func NewPool[T any](size int, newUnitFn func() T) PoolOf[T] {
	units := make([]unit[T], size)

	for i := 0; i < size; i++ {
		units[i] = unit[T]{
			value:     newUnitFn(),
			available: true,
		}
	}

	return &poolOf[T]{units, newUnitFn}
}

func (p *poolOf[T]) TakeOne() (T, func()) {

	getReleaseFn := func(i int) func() {
		return func() {
			p.units[i] = unit[T]{
				value:     p.newUnitFn(),
				available: true,
			}
		}
	}

	for i, unit := range p.units {
		if unit.available {
			p.units[i].available = false
			return unit.value, getReleaseFn(i)
		}
	}

	// We only reach this point if there are not available units. Create a new one:

	newUnit := unit[T]{
		value:     p.newUnitFn(),
		available: false,
	}

	p.units = append(p.units, newUnit)

	return newUnit.value, getReleaseFn(len(p.units) - 1)
}

func (p *poolOf[T]) Len() int {
	return len(p.units)
}
