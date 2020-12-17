package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

type ranking struct {
	Company string  `json:"company"`
	Image   string  `json:"image"`
	Matches int     `json:"matches"`
	Wins    int     `json:"wins"`
	Losses  int     `json:"losses"`
	Winrate float64 `json:"winrate"`
	Ranking int     `json:"ranking"`
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return show(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	rks, err := getItems()
	if err != nil {
		return serverError(err)
	}
	if rks == nil {
		return clientError(http.StatusNotFound)
	}

	sort.Slice(rks, func(i, j int) bool {
		return rks[i].Winrate > rks[j].Winrate
	})

	currentRank := 1

	for i := range rks {
		rk := rks[i]
		rk.Ranking = currentRank
		currentRank++
		rks[i] = rk
	}

	js, err := json.Marshal(rks)
	if err != nil {
		return serverError(err)
	}

	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
		Headers:    headers,
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(router)
}
