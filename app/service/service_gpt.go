package service

import (
	"context"
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/apimodels"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
)

type GPTSvc struct {
	pbs.UnimplementedGPTServiceServer
	Clients core.Clients
	Tools   core.Tools
}

func (svc *GPTSvc) NewGPTChat(ctx context.Context, req *pbs.NewGPTChatRequest) (*pbs.NewGPTChatResponse, error) {

	gptResponse, err := svc.Clients.SendRequestToGPT(ctx, req.Message)
	if err != nil {
		return nil, fmt.Errorf("error calling GPT API: %w", err)
	}

	dbGPTChat, err := svc.Clients.DBCreateGPTChat(ctx, req.Message)
	if err != nil {
		return nil, errs.GRPCFromDB(err, core.GetRouteFromCtx(ctx).Name)
	}

	dbMessages := []*models.GPTMessage{
		{Title: "Instructions", From: "user", Content: "You are a highly...", ChatID: dbGPTChat.ID},
		{Title: "User prompt", From: "user", Content: req.Message, ChatID: dbGPTChat.ID},
		{Title: "GPT response", From: "assistant", Content: gptResponse, ChatID: dbGPTChat.ID},
	}

	for _, msg := range dbMessages {
		if _, err := svc.Clients.DBCreateGPTMessage(ctx, msg); err != nil {
			return nil, errs.GRPCFromDB(err, core.GetRouteFromCtx(ctx).Name)
		}
	}

	return &pbs.NewGPTChatResponse{GptMessage: gptResponse, Chat: &pbs.GPTChatInfo{Id: int32(dbGPTChat.ID), Title: dbGPTChat.Title}}, nil
}

func (svc *GPTSvc) ReplyToGPTChat(ctx context.Context, req *pbs.ReplyToGPTChatRequest) (*pbs.ReplyToGPTChatResponse, error) {

	dbGPTChat, err := svc.Clients.DBGetGPTChat(ctx, core.WithID(req.ChatId))
	if err != nil {
		if errs.IsDBNotFound(err) {
			return nil, errs.GRPCNotFound("GPT Chat", int(req.ChatId))
		}
		return nil, errs.GRPCFromDB(err, core.GetRouteFromCtx(ctx).Name)
	}

	var prevMsgs []apimodels.GPTChatMsg
	for _, msg := range dbGPTChat.Messages {
		prevMsgs = append(prevMsgs, apimodels.GPTChatMsg{Role: msg.From, Content: msg.Content})
	}

	gptResponse, err := svc.Clients.SendRequestToGPT(ctx, req.Message, prevMsgs...)
	if err != nil {
		return nil, fmt.Errorf("error calling GPT API: %w", err)
	}

	dbMessages := []*models.GPTMessage{
		{Title: "User response", From: "user", Content: req.Message, ChatID: dbGPTChat.ID},
		{Title: "GPT response", From: "assistant", Content: gptResponse, ChatID: dbGPTChat.ID},
	}

	for _, msg := range dbMessages {
		if _, err := svc.Clients.DBCreateGPTMessage(ctx, msg); err != nil {
			return nil, errs.GRPCFromDB(err, core.GetRouteFromCtx(ctx).Name)
		}
	}

	return &pbs.ReplyToGPTChatResponse{GptMessage: gptResponse, Chat: &pbs.GPTChatInfo{Id: int32(dbGPTChat.ID), Title: dbGPTChat.Title}}, nil
}

func (svc *GPTSvc) NewGPTImage(ctx context.Context, req *pbs.NewGPTImageRequest) (*pbs.NewGPTImageResponse, error) {

	dalleResponse, err := svc.Clients.SendRequestToDallE(ctx, req.Message, req.Size)
	if err != nil {
		return nil, fmt.Errorf("error calling Dall-E API: %w", err)
	}

	dbGPTChat, err := svc.Clients.DBCreateGPTChat(ctx, req.Message)
	if err != nil {
		return nil, errs.GRPCFromDB(err, core.GetRouteFromCtx(ctx).Name)
	}

	dbMessages := []*models.GPTMessage{
		{Title: "Image prompt", From: "user", Content: req.Message, ChatID: dbGPTChat.ID},
		{Title: "Image", From: "assistant", Content: dalleResponse.URL, ChatID: dbGPTChat.ID},
	}

	for _, msg := range dbMessages {
		if _, err := svc.Clients.DBCreateGPTMessage(ctx, msg); err != nil {
			return nil, errs.GRPCFromDB(err, core.GetRouteFromCtx(ctx).Name)
		}
	}

	return &pbs.NewGPTImageResponse{ImageUrl: dalleResponse.URL, Chat: &pbs.GPTChatInfo{Id: int32(dbGPTChat.ID), Title: dbGPTChat.Title}}, nil
}
