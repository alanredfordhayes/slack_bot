package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(context context.Context, event events.DynamoDBStreamRecord) {
	fmt.Println(context)
	fmt.Printf("Context Type: %T\n", context)
	fmt.Println(event)
	fmt.Printf("Event Record Type: %T\n", event.NewImage)
}

func main() {
	lambda.Start(Handler)
}
