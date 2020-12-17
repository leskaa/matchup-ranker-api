package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
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

type matchup struct {
	Company1         string `json:"company1"`
	Company2         string `json:"company2"`
	Image1           string `json:"image1"`
	Image2           string `json:"image2"`
	VerificationCode string `json:"verificationCode"`
	Voted            string `json:"voted"`
}

type winBody struct {
	VerificationCode string `json:"verificationCode"`
	Winner           int    `json:"winner"`
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return show(req)
	case "POST":
		return create(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func show(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	rks, err := getRankings()
	if err != nil {
		return serverError(err)
	}
	if rks == nil {
		return clientError(http.StatusNotFound)
	}

	// Create two different random numbers and get the companies at those indexes
	rand.Seed(time.Now().UnixNano())
	p := rand.Perm(len(rks))
	rk1 := rks[p[0]]
	rk2 := rks[p[1]]

	mu := new(matchup)
	mu.Company1 = rk1.Company
	mu.Company2 = rk2.Company
	mu.Image1 = rk1.Image
	mu.Image2 = rk2.Image
	mu.VerificationCode = uuid.New().String()
	mu.Voted = "undecided"

	err = putMatchup(mu)
	if err != nil {
		return serverError(err)
	}

	js, err := json.Marshal(mu)
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

func create(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["content-type"] != "application/json" && req.Headers["Content-Type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	wb := new(winBody)
	err := json.Unmarshal([]byte(req.Body), wb)
	if err != nil {
		return clientError(http.StatusUnprocessableEntity)
	}
	if wb.Winner != 1 && wb.Winner != 2 {
		return clientError(http.StatusBadRequest)
	}

	mu, err := getMatchup(wb.VerificationCode)
	if err != nil {
		return serverError(err)
	}
	if mu == nil {
		return clientError(http.StatusForbidden)
	}
	if strings.Compare(mu.Voted, "decided") == 0 {
		return clientError(http.StatusForbidden)
	}

	rk1, err := getRanking(mu.Company1)
	if err != nil {
		return serverError(err)
	}
	if rk1 == nil {
		return clientError(http.StatusNotFound)
	}

	rk2, err := getRanking(mu.Company2)
	if err != nil {
		return serverError(err)
	}
	if rk2 == nil {
		return clientError(http.StatusNotFound)
	}

	err = updateMatchup(wb.VerificationCode)
	if err != nil {
		return serverError(err)
	}

	if wb.Winner == 1 {
		err = updateRanking(rk1.Company, true, float64(rk1.Wins+1)/float64(rk1.Matches+1))
		if err != nil {
			return serverError(err)
		}
		err = updateRanking(rk2.Company, false, float64(rk2.Wins)/float64(rk2.Matches+1))
		if err != nil {
			return serverError(err)
		}
	} else {
		err = updateRanking(rk1.Company, false, float64(rk1.Wins)/float64(rk1.Matches+1))
		if err != nil {
			return serverError(err)
		}
		err = updateRanking(rk2.Company, true, float64(rk2.Wins+1)/float64(rk2.Matches+1))
		if err != nil {
			return serverError(err)
		}
	}

	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
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
	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(router)
}
