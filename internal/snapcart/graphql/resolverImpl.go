package graphql

import (
	"github.com/graphql-go/graphql"
	"log"
	"snapcart/internal/snapcart"
	"snapcart/model"
)

type graphqlResolverImpl struct {
	service snapcart.Service
}

func NewGraphqlResolverImpl(service snapcart.Service) *graphqlResolverImpl {
	return &graphqlResolverImpl{
		service: service,
	}
}

func (resolver *graphqlResolverImpl) ViewMessage() graphql.FieldResolveFn  {
	return func(p graphql.ResolveParams) (interface{}, error) {
		chat := model.ChatRequest{
			Chat: p.Args["message"].(string),
		}
		viewMessage, err := resolver.service.ViewMessage(p.Context, chat)

		if err != nil {
			return model.ChatResponse{
				Status: 404,
			}, nil
		}

		return model.ChatResponse{
			Status: 200,
			ID: viewMessage.ID,
			Chat: viewMessage.Chat,
			CreatedAt: viewMessage.CreatedAt,
			UpdatedAt: viewMessage.UpdatedAt,
			DeletedAt: viewMessage.DeletedAt.Time,
		}, nil
	}


}


func (resolver *graphqlResolverImpl) AddMessage() graphql.FieldResolveFn  {
	return func(p graphql.ResolveParams) (interface{}, error) {

		chatRequest := model.ChatRequest{
			Chat: p.Args["message"].(string),
		}

		err := resolver.service.AddMessage(p.Context, chatRequest)
		if err != nil {
			return model.ChatResponse{
				Status: 404,
			}, nil
		}

		return model.ChatResponse{
			Status: 202,
		}, nil
	}
}


func (resolver *graphqlResolverImpl) Messages() graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {

		messagesData := make(chan interface{})

		go func(id uint) {
				for {
					messages, err := resolver.service.Messages(p.Context, id)
					if err != nil {
						close(messagesData)
					}
					select {
						   case <-p.Context.Done():
							  log.Println("[RootSubscription] [Subscribe] subscription canceled")
							  close(messagesData)
							  return
						   default:

							   for _, value := range *messages {
								///  fmt.Println(value)
								  messagesData <- model.ChatResponse{
									   ID:        value.ID,
									   Chat:      value.Chat,
									   CreatedAt: value.CreatedAt,
									   UpdatedAt: value.UpdatedAt,
									   DeletedAt: value.DeletedAt.Time,
								  }
							   }

					  }

					  if len(*messages) == 0 {
						  id = (*messages)[len(*messages)-1].ID
					  }else{
					  	id = 0
					  }

		  }
		}( uint(p.Args["id"].(int)))

			return messagesData, nil

	}
}



