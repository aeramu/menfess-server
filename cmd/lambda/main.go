package main

//
//import (
//	"context"
//	"encoding/json"
//	"github.com/aeramu/menfess-server/internal/gateway/server"
//	postRepo "github.com/aeramu/menfess-server/internal/post/repository"
//	post "github.com/aeramu/menfess-server/internal/post/service"
//	roomRepo "github.com/aeramu/menfess-server/room/repository"
//	room "github.com/aeramu/menfess-server/room/service"
//	"github.com/aws/aws-lambda-go/events"
//	"github.com/aws/aws-lambda-go/lambda"
//)
//
//func main() {
//	lambda.Start(handler)
//}
//
//func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//	//convert request body to json
//	var parameter struct {
//		Query         string
//		OperationName string
//		Variables     map[string]interface{}
//	}
//	json.Unmarshal([]byte(request.Body), &parameter)
//
//	//add token from header
//	context := context.WithValue(ctx, "request", request.Headers)
//
//	postRepo := postRepo.NewRepository()
//	postService := post.NewService(postRepo)
//
//	roomRepo := roomRepo.NewRepository()
//	roomService := room.NewService(roomRepo)
//
//	handler := server.NewServer(context, postService, roomService)
//
//	response := handler.Response(ctx, parameter.Query, parameter.OperationName, parameter.Variables)
//	responseJSON, _ := json.Marshal(response)
//
//	//response
//	return events.APIGatewayProxyResponse{
//		Headers: map[string]string{
//			"Access-Control-Allow-Origin":  "*",
//			"Access-Control-Allow-Headers": "Content-Type",
//		},
//		Body:       string(responseJSON),
//		StatusCode: 200,
//	}, nil
//}
