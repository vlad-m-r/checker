package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vlad-m-r/checker/api/controllers"
	"github.com/vlad-m-r/checker/api/models"
)

func HandleRequest(ctx context.Context, name models.LambdaEvent) (string, error) {
	controllers.RunChecks("config.yaml")
	return "Done", nil
}

func main() {
	lambda.Start(HandleRequest)
}
