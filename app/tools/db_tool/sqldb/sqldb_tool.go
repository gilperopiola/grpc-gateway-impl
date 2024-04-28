package sqldb

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

var _ core.DBTool = (*sqlDBTool)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - SQL Database Tool -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// The SQL DB Tool holds a SQL Database object/connection.
//
// -> DB Tool = High Level Operations (e.g. CreateUser, GetUser, GetUsers)
// -> DB = Low Level Operations (e.g. Insert, Find, Count)

type sqlDBTool struct {
	DB core.SQLDB
}

func NewDBTool(db core.SQLDB) core.DBTool {
	return &sqlDBTool{db}
}

func (sdbt *sqlDBTool) GetDBTool() core.DBTool { return sdbt }

func (sdbt *sqlDBTool) GetDB() core.DB { return sdbt.DB }

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
