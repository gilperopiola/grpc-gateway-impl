package core

import (
	"fmt"
	"strconv"
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

func WithID(id int32) SqlDBOpt {
	return WithCondition(Where, "id", strconv.Itoa(int(id)))
}

func WithUsername(username string) SqlDBOpt {
	return WithCondition(Where, "username", username)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*      - Low Level SQL Options -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func WithCondition(operation Operation, field, value string) SqlDBOpt {
	if field == "" {
		return func(db InnerSqlDB) {} // No-op
	}

	return func(db InnerSqlDB) {
		if operation == Where || operation == And {
			db.Where(fmt.Sprintf("%s = ?", field), value)
			return
		}

		if operation == Or {
			db.Or(fmt.Sprintf("%s = ?", field), value)
			return
		}

		if operation == Like {
			db.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value+"%")
			return
		}
	}
}
