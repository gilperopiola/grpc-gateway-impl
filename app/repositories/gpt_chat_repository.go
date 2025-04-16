package repositories

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*       - GPT Chat Repository -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// GormGPTChatRepository implements the GPTChatRepository interface using GORM
type GormGPTChatRepository struct {
	db core.DBOperations
}

// Verify that GormGPTChatRepository implements the core.GPTChatRepository interface
var _ core.GPTChatRepository = (*GormGPTChatRepository)(nil)

// NewGormGPTChatRepository creates a new GormGPTChatRepository
func NewGormGPTChatRepository(db core.DBOperations) *GormGPTChatRepository {
	return &GormGPTChatRepository{db: db}
}

// GetChatByID retrieves a GPT chat by its ID
func (r *GormGPTChatRepository) GetChatByID(ctx god.Ctx, id int) (*models.GPTChat, error) {
	var chat models.GPTChat

	// Get the chat along with its messages
	err := r.db.WithContext(ctx).(interface {
		First(out interface{}, where ...interface{}) error
		Preload(query string, args ...interface{}) interface {
			First(out interface{}, where ...interface{}) error
		}
	}).Preload("Messages").First(&chat, id)

	if err != nil {
		return nil, &errs.DBErr{Err: err, Context: errs.ChatNotFound}
	}

	return &chat, nil
}

// CreateChat creates a new GPT chat with the specified title
func (r *GormGPTChatRepository) CreateChat(ctx god.Ctx, title string) (*models.GPTChat, error) {
	chat := models.GPTChat{
		Title: title,
	}

	err := r.db.WithContext(ctx).Create(&chat)
	if err != nil {
		return nil, &errs.DBErr{Err: err, Context: errs.FailedToCreateChat}
	}

	return &chat, nil
}

// CreateMessage creates a new GPT message in a chat
func (r *GormGPTChatRepository) CreateMessage(ctx god.Ctx, message *models.GPTMessage) (*models.GPTMessage, error) {
	err := r.db.WithContext(ctx).Create(message)
	if err != nil {
		return nil, &errs.DBErr{Err: err, Context: errs.FailedToCreateMessage}
	}

	return message, nil
}
