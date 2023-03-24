package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Models   data.Models
	WebPort  string
	RPCPort  string
	GRPCPort string
	MongoURL string
}

func main() {
	app := Config{
		WebPort:  "80",
		RPCPort:  "5001",
		GRPCPort: "50001",
		MongoURL: "mongodb://mongo:27017",
	}

	// connect to mongo
	mongoClient, err := connectToMongo(app.MongoURL)
	if err != nil {
		log.Panic(err)
	}
	app.Models = data.New(mongoClient)

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Register RPC Server
	rpcServer := &RPCServer{Client: mongoClient}
	err = rpc.Register(rpcServer)
	go app.rpcListen()

	// start web server
	log.Println("Starting service on port", app.WebPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.WebPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port ", app.RPCPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", app.RPCPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo(mongoURL string) (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	connect, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return connect, nil
}
