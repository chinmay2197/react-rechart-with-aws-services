package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	elastic "gopkg.in/olivere/elastic.v7"
)

// Item Record struct
type Item struct {
	ID       int    `json:"ID"`
	Time     string `json:"Time"`
	Interest int    `json:"Interest"`
}

// Query Response struct
type RecordResponse struct {
	StatusCode int64  `json:statuscode`
	Message    string `json:string`
	Count      int    `json:conut`
	Items      []Item `json:items`
}

// Get ES Client
func GetESClient() (*elastic.Client, error) {

	client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("https://%s", os.Getenv("ELASTICSEARCH_URL"))),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetBasicAuth(os.Getenv("ELASTICSEARCH_USERNAME"), os.Getenv("ELASTICSEARCH_PASSWORD")),
	)

	return client, err

}
func searchHandler() (RecordResponse, error) {

	ctx := context.Background()
	esclient, err := GetESClient()
	if err != nil {
		log.Fatal("Error in initializing ES client", err)
		panic("Client fail")
	}

	var records []Item

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("_index", "interestovertime"))

	queryStr, err1 := searchSource.Source()
	queryJs, err2 := json.Marshal(queryStr)

	if err1 != nil || err2 != nil {
		log.Fatal("Error occurred during query marshal=", err1, err2)
		panic("Client failed while processing query")
	}
	log.Print("ElasticSerach query=", string(queryJs))

	searchService := esclient.Search().Index("interestovertime").SearchSource(searchSource).Sort("Time", false).Size(100)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		log.Fatal("Error occurred while processing search query", err)
		panic("Query returned error")
	}

	for _, hit := range searchResult.Hits.Hits {
		var record Item
		err := json.Unmarshal(hit.Source, &record)
		if err != nil {
			log.Fatal("Error occurred while json unmarshalling of record ", err)
		}

		records = append(records, record)
	}
	jsonInfo, err := json.Marshal(records)
	if err != nil {
		log.Fatal("Error occurred while json marshalling of record ", err)
		panic("Marshalling of records failed")
	}

	resp := RecordResponse{
		StatusCode: 200,
		Items:      records,
		Count:      len(records),
		Message:    "200 success",
	}
	return resp, nil

}

func main() {
	lambda.Start(searchHandler)
}
