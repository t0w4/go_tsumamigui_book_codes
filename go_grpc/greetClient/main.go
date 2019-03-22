package main

import (
	"context"
	"fmt"
	"go_grpc/proto"
	"io"
	"os"

	"google.golang.org/grpc"
)

const address = "localhost:50051"

func main() {
	con, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("dial error: %v\n", err)
		os.Exit(1)
	}
	defer con.Close()

	cli := greet.NewGreetServiceClient(con)

	callSimple(cli)
	//callServerStreaming(cli)
	//callClientStreaming(cli)
	//callBidirectionalStreaming(cli)
}

func callSimple(cli greet.GreetServiceClient) {
	user := &greet.User{
		Name: "bob",
		Age:  18,
	}
	req := &greet.UserRequest{
		User: user,
	}
	res, err := cli.Greet(context.Background(), req)
	if err != nil {
		fmt.Printf("greet error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("greet response: %v\n", res.Result)
	os.Exit(0)
}

func callServerStreaming(cli greet.GreetServiceClient) {
	user1 := &greet.User{
		Name: "bob",
		Age:  18,
	}
	user2 := &greet.User{
		Name: "alice",
		Age:  12,
	}
	req := &greet.UsersRequest{
		Users: []*greet.User{
			user1,
			user2,
		},
	}

	stream, err := cli.GreetServerSideStreaming(context.Background(), req)
	if err != nil {
		fmt.Printf("greet server side stream create error: %v\n", err)
		os.Exit(1)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("during streaming receive error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("greet server side streaming response: %v\n", res.Result)
	}
	os.Exit(0)
}

func callClientStreaming(cli greet.GreetServiceClient) {
	req1 := &greet.UserRequest{
		User: &greet.User{
			Name: "bob",
			Age:  18,
		},
	}
	req2 := &greet.UserRequest{
		User: &greet.User{
			Name: "alice",
			Age:  12,
		},
	}
	reqs := []*greet.UserRequest{
		req1,
		req2,
	}
	stream, err := cli.GreetClientSideStreaming(context.Background())
	if err != nil {
		fmt.Printf("greet client side stream create error: %v\n", err)
		os.Exit(1)
	}
	for _, req := range reqs {
		err := stream.Send(req)
		if err != nil {
			fmt.Printf("during streaming send error: %v\n", err)
			os.Exit(1)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("close and receive response error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("greet client side streaming response: %v\n", res.Result)
	os.Exit(0)
}

func callBidirectionalStreaming(cli greet.GreetServiceClient) {
	req1 := &greet.UserRequest{
		User: &greet.User{
			Name: "bob",
			Age:  18,
		},
	}
	req2 := &greet.UserRequest{
		User: &greet.User{
			Name: "alice",
			Age:  12,
		},
	}
	reqs := []*greet.UserRequest{
		req1,
		req2,
	}

	stream, err := cli.GreetBidirectionalStreaming(context.Background())
	if err != nil {
		fmt.Printf("greet bidirectional streame create error: %v\n", err)
		os.Exit(1)
	}

	waitc := make(chan struct{})
	go func() {
		for _, req := range reqs {
			err := stream.Send(req)
			if err != nil {
				fmt.Printf("during streaming send error: %v\n", err)
				os.Exit(1)
			}
		}
		err := stream.CloseSend()
		if err != nil {
			fmt.Printf("close send error: %v\n", err)
			os.Exit(1)
		}
	}()
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("during streaming receive error: %v\n", err)
				os.Exit(1)
			}
			res.GetResult()
			fmt.Printf("greet bidirectional streaming response: %v\n", res.Result)
		}
		close(waitc)
	}()

	<-waitc
	os.Exit(0)
}
