package main

import (
	"context"
	"fmt"
	"log-service/data"
	"log-service/logs"
)

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (logServer *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	//write a log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := logServer.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: fmt.Sprintf("Failed. Details: %s", err),
		}
		return res, err
	}

	res := &logs.LogResponse{
		Result: "Processed payload via gRPC",
	}

	return res, nil
}
