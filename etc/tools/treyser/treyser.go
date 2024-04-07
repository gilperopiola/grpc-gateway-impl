package treyser

import (
	"fmt"
	"time"
)

type Treyser struct {
	name  string
	level uint
	start time.Time
}

func NewTreyser(name string, treysLevel uint) *Treyser {
	return &Treyser{
		name:  name,
		level: treysLevel,
		start: time.Now(),
	}
}

func (t *Treyser) Elapsed() time.Duration {
	return time.Since(t.start)
}

func (t *Treyser) Treys() {
	for i := 0; i < int(t.level); i++ {
		fmt.Print("\t")
	}
	fmt.Printf("%s took %dms\n", t.name, t.Elapsed().Milliseconds())
}
