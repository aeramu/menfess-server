package main

import (
	"context"
	"github.com/aeramu/menfess-server/gateway/server"
	postRepo "github.com/aeramu/menfess-server/post/repository"
	roomRepo "github.com/aeramu/menfess-server/room/repository"
	"log"
	"net/http"
	"os"

	post "github.com/aeramu/menfess-server/post/service"
	room "github.com/aeramu/menfess-server/room/service"
	"github.com/friendsofgo/graphiql"
)

func main() {
	//TODO: catch header, not test header
	ctx := context.WithValue(context.Background(), "request", map[string]string{
		"id": "5ef89baaec8ff2af8b9934c1",
	})

	postRepo := postRepo.NewRepository()
	if postRepo == nil{
		return
	}
	postService := post.NewService(postRepo)

	roomRepo := roomRepo.NewRepository()
	if roomRepo == nil{
		return
	}
	roomService := room.NewService(roomRepo)

	handler := server.NewServer(ctx, postService, roomService)
	http.Handle("/", handler)

	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/")
	if err != nil {
		panic(err)
	}
	http.Handle("/graphiql", graphiqlHandler)

	port := getPort()
	log.Println("Server ready at " + port)
	log.Println("Graphiql ready at " + port + "/graphiql")
	log.Fatal(http.ListenAndServe(port, nil))
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	return ":" + port
}
