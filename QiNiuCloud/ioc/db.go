package ioc

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InitMongoDB() *mongo.Client {
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
	}
	opts := options.Client().
		ApplyURI("mongodb://127.0.0.1:27017/").
		SetMonitor(monitor)
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}
	return client
}
