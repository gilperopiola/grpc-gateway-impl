package mongo

import (
	"context"
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	UsersCollection Collections = "users"
)

func SetupDBTool(db core.MongoDB) *mongoDBTool {
	return &mongoDBTool{db}
}

func (s *mongoDBTool) GetDB() core.DB {
	return s.DB
}

func (s *mongoDBTool) GetDBTool() core.DBTool {
	return s
}

func (s *mongoDBTool) CreateUser(ctx context.Context, username, hashedPwd string) (*core.User, error) {
	user := &core.User{Username: username, Password: hashedPwd}

	result, err := s.DB.InsertOne(ctx, string(UsersCollection), user)
	if err != nil || result.InsertedID == nil {
		return nil, errs.DBErr{err, CreateUserErr}
	}

	return user, nil
}

func (s *mongoDBTool) GetUser(ctx context.Context, opts ...any) (*core.User, error) {

	if len(opts) == 0 {
		return nil, errs.DBErr{nil, NoOptionsErr}
	}

	filter := &bson.D{}
	for _, opt := range opts {
		opt.(core.MongoDBOpt)(filter)
	}

	result := s.DB.FindOne(ctx, string(UsersCollection), filter)
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

func (s *mongoDBTool) GetUsers(ctx context.Context, page, pageSize int, opts ...any) (core.Users, int, error) {
	filter := &bson.D{}
	for _, opt := range opts {
		opt.(core.MongoDBOpt)(filter)
	}

	matches, err := s.DB.Count(ctx, string(UsersCollection), filter)
	if err != nil {
		return nil, 0, errs.DBErr{err, CountUsersErr}
	}
	if matches == 0 {
		return nil, 0, nil
	}

	result, err := s.DB.Find(ctx, string(UsersCollection), filter, page, pageSize)
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

// T0D0 this can probably be done in a way that ExternalLayer can hold many Storage(s),
// each of those being the same Storage struct but with a generic type that would be sql or MongoDB or so. DRY overload.
