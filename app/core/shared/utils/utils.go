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

/* -~-~-~-~-~ Logger Utils -~-~-~-~-~- */
