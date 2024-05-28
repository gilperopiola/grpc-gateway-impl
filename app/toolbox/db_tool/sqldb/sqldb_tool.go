package sqldb

import (
	"errors"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
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

func (sdbt sqlDBTool) GetDB() core.DB {
	return sdbt.DB
}

func (sdbt sqlDBTool) CloseDB() {
	sdbt.DB.Close()
}

func (sdbt sqlDBTool) IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, mongo.ErrNoDocuments)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
