package utils

func FmtForRetrier(f func() (any, error)) func() error {
	return func() error {
		_, err := f()
		return err
	}
}
