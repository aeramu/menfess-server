package main

import (
	"context"
	"log"
	"net/http"
	"os"

	authClient "github.com/aeramu/menfess-server/internal/auth/client"
	auth "github.com/aeramu/menfess-server/internal/auth/service"
	"github.com/aeramu/menfess-server/internal/gateway/server"
	postClient "github.com/aeramu/menfess-server/internal/post/client"
	postRepository "github.com/aeramu/menfess-server/internal/post/repository"
	post "github.com/aeramu/menfess-server/internal/post/service"
	userRepository "github.com/aeramu/menfess-server/internal/user/repository"
	user "github.com/aeramu/menfess-server/internal/user/service"
	"github.com/friendsofgo/graphiql"
)

func main() {
	//TODO: catch header, not test header
	ctx := context.WithValue(context.Background(), "request", map[string]string{
		"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYXlsb2FkIjp7IklEIjoiNWZhZTc5ODU4YjlhYmEzNTVmNDQ5MmY3In19.f7XPo_Jj20Lpt3xLv8nkC0NTpemoKaEcO7UsFvcsQ0A",
	})

	userRepo := userRepository.NewRepository()
	userService := user.NewService(userRepo)

	postNotifClient := postClient.NewNotificationClient(userService)
	postRepo := postRepository.NewRepository()
	postService := post.NewService(postRepo, postNotifClient)

	authUserClient := authClient.NewUserClient(userService)
	authNotifClient := authClient.NewNotificationClient(userService)
	authService := auth.NewService(authUserClient, authNotifClient)

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
