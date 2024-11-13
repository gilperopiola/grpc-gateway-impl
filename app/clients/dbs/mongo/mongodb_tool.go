package mongo

import (
	"fmt"
	"time"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ core.DB = &mongoDBConn{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*     - External Layer: Storage -     */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type mongoDBConn struct {
	DB core.InnerMongoDB
}

type Collections string

const (
	UsersCollection    Collections = "users"
	GroupsCollection   Collections = "groups"
	GPTChatsCollection Collections = "gpt_chats"
)

func SetupDBConn(db core.InnerMongoDB) *mongoDBConn {
	return &mongoDBConn{db}
}

func (dbt *mongoDBConn) GetDB() any { return dbt.DB }

func (dbt *mongoDBConn) CloseDB() {
	ctx, cancel := god.NewCtxWithTimeout(5 * time.Second)
	dbt.DB.Close(ctx)
	cancel()
}

func (dbt *mongoDBConn) DBCreateGroup(ctx god.Ctx, name string, ownerID int, invitedUserIDs []int) (*models.Group, error) {
	group := &models.Group{Name: name, OwnerID: ownerID}

	result, err := dbt.DB.InsertOne(ctx, string(GroupsCollection), group)
	if err != nil || result.InsertedID == nil {
		return nil, errs.DBErr{err, "T0D0"}
	}

	return group, nil
}

func (dbt *mongoDBConn) DBGetGroup(ctx god.Ctx, opts ...any) (*models.Group, error) {
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

	var group models.Group
	if err := result.Decode(&group); err != nil {
		return nil, errs.DBErr{err, fmt.Sprintf("error decoding group: %v", err)}
	}

	return &group, nil
}

func (dbt *mongoDBConn) DBCreateUser(ctx god.Ctx, username, hashedPwd string) (*models.User, error) {
	user := &models.User{Username: username, Password: hashedPwd}

	result, err := dbt.DB.InsertOne(ctx, string(UsersCollection), user)
	if err != nil || result.InsertedID == nil {
		return nil, errs.DBErr{err, CreateUserErr}
	}

	return user, nil
}

func (dbt *mongoDBConn) DBGetUser(ctx god.Ctx, opts ...any) (*models.User, error) {

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

	var user models.User
	if err := result.Decode(&user); err != nil {
		return nil, errs.DBErr{err, fmt.Sprintf("error decoding user: %v", err)}
	}

	return &user, nil
}

func (dbt *mongoDBConn) DBGetUsers(ctx god.Ctx, page, pageSize int, opts ...any) ([]*models.User, int, error) {
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

	var users []*models.User
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

func (db *mongoDBConn) DBGetGPTChat(ctx god.Ctx, opts ...any) (*models.GPTChat, error) {
	return nil, nil
}
func (db *mongoDBConn) DBCreateGPTChat(ctx god.Ctx, title string) (*models.GPTChat, error) {
	return nil, nil
}
func (db *mongoDBConn) DBCreateGPTMessage(ctx god.Ctx, message *models.GPTMessage) (*models.GPTMessage, error) {
	return nil, nil
}
