package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aeramu/menfess-server/pkg/playground"

	"github.com/aeramu/menfess-server/internal/gateway/handler/graphql"

	"github.com/aeramu/menfess-server/internal/database/mongodb"
	"go.mongodb.org/mongo-driver/mongo"

	authNotifClient "github.com/aeramu/menfess-server/internal/auth/client/notification"
	authUserClient "github.com/aeramu/menfess-server/internal/auth/client/user"
	authRepo "github.com/aeramu/menfess-server/internal/auth/repository/mongodb"
	auth "github.com/aeramu/menfess-server/internal/auth/service"
	notifExpoClient "github.com/aeramu/menfess-server/internal/notification/client/expo"
	notifRepo "github.com/aeramu/menfess-server/internal/notification/repository/mongodb"
	notif "github.com/aeramu/menfess-server/internal/notification/service"
	postNotifClient "github.com/aeramu/menfess-server/internal/post/client/notification"
	postRepo "github.com/aeramu/menfess-server/internal/post/repository/mongodb"
	post "github.com/aeramu/menfess-server/internal/post/service"
	userRepo "github.com/aeramu/menfess-server/internal/user/repository/mongodb"
	user "github.com/aeramu/menfess-server/internal/user/service"
)

func newNotifService(db *mongo.Database) notif.Service {
	repo := notifRepo.NewRepository(db)
	expoClient := notifExpoClient.NewClient()
	return notif.NewService(repo, expoClient)
}

func newUserService(db *mongo.Database) user.Service {
	repo := userRepo.NewRepository(db)
	return user.NewService(repo)
}

func newPostService(db *mongo.Database, userService user.Service, notifService notif.Service) post.Service {
	repo := postRepo.NewRepository(db)
	notifClient := postNotifClient.NewClient(notifService, userService)
	return post.NewService(repo, notifClient)
}

func newAuthService(db *mongo.Database, userService user.Service, notifService notif.Service) auth.Service {
	repo := authRepo.NewRepository(db)
	userClient := authUserClient.NewClient(userService)
	notifClient := authNotifClient.NewClient(notifService)
	return auth.NewService(repo, userClient, notifClient)
}

func main() {
	db := mongodb.NewDatabase()

	notifService := newNotifService(db)
	userService := newUserService(db)
	postService := newPostService(db, userService, notifService)
	authService := newAuthService(db, userService, notifService)

	handler := graphql.NewHandler(postService, authService, userService)
	http.Handle("/", handler)

	playgorundHandler := playground.Handler("Playground", "/")
	http.HandleFunc("/playground", playgorundHandler)

	port := getPort()
	log.Println("Server ready at " + port)
	log.Println("Playground ready at " + port + "/playground")
	log.Fatal(http.ListenAndServe(port, nil))
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	return ":" + port
}
