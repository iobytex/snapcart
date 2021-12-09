package snapcart

import (
	"context"
	"snapcart/model"
)

type Repository interface {
	ViewMessage(ctx context.Context,messageRequest model.ChatRequest) (model.Message,error)
	AddMessage(ctx context.Context,messageRequest model.ChatRequest) error
	Messages(ctx context.Context,id uint) (*[]model.Message,error)
}
