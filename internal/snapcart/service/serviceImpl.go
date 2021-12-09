package service

import (
	"context"
	"snapcart/internal/snapcart"
	"snapcart/model"
)

type serviceImpl struct {
	repo snapcart.Repository
}

func NewServiceImpl(repo snapcart.Repository) *serviceImpl {
	return &serviceImpl{
		repo: repo,
	}
}



func (serviceI *serviceImpl) ViewMessage(ctx context.Context,messageRequest model.ChatRequest) (model.Message,error){
	message, err := serviceI.repo.ViewMessage(ctx,messageRequest)
	if err != nil {
		return model.Message{}, err
	}
	return message,nil
}


func (serviceI *serviceImpl)  AddMessage(ctx context.Context,messageRequest model.ChatRequest) error{
	err := serviceI.repo.AddMessage(ctx, messageRequest)
	if err != nil {
		return err
	}
	return nil
}


func (serviceI *serviceImpl)  Messages(ctx context.Context,id int) (*[]model.Message,error){
	messages, err := serviceI.repo.Messages(ctx, id)
	if err != nil {
		return nil, err
	}
	return messages,nil
}
