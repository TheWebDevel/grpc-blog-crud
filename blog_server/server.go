package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/thewebdevel/grpc-blog/blogpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type server struct{}

// Blog item structure
// Use bson to map our struct elements with mongodb attributes
type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

// Make the MongoDB collection globally available
var collection *mongo.Collection

func main() {
	// If we get the crash code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Connect with MongoDB
	fmt.Println("Connecting to MongoDB...")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Use the collection
	collection = client.Database("mydb").Collection("blog")

	// Start Service
	fmt.Println("Blog service started...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for ctrl c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until the signal is received
	<-ch
	fmt.Println("Stopping the server...")
	s.Stop()
	fmt.Println("Closing the listener...")
	lis.Close()
	fmt.Println("Closing MongoDB connection...")
	client.Disconnect(context.TODO())
	fmt.Println("End of Program")
}
