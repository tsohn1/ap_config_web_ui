package main

import (
	"net"
	"google.golang.org/grpc/reflection"
	"log"
	"google.golang.org/grpc"
	"context"
	validate "ap_config_web_ui/validate"
)

const (
	YAML_FOLDER = "config_files/"
	NETWORK_ENV = YAML_FOLDER + "network.yaml"
	VALIDATE_YAML_CHANGES = false
	GRPC_SUCCESS_TOKEN_NETWORK = 1
	GRPC_SUCCESS_TOKEN_OPERATION = 2
	GRPC_FAIL_TOKEN = 0
)


type ValidateServer struct {
	validate.UnimplementedValidateServer
}

func (s *ValidateServer) Verify(ctx context.Context, req *validate.VerifyRequest) (*validate.VerifyResponse, error) {
	if req.Token == GRPC_SUCCESS_TOKEN_NETWORK || req.Token == GRPC_SUCCESS_TOKEN_OPERATION {
		return &validate.VerifyResponse{IsValid: true}, nil
	} else{
		return &validate.VerifyResponse{IsValid: false}, nil
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)

	}

	s := grpc.NewServer()
	validate.RegisterValidateServer(s, &ValidateServer{})
	reflection.Register(s)

	if err:= s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}