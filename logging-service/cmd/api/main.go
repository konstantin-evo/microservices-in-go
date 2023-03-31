package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
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
		log.Panicf("Failed to establish connection to mongoDB. URL: %s. Error: %v", app.MongoURL, err)
	}
	app.Models = data.New(mongoClient)

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			log.Panicf("Failed to close mongoDB connection. Error: %v", err)
		}
	}()

	// Register RPC & gPRC Server
	rpcServer := &RPCServer{Client: mongoClient}
	err = rpc.Register(rpcServer)
	go app.RPCListen()
	go app.gRPCListen()

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

func (app *Config) RPCListen() error {
	log.Println("Starting RPC server on port ", app.RPCPort)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", app.RPCPort))
	if err != nil {
		log.Printf("Failed to establish TCP connection for RPC on port: %s. Error: %v", app.RPCPort, err)
		return err
	}
	defer listener.Close()

	for {
		rpcConn, err := listener.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(rpcConn)
	}
}

func (app *Config) gRPCListen() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", app.GRPCPort))
	if err != nil {
		log.Panicf("Failed to establish TCP connection for gRPC on port: %s. Error: %v", app.GRPCPort, err)
	}
	defer listener.Close()

	gRPCServer := grpc.NewServer()
	logs.RegisterLogServiceServer(gRPCServer, &LogServer{Models: app.Models})

	if err := gRPCServer.Serve(listener); err != nil {
		log.Panicf("Failed to listen gRPC: %v", err)
	}
	log.Printf("gRPCServer starting on port: %s", app.GRPCPort)

	return nil
}

func connectToMongo(mongoURL string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to mongo!")

	return conn, nil
}
