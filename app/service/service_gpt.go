package service

import (
	"context"
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/types/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/utils"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/apis/apimodels"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools/dbs/sqldb"
)

type GPTSubService struct {
	pbs.UnimplementedGPTServiceServer
	Tools core.Tools
}

func (svc *GPTSubService) NewGPTChat(ctx context.Context, req *pbs.NewGPTChatRequest) (*pbs.NewGPTChatResponse, error) {

	// Call the GPT API.
	gptResponse, err := svc.Tools.SendToGPT(ctx, req.Message)
	if err != nil {
		return nil, fmt.Errorf("error calling GPT API: %w", err)
	}

	// Create a new GPTChat in the database
	gptChat, err := svc.Tools.CreateGPTChat(ctx, req.Message)
	if err != nil {
		return nil, errs.GRPCFromDB(err, utils.RouteNameFromCtx(ctx))
	}

	// Define the initial messages
	messages := []*models.GPTMessage{
		{Title: "System default message", From: "system", Content: "You are a highly...", ChatID: gptChat.ID},
		{Title: "User prompt", From: "user", Content: req.Message, ChatID: gptChat.ID},
		{Title: "GPT response", From: "assistant", Content: gptResponse, ChatID: gptChat.ID},
	}

	// Add all messages to the database
	for _, msg := range messages {
		if _, err := svc.Tools.CreateGPTMessage(ctx, msg); err != nil {
			return nil, errs.GRPCFromDB(err, utils.RouteNameFromCtx(ctx))
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

func (svc *GPTSubService) ReplyToGPTChat(ctx context.Context, req *pbs.ReplyToGPTChatRequest) (*pbs.ReplyToGPTChatResponse, error) {

	// Get existing chat from DB.
	chat, err := svc.Tools.GetGPTChat(ctx, sqldb.WithID(req.ChatId))
	if err != nil {
		if svc.Tools.IsNotFound(err) {
			return nil, errs.GRPCNotFound("GPT Chat", int(req.ChatId))
		}
		return nil, errs.GRPCFromDB(err, utils.RouteNameFromCtx(ctx))
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
	gptResponse, err := svc.Tools.SendToGPT(ctx, req.Message, previousChatMsgs...)
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
		if _, err := svc.Tools.CreateGPTMessage(ctx, msg); err != nil {
			return nil, errs.GRPCFromDB(err, utils.RouteNameFromCtx(ctx))
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
