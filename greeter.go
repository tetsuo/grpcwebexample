package main

import (
	"context"
	"fmt"
	"strings"

	"example.com/grpcwebexample/internal/hellopb"
)

type greeterService struct {
	hellopb.UnimplementedGreeterServer
}

func newGreeterService() *greeterService {
	return &greeterService{}
}

func (g *greeterService) SayHello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloReply, error) {
	return &hellopb.HelloReply{Message: g.greeting(req.GetName())}, nil
}

func (g *greeterService) greeting(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "world"
	}
	return fmt.Sprintf("Hello, %s!", name)
}
