package clients

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/clients/apis"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/logs"
	"github.com/gilperopiola/grpc-gateway-impl/app/db"
	"github.com/gilperopiola/grpc-gateway-impl/app/repositories"
)

var _ core.Clients = &Clients{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Clients -            */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v2 */

// Clients provides access to external services (APIs) and data repositories
type Clients struct {
	core.APIClients
	DB           core.DBOperations
	Repositories *repositories.RepositoryRegistry
}

// Setup initializes all clients and repositories
func Setup(cfg *core.Config, tools core.Tools) (*Clients, error) {
	// Initialize database
	dbOperator, err := db.NewGormDB(&cfg.DBCfg, tools.HashPassword)
	if err != nil {
		return nil, err
	}

	// Initialize repositories
	repos := repositories.NewRepositoryRegistry(dbOperator)

	// Initialize API clients
	apiClients := apis.NewAPIs(&cfg.APIsCfg)

	clients := Clients{
		APIClients:   apiClients,
		DB:           dbOperator,
		Repositories: repos,
	}

	logs.InitModuleOK("Clients", "��")
	return &clients, nil
}

// GetDB returns the underlying database instance
// This is maintained for backward compatibility
func (c *Clients) GetDB() any {
	return c.DB
}

// CloseDB closes the database connection
// This is maintained for backward compatibility
func (c *Clients) CloseDB() error {
	return c.DB.CloseDB()
}

// UserRepository returns the user repository
func (c *Clients) UserRepository() core.UserRepository {
	return c.Repositories.UserRepository
}

// GroupRepository returns the group repository
func (c *Clients) GroupRepository() core.GroupRepository {
	return c.Repositories.GroupRepository
}

// GPTChatRepository returns the GPT chat repository
func (c *Clients) GPTChatRepository() core.GPTChatRepository {
	return c.Repositories.GPTChatRepository
}
