package snapcart

import (
	"context"
	"snapcart/model"
)

type Service interface {
	ViewMessage(ctx context.Context,messageRequest model.ChatRequest) (model.Message,error)
	AddMessage(ctx context.Context,messageRequest model.ChatRequest) error
	Messages(ctx context.Context,id int) (*[]model.Message,error)
}
