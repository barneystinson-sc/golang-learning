package main

import (
	"fmt"
	"golang-learning/pb"
	"golang-learning/sample"
	"golang-learning/serializer"
)

func main() {
	processor1 := sample.NewCPU()
	serializer.WriteProtobufToFile(processor1, "processor.txt")
	processor2 := pb.CPUMessage{}
	serializer.ReadBinaryToProtobuf("processor.txt", &processor2)
	fmt.Println(processor2.Brand)
	serializer.WriteProtobufToJson(processor1, "processor.json")
}
