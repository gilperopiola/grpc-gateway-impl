package utils

// Oh no, a utils package! We're all gonna die!

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Utils -               */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* -~-~-~-~-~ General Utils -~-~-~-~-~- */

// Returns the first element of a string slice, or a fallback if the slice is empty.
func FirstOrDefault(slice []string, fallback string) string {
	if len(slice) > 0 {
		return slice[0]
	}
	return fallback
}

/* -~-~-~-~ Types & Conversions -~-~-~-~ */

type Int32Slice []int32

func (int32Slice Int32Slice) ToIntSlice() []int {
	intSlice := make([]int, len(int32Slice))
	for index, int32Val := range int32Slice {
		intSlice[index] = int(int32Val)
	}
	return intSlice
}
