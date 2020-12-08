package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/olivere/elastic"
	"github.com/serverless/examples/aws-golang-dynamo-stream-to-elasticsearch/dstream"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func handler(e events.DynamoDBEvent) error {
	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""),
	},
	)
	var dynamoSvc = dynamodb.New(awsSession)
	var esclient = new(dstream.Elasticsearch)

	var item map[string]events.DynamoDBAttributeValue
	log.Print("ES sync starting...")
	for _, v := range e.Records {
		switch v.EventName {
		case "INSERT":
			fallthrough
		case "MODIFY":
			tableName := strings.Split(v.EventSourceArn, "/")[1]
			item = v.Change.NewImage

			details, err := (&dstream.DynamoDetails{
				DynamoDBAPI: dynamoSvc,
			}).Get(tableName)

			if err != nil {
				log.Fatal("Unable to get record details from dynamodb stream", err)
				return err
			}

			svc, err := elastic.NewClient(
				elastic.SetSniff(false),
				elastic.SetURL(fmt.Sprintf("https://%s", os.Getenv("ELASTICSEARCH_URL"))),
				elastic.SetBasicAuth(os.Getenv("ELASTICSEARCH_USERNAME"), os.Getenv("ELASTICSEARCH_PASSWORD")),
			)
			if err != nil {
				log.Fatal("Error occurred while connecting to ElasticSearch", err)
				return err
			}
			esclient.Client = svc
			resp, err := esclient.Update(details, item)
			if err != nil {
				log.Fatal("Error occurred while pushing to ElasticSearch", err)
				return err
			}
			log.Print("Result=", resp.Result)
		default:
		}
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
