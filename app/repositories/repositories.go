package repositories

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - Repository Registry -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// RepositoryRegistry provides access to all repositories
type RepositoryRegistry struct {
	UserRepository    core.UserRepository
	GroupRepository   core.GroupRepository
	GPTChatRepository core.GPTChatRepository
}

// NewRepositoryRegistry creates a new RepositoryRegistry with all repositories
func NewRepositoryRegistry(db core.DBOperations) *RepositoryRegistry {
	return &RepositoryRegistry{
		UserRepository:    NewGormUserRepository(db),
		GroupRepository:   NewGormGroupRepository(db),
		GPTChatRepository: NewGormGPTChatRepository(db),
	}
}
