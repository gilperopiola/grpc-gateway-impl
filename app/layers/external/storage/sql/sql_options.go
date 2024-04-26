package sql

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

func WithUserID(userID int) core.SQLQueryOpt {
	return WithCondition(Where, "id", strconv.Itoa(userID))
}

func WithUsername(username string) core.SQLQueryOpt {
	return WithCondition(Where, "username", username)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*      - Low Level SQL Options -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func WithCondition(operation Operation, field, value string) core.SQLQueryOpt {
	if field == "" {
		core.LogWeirdBehaviour("Empty field in SQL condition -> value = " + value)
		return func(db core.SQLDatabaseAPI) {} // No-op
	}

	return func(db core.SQLDatabaseAPI) {
		if operation == Where || operation == And { // Where / And
			db.Where(fmt.Sprintf("%s = ?", field), value)
			return
		}

		if operation == Or { // Or
			db.Or(fmt.Sprintf("%s = ?", field), value)
			return
		}

		if operation == Like { // Like
			if value == "" {
				core.LogWeirdBehaviour("Empty value in SQL condition -> field = " + field)
				return
			}

			db.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value+"%")
			return
		}
	}
}
