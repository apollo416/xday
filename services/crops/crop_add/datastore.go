package main

import (
	"fmt"

	"github.com/apollo416/xday/pkg/crops"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type datastore struct {
	dynamocli *dynamodb.DynamoDB
}

func (d *datastore) Get(id string) (crops.Crop, error) {
	return crops.Crop{}, nil
}

func (d *datastore) Put(c crops.Crop) error {
	av, err := dynamodbattribute.MarshalMap(c)
	if err != nil {
		return fmt.Errorf("got error marshalling item: %v", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("crops"),
	}

	_, err = d.dynamocli.PutItem(input)
	if err != nil {
		return fmt.Errorf("got error calling PutItem: %v", err)
	}

	return nil
}

func (d *datastore) Update(c crops.Crop) error {
	return nil
}
