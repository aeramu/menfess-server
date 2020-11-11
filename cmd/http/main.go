package main

import (
	"context"
	"log"
	"net/http"
	"os"

	authClient "github.com/aeramu/menfess-server/auth/client"
	auth "github.com/aeramu/menfess-server/auth/service"
	"github.com/aeramu/menfess-server/gateway/server"
	postClient "github.com/aeramu/menfess-server/post/client"
	postRepository "github.com/aeramu/menfess-server/post/repository"
	post "github.com/aeramu/menfess-server/post/service"
	userRepository "github.com/aeramu/menfess-server/user/repository"
	user "github.com/aeramu/menfess-server/user/service"
	"github.com/friendsofgo/graphiql"
)

func main() {
	//TODO: catch header, not test header
	ctx := context.WithValue(context.Background(), "request", map[string]string{
		"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYXlsb2FkIjp7IklEIjoiNWY5ZTkyNGYzZTU4N2E1ZTc1OTYyZWU2In19.9nxUGhIrhk_reePdqtg_hspD1ab64PX6gmaZPtodmwU",
	})

	userRepo := userRepository.NewRepository()
	userService := user.NewService(userRepo)

	notifClient := postClient.NewNotificationClient(userService)
	postRepo := postRepository.NewRepository()
	postService := post.NewService(postRepo, notifClient)

	userClient := authClient.NewUserClient(userService)
	authService := auth.NewService(userClient)

	handler := server.NewServer(ctx, postService, authService, userService)
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
