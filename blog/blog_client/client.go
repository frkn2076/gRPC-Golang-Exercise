package main

import(
	"context"
	"fmt"
	"io"
	"log"

	"app/gRPC-Golang-Exercise/blog/blogpb"

	"google.golang.org/grpc"
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

	fmt.Printf("Blog was read: %v \n", readBlogRes)


	//update blog
	newBlog := &blogpb.Blog{
		Id: blogId,
		AuthorId: "Changed Author",
		Title: "Changed Title",
		Content: "Changed Content",
	}

	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil{
		fmt.Printf("Error happened while updating: %v \n", readBlogErr)
		return
	}
	fmt.Printf("Blog was updated: %v \n", updateRes)


	//delete blog

	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogId})
	if deleteErr != nil{
		fmt.Printf("Error happened while deleting: %v \n", readBlogErr)
		return
	}
	fmt.Printf("Blog was deleted: %v \n", deleteRes)



	//list blogs
	stream, streamErr := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})

	if streamErr != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", streamErr)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			//We've reached the end of the stream
			break;
		}
		if err != nil {
			log.Fatalf("Error occured while stream %v", err)
		}
		fmt.Println(res.GetBlog())
	}
}