package mongo

import (
	"errors"
	"fmt"
	"time"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var _ core.DBTool = (*mongoDBTool)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*     - External Layer: Storage -     */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type mongoDBTool struct {
	DB core.MongoDB
}

type Collections string

const (
	UsersCollection  Collections = "users"
	GroupsCollection Collections = "groups"
)

func SetupDBTool(db core.MongoDB) *mongoDBTool {
	return &mongoDBTool{db}
}

func (dbt *mongoDBTool) GetDB() core.DB {
	return dbt.DB
}

func (dbt *mongoDBTool) CloseDB() {
	ctx, cancel := god.NewCtxWithTimeout(5 * time.Second)
	dbt.DB.Close(ctx)
	cancel()
}

func (dbt *mongoDBTool) IsNotFound(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, gorm.ErrRecordNotFound)
}

func (dbt *mongoDBTool) CreateGroup(ctx god.Ctx, name string, ownerID int, invitedUserIDs []int) (*core.Group, error) {
	group := &core.Group{Name: name, OwnerID: ownerID}

	result, err := dbt.DB.InsertOne(ctx, string(GroupsCollection), group)
	if err != nil || result.InsertedID == nil {
		return nil, errs.DBErr{err, "T0D0"}
	}

	return group, nil
}

func (dbt *mongoDBTool) GetGroup(ctx god.Ctx, opts ...any) (*core.Group, error) {
	if len(opts) == 0 {
		return nil, errs.DBErr{nil, NoOptionsErr}
	}

	filter := &bson.D{}
	for _, opt := range opts {
		opt.(core.MongoDBOpt)(filter)
	}

	result := dbt.DB.FindOne(ctx, string(GroupsCollection), filter)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.DBErr{mongo.ErrNoDocuments, "group not found"}
		}

		return nil, errs.DBErr{mongo.ErrNoDocuments, fmt.Sprintf("error finding group: %v", err)}
	}

	var group core.Group
	if err := result.Decode(&group); err != nil {
		return nil, errs.DBErr{err, fmt.Sprintf("error decoding group: %v", err)}
	}

	return &group, nil
}

func (dbt *mongoDBTool) CreateUser(ctx god.Ctx, username, hashedPwd string) (*core.User, error) {
	user := &core.User{Username: username, Password: hashedPwd}

	result, err := dbt.DB.InsertOne(ctx, string(UsersCollection), user)
	if err != nil || result.InsertedID == nil {
		return nil, errs.DBErr{err, CreateUserErr}
	}

	return user, nil
}

func (dbt *mongoDBTool) GetUser(ctx god.Ctx, opts ...any) (*core.User, error) {

	if len(opts) == 0 {
		return nil, errs.DBErr{nil, NoOptionsErr}
	}

	filter := &bson.D{}
	for _, opt := range opts {
		opt.(core.MongoDBOpt)(filter)
	}

	result := dbt.DB.FindOne(ctx, string(UsersCollection), filter)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errs.DBErr{mongo.ErrNoDocuments, "user not found"}
		}
		return nil, errs.DBErr{mongo.ErrNoDocuments, fmt.Sprintf("error finding user: %v", err)}
	}

	var user core.User
	if err := result.Decode(&user); err != nil {
		return nil, errs.DBErr{err, fmt.Sprintf("error decoding user: %v", err)}
	}

	return &user, nil
}

func (dbt *mongoDBTool) GetUsers(ctx god.Ctx, page, pageSize int, opts ...any) (core.Users, int, error) {
	filter := &bson.D{}
	for _, opt := range opts {
		opt.(core.MongoDBOpt)(filter)
	}

	matches, err := dbt.DB.Count(ctx, string(UsersCollection), filter)
	if err != nil {
		return nil, 0, errs.DBErr{err, CountUsersErr}
	}
	if matches == 0 {
		return nil, 0, nil
	}

	result, err := dbt.DB.Find(ctx, string(UsersCollection), filter, page, pageSize)
	if err != nil {
		return nil, 0, errs.DBErr{mongo.ErrNoDocuments, fmt.Sprintf("error finding users: %v", err)}
	}

	var users core.Users
	if err := result.Decode(&users); err != nil {
		return nil, 0, errs.DBErr{err, fmt.Sprintf("error decoding users: %v", err)}
	}

	return users, int(matches), nil
}

var (
	CreateUserErr = errs.DBCreatingUser
	GetUserErr    = errs.DBGettingUser
	GetUsersErr   = errs.DBGettingUsers
	CountUsersErr = errs.DBCountingUsers
	NoOptionsErr  = errs.DBNoQueryOpts
)
