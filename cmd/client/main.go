package main

import (
	"context"
	"flag"
	"fmt"
	"golang-learning/pb"
	"io"
	"time"

	"google.golang.org/grpc"
)

func main() {
	serverAddr := flag.String("address", "localhost:5000", "localhost:5000")
	flag.Parse()
	fmt.Println("Dial server", *serverAddr)
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Error while grpc client", err.Error())
	}
	laptopClient := pb.NewLaptopCPUServiceClient(conn)
	reqMsg := &pb.NewLaptopCPU{PurchaseToken: "p1"}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	res, err := laptopClient.GetCPU(ctx, reqMsg)
	if err != nil {
		fmt.Println("error while getting CPU", err.Error())
	} else {
		fmt.Println("CPU Brand is", res.Brand)
	}

	res, err = laptopClient.GetCPU(ctx, &pb.NewLaptopCPU{}) //Emulating bad Request
	if err != nil {
		fmt.Println("error while getting CPU", err.Error())
	} else {
		fmt.Println("CPU Brand is", res.Brand)
	}

	//Server streaming below

	stream, err := laptopClient.GetLaptopStreamService(ctx, &pb.Empty{})
	if err != nil {
		fmt.Println("error while getting streaming CPU", err.Error())
	} else {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("streaming response finished", err.Error())
				return
			}
			if err != nil {
				fmt.Println("Cannot receive response", err.Error())
				return
			}
			laptop := res.GetCPUMessage()
			fmt.Println("Got cpu of brand", laptop.Brand)
		}
	}
}
