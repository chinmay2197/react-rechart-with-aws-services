# react-rechart-with-aws-services

AWS services used:
* Dynamodb and stream
* Elastic Search
* CloudWatch Scheduler

To populate dynamodb with multiTimeline.csv records, Run command
* go run csvtoddb.go

Below are golang lambda handler which will be used to generate Dynamodb stream and sync with ElasticSearch
* updateddb.go -> It will randomly update 100 records from the db
* essync.go -> It will push any change in the db to ES using stream events
* essearch.go -> It will run elasticsearch query and provide json response for rechart


To install node modules deps inside borderfreechart/
* npm install

To run react app and visuliaze Line chart, run below command inside borderfreechart/
* npm start

which will render a Line chart graph in your browser https://borderfreecs.s3.ap-south-1.amazonaws.com/index.html


To generate production build, run below command inside borderfreechart/
* npm run build

Reference:
https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/using-dynamodb-with-go-sdk.html
https://medium.com/@kevinlohier.kl/how-to-fetch-apis-in-react-and-effectively-use-data-responses-to-create-graphs-using-recharts-5a4eea4b5184
https://github.com/olivere/elastic
https://github.com/serverless/examples
