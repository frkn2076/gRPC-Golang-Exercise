package main

import(
	"fmt"
	"log"
	"context"
	"io"

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

	doServerStreaming(c)
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