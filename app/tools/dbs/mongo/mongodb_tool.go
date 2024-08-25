package mongo

import (
	"errors"
	"fmt"
	"time"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/types/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ core.DBTool = &mongoDBTool{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*     - External Layer: Storage -     */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type mongoDBTool struct {
	DB core.MongoDB
}

type Collections string

const (
	UsersCollection    Collections = "users"
	GroupsCollection   Collections = "groups"
	GPTChatsCollection Collections = "gpt_chats"
)

func SetupDBTool(db core.MongoDB) *mongoDBTool {
	return &mongoDBTool{db}
}

func (dbt *mongoDBTool) GetDB() core.AnyDB {
	return dbt.DB
}

func (dbt *mongoDBTool) CloseDB() {
	ctx, cancel := god.NewCtxWithTimeout(5 * time.Second)
	dbt.DB.Close(ctx)
	cancel()
}

func (dbt *mongoDBTool) IsNotFound(err error) bool {
	return errors.Is(err, mongo.ErrNoDocuments)
}

func (dbt *mongoDBTool) CreateGroup(ctx god.Ctx, name string, ownerID int, invitedUserIDs []int) (*models.Group, error) {
	group := &models.Group{Name: name, OwnerID: ownerID}

	result, err := dbt.DB.InsertOne(ctx, string(GroupsCollection), group)
	if err != nil || result.InsertedID == nil {
		return nil, errs.DBErr{err, "T0D0"}
	}

	return group, nil
}

func (dbt *mongoDBTool) GetGroup(ctx god.Ctx, opts ...any) (*models.Group, error) {
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

func (dbt *mongoDBTool) CreateUser(ctx god.Ctx, username, hashedPwd string) (*models.User, error) {
	user := &models.User{Username: username, Password: hashedPwd}

	result, err := dbt.DB.InsertOne(ctx, string(UsersCollection), user)
	if err != nil || result.InsertedID == nil {
		return nil, errs.DBErr{err, CreateUserErr}
	}

	return user, nil
}

func (dbt *mongoDBTool) GetUser(ctx god.Ctx, opts ...any) (*models.User, error) {

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

func (dbt *mongoDBTool) GetUsers(ctx god.Ctx, page, pageSize int, opts ...any) (models.Users, int, error) {
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

	var users models.Users
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

func (db *mongoDBTool) GetGPTChat(ctx god.Ctx, opts ...any) (*models.GPTChat, error) {
	return nil, nil
}
func (db *mongoDBTool) CreateGPTChat(ctx god.Ctx, title string) (*models.GPTChat, error) {
	return nil, nil
}
func (db *mongoDBTool) CreateGPTMessage(ctx god.Ctx, message *models.GPTMessage) (*models.GPTMessage, error) {
	return nil, nil
}
