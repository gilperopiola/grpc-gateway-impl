package db

import (
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/cfg"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/models"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/* ----------------------------------- */
/*             - Database -            */
/* ----------------------------------- */

// NewDB returns a new instance of the Database, either by connecting with a given configuration
// or by using an already configured database. This is useful for testing.
func NewDB(c *cfg.DBConfig) GormAdapter {
	dsn := connectionString(c)

	gormDB, err := gormConnect(dsn, gormConfigOption(c.GormLogLevel))
	if err != nil {
		zap.S().Fatalf(errs.FatalErrMsgConnectingDB, err)
	}

	// Migrate the models to the database.
	gormDB.AutoMigrate(models.AllModels...)

	// Insert the admin user if it doesn't exist.
	if c.InsertAdmin && c.AdminPwd != "" {
		admin := models.User{
			Username: "admin",
			Password: c.AdminPwd,
			Role:     models.AdminRole,
		}

		if err := gormDB.FirstOrCreate(&admin).Error(); err != nil {
			zap.S().Warnf(errs.ErrMsgInsertingAdmin, err)
		}
	}

	// database := &DBWrapper{}
	// database.Connect(c)
	// return database
	return gormDB
}

//// Connect connects to the database using the given configuration.
//func (dbw *DBWrapper) Connect(c *cfg.DBConfig) {
//	dsn := connectionString(c)
//
//	var err error
//	if dbw.DB, err = gormConnect(dsn, gormConfigOption(c.GormLogLevel)); err != nil {
//		zap.S().Fatalf(errs.FatalErrMsgConnectingDB, err)
//	}
//
//	// Migrate the models to the database.
//	dbw.DB.AutoMigrate(models.AllModels...)
//
//	// Insert the admin user if it doesn't exist.
//	if c.InsertAdmin && c.AdminPassword != "" {
//		dbw.InsertAdmin("admin", c.AdminPassword)
//	}
//}
//
//// InsertAdmin inserts the admin user if it doesn't exist.
//func (dbw *DBWrapper) InsertAdmin(adminUsername, adminPwd string) {
//	admin := models.User{
//		Username: adminUsername,
//		Password: adminPwd,
//		Role:     models.AdminRole,
//	}
//
//	if err := dbw.DB.FirstOrCreate(&admin).Error(); err != nil {
//		zap.S().Warnf(errs.ErrMsgInsertingAdmin, err)
//	}
//}
//
//// Close closes the database connection.
//func (dbw *DBWrapper) Close() {
//	sqlDB := dbw.DB.GetSQL()
//	if sqlDB != nil {
//		sqlDB.Close()
//	}
//}

// gormConnect calls gorm.Open and wraps the returned *gorm.DB with our concrete type that implements GormAdapter.
func gormConnect(dsn string, opts ...gorm.Option) (GormAdapter, error) {
	gormDB, err := gorm.Open(mysql.Open(dsn), opts...)
	return newGormAdapter(gormDB), err
}

// gormConfigOption returns a new GormConfig with the given log level.
func gormConfigOption(logLevel int) *gorm.Config {
	return &gorm.Config{
		Logger:         newGormLoggerAdapter(zap.L(), logLevel),
		TranslateError: true,
	}
}

// connectionString returns the connection string for the database.
// Only MySQL is supported for now.
func connectionString(c *cfg.DBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s%s", c.Username, c.Password, c.Hostname, c.Port, c.Schema, c.Params)
}
