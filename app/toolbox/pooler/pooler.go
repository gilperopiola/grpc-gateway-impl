package pooler

import (
	"golang.org/x/time/rate"
)

type pooler struct {
	userRateLimitingPool PoolOf[unit[*rate.Limiter]]
}
