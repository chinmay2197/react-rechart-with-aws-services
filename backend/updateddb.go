package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Dynamodb Item
type Item struct {
	ID       int
	Time     string
	Interest int
}

func generateRandomID(min int, max int) int {
	return min + rand.Intn(max-min)
}

func lambdaHandler() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""),
	},
	)

	if err != nil {
		log.Fatal("Unable to create AWS session", err)
		return
	}

	svc := dynamodb.New(sess)
	tableName := "InterestOverTime"
	rand.Seed(time.Now().UTC().UnixNano())

	params := &dynamodb.ScanInput{
		ProjectionExpression: aws.String("ID"),
		TableName:            aws.String(tableName),
	}

	result, err := svc.Scan(params)
	if err != nil {
		log.Fatal("Error occurred while running Dynamodb query ", err)
		return
	}

	records := len(result.Items)
	log.Print("Total Records in Dynamodb =", records)

	for i := 0; i < 100; i++ {
		randomID := generateRandomID(0, records-1)
		currTime := time.Now()
		randomInterest := generateRandomID(0, 100)
		updateTime := currTime.Format("2006-01-02T15:04:05Z")
		input := &dynamodb.UpdateItemInput{
			ExpressionAttributeNames: map[string]*string{
				"#I":  aws.String("Interest"),
				"#UT": aws.String("UpdateTime"),
			},
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":r": {
					N: aws.String(strconv.Itoa(randomInterest)),
				},
				":s": {
					S: aws.String(updateTime),
				},
			},
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"ID": {
					N: aws.String(strconv.Itoa(randomID)),
				},
			},
			ReturnValues:     aws.String("UPDATED_NEW"),
			UpdateExpression: aws.String("set #I = :r, #UT = :s"),
		}

		_, err = svc.UpdateItem(input)
		if err != nil {
			log.Fatal("Error occurred while updating record in Dynamodb with record ID=", randomID)
		}
		log.Print("Record Updated. Record ID=", randomID)
	}

}

func main() {
	lambda.Start(lambdaHandler)
}
