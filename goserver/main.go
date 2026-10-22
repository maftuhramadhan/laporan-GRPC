package main

import (
	"context"
	"log"
	"net"

	"goserver/pb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) GetUserProfile(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	log.Printf("Menerima permintaan untuk user ID: %d", req.GetUserId())

	return &pb.UserResponse{
		Name:  "Maftuh",
		Email: "maftuh@example.com",
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Gagal membuka port: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &server{})

	log.Println("gRPC server berjalan pada port 50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
