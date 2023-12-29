package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type program_secrets struct {
	OPEN_AI_API_KEY      string `json:"OPEN_AI_API_KEY"`
	OPEN_AI_ORG_ID       string `json:"OPEN_AI_ORG_ID"`
	OPEN_AI_ASSISTANT_ID string `json:"OPEN_AI_ASSISTANT_ID"`
}

func http_client() *http.Client {
	client := &http.Client{}
	return client
}

func http_get_raw(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return req, err
}

func http_end_request(client *http.Client, req *http.Request) (string, error) {
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	return string(b), err
}

func http_get(url string, headers map[string][]string) (string, error) {
	client := http_client()
	req, err := http_get_raw(url)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header = http.Header(headers)
	str, err := http_end_request(client, req)
	return str, err
}

func GetSecret(secretName string) (string, error) {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	conn := secretsmanager.NewFromConfig(config)
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	}
	result, err := conn.GetSecretValue(context.TODO(), input)
	if err != nil {
		return "", err
	}
	return *result.SecretString, err
}

func slack_bot_secrets() program_secrets {
	var secretName string = "slack_bot"
	secret, err := GetSecret(secretName)
	if err != nil {
		log.Fatalf("unable to get secret, %v", err)
	}
	secret_data := program_secrets{}
	json.Unmarshal([]byte(secret), &secret_data)
	return secret_data
}

func open_ai_assistant_headers(secrets program_secrets) map[string][]string {
	var headers = map[string][]string{
		"Authorization":       {fmt.Sprintf("Bearer %s", secrets.OPEN_AI_API_KEY)},
		"OpenAI-Organization": {secrets.OPEN_AI_ORG_ID},
		"OpenAI-Beta":         {"assistants=v1"},
	}
	return headers
}

func open_ai_retrieve_assistant(secrets program_secrets) (string, error) {
	headers := open_ai_assistant_headers(secrets)
	res, err := http_get(fmt.Sprintf("https://api.openai.com/v1/assistants/%s", secrets.OPEN_AI_ASSISTANT_ID), headers)
	return res, err
}

func Handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	secrets := slack_bot_secrets()
	my_assistant, err := open_ai_retrieve_assistant(secrets)
	if err != nil {
		log.Fatalf("unable to get secret, %v", err)
	} else {
		fmt.Printf("%s", my_assistant)
	}
	return events.APIGatewayProxyResponse{
		Body:       "Hello, World! Your request was received at " + request.RequestContext.Time + ". And here is my assistant information:  " + my_assistant,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
