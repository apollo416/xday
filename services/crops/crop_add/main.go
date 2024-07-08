package main

import (
	"context"
	"time"

	"github.com/apollo416/xday/pkg/crops"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

// dyna.PutItem

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	service := getService()

	crop := crops.Crop{
		Id:      uuid.New(),
		Created: time.Now(),
	}

	if err := service.Add(crop); err != nil {
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 201}, nil
}

func main() {
	lambda.Start(handler)
}
