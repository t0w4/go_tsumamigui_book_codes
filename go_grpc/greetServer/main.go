package main

import (
	"context"
	"fmt"
	"go_grpc/proto"
	"io"
	"net"
	"os"

	"google.golang.org/grpc"
)

const address = "localhost:50051"

type server struct{}

func main() {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("listen error: %v\n", err)
		os.Exit(1)
	}

	s := grpc.NewServer()
	greet.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		fmt.Printf("serve error: %v\n", err)
		os.Exit(1)
	}
}

func (s *server) Greet(ctx context.Context, req *greet.UserRequest) (*greet.GreetResponse, error) {
	result := fmt.Sprintf("Hello! I'm %s, %d.", req.User.Name, req.User.Age)
	res := &greet.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (s *server) GreetServerSideStreaming(req *greet.UsersRequest, stream greet.GreetService_GreetServerSideStreamingServer) error {
	users := req.Users
	for _, user := range users {
		result := fmt.Sprintf("Hello! I'm %s, %d.", user.Name, user.Age)
		res := &greet.GreetResponse{
			Result: result,
		}
		err := stream.Send(res)
		if err != nil {
			fmt.Printf("during streaming send error: %v\n", err)
			return err
		}
	}
	return nil
}

func (s *server) GreetClientSideStreaming(stream greet.GreetService_GreetClientSideStreamingServer) error {
	result := "Hello!"
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			result += "."
			return stream.SendAndClose(&greet.GreetResponse{
				Result: result,
			})
		}
		if err != nil {
			fmt.Printf("during streaming receive error: %v\n", err)
			return err
		}

		result += fmt.Sprintf(" %s", req.User.Name)
	}
}

func (s *server) GreetBidirectionalStreaming(stream greet.GreetService_GreetBidirectionalStreamingServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Printf("during streaming receive error: %v\n", err)
			return err
		}
		result := fmt.Sprintf("Hello! I'm %s, %d.", req.User.Name, req.User.Age)

		err = stream.Send(&greet.GreetResponse{
			Result: result,
		})
		if err != nil {
			fmt.Printf("during streaming send error: %v\n", err)
			return err
		}
	}

}
