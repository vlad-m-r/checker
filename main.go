package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/vlad-m-r/checker/api/controllers"
	"github.com/vlad-m-r/checker/api/models"
	"log"
	"time"
)

func HandleRequest(ctx context.Context, name models.LambdaEvent) (string, error) {
	main()
	return "Done", nil
}

func main() {
	start := time.Now()

	log.Println("Tests started")
	controllers.RunChecks("config.yaml")
	log.Println("Tests passed")
	elapsed := time.Since(start).Seconds()
	log.Printf("Tests took %f\n", elapsed)
}

func mainCloud() {
	lambda.Start(HandleRequest)
}

func mainLocal() {
	main()
}
