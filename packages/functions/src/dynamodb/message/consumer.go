package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(context context.Context, event events.DynamoDBEvent) {
	fmt.Println(context)
	fmt.Printf("Context Type: %T\n", context)
	fmt.Println(event)
	fmt.Printf("Event Type: %T\n", event.Records)
}

func main() {
	lambda.Start(Handler)
}
