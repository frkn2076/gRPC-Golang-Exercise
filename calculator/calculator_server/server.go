package main

import (
	"context"
	"log"
	"net"
	"time"

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

func (*server) PrimeNumberDecompositon(req *calculatorpb.PrimeNumberDecompositonRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositonServer) error {
	number := req.GetPrimeNumberDecompositing().GetNumber()
	for {
		for i := int32(2); i <= number; i++ {
			if number % i == 0 {
				number = number / i

				res := &calculatorpb.PrimeNumberDecompositonResponse{
					Result: i,
				}
				stream.Send(res)
				time.Sleep(1000 * time.Millisecond)
				break;
			}
		}
		if number == 1 {
			break;
		}
	}
	return nil
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