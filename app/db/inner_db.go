package db

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"

	"gorm.io/gorm"
)

var _ core.InnerDB = &innerDB{}

// Embeds a *gorm.DB object which itself holds a connection to an SQL DB.
// Implements all basic low-level methods to interact with an SQL DB.
type innerDB struct {
	*gorm.DB
}

func (this *innerDB) GetInnerDB() any { return this.DB }

func (this *innerDB) InsertAdmin(hashedPwd string) {
	admin := models.User{
		Username: "admin",
		Password: hashedPwd,
		Role:     models.AdminRole,
	}
	logs.WarnIfErr(this.DB.FirstOrCreate(&admin).Error, errs.FailedToInsertDBAdmin)
}

func (this *innerDB) Association(column string) core.SqlDBAssoc { return this.DB.Association(column) }

func (this *innerDB) Count(value *int64) core.InnerDB { return &innerDB{this.DB.Count(value)} }

func (this *innerDB) Create(value any) core.InnerDB { return &innerDB{this.DB.Create(value)} }

func (this *innerDB) Debug() core.InnerDB { return &innerDB{this.DB.Debug()} }

func (this *innerDB) Error() error { return this.DB.Error }

func (this *innerDB) Group(query string) core.InnerDB { return &innerDB{this.DB.Group(query)} }

func (this *innerDB) Limit(value int) core.InnerDB { return &innerDB{this.DB.Limit(value)} }

func (this *innerDB) Model(value any) core.InnerDB { return &innerDB{this.DB.Model(value)} }

func (this *innerDB) Offset(value int) core.InnerDB { return &innerDB{this.DB.Offset(value)} }

func (this *innerDB) Order(value string) core.InnerDB { return &innerDB{this.DB.Order(value)} }

func (this *innerDB) RowsAffected() int64 { return this.DB.RowsAffected }

func (this *innerDB) Save(value any) core.InnerDB { return &innerDB{this.DB.Save(value)} }

func (this *innerDB) Scan(to any) core.InnerDB { return &innerDB{this.DB.Scan(to)} }

func (this *innerDB) Close() {
	sqlDB, err := this.DB.DB()
	logs.LogIfErr(err, errs.FailedToGetSQLDB)
	if sqlDB != nil {
		logs.LogIfErr(sqlDB.Close(), errs.FailedToCloseSQLDB)
	}
}

func (this *innerDB) Delete(val any, where ...any) core.InnerDB {
	return &innerDB{this.DB.Delete(val, where)}
}

func (this *innerDB) Find(out any, where ...any) core.InnerDB {
	return &innerDB{this.DB.Find(out, where...)}
}

func (this *innerDB) First(out any, where ...any) core.InnerDB {
	return &innerDB{this.DB.First(out, where...)}
}

func (this *innerDB) FirstOrCreate(out any, where ...any) core.InnerDB {
	return &innerDB{this.DB.FirstOrCreate(out, where...)}
}

func (this *innerDB) Joins(qry string, args ...any) core.InnerDB {
	return &innerDB{this.DB.Joins(qry, args)}
}

func (this *innerDB) Or(query any, args ...any) core.InnerDB {
	return &innerDB{this.DB.Or(query, args...)}
}

func (this *innerDB) Pluck(col string, val any) core.InnerDB {
	return &innerDB{this.DB.Pluck(col, val)}
}

func (this *innerDB) Preload(query string, args ...any) core.InnerDB {
	return &innerDB{this.DB.Preload(query, args...)}
}

func (this *innerDB) Raw(sql string, vals ...any) core.InnerDB {
	return &innerDB{this.DB.Raw(sql, vals...)}
}

func (this *innerDB) Scopes(fns ...func(core.InnerDB) core.InnerDB) core.InnerDB {
	adaptedFns := make([]func(*gorm.DB) *gorm.DB, len(fns))
	for i, fn := range fns {
		adaptedFns[i] = func(db *gorm.DB) *gorm.DB {
			return fn(&innerDB{db}).(*innerDB).DB // Messy. T0D0.
		}
	}

	return &innerDB{this.DB.Scopes(adaptedFns...)}
}

func (this *innerDB) Unscoped() core.InnerDB {
	return &innerDB{this.DB.Unscoped()}
}

// Calling the actual gorm WithContext func makes our SQLOptions fail to apply for some reason. T0D0.
func (this *innerDB) WithContext(ctx god.Ctx) core.InnerDB {
	return &innerDB{this.DB}
}

func (this *innerDB) Where(q any, args ...any) core.InnerDB {
	return &innerDB{this.DB.Where(q, args...)}
}
