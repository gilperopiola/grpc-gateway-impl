package db

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
)

var _ = &models.GPTChat{}
var _ = &models.GPTMessage{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*     - SQL DB Tool: GPT Chats -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Deprecated: Use repositories.GPTChatRepository instead
func (db *LegacyDB) DBCreateGPTChat(ctx god.Ctx, title string) (*models.GPTChat, error) {
	gptChat := models.GPTChat{Title: title}
	if err := db.InnerDB.WithContext(ctx).Create(&gptChat).Error(); err != nil {
		return nil, &errs.DBErr{Err: err, Context: "db error -> creating gpt chat"}
	}
	return &gptChat, nil
}

// Deprecated: Use repositories.GPTChatRepository instead
func (db *LegacyDB) DBGetGPTChat(ctx god.Ctx, opts ...any) (*models.GPTChat, error) {
	if len(opts) == 0 {
		return nil, &errs.DBErr{Err: nil, Context: NoOptionsErr}
	}

	query := db.InnerDB.Model(&models.GPTChat{}).WithContext(ctx)
	for _, opt := range opts {
		opt.(core.SqlDBOpt)(query)
	}

	var gptChat models.GPTChat
	if err := query.Preload("Messages").First(&gptChat).Error(); err != nil {
		return nil, &errs.DBErr{Err: err, Context: "db error -> getting gpt chat"}
	}
	return &gptChat, nil
}

/* -~-~-~-~-~- GPT Messages -~-~-~-~-~- */

// Deprecated: Use repositories.GPTChatRepository instead
func (db *LegacyDB) DBCreateGPTMessage(ctx god.Ctx, message *models.GPTMessage) (*models.GPTMessage, error) {
	if err := db.InnerDB.WithContext(ctx).Create(&message).Error(); err != nil {
		return nil, &errs.DBErr{Err: err, Context: "db error -> creating gpt message"}
	}
	return message, nil
}
