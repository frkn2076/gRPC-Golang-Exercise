package main

import(
	"context"
	"fmt"
	// "io"
	"log"
	// "time"

	"app/gRPC-Golang-Exercise/blog/blogpb"

	"google.golang.org/grpc"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/credentials"
	// "google.golang.org/grpc/status"
)

func main() {

	fmt.Println("Blog Client")
	
	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect : %v", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	//Create blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Furkan",
		Title: "My First Blog",
		Content: "Content of the first blog",
	}
	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: &v", err)
		return
	}
	fmt.Printf("Blog has created: &v", res)
	blogId := res.GetBlog().GetId()



	//read blog
	fmt.Println("Reading the blog")
	
	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId:  "asdasd",
	})
	if err2 != nil {
		fmt.Printf("Error happened while reading: %v \n", err2)
	}

	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId:  blogId,
	})
	if readBlogErr != nil {
		fmt.Printf("Error happened while reading: %v \n", readBlogErr)
		return
	}

	fmt.Printf("Blog was read: %v", readBlogRes)
}