package main

import (
	"context"
	"io"
	"log"
	"math"
	"net"
	"time"
	"fmt"

	"app/gRPC-Golang-Exercise/calculator/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (*server) Calculate(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatingResponse, error) {
	num1 := req.GetCalculating().GetNumber1()
	num2 := req.GetCalculating().GetNumber2()
	result := num1 + num2
	res := &calculatorpb.CalculatingResponse{
		Sum: result,
	}
	return res, nil
}

func (*server) PrimeNumberDecompositon(req *calculatorpb.PrimeNumberDecompositonRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositonServer) error {
	number := req.GetPrimeNumberDecompositing().GetNumber()
	for {
		for i := int32(2); i <= number; i++ {
			if number%i == 0 {
				number = number / i

				res := &calculatorpb.PrimeNumberDecompositonResponse{
					Result: i,
				}
				stream.Send(res)
				time.Sleep(1000 * time.Millisecond)
				break
			}
		}
		if number == 1 {
			break
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
			log.Fatalf("Error while reading client stream %v", err)
		}

		number := req.GetAveraging().GetNumber()
		sum += number
		i++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	var max int32 = math.MinInt32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		number := req.GetNumber()
		if number > max {
			max = number
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Result: number,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", sendErr)
				return sendErr
			}
		}

	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Received a negative  number: %v", number),)
	}

	return &calculatorpb.SquareRootResponse {
		SquareRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
