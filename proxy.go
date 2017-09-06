package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func handler(w http.ResponseWriter, r *http.Request) {

	// svc := dynamodb.New(session.New())
	svc := dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"sessionId": {
				S: aws.String("bb043bb6-2118-47fb-8206-1539be03a57f"),
			},
		},
		TableName: aws.String("filmstruckdev_dev_EngineSession"),
	}
	result, err := svc.GetItem(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)

	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = "127.0.0.1:8081"
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/receiveMsg", handler)
	panic(http.ListenAndServe(":8080", nil))

}
