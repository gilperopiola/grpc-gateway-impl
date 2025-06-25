package core

import (
	"context"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - Database Interfaces -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// DBOperations defines a simplified interface for database operations
// This replaces the previous InnerDB interface
type DBOperations interface {
	// Core database operations
	FindError(out any, where ...any) error
	FirstError(out any, where ...any) error
	CreateError(value any) error
	SaveError(value any) error
	DeleteError(value any, where ...any) error
	CountError(value *int64) error

	WithContext(ctx context.Context) DBOperations
	Transaction(fn func(tx DBOperations) error) error
	CloseDB() error

	Count(value *int64) error
	Model(value any) DBOperations
	Offset(value int) DBOperations
	Limit(value int) DBOperations
	Preload(query string, args ...any) DBOperations
}

// UserRepository handles user-related database operations
type UserRepository interface {
	CreateUser(ctx god.Ctx, username, hashedPwd string) (*models.User, error)
	GetUserByID(ctx god.Ctx, id int) (*models.User, error)
	GetUserByUsername(ctx god.Ctx, username string) (*models.User, error)
	GetUsers(ctx god.Ctx, page, pageSize int) ([]*models.User, int, error)
}

// GroupRepository handles group-related database operations
type GroupRepository interface {
	CreateGroup(ctx god.Ctx, name string, ownerID int, invitedUserIDs []int) (*models.Group, error)
	GetGroupByID(ctx god.Ctx, id int) (*models.Group, error)
	GetGroupsByUserID(ctx god.Ctx, userID int) ([]*models.Group, error)
}

// GPTChatRepository handles GPT chat-related database operations
type GPTChatRepository interface {
	GetChatByID(ctx god.Ctx, id int) (*models.GPTChat, error)
	CreateChat(ctx god.Ctx, title string) (*models.GPTChat, error)
	CreateMessage(ctx god.Ctx, message *models.GPTMessage) (*models.GPTMessage, error)
}
