package server

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	graphqlInternal "snapcart/internal/snapcart/graphql"
	"snapcart/internal/snapcart/repository"
	"snapcart/internal/snapcart/service"
	"sync"
)

type server struct {
	db *gorm.DB
}

func NewServer(db *gorm.DB) *server {
	return &server{
		db: db,
	}
}

var schemaQL graphql.Schema

func (s server) Run()  error {


	repo := repository.NewRepositoryImpl(s.db)
	serviceImpl := service.NewServiceImpl(repo)

	graphqlResolver := graphqlInternal.NewGraphqlResolverImpl(serviceImpl)
	graphqlType := graphqlInternal.NewGraphqlType(graphqlResolver)

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphqlType.Query(),
		Mutation: graphqlType.Mutation(),
		Subscription: graphqlType.Subscription(),
	})
	if err != nil {
		return 	errors.Wrap(err,"Server: Run")
	}

	schemaQL= schema

	h := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
		GraphiQL: false,
		Playground: true,
	})

	http.Handle("/snap", h)
	http.HandleFunc("/subscriptions", SubscriptionsHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ginRunErr := make(chan error)
	go func() {
		httpErr := http.ListenAndServe(":8080", nil)
		if httpErr != nil {
			ginRunErr <- errors.Wrap(httpErr,"")
		}
	}()


	if ginRunErr!=nil{
		msg := <- ginRunErr
		return msg
	}

	return nil
}


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"graphql-ws"},
}

type ConnectionACKMessage struct {
	OperationID string `json:"id,omitempty"`
	Type        string `json:"type"`
	Payload     struct {
		Query string `json:"query"`
	} `json:"payload,omitempty"`
}

func SubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to do websocket upgrade: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	connectionACK, err := json.Marshal(map[string]string{
		"type": "connection_ack",
	})
	if err != nil {
		log.Printf("failed to marshal ws connection ack: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, connectionACK); err != nil {
		log.Printf("failed to write to ws connection: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go handleSubscription(conn)
}

func handleSubscription(conn *websocket.Conn) {
	var subscriber *Subscriber
	subscriptionCtx, subscriptionCancelFn := context.WithCancel(context.Background())

	handleClosedConnection := func() {
		log.Println("[SubscriptionsHandler] subscriber closed connection")
		unsubscribe(subscriptionCancelFn, subscriber)
		return
	}

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("failed to read websocket message: %v", err)
			return
		}

		var msg ConnectionACKMessage
		if err := json.Unmarshal(p, &msg); err != nil {
			log.Printf("failed to unmarshal websocket message: %v", err)
			continue
		}

		if msg.Type == "stop" {
			handleClosedConnection()
			return
		}

		if msg.Type == "start" {
			subscriber = subscribe(subscriptionCtx, subscriptionCancelFn, conn, msg)
		}
	}
}

type Subscriber struct {
	UUID          string
	Conn          *websocket.Conn
	RequestString string
	OperationID   string
}

var subscribers sync.Map

func subscribersSize() uint64 {
	var size uint64
	subscribers.Range(func(_, _ interface{}) bool {
		size++
		return true
	})
	return size
}

func unsubscribe(subscriptionCancelFn context.CancelFunc, subscriber *Subscriber) {
	subscriptionCancelFn()
	if subscriber != nil {
		subscriber.Conn.Close()
		subscribers.Delete(subscriber.UUID)
	}
	log.Printf("[SubscriptionsHandler] subscribers size: %+v", subscribersSize())
}

func subscribe(ctx context.Context, subscriptionCancelFn context.CancelFunc, conn *websocket.Conn, msg ConnectionACKMessage) *Subscriber {
	subscriber := &Subscriber{
		UUID:          uuid.New().String(),
		Conn:          conn,
		RequestString: msg.Payload.Query,
		OperationID:   msg.OperationID,
	}
	subscribers.Store(subscriber.UUID, &subscriber)

	log.Printf("[SubscriptionsHandler] subscribers size: %+v", subscribersSize())

	sendMessage := func(r *graphql.Result) error {
		message, err := json.Marshal(map[string]interface{}{
			"type":    "data",
			"id":      subscriber.OperationID,
			"payload": r.Data,
		})
		if err != nil {
			return err
		}

		if err := subscriber.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return err
		}

		return nil
	}

	go func() {
		subscribeParams := graphql.Params{
			Context:       ctx,
			RequestString: msg.Payload.Query,
			Schema:        schemaQL,
		}

		subscribeChannel := graphql.Subscribe(subscribeParams)

		for {
			select {
			case <-ctx.Done():
				log.Printf("[SubscriptionsHandler] subscription ctx done")
				return
			case r, isOpen := <-subscribeChannel:
				if !isOpen {
					log.Printf("[SubscriptionsHandler] subscription channel closed")
					unsubscribe(subscriptionCancelFn, subscriber)
					return
				}
				if err := sendMessage(r); err != nil {
					if err == websocket.ErrCloseSent {
						unsubscribe(subscriptionCancelFn, subscriber)
					}
					log.Printf("failed to send message: %v", err)
				}
			}
		}
	}()

	return subscriber
}