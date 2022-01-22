package sample

import (
	"golang-learning/pb"
)

func NewCPU() *pb.CPUMessage {
	cpu := pb.CPUMessage{
		Brand:         "AMD",
		Name:          "P1",
		NumberThreads: 4,
		NumberCores:   4,
	}
	return &cpu
}
