package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))

// Get an unsorted list of companies from dynamodb by rankings
func getItems() ([]ranking, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("prestige-companies"),
	}

	result, err := db.Scan(input)
	if err != nil {
		return nil, err
	}
	if result.Items == nil {
		return nil, nil
	}

	rks := make([]ranking, 0)

	for _, i := range result.Items {
		rk := new(ranking)
		err = dynamodbattribute.UnmarshalMap(i, rk)
		if err != nil {
			return nil, err
		}
		rks = append(rks, *rk)
	}

	return rks, nil
}
