package main

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	protoc "github.com/k-shir0/connect-go-example/pkg/gen/api/grpc/v1"
	"github.com/k-shir0/connect-go-example/pkg/gen/api/grpc/v1/examplev1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

type ExampleServer struct{}

var tasks []*protoc.Task

func (e ExampleServer) CreateTask(_ context.Context, req *connect.Request[protoc.CreateTaskRequest]) (*connect.Response[protoc.CreateTaskResponse], error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	tasks = append(tasks, &protoc.Task{
		// id is uuid
		Id:          id.String(),
		Name:        req.Msg.Name,
		Description: req.Msg.Description,
	})

	res := connect.NewResponse(&protoc.CreateTaskResponse{
		Id: id.String(),
	})

	return res, nil
}

func (e ExampleServer) ReadAllTask(_ context.Context, _ *connect.Request[emptypb.Empty]) (*connect.Response[protoc.ReadAllTaskResponse], error) {
	res := connect.NewResponse(&protoc.ReadAllTaskResponse{
		Tasks: tasks,
	})

	return res, nil
}

func main() {
	example := &ExampleServer{}
	mux := http.NewServeMux()
	path, handler := examplev1connect.NewExampleServiceHandler(example)
	mux.Handle(path, handler)
	http.ListenAndServe(":8080", h2c.NewHandler(mux, &http2.Server{}))
}
