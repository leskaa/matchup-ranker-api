package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))

// Get an unsorted list of companies from dynamodb by rankings
func getRankings() ([]ranking, error) {
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

	var rks []ranking

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

// Get a company from dynamodb by company name
func getRanking(company string) (*ranking, error) {
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

// Update a company ranking by company name
func updateRanking(company string, result bool, winrate float64) error {

	// Set win or loss increment
	w := "0"
	l := "1"
	if result {
		w = "1"
		l = "0"
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("prestige-companies"),
		Key: map[string]*dynamodb.AttributeValue{
			"Company": {
				S: aws.String(company),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				N: aws.String(fmt.Sprintf("%f", winrate)),
			},
			":w": {
				N: aws.String(w),
			},
			":l": {
				N: aws.String(l),
			},
			":m": {
				N: aws.String("1"),
			},
		},
		UpdateExpression: aws.String("SET Matches = Matches + :m, Wins = Wins + :w, Losses = Losses + :l, Winrate = :r"),
	}

	_, err := db.UpdateItem(input)
	return err
}

// Get a matchup by VerificationCode
func getMatchup(verificationCode string) (*matchup, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String("prestige-matchups"),
		Key: map[string]*dynamodb.AttributeValue{
			"VerificationCode": {
				S: aws.String(verificationCode),
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

	mu := new(matchup)
	err = dynamodbattribute.UnmarshalMap(result.Item, mu)
	if err != nil {
		return nil, err
	}

	return mu, nil
}

// Add a new matchup to the table
func putMatchup(mu *matchup) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("prestige-matchups"),
		Item: map[string]*dynamodb.AttributeValue{
			"Company1": {
				S: aws.String(mu.Company1),
			},
			"Company2": {
				S: aws.String(mu.Company2),
			},
			"Image1": {
				S: aws.String(mu.Image1),
			},
			"Image2": {
				S: aws.String(mu.Image2),
			},
			"VerificationCode": {
				S: aws.String(mu.VerificationCode),
			},
			"Voted": {
				S: aws.String(mu.Voted),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}

// Update a matchup to voted state
func updateMatchup(verificationCode string) error {

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("prestige-matchups"),
		Key: map[string]*dynamodb.AttributeValue{
			"VerificationCode": {
				S: aws.String(verificationCode),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":d": {
				S: aws.String("decided"),
			},
		},
		UpdateExpression: aws.String("SET Voted = :d"),
	}

	_, err := db.UpdateItem(input)
	return err
}
