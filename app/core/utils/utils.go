package utils

import (
	"errors"
	"io/fs"

	"go.uber.org/zap"
)

// Oh no, a utils package! We're all gonna die!

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Utils -               */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns the first element of a string slice, or a fallback if the slice is empty.
func FirstOrDefault(slice []string, fallback string) string {
	if len(slice) > 0 {
		return slice[0]
	}
	return fallback
}

// On Windows, I'm getting a *fs.PathError when calling zap.L().Sync() to flush
// the Logger on shutdown.
// This just calls zap.L().Sync() and ignores that specific error.
// See https://github.com/uber-go/zap/issues/991
func SyncLogger() error {
	var pathErr *fs.PathError
	if err := zap.L().Sync(); err != nil && !errors.As(err, &pathErr) {
		return err
	}
	return nil
}
