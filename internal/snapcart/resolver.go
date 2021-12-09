package snapcart

import "github.com/graphql-go/graphql"

type Resolver interface {
	NewMessage() graphql.FieldResolveFn
	AddMessage() graphql.FieldResolveFn
	Messages() graphql.FieldResolveFn
}
