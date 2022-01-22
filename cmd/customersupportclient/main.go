package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"golang-learning/pb"
	"log"
	"os"
	"strings"

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
	stream, err := laptopClient.CustomerSupportService(context.Background())
	waitResponse := make(chan error)
	waitforReceiver := make(chan bool)
	go func() {
		responseMsg, err := stream.Recv()
		if err != nil {
			waitforReceiver <- true
			log.Println("Error in receiving message", err.Error())
			waitResponse <- err
		} else {
			fmt.Println("Message from Customer Support :", responseMsg.GetMessage())
		}
		for true {
			responseMsg, err := stream.Recv()
			if err != nil {
				log.Println("Error in receiving message", err.Error())
				waitResponse <- err
				break
			} else {
				fmt.Println("Message from Customer Support :", responseMsg.GetMessage())
			}
		}
	}()
	for true {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")
		clientResponse := &pb.CustomerSupportRequest{Message: text}
		if text == "EOF" {
			break
		}
		stream.Send(clientResponse)
	}
	<-waitResponse
}
