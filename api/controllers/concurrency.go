package controllers

import (
	"github.com/vlad-m-r/checker/api/models"
	"gopkg.in/yaml.v2"
	"log"
)

func NewConcurrencyController(yamlContent []byte) *ConcurrencyController {
	concurrencyController := ConcurrencyController{}
	concurrencyController.UnmarshalYaml(yamlContent)
	return &concurrencyController
}

type ConcurrencyController struct {
	Yaml models.ConcurrencyConfig `yaml:"Concurrency"`
}

func (c *ConcurrencyController) UnmarshalYaml(yamlContent []byte) {
	err := yaml.Unmarshal(yamlContent, &c)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}
}
