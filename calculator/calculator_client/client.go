package main

import(
	"fmt"
	"log"
	"context"
	"io"
	"time"

	"app/gRPC-Golang-Exercise/calculator/calculatorpb"


	"google.golang.org/grpc"
)

func main()  {
	fmt.Println("Hello I'm a calculator client")
	cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnary(c)

	// doServerStreaming(c)

	doClientStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient){
	fmt.Println("Starting to do a Unary RPC...")
	req := &calculatorpb.CalculatorRequest{
		Calculating: &calculatorpb.Calculating{
			Number1: 4,
			Number2: 9,
		},
	}
	res, err := c.Calculate(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Calculator: %v", res.Sum)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient){
	fmt.Println("Starting to do a Unary RPC...")
	req := &calculatorpb.PrimeNumberDecompositonRequest{
		PrimeNumberDecompositing: &calculatorpb.PrimeNumberDecompositing{
			Number: 120,
		},
	}
	resStream, err := c.PrimeNumberDecompositon(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			//We've reached the end of the stream
			break;
		}
		if err != nil {
			log.Fatalf("Error occured while stream %v", err)
		}
		log.Printf("Response from PrimeNumberDecompositon %v", msg.GetResult())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient){
	fmt.Println("Starting to do a Client Streaming RPC...")
	
	requests := []*calculatorpb.AverageRequest{
		&calculatorpb.AverageRequest{
			Averaging: &calculatorpb.Averaging{
				Number: 3,
			},
		},
		&calculatorpb.AverageRequest{
			Averaging: &calculatorpb.Averaging{
				Number: 5,
			},
		},
		&calculatorpb.AverageRequest{
			Averaging: &calculatorpb.Averaging{
				Number: 9,
			},
		},
		&calculatorpb.AverageRequest{
			Averaging: &calculatorpb.Averaging{
				Number: 54,
			},
		},
		&calculatorpb.AverageRequest{
			Averaging: &calculatorpb.Averaging{
				Number: 23,
			},
		},
	}

	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalf("error while calling Average: %v", err)
	}

	//we iterate over our slice and send each message individually 
	for _, req := range requests {
		fmt.Printf("Sending request %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from Average: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}