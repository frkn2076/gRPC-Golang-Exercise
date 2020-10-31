package main

import (
	"context"
	"log"
	"net"

	"app/gRPC-Golang-Exercise/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Calculate(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatingResponse, error){
	num1 := req.GetCalculating().GetNumber1()
	num2 := req.GetCalculating().GetNumber2()
	result := num1 + num2
	res := &calculatorpb.CalculatingResponse{
		Sum: result,
	}
	return res,nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}