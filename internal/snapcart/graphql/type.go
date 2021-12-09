package graphql

import (
	"github.com/graphql-go/graphql"
	"snapcart/model"
)

type graphqlType struct {
	resolverImpl *graphqlResolverImpl
}

func NewGraphqlType(resolverImpl *graphqlResolverImpl) *graphqlType {
	return &graphqlType{
		resolverImpl: resolverImpl,
	}
}

func (g *graphqlType) Query()  *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"message": &graphql.Field{
					Type: model.ResponseType,
					Description: "New messages",
					Args: graphql.FieldConfigArgument{
						"message": 	&graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: g.resolverImpl.ViewMessage(),
				},
			},
		})
}

func (g *graphqlType) Mutation()  *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"message": &graphql.Field{
					Type: model.ResponseType,
					Description: "Add new message",
					Args: graphql.FieldConfigArgument{
						"message": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
					},
					Resolve: g.resolverImpl.AddMessage(),
				},
			},
		})
}

func (g *graphqlType) Subscription() *graphql.Object {
	return graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Subscription",
			Fields: graphql.Fields{
				"messages": &graphql.Field{
					Type: model.ResponseType,
					Description: "Read recorded messages",
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.Int),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return p.Source,nil
					},
					Subscribe: g.resolverImpl.Messages(),
				},
			},
		})
}