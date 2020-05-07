package controllers

import (
	"github.com/vlad-m-r/checker/api/models"
	"gopkg.in/yaml.v2"
	"log"
)

func NewEmailClient(yamlContent []byte) *EmailClient {
	emailClient := EmailClient{}
	emailClient.UnmarshalYaml(yamlContent)
	return &emailClient
}

type EmailClient struct {
	Yaml models.EmailClientConfig `yaml:"EmailClient"`
}

func (e *EmailClient) UnmarshalYaml(yamlContent []byte) {
	err := yaml.Unmarshal(yamlContent, &e)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}
}
