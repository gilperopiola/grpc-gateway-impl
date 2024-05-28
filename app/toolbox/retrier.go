package toolbox

import (
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

var _ core.Retrier = (*retrier)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Retrier -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type retrier struct {
	config *core.RetrierCfg
}

func NewRetrier(cfg *core.RetrierCfg) core.Retrier {
	return &retrier{cfg}
}

func (r retrier) TryToConnectToDB(connectToDB func() (any, error), execOnFailure func()) (any, error) {
	var resp any
	var err error

	for i := 0; i < r.config.DBConnRetries; i++ {
		if resp, err = connectToDB(); err == nil {
			return resp, nil
		}

		execOnFailure()
		time.Sleep(time.Second)
	}

	return resp, err
}
