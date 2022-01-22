package service

import (
	"bytes"
	"context"
	"fmt"
	"golang-learning/pb"
	"golang-learning/sample"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const MaxImageSize = 1 << 20 //1 MB
type LaptopServer struct {
	imageStore ImageStore
}

func NewLaptopServer(imageStore ImageStore) *LaptopServer {
	return &LaptopServer{imageStore: imageStore}
}

func (server *LaptopServer) GetCPU(ctx context.Context, request *pb.NewLaptopCPU) (*pb.CPUMessage, error) {
	purchase_token := request.GetPurchaseToken()
	// time.Sleep(6 * time.Second) uncomment to check how a service behaves in case load is heavy
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Deadline exceeded")
		return nil, status.Error(codes.DeadlineExceeded, "Deadline exceeded")
	}
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.DeadlineExceeded, "Context Cancelled")
	}
	if purchase_token != "" {
		return sample.NewCPU(), nil
	} else {
		return nil, status.Error(codes.Internal, "Invalid argument")
	}
}

func (server *LaptopServer) GetLaptopStreamService(req *pb.Empty, stream pb.LaptopCPUService_GetLaptopStreamServiceServer) error {
	// time.Sleep(6 * time.Second) uncomment to check how a service behaves in case load is heavy
	for i := 0; i < 10; i++ {
		if stream.Context().Err() == context.DeadlineExceeded {
			fmt.Println("Deadline exceeded")
			return status.Error(codes.DeadlineExceeded, "Deadline exceeded")
		}
		if stream.Context().Err() == context.Canceled {
			return status.Error(codes.DeadlineExceeded, "Context Cancelled")
		}
		time.Sleep(1 * time.Second) //uncomment to emulate load
		err := stream.Send(&pb.GetLaptopStream{CPUMessage: sample.NewCPU()})
		if err != nil {
			fmt.Println("Error in send ", err.Error())
		}
	}
	return nil
}

func (server *LaptopServer) UploadImageService(stream pb.LaptopCPUService_UploadImageServiceServer) error {
	req, err := stream.Recv()
	ctx := stream.Context()
	if err != nil {
		log.Println("Cannot receive image info", err.Error())
		return status.Error(codes.Unknown, "cannot receive image info")
	}

	laptopId := req.GetImageData().LaptopId
	imageType := req.GetImageData().ImageType

	log.Printf("Received laptop with Id %s and image type %s\n", laptopId, imageType)

	imageData := bytes.Buffer{}
	imageSize := 0
	for {
		log.Println("Waiting to receive more data")
		req, err := stream.Recv()
		// time.Sleep(10 * time.Second) Un comment to check context deadline flow
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Deadline exceeded")
			return status.Error(codes.DeadlineExceeded, "Deadline exceeded")
		}
		if err != io.EOF {
			log.Println("No more data")
			break
		}
		if err != nil {
			log.Println("cannot receive the chunked data ", err.Error())
			return status.Error(codes.Unknown, "cannot receive the chunked data ")
		}
		chunk := req.GetChunkData()
		size := len(chunk)

		imageSize += size
		if imageSize > MaxImageSize {
			return status.Error(codes.InvalidArgument, "Image size is too large")
		}

		_, err = imageData.Write(chunk)

		if err != nil {
			log.Println("cannot write the chunked data ", err.Error())
			return status.Error(codes.Internal, "cannot write chunk data")
		}
	}
	imageId, err := server.imageStore.Save(laptopId, imageType, imageData)
	if err != nil {
		log.Println("cannot save the image data", err.Error())
		return status.Error(codes.Internal, "cannot save the image data")
	}
	res := &pb.UploadImageResponse{
		Id:   imageId,
		Size: uint32(imageSize),
	}
	err = stream.SendAndClose(res)
	if err != nil {
		log.Println("cannot close the stream", err.Error())
		return status.Error(codes.Internal, "cannot close the stream")
	}
	return nil
}
