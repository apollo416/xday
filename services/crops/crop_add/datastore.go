package main

import (
	"fmt"
	"time"

	"github.com/apollo416/xday/pkg/crops"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type UUID uuid.UUID

func (u UUID) String() string {
	return uuid.UUID(u).String()
}

func (u UUID) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	s := u.String()
	av.S = &s
	return nil
}

func (u *UUID) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	if av.S == nil {
		return nil
	}

	id, err := uuid.Parse(*av.S)
	if err != nil {
		return err
	}

	*u = UUID(id)

	return nil
}

type CropItem struct {
	ID             UUID      `json:"id" dynamodbav:"id"`
	Status         string    `json:"status" dynamodbav:"status"`
	Cultivar       UUID      `json:"cultivar" dynamodbav:"cultivar"`
	CultivarStart  time.Time `json:"cultivar_start" dynamodbav:"cultivar_start"`
	CultivarEnd    time.Time `json:"cultivar_end" dynamodbav:"cultivar_end"`
	MaturationTime int       `json:"maturation_time" dynamodbav:"maturation_time"`
	Created        time.Time `json:"created" dynamodbav:"created"`
	Generation     int       `json:"generation" dynamodbav:"generation"`
}

func FromCrop(c crops.Crop) CropItem {
	return CropItem{
		ID:             UUID(c.ID),
		Status:         c.Status,
		Cultivar:       UUID(c.Cultivar),
		CultivarStart:  c.CultivarStart,
		CultivarEnd:    c.CultivarEnd,
		MaturationTime: c.MaturationTime,
		Created:        c.Created,
		Generation:     c.Generation,
	}
}

type datastore struct {
	dynamocli *dynamodb.DynamoDB
}

func (d *datastore) Get(id uuid.UUID) (crops.Crop, error) {
	return crops.Crop{}, nil
}

func (d *datastore) Put(c crops.Crop) error {
	cropItem := FromCrop(c)

	av, err := dynamodbattribute.MarshalMap(cropItem)
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
