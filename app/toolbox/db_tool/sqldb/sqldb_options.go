package sqldb

import (
	"fmt"
	"strconv"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

type Operation string

const (
	Where Operation = "where"
	And   Operation = "and"
	Or    Operation = "or"
	Like  Operation = "like"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*      - High Level SQL Options -     */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func WithID(id int32) core.SQLDBOpt {
	return WithCondition(Where, "id", strconv.Itoa(int(id)))
}

func WithUsername(username string) core.SQLDBOpt {
	return WithCondition(Where, "username", username)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*      - Low Level SQL Options -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func WithCondition(operation Operation, field, value string) core.SQLDBOpt {
	if field == "" {
		core.LogWeirdBehaviour("Empty field in SQL condition -> value = " + value)
		return func(db core.SQLDB) {} // No-op
	}

	return func(db core.SQLDB) {
		if operation == Where || operation == And { // Where / And
			db.Where(fmt.Sprintf("%s = ?", field), value)
			return
		}

		if operation == Or { // Or
			db.Or(fmt.Sprintf("%s = ?", field), value)
			return
		}

		if operation == Like { // Like
			db.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value+"%")
			return
		}
	}
}
