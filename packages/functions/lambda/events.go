package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type url_verification struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

func create_url_verification(request_body string) url_verification {
	url_verification := url_verification{}
	json.Unmarshal([]byte(request_body), &url_verification)
	return url_verification
}

func events_handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	url_verification := create_url_verification(request.Body)
	switch url_verification.Challenge {
	case "":
		return events.APIGatewayProxyResponse{
			Body:       request.Body,
			StatusCode: 200,
		}, nil
	}

	fmt.Println(request.Body)

	return events.APIGatewayProxyResponse{
		Body:       request.Body,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(events_handler)
}
