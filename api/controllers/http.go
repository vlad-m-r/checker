package controllers

import (
	"crypto/tls"
	"github.com/vlad-m-r/checker/api/models"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"time"
)

func NewHttpClient(yamlContent []byte) *HttpClient {
	httpClient := HttpClient{}
	httpClient.UnmarshalYaml(yamlContent)
	httpClient.setupClient()
	return &httpClient
}

type HttpClient struct {
	Yaml models.HttpClientConfig `yaml:"HttpClient"`
}

func (hc *HttpClient) UnmarshalYaml(yamlContent []byte) {
	err := yaml.Unmarshal(yamlContent, &hc)
	if err != nil {
		log.Fatalf("HttpClient: Unmarshal error: %v", err)
	}
}

func (hc HttpClient) setupClient() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: hc.Yaml.InsecureSkipVerify}
	http.DefaultTransport.(*http.Transport).ResponseHeaderTimeout = time.Second * hc.Yaml.ResponseHeaderTimeout
}
