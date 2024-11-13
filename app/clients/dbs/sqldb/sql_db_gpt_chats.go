package sqldb

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/models"
)

var _ = &models.GPTChat{}
var _ = &models.GPTMessage{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*     - SQL DB Tool: GPT Chats -      */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (db *DB) DBCreateGPTChat(ctx god.Ctx, title string) (*models.GPTChat, error) {
	gptChat := models.GPTChat{Title: title}
	if err := db.InnerDB.WithContext(ctx).Create(&gptChat).Error(); err != nil {
		return nil, &errs.DBErr{err, "db error -> creating gpt chat"}
	}
	return &gptChat, nil
}

func (db *DB) DBGetGPTChat(ctx god.Ctx, opts ...any) (*models.GPTChat, error) {
	if len(opts) == 0 {
		return nil, &errs.DBErr{nil, NoOptionsErr}
	}

	query := db.InnerDB.Model(&models.GPTChat{}).WithContext(ctx)
	for _, opt := range opts {
		opt.(core.SqlDBOpt)(query)
	}

	var gptChat models.GPTChat
	if err := query.Preload("Messages").First(&gptChat).Error(); err != nil {
		return nil, &errs.DBErr{err, "db error -> getting gpt chat"}
	}
	return &gptChat, nil
}

/* -~-~-~-~-~- GPT Messages -~-~-~-~-~- */

func (db *DB) DBCreateGPTMessage(ctx god.Ctx, message *models.GPTMessage) (*models.GPTMessage, error) {
	if err := db.InnerDB.WithContext(ctx).Create(&message).Error(); err != nil {
		return nil, &errs.DBErr{err, "db error -> creating gpt message"}
	}
	return message, nil
}
