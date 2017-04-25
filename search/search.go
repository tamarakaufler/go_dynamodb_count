package search

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var awsAccessKeyID, awsSecretAccessKey string

// Result is a data structure,that will be used for display on the web page
type Result struct {
	Table string
	Count int64
	Error error
}

// Set up the AWS credentials if required
func init() {
	flag.StringVar(&awsAccessKeyID, "aws_a", os.Getenv("AWS_ACCESS_KEY_ID"), "AWS access key credential")
	flag.StringVar(&awsSecretAccessKey, "aws_s", os.Getenv("AWS_SECRET_ACCESS_KEY"), "AWS secret key credential")

	//	if awsAccessKeyID != "" || awsSecretAccessKey == "" {
	//		log.Panic("AWS credentials missing")
	//	}

	os.Setenv("AWS_ACCESS_KEY_ID", awsAccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", awsSecretAccessKey)
}

// The date input needs to satisfy the dynamoDb format
func checkDate(date, tm string) (string, error) {
	datetime := date + "T" + tm
	format := "2006-01-02T15:04:05.000Z"
	_, err := time.Parse(format, datetime)

	if err != nil {
		return "", errors.New("Start date and time must be provided and be yyy-MM-dd and hh:mm:ss")
	}

	return datetime, nil
}

// Process the form data
func getSearchTerms(r *http.Request) (string, map[string]string, error) {
	flag.Parse()
	r.ParseForm()

	var fromDatetime, toDatetime string
	var err error

	attributes := make(map[string]string)

	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	table := r.Form.Get("table")
	if table == "" {
		return "", nil, err
	}

	namespace := r.Form.Get("namespace")
	if namespace == "" {
		return "", nil, err
	}

	fromDate := r.Form.Get("from_date")
	fromTime := r.Form.Get("from_time")
	fromDatetime, err = checkDate(fromDate, fromTime)
	if err != nil {
		return "", nil, err
	}

	toDate := r.Form.Get("to_date")
	toTime := r.Form.Get("to_time")
	toDatetime, err = checkDate(toDate, toTime)
	if err != nil {
		return "", nil, err
	}

	attributes["namespace"] = namespace
	attributes["from_date"] = fromDatetime
	attributes["to_date"] = toDatetime

	return table, attributes, nil

}

// Set up the DynamoDB client for performing the scan
func getDynamoDbClient() *dynamodb.DynamoDB {

	endpoint := "http://localhost:8000"
	region := "eu-west-1"
	awsCfg := &aws.Config{
		Endpoint: &endpoint,
		Region:   &region,
	}
	sess := session.Must(session.NewSession(awsCfg))

	return dynamodb.New(sess)
}

// Do the scanning and get the number of records satisfying the search criteria
func getScanResult(client *dynamodb.DynamoDB, table string, attributes map[string]string) (*Result, error) {

	dbTable := "nectar-service-" + table + "-production"

	var dateField string

	switch table {
	case "order":
		dateField = "createDate"
	case "member":
		dateField = "date"
	}

	namespace := attributes["namespace"]
	fromDate := attributes["from_date"]
	toDate := attributes["to_date"]

	// Set up filtering
	params := &dynamodb.ScanInput{
		TableName:        aws.String(dbTable),
		FilterExpression: aws.String("namespace = :namespace and #dt between :from_date and :to_date"),
		ExpressionAttributeNames: map[string]*string{
			"#dt": aws.String(dateField),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":namespace": &dynamodb.AttributeValue{S: aws.String(namespace)},
			":from_date": &dynamodb.AttributeValue{S: aws.String(fromDate)},
			":to_date":   &dynamodb.AttributeValue{S: aws.String(toDate)},
		},
	}

	// Do the scanning
	var count int64
	for {
		result, err := client.Scan(params)
		if err != nil {
			return nil, err
		}

		count = count + *result.Count

		if result.LastEvaluatedKey != nil {
			params.ExclusiveStartKey = result.LastEvaluatedKey
		} else {
			break
		}
	}

	// Store the result for web page display
	sr := &Result{
		Table: table,
		Count: count,
	}

	return sr, nil
}

// DoDynamoDBScan function provides search results for display on the web page
// Input: *http.Request
// Output: *Result:		&struct{
//							Table string
//							Count int64
//							Error error
//						}
//			error
func DoDynamoDBScan(r *http.Request) (*Result, error) {

	// Extract Form search parameters
	table, attributes, err := getSearchTerms(r)
	if err != nil {
		return nil, err
	}

	dbClient := getDynamoDbClient()

	result, err := getScanResult(dbClient, table, attributes)
	if err != nil {
		return nil, err
	}

	return result, nil
}
