package main

import (
	"context"
	"encoding/json"

	graphql "github.com/aeramu/menfess-server/implementation/graphql.handler"
	mongodb "github.com/aeramu/menfess-server/implementation/mongodb/repository"
	"github.com/aeramu/menfess-server/usecase"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//convert request body to json
	var parameter struct {
		Query         string
		OperationName string
		Variables     map[string]interface{}
	}
	json.Unmarshal([]byte(request.Body), &parameter)

	//add token from header
	context := context.WithValue(ctx, "request", request.Headers)

	repository := mongodb.New()
	defer mongodb.Disconnect()
	interactor := usecase.InteractorConstructor{
		Repository: repository,
	}.New()
	handler := graphql.New(context, interactor)

	response := handler.Response(context, parameter.Query, parameter.OperationName, parameter.Variables)
	responseJSON, _ := json.Marshal(response)

	//response
	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		Body:       string(responseJSON),
		StatusCode: 200,
	}, nil
}
