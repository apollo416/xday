package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// ! dyna.PutItem
// ! api.POST 201
// ! api.POST 500

type Item struct {
	Year   int
	Title  string
	Plot   string
	Rating float64
}

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))

var tableName = "test_table"
var dyna = dynamodb.New(sess)

func HandleRequest(ctx context.Context) error {
	dyna.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"Year": {
				N: aws.String("year"),
			},
			"Title": {
				S: aws.String("title"),
			},
		},
	})

	item := Item{
		Year:   2015,
		Title:  "The Big New Movie",
		Plot:   "Nothing happens at all.",
		Rating: 0.0,
	}

	av, _ := dynamodbattribute.MarshalMap(item)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err := dyna.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
