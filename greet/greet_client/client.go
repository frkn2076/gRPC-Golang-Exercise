package main

import(
	"fmt"
	"log"
	"context"

	"app/gRPC-Golang-Exercise/greet/greetpb"


	"google.golang.org/grpc"
)

func main()  {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client &f:", c)

	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient){
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Furkan",
			LastName: "Öztürk",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v",res.Result)
}