package main

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var tableName = "SerialHostLookup"
var region = "us-east-1"

func awsDynamo(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query()["deviceId"][0]

	var host = getHost(deviceID)
	fmt.Println("get host: ", host)

	deleteHost(deviceID)
	setHost(deviceID)
}

func getHost(serial string) string {
	host := ""
	svc := dynamodb.New(session.New(), aws.NewConfig().WithRegion(region))
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(serial),
			},
		},
		TableName: aws.String(tableName),
	}
	result, err := svc.GetItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Code(), aerr.Error())
		} else {
			// Print the error, cast err to awserr.Error to get the Code and Message from an error.
			fmt.Println(err.Error())
		}
	} else {
		if result.Item != nil {
			host = *result.Item["host"].S
		}
	}
	return host
}

func deleteHost(serial string) {
	svc := dynamodb.New(session.New(), aws.NewConfig().WithRegion(region))
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(serial),
			},
		},
		TableName: aws.String(tableName),
	}
	result, err := svc.DeleteItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Code(), aerr.Error())
		} else {
			// Print the error, cast err to awserr.Error to get the Code and Message from an error.
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(result)
	}
}

func setHost(serial string) {
	host := "127.0.0.1:8081"
	svc := dynamodb.New(session.New(), aws.NewConfig().WithRegion(region))
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(serial),
			},
			"host": {
				S: aws.String(host),
			},
		},
		TableName: aws.String(tableName),
	}
	result, err := svc.PutItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Code(), aerr.Error())
		} else {
			// Print the error, cast err to awserr.Error to get the Code and Message from an error.
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(result)
	}
}

func main() {
	http.HandleFunc("/awsDynamo", awsDynamo)
	panic(http.ListenAndServe(":8080", nil))
}
