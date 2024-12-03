package etc

import (
	"fmt"
	"strconv"
	"testing"
)

// Use these instead of a blank identifier to prevent the compiler from optimizing out the variable assignment.
var (
	dontOptimizeInt    int
	dontOptimizeString string
)

func BenchmarkConditions(b *testing.B) {
	const x = 10
	const z = "Hola"

	b.Run("If true", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if true {
				_ = 1
			}
		}
	})

	b.Run("If true or false", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if true || false {
				_ = 1
			}
		}
	})

	b.Run("If int == 10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if x == 10 {
				_ = 1
			}
		}
	})

	b.Run("If string == 'Hola' -> true", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if z == "Hola" {
				_ = 1
			}
		}
	})

	b.Run("If string == 'Chau' -> false", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if z == "Chau" {
				_ = 1
			}
		}
	})

	b.Run("Switch case", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			switch x {
			case 10:
				_ = 1
			}
		}
	})

	b.Run("Switch case default", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			switch x {
			case 9:
				_ = 1
			default:
				_ = 1
			}
		}
	})
}

func BenchmarkTypeAssertion(b *testing.B) {
	var x any = 10
	for i := 0; i < b.N; i++ { // ➤ 0.6 ns/op
		_ = x.(int)
	}
}

func BenchmarkIntToString(b *testing.B) {
	var x int = 9

	b.Run("Sprintf", func(b *testing.B) { // ➤ 115 ns/op
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%d", x)
		}
	})

	b.Run("Strconv", func(b *testing.B) { // ➤ 5 ns/op 🏆
		for i := 0; i < b.N; i++ {
			_ = strconv.Itoa(x)
		}
	})

	b.Run("Casting", func(b *testing.B) { // ➤ 7 ns/op
		for i := 0; i < b.N; i++ {
			_ = string(rune(x))
		}
	})

	b.Run("Sprintf_2", func(b *testing.B) { // ➤ 123 ns/op
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%v", x)
		}
	})

	b.Run("Sprintf_3", func(b *testing.B) { // ➤ 179 ns/op
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%q", x)
		}
	})
}

func BenchmarkStringFormation(b *testing.B) {
	const msg = "Hello World 4"

	b.Run("String Assignment", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dontOptimizeString = "Hello World 1"
		}
	})
	b.Run("Sprintf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dontOptimizeString = fmt.Sprintf("%s", "Hello World 2")
		}
	})
	b.Run("Sprintf Many", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dontOptimizeString = fmt.Sprintf("%s%s", "Hello ", "World 2")
		}
	})
	b.Run("String Concatenation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dontOptimizeString = "Hello" + " World 3"
		}
	})
	b.Run("String Concatenation Many", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dontOptimizeString = "He" + "ll" + "o " + "Wo" + "rl" + "d " + "3"
		}
	})
	b.Run("Const String Assignment", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			dontOptimizeString = msg
		}
	})
}

func BenchmarkFloatToString(b *testing.B) {
	var x float64 = 9.9

	b.Run("Sprintf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%f", x)
		}
	})
	b.Run("Strconv", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = strconv.FormatFloat(x, 'f', -1, 64)
		}
	})
	b.Run("Sprintf_2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%v", x)
		}
	})
	b.Run("Sprintf_3", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = fmt.Sprintf("%.2f", x)
		}
	})
}

type A struct {
	Str1    string
	Str2    string
	Int     int
	Int2    string
	Int3    string
	Int4    string
	Int5    string
	Int6    string
	Int7    string
	Int8    string
	Int9    string
	Int10   string
	String3 string
	String4 string
	String5 string
	String6 string
	String7 string
	String8 string
}

func (a *A) GetFieldsPtrReceiver() (string, string, int) {
	return a.Str1, a.Str2, a.Int
}

func (a A) GetFieldsValueReceiver() (string, string, int) {
	return a.Str1, a.Str2, a.Int
}

type B struct {
	Str1 string
	Str2 string
	Int  int
}

func (a *B) GetFieldsPtrReceiver() (string, string, int) {
	return a.Str1, a.Str2, a.Int
}

func (a B) GetFieldsValueReceiver() (string, string, int) {
	return a.Str1, a.Str2, a.Int
}

func BenchmarkPointerVsValueReceiver(b *testing.B) {
	var ptr *A = &A{"A", "B ", 1, "D", "E", "F", "G", "H", "I", "J", "K", "K", "K", "K", "K", "K", "K", "K"}
	var val A = A{"A", "B ", 1, "D", "E", "F", "G", "H", "I", "J", "K", "K", "K", "K", "K", "K", "K", "K"}
	var ptrB *B = &B{"A", "B ", 1}
	var valB B = B{"A", "B ", 1}

	b.Run("PointerReceiver", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = ptr.GetFieldsPtrReceiver()
		}
	})
	b.Run("ValueReceiver", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = val.GetFieldsValueReceiver()
		}
	})
	b.Run("PointerReceiverB", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = ptrB.GetFieldsPtrReceiver()
		}
	})
	b.Run("ValueReceiverB", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = valB.GetFieldsValueReceiver()
		}
	})
}
