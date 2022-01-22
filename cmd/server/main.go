package main

import (
	"fmt"
	"golang-learning/pb"
	"golang-learning/service"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	grpcServer := grpc.NewServer()
	newImageStore := service.NewDiskImageStore("/Users/mayankparmar/moj/golang-learning/images")
	laptopService := service.NewLaptopServer(newImageStore)

	pb.RegisterLaptopCPUServiceServer(grpcServer, laptopService)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		fmt.Println("Error while starting grpc server", err.Error())
	}
	fmt.Println("Listener initialized")
	c := make(chan bool)
	go func() {
		err = grpcServer.Serve(listener)
		if err != nil {
			fmt.Println("Error while starting grpc server", err.Error())
		}
		c <- true
	}()
	fmt.Println("GRPC Server started at port 5000")
	<-c
}
