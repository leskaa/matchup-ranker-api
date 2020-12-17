package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))

// Get a company from dynamodb by company name
func getItem(company string) (*ranking, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("prestige-companies"),
		Key: map[string]*dynamodb.AttributeValue{
			"Company": {
				S: aws.String(company),
			},
		},
	}

	result, err := db.GetItem(input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	rk := new(ranking)
	err = dynamodbattribute.UnmarshalMap(result.Item, rk)
	if err != nil {
		return nil, err
	}

	return rk, nil
}
