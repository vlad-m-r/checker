package controllers

import (
	"github.com/vlad-m-r/checker/api/utils"
	"log"
)

func RunChecks(yamlFile string) {

	// Load yaml file
	yamlContent := utils.ReadYaml(yamlFile)

	NewHttpClient(yamlContent)
	awsClient := NewAwsClient(yamlContent)
	emailClient := NewEmailClient(yamlContent)
	concurrencyController := NewConcurrencyController(yamlContent)

	checksController := NewChecksController(yamlContent)
	checksController.concurrencyLimit = concurrencyController.Yaml.Limit

	// Run checks
	results := checksController.getCheckResults()

	// Notify
	if checksController.ErrorOccurred {
		awsClient.sendSESMail(results, emailClient)
	} else {
		log.Println("No errors found")
	}

}
