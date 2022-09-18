package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/munizigor/grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// AddUser(client)
	// AddUserVerbose(client)
	AddUsers(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "joao@joao.com.br",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC Request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Joao",
		Email: "joao@joao.com.br",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make gRPC Request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}
		fmt.Println("Status", stream.Status, " - ", stream.GetUser())
	}

}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "w1",
			Name:  "Igor",
			Email: "igor.mm@uol.com",
		},
		&pb.User{
			Id:    "w2",
			Name:  "Igor 2",
			Email: "igor2.mm@uol.com",
		},
		&pb.User{
			Id:    "w3",
			Name:  "Igor 3",
			Email: "igor3.mm@uol.com",
		},
		&pb.User{
			Id:    "w4",
			Name:  "Igor 4",
			Email: "igor4.mm@uol.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 2)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error creating response: %v", err)
	}

	fmt.Println("Dados adicionados:\n", res)
}
