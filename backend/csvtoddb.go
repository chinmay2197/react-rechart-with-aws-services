package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Dynamodb Record
type Item struct {
	ID       int
	Time     string
	Interest int
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func main() {
	records := readCsvFile("./multiTimeline.csv")

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

	for i, record := range records {
		interestRate, err := strconv.Atoi(record[1])
		item := Item{
			ID:       i,
			Time:     record[0],
			Interest: interestRate,
		}

		av, err := dynamodbattribute.MarshalMap(item)
		if err != nil {
			log.Fatal("Error occurred while marshalling new record", err)
			return
		}
		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}
		ok, err := svc.PutItem(input)
		if err != nil {
			log.Fatal("Error occurred while creating new record in Dynamodb", err)
			return
		}
		log.Print("Record created with ID=", item.ID, " status=", ok)
	}
}
