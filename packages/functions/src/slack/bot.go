package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type slackBotSecrets struct {
	Slack_token              string `json:"SLACK_TOKEN"`
	Slack_verification_token string `json:"SLACK_VERIFICATION_TOKEN"`
	Slack_signing_secret     string `json:"SLACK_SIGNING_SECRET"`
}

type slackDynamoDbItem struct {
	Text       string
	User       string
	Ts         string
	Team       string
	Channel    string
	Event_id   string
	Event_time float64
}

type slackRequestBodyEvent struct {
	Client_msg_id string `json:"client_msg_id"`
	Type          string `json:"type"`
	Text          string `json:"text"`
	User          string `json:"user"`
	Ts            string `json:"ts"`
	Team          string `json:"team"`
	Channel       string `json:"channel"`
	Event_ts      string `json:"event_ts"`
	Channel_type  string `json:"channel_type"`
}

type slackRequestBodyAuthorizations struct {
	Enterprise_id         string `json:"enterprise_id"`
	Team_id               string `json:"team_id"`
	User_id               string `json:"user_id"`
	Is_bot                string `json:"is_bot"`
	Is_enterprise_install string `json:"is_enterprise_install"`
}

type slackRequestBody struct {
	Token                 string                           `json:"token"`
	Team_id               string                           `json:"team_id"`
	Context_team_id       string                           `json:"context_team_id"`
	Context_enterprise_id string                           `json:"context_enterprise_id"`
	Api_app_id            string                           `json:"api_app_id"`
	Event                 slackRequestBodyEvent            `json:"event"`
	Authorizations        []slackRequestBodyAuthorizations `json:"authorizations"`
	Type                  string                           `json:"type"`
	Event_id              string                           `json:"event_id"`
	Event_time            float64                          `json:"event_time"`
	Is_ext_shared_channel string                           `json:"is_ext_shared_channel"`
	Event_context         string                           `json:"event_context"`
	Challenge             string                           `json:"challenge"`
}

func putItemIntoDynamoDb(dynamodbServiceClient *dynamodb.DynamoDB, dynamoDbPutItemInput *dynamodb.PutItemInput) string {
	_, err := dynamodbServiceClient.PutItem(dynamoDbPutItemInput)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}
	return "processed"
}

func createDynamoDbPutItemInput(marshalMap map[string]*dynamodb.AttributeValue, dynamoDbTableName string) *dynamodb.PutItemInput {
	dynamoDbPutItemInput := &dynamodb.PutItemInput{
		Item:      marshalMap,
		TableName: aws.String(dynamoDbTableName),
	}
	return dynamoDbPutItemInput
}

func createDynamoDbTableName(slackRequestBodyJsonEventType string) string {
	var dynamodbTableName string
	switch slackRequestBodyJsonEventType {
	case "message":
		dynamodbTableName = "ahayes-slack-bot-slack_event_message"
	case "app_mention":
		dynamodbTableName = "ahayes-slack-bot-slack_event_app_mention"
	default:
		dynamodbTableName = "ahayes-slack-bot-slack_event_app_not_found"
	}
	return dynamodbTableName
}

func createDynamoDbMarshalMap(item slackDynamoDbItem) map[string]*dynamodb.AttributeValue {
	marshalMap, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Fatalf("Got error marshalling new dynamodb item: %s", err)
	}
	return marshalMap
}

func createDynamoDbItem(slackRequestBodyJson slackRequestBody) slackDynamoDbItem {
	dynamoDbItem := slackDynamoDbItem{
		Text:       slackRequestBodyJson.Event.Text,
		User:       slackRequestBodyJson.Event.User,
		Ts:         slackRequestBodyJson.Event.Ts,
		Team:       slackRequestBodyJson.Event.Team,
		Channel:    slackRequestBodyJson.Event.Channel,
		Event_id:   slackRequestBodyJson.Event_id,
		Event_time: slackRequestBodyJson.Event_time,
	}
	return dynamoDbItem
}

func createDynamoDbClient(awsSession *session.Session) *dynamodb.DynamoDB {
	dynamodbServiceClient := dynamodb.New(awsSession)
	return dynamodbServiceClient
}

func createAwsSession() *session.Session {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return awsSession
}

func compareCalculatedSigningSignatureWithHeaderSigningSignature(calculatedSignatureHash []byte, headerSignatureHash []byte) bool {
	comparisonBool := hmac.Equal(calculatedSignatureHash, headerSignatureHash)
	return comparisonBool
}

