package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(context context.Context, event events.DynamoDBStreamRecord) {
	fmt.Println(context)
	fmt.Println(event)
}

func main() {
	lambda.Start(Handler)
}
