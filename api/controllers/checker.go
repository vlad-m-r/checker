package controllers

import (
	"github.com/vlad-m-r/checker/api/models"
	"gopkg.in/yaml.v2"
	"log"
)

func NewChecksController(yamlContent []byte) *ChecksController {
	checkController := ChecksController{}
	checkController.UnmarshalYaml(yamlContent)
	return &checkController
}

type CheckResult struct {
	Name          string
	URL           string
	err           []error
	seconds       float64
	ErrorOccurred bool
}

type ChecksController struct {
	Yaml             []models.Check `yaml:"Checks"`
	results          []CheckResult
	ErrorOccurred    bool
	concurrencyLimit int
}

func (c *ChecksController) UnmarshalYaml(yamlContent []byte) {
	err := yaml.Unmarshal(yamlContent, &c)
	if err != nil {
		log.Fatalf("NotifiersController: Unmarshal error: %v", err)
	}
}

func (c *ChecksController) getCheckResults() []*CheckResult {
	var results []*CheckResult

	sChan := make(chan struct{}, c.concurrencyLimit)
	rChan := make(chan *CheckResult)

	defer func() {
		close(sChan)
		close(rChan)
	}()

	// Run goroutines with channels: buffered and unbuffered
	for _, check := range c.Yaml {
		go RequestControllerFactory(check, sChan, rChan)
	}

	// Start reading results from the result channel
	for range c.Yaml {
		r := <-rChan
		// Dump result
		results = append(results, r)
		if r.ErrorOccurred {
			c.ErrorOccurred = true
		}
	}

	return results
}
