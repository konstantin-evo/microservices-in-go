package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"log-service/data"
	"time"
)

type RPCServer struct {
	Client *mongo.Client
}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := r.Client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	*resp = "Processed payload via RPC"
	return nil
}
