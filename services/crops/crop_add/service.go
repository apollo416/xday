package main

import (
	"github.com/apollo416/xday/pkg/crops"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var service *crops.CropsService

func getService() *crops.CropsService {
	if service == nil {
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		dynamocli := dynamodb.New(sess)
		datastore := &datastore{dynamocli: dynamocli}
		service = crops.NewCropsService(datastore)
	}
	return service
}
