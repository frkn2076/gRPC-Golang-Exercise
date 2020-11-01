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

	// doClientStreaming(c)

	doBiDirectionalStreaming(c)
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

func doBiDirectionalStreaming(c calculatorpb.CalculatorServiceClient){

	fmt.Println("Starting to do a Bi-Directional Streaming RPC...")

	//we create a stream by invoking the client
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*calculatorpb.FindMaximumRequest{
		&calculatorpb.FindMaximumRequest{
			Number: 1,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 5,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 3,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 6,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 2,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 20,
		},
	}

	waitChannel := make(chan struct{})

	//we send a bunch of messages to the client (go routine)
	go func() {
		//function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	//we receive a bunch of messages from the client (go routine)
	go func() {
		//function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break;
			}
			if err != nil {
				log.Fatalf("Error while receiving: &v", err)
				break;
			}
			fmt.Printf("Received &v\n", res.GetResult())
		}
		close(waitChannel)
	}()

	//block until everything is done
	//waits channel to be closed. If you removed the <- line , the program would exit before the go func even started.
	<-waitChannel

}