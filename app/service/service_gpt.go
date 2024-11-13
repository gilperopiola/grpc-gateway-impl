package service

import (
	"context"
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/app/clients/apis/apimodels"
	"github.com/gilperopiola/grpc-gateway-impl/app/clients/dbs/sqldb"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/utils"
)

type GPTSvc struct {
	pbs.UnimplementedGPTServiceServer
	Clients core.Clients
	Tools   core.Tools
}

func (svc *GPTSvc) NewGPTChat(ctx context.Context, req *pbs.NewGPTChatRequest) (*pbs.NewGPTChatResponse, error) {

	// Call the GPT API.
	gptResponse, err := svc.Clients.SendToGPT(ctx, req.Message)
	if err != nil {
		return nil, fmt.Errorf("error calling GPT API: %w", err)
	}

	// Create a new GPTChat in the database
	gptChat, err := svc.Clients.DBCreateGPTChat(ctx, req.Message)
	if err != nil {
		return nil, errs.GRPCFromDB(err, shared.GetRouteFromCtx(ctx).Name)
	}

	// Define the initial messages
	messages := []*models.GPTMessage{
		{Title: "System default message", From: "system", Content: "You are a highly...", ChatID: gptChat.ID},
		{Title: "User prompt", From: "user", Content: req.Message, ChatID: gptChat.ID},
		{Title: "GPT response", From: "assistant", Content: gptResponse, ChatID: gptChat.ID},
	}

	// Add all messages to the database
	for _, msg := range messages {
		if _, err := svc.Clients.DBCreateGPTMessage(ctx, msg); err != nil {
			return nil, errs.GRPCFromDB(err, shared.GetRouteFromCtx(ctx).Name)
		}
	}

	return &pbs.NewGPTChatResponse{
		GptMessage: gptResponse,
		Chat: &pbs.GPTChatInfo{
			Id:    int32(gptChat.ID),
			Title: gptChat.Title,
		},
	}, nil
}

func (svc *GPTSvc) ReplyToGPTChat(ctx context.Context, req *pbs.ReplyToGPTChatRequest) (*pbs.ReplyToGPTChatResponse, error) {

	// Get existing chat from DB.
	chat, err := svc.Clients.DBGetGPTChat(ctx, sqldb.WithID(req.ChatId))
	if err != nil {
		if utils.IsNotFound(err) {
			return nil, errs.GRPCNotFound("GPT Chat", int(req.ChatId))
		}
		return nil, errs.GRPCFromDB(err, shared.GetRouteFromCtx(ctx).Name)
	}

	// Prepare the previous messages for the GPT API call
	var previousChatMsgs []apimodels.GPTMessage
	for _, msg := range chat.Messages {
		previousChatMsgs = append(previousChatMsgs, apimodels.GPTMessage{
			Role:    msg.From,
			Content: msg.Content,
		})
	}

	// Call the GPT API to generate a response
	gptResponse, err := svc.Clients.SendToGPT(ctx, req.Message, previousChatMsgs...)
	if err != nil {
		return nil, fmt.Errorf("error calling GPT API: %w", err)
	}

	// Create the user message and the assistant's response in the database
	messages := []*models.GPTMessage{
		{Title: "User response", From: "user", Content: req.Message, ChatID: chat.ID},
		{Title: "Assistant response", From: "assistant", Content: gptResponse, ChatID: chat.ID},
	}

	// Store all messages in the database
	for _, msg := range messages {
		if _, err := svc.Clients.DBCreateGPTMessage(ctx, msg); err != nil {
			return nil, errs.GRPCFromDB(err, shared.GetRouteFromCtx(ctx).Name)
		}
	}

	return &pbs.ReplyToGPTChatResponse{
		GptMessage: gptResponse,
		Chat: &pbs.GPTChatInfo{
			Id:    int32(chat.ID),
			Title: chat.Title,
		},
	}, nil
}