func createCalculatedSlackRequestSignatureHashUsingSigningSecretAsKey(slackBotSecretsJsonSlackSigningSecret string, requestSignature []byte, xSlackSignature string) string {
	secretKey := []byte(slackBotSecretsJsonSlackSigningSecret)
	hash := hmac.New(sha256.New, secretKey)
	hash.Write(requestSignature)
	signatureHash := hash.Sum(nil)
	hexEncodedSignatureHash := hex.EncodeToString(signatureHash)
	finalHexEncodedSignatureHash := fmt.Sprintf("v0=%v", hexEncodedSignatureHash)
	return finalHexEncodedSignatureHash
}

func createRequestSignature(requestSignatureVersion string, xSlackRequestTimestamp string, requestBody string) []byte {
	requestSignature := []byte(fmt.Sprintf("%v:%v:%v", requestSignatureVersion, xSlackRequestTimestamp, requestBody))
	return requestSignature
}

func jsonUnmarshalSlackRequestBodyToJson(SlackRequestBody string) slackRequestBody {
	slackRequestBodyJson := slackRequestBody{}
	json.Unmarshal([]byte(SlackRequestBody), &slackRequestBodyJson)
	return slackRequestBodyJson
}

func getSecretsManagerSecretValueOutput(secretsManagerClient *secretsmanager.Client, secretsManagerSecretValueInput *secretsmanager.GetSecretValueInput) *secretsmanager.GetSecretValueOutput {
	secretsManagerSecretValueOutput, err := secretsManagerClient.GetSecretValue(context.TODO(), secretsManagerSecretValueInput)
	if err != nil {
		log.Fatal("unable to get the secrets manager secret value output")
	}
	return secretsManagerSecretValueOutput
}

func getSecretsManagerSecretValueInput(secretId string) *secretsmanager.GetSecretValueInput {
	secretsManagerSecretValueInput := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretId),
		VersionStage: aws.String("AWSCURRENT"),
	}
	return secretsManagerSecretValueInput
}

func createSecretManagerClient() *secretsmanager.Client {
	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	secretsManagerClient := secretsmanager.NewFromConfig(config)
	return secretsManagerClient
}

func jsonUnmarshalSecretsManagerSecretValueOutput(secretsManagerSecretValueOutputSecretString string) slackBotSecrets {
	slackBotSecretsJson := slackBotSecrets{}
	json.Unmarshal([]byte(secretsManagerSecretValueOutputSecretString), &slackBotSecretsJson)
	return slackBotSecretsJson
}

func handler(request events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	slackRequestBodyJson := jsonUnmarshalSlackRequestBodyToJson(request.Body)
	if slackRequestBodyJson.Challenge != "" {
		return events.APIGatewayProxyResponse{
			Body:       request.Body,
			StatusCode: 200,
		}, nil
	}
	secretsManagerClient := createSecretManagerClient()
	secretsManagerSecretValueInput := getSecretsManagerSecretValueInput("slack_bot")
	secretsManagerSecretValueOutput := getSecretsManagerSecretValueOutput(secretsManagerClient, secretsManagerSecretValueInput)
	slackBotSecretsJson := jsonUnmarshalSecretsManagerSecretValueOutput(*secretsManagerSecretValueOutput.SecretString)
	requestSignature := createRequestSignature("v0", request.Headers["x-slack-request-timestamp"], request.Body)
	finalHexEncodedSignatureHash := createCalculatedSlackRequestSignatureHashUsingSigningSecretAsKey(slackBotSecretsJson.Slack_signing_secret, requestSignature, request.Headers["x-slack-signature"])
	comparisonBool := compareCalculatedSigningSignatureWithHeaderSigningSignature([]byte(finalHexEncodedSignatureHash), []byte(request.Headers["x-slack-signature"]))
	var processed string
	if comparisonBool {
		awsSession := createAwsSession()
		dynamodbServiceClient := createDynamoDbClient(awsSession)
		dynamoDbItem := createDynamoDbItem(slackRequestBodyJson)
		marshalMap := createDynamoDbMarshalMap(dynamoDbItem)
		dynamoDbTableName := createDynamoDbTableName(slackRequestBodyJson.Event.Type)
		dynamoDbPutItemInput := createDynamoDbPutItemInput(marshalMap, dynamoDbTableName)
		processed = putItemIntoDynamoDb(dynamodbServiceClient, dynamoDbPutItemInput)
	} else {
		processed = "unprocessed"
	}
	return events.APIGatewayProxyResponse{
		Body:       processed,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
