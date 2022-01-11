package main

import (
	"context"
	"fmt"
	"grpc-go-course/blog/blogpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Junya",
		Title:    "My First Blog",
		Content:  "Content of the first blog.",
	}
	creatingBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v\n", creatingBlogRes.Blog)
	blogID := creatingBlogRes.GetBlog().GetId()

	// read blog
	fmt.Println("Reading the blog")
	if _, err = c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: ""}); err != nil {
		fmt.Printf("error happened while reading: %v\n", err)
	}

	readBlogRes, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blogID})
	if err != nil {
		fmt.Printf("error happened while reading: %v\n", err)
	}
	fmt.Printf("Blog was read: %v\n", readBlogRes.Blog)
}
