package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"app/gRPC-Golang-Exercise/blog/blogpb"

	"gopkg.in/mgo.v2/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type server struct{}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {

	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}

	objectID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot convert to OID error: %v", err))
	}

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id: objectID.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title: blog.GetTitle(),
			Content: blog.GetContent(),
		},
	}, nil
}

func (*server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {

	blogID := req.GetBlogId()

	objectID, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Cannot parse ID: %v", err))
	}

	//create an empty struct
	data := &blogItem{}
	filter := bson.M{"_id": objectID}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Cannot find blog with specified ID: %v", err))
	}

	return &blogpb.ReadBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil

}

func (*server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {

	blog  := req.GetBlog()

	iod, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Cannot parse ID: %v", err))
	}

	//create an empty struct
	data := &blogItem{}
	filter := bson.M{"_id": iod}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Cannot find blog with specified ID: %v", err))
	}

	//we update our internal struct
	data.AuthorID = blog.GetAuthorId()
	data.Content = blog.GetContent()
	data.Title = blog.GetTitle()

	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot update object in MongoDB: %v", err))
	}

	return &blogpb.UpdateBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil

}

func (*server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {

	blogID := req.GetBlogId()

	iod, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Cannot parse ID: %v", err))
	}

	filter := bson.M{"_id": iod}

	res, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot delete object in MongoDB: %v", err))
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Cannot find blog in MongoDB: %v", err))
	}

	return &blogpb.DeleteBlogResponse{
		BlogId: req.GetBlogId(),
	}, nil

}

func dataToBlogPb(data *blogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id: data.ID.Hex(),
		AuthorId: data.AuthorID,
		Content: data.Content,
		Title: data.Title,
	}
}

func (*server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {

	cur, err := collection.Find(context.Background(), primitive.D{{}}) //filter is nil -> primitive.D{{}}
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		data := &blogItem{}
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("Error while decoding data from MongoDB: %v", err))
		}
		stream.Send(&blogpb.ListBlogResponse{Blog: dataToBlogPb(data)})
	}

	if err := cur.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}

	return nil
}

func main() {
	//if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Connecting to MongoDB")

	//connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("mydb").Collection("blog")

	fmt.Println("Blog Service started")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {
		fmt.Println("Starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to server: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//Block until a signal is receieved
	<-ch
	fmt.Println("Stopping the server...")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("Closing MongoDB connection")
	client.Disconnect(context.TODO())
	fmt.Println("End of program")

}
