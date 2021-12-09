package repository

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"snapcart/model"
	"snapcart/pkg/gorm_errors"
)

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepositoryImpl(db *gorm.DB) *repositoryImpl {
	return &repositoryImpl{
		db:db,
	}
}


func (repositoryI *repositoryImpl) ViewMessage(ctx context.Context,messageRequest model.ChatRequest) (model.Message,error){

	var message model.Message

	if result :=  repositoryI.db.WithContext(ctx).Where("chat LIKE ?", fmt.Sprintf("%s%s%s","%",messageRequest.Chat,"%")).First(&message) ; result.Error != nil {
		return model.Message{},gorm_errors.GormError(result.Error)
	}
	return message,nil
}


func (repositoryI *repositoryImpl)  AddMessage(ctx context.Context,messageRequest model.ChatRequest) error{

	message := model.Message{
		Chat: messageRequest.Chat,
	}

	if result := repositoryI.db.WithContext(ctx).Create(&message); result.Error != nil {
		return gorm_errors.GormError(result.Error)
	}
	return nil
}


func (repositoryI *repositoryImpl)  Messages(ctx context.Context,id uint) (*[]model.Message,error){

	var messages []model.Message

	fmt.Print(id)
	if id == 0 {
		//Get new record
		if result :=  repositoryI.db.WithContext(ctx).Order("id desc").Limit(1).Find(&messages) ; result.Error != nil {
			return  nil,gorm_errors.GormError(result.Error)
		}
	}else{
		if result :=  repositoryI.db.WithContext(ctx).Order("id asc").Find(&messages,"id > ?", id) ; result.Error != nil {
			fmt.Print("Error!!!!!!!!!")
			return  nil,gorm_errors.GormError(result.Error)
		}
	}


	return &messages,nil
}

