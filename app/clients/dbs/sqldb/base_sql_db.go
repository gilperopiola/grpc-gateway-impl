package sqldb

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/models"

	"gorm.io/gorm"
)

var _ core.InnerSqlDB = &baseSQLDB{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*           - SQL Database -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Embeds a *gorm.DB object which itself holds a connection to an SQL DB.
// Implements all basic low-level methods to interact with an SQL DB.
type baseSQLDB struct {
	*gorm.DB
}

func (this *baseSQLDB) GetInnerDB() any { return this.DB }

func (this *baseSQLDB) InsertAdmin(hashedPwd string) {
	admin := models.User{
		Username: "admin",
		Password: hashedPwd,
		Role:     models.AdminRole,
	}
	logs.WarnIfErr(this.DB.FirstOrCreate(&admin).Error, errs.FailedToInsertDBAdmin)
}

func (this *baseSQLDB) Association(column string) core.SqlDBAssoc { return this.DB.Association(column) }

func (this *baseSQLDB) Count(value *int64) core.InnerSqlDB { return &baseSQLDB{this.DB.Count(value)} }

func (this *baseSQLDB) Create(value any) core.InnerSqlDB { return &baseSQLDB{this.DB.Create(value)} }

func (this *baseSQLDB) Debug() core.InnerSqlDB { return &baseSQLDB{this.DB.Debug()} }

func (this *baseSQLDB) Error() error { return this.DB.Error }

func (this *baseSQLDB) Group(query string) core.InnerSqlDB { return &baseSQLDB{this.DB.Group(query)} }

func (this *baseSQLDB) Limit(value int) core.InnerSqlDB { return &baseSQLDB{this.DB.Limit(value)} }

func (this *baseSQLDB) Model(value any) core.InnerSqlDB { return &baseSQLDB{this.DB.Model(value)} }

func (this *baseSQLDB) Offset(value int) core.InnerSqlDB { return &baseSQLDB{this.DB.Offset(value)} }

func (this *baseSQLDB) Order(value string) core.InnerSqlDB { return &baseSQLDB{this.DB.Order(value)} }

func (this *baseSQLDB) Row() core.SqlRow { return this.DB.Row() }

func (this *baseSQLDB) Rows() (core.SqlRows, error) { return this.DB.Rows() }

func (this *baseSQLDB) RowsAffected() int64 { return this.DB.RowsAffected }

func (this *baseSQLDB) Save(value any) core.InnerSqlDB { return &baseSQLDB{this.DB.Save(value)} }

func (this *baseSQLDB) Scan(to any) core.InnerSqlDB { return &baseSQLDB{this.DB.Scan(to)} }

func (this *baseSQLDB) Close() {
	innerDB, err := this.DB.DB()
	logs.LogIfErr(err, errs.FailedToGetSQLDB)

	if innerDB != nil {
		logs.LogIfErr(innerDB.Close(), errs.FailedToCloseSqlDB)
	}
}

func (this *baseSQLDB) Delete(val any, where ...any) core.InnerSqlDB {
	return &baseSQLDB{this.DB.Delete(val, where)}
}

func (this *baseSQLDB) Find(out any, where ...any) core.InnerSqlDB {
	return &baseSQLDB{this.DB.Find(out, where...)}
}

func (this *baseSQLDB) First(out any, where ...any) core.InnerSqlDB {
	return &baseSQLDB{this.DB.First(out, where...)}
}

func (this *baseSQLDB) FirstOrCreate(out any, where ...any) core.InnerSqlDB {
	return &baseSQLDB{this.DB.FirstOrCreate(out, where...)}
}

func (this *baseSQLDB) Joins(qry string, args ...any) core.InnerSqlDB {
	return &baseSQLDB{this.DB.Joins(qry, args)}
}

func (this *baseSQLDB) Or(query any, args ...any) core.InnerSqlDB {
	return &baseSQLDB{this.DB.Or(query, args...)}
}

func (this *baseSQLDB) Pluck(col string, val any) core.InnerSqlDB {
	return &baseSQLDB{this.DB.Pluck(col, val)}
}

func (this *baseSQLDB) Preload(query string, args ...any) core.InnerSqlDB {
	return &baseSQLDB{this.DB.Preload(query, args...)}
}

func (this *baseSQLDB) Raw(sql string, vals ...any) core.InnerSqlDB {
	return &baseSQLDB{this.DB.Raw(sql, vals...)}
}

func (this *baseSQLDB) Scopes(fns ...func(core.InnerSqlDB) core.InnerSqlDB) core.InnerSqlDB {
	adaptedFns := make([]func(*gorm.DB) *gorm.DB, len(fns))
	for i, fn := range fns {
		adaptedFns[i] = func(db *gorm.DB) *gorm.DB {
			return fn(&baseSQLDB{db}).(*baseSQLDB).DB // Messy. T0D0.
		}
	}

	return &baseSQLDB{this.DB.Scopes(adaptedFns...)}
}

func (this *baseSQLDB) Unscoped() core.InnerSqlDB {
	return &baseSQLDB{this.DB.Unscoped()}
}

// Calling the actual gorm WithContext func makes our SQLOptions fail to apply for some reason. T0D0.
func (this *baseSQLDB) WithContext(ctx god.Ctx) core.InnerSqlDB {
	return &baseSQLDB{this.DB}
}

func (this *baseSQLDB) Where(q any, args ...any) core.InnerSqlDB {
	return &baseSQLDB{this.DB.Where(q, args...)}
}
