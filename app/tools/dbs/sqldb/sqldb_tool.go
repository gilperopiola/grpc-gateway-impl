package sqldb

import (
	"errors"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"gorm.io/gorm"
)

var _ core.DBTool = &sqlDBTool{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - SQL Database Tool -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// The SQL DB Tool holds a SQL Database object/connection.
//
// -> DB Tool = High Level Operations (e.g. CreateUser, GetUser, GetUsers)
// -> DB = Low Level Operations (e.g. Insert, Find, Count)

type sqlDBTool struct {
	DB core.SqlDB
}

func NewDBTool(db core.SqlDB) core.DBTool {
	return &sqlDBTool{db}
}

func (this sqlDBTool) GetDB() core.AnyDB         { return this.DB }
func (this sqlDBTool) CloseDB()                  { this.DB.Close() }
func (this sqlDBTool) IsNotFound(err error) bool { return errors.Is(err, gorm.ErrRecordNotFound) }
