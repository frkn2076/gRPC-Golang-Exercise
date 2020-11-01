package main

import (
	"context"
	"log"
	"net"
	"time"
	"io"

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

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	sum := int32(0)
	i := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(i)
			//we have finished reading the client stream
			return stream.SendAndClose(&calculatorpb.AverageResponse{
				Res: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v" ,err)
		}

		number := req.GetAveraging().GetNumber()
		sum += number
		i++
	}
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