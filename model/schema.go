package model

import "github.com/graphql-go/graphql"

var ResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Response",
	Fields: graphql.Fields{
		"status": &graphql.Field{
			Type: graphql.Int,
		},
		"chat": &graphql.Field{
			Type: graphql.String,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"updated_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"deleted_at": &graphql.Field{
			Type: graphql.DateTime,
		},

	},
})

var ChatType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Chat",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

