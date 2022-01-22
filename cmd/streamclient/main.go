package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"golang-learning/pb"
	"io"
	"log"
	"os"
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

	file, err := os.Open("/Users/mayankparmar/Downloads/test.png")
	if err != nil {
		log.Fatalf("Cannot open the image %s")
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := laptopClient.UploadImageService(ctx)

	if err != nil {
		log.Fatal("Cannot upload the image", err.Error())
	}

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_ImageData{
			ImageData: &pb.ImageInfo{
				LaptopId:  "TestLaptopId",
				ImageType: ".png",
			},
		},
	}

	err = stream.Send(req)

	if err != nil {
		log.Fatal("Cannot send the file", err.Error())
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Error in reading file", err)
		}
		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}
		err = stream.Send(req)
		if err != nil {
			errStream := stream.RecvMsg(nil)
			log.Fatal("Error in sending data on stream", errStream)
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("Error in uploading file file", err)
	}
	log.Printf("Image successfully upload to server with the Id %s  and size %s \n", res.GetId(), res.GetSize())

}
